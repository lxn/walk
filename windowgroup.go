// Copyright 2019 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"sync"
	"unsafe"

	"github.com/lxn/win"
)

// The global window group manager instance.
var wgm windowGroupManager

// windowGroupManager manages window groups for each thread with one or
// more windows.
type windowGroupManager struct {
	mutex  sync.RWMutex
	groups map[uint32]*WindowGroup
}

// Group returns a window group for the given thread ID, if one exists.
// If a group does not already exist it returns nil.
func (m *windowGroupManager) Group(threadID uint32) *WindowGroup {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.groups == nil {
		return nil
	}
	return m.groups[threadID]
}

// CreateGroup returns a window group for the given thread ID. If one does
// not already exist, it will be created.
//
// The group will have its counter incremented as a result of this call.
// It is the caller's responsibility to call Done when finished with the
// group.
func (m *windowGroupManager) CreateGroup(threadID uint32) *WindowGroup {
	// Fast path with read lock
	m.mutex.RLock()
	if m.groups != nil {
		if group := m.groups[threadID]; group != nil {
			m.mutex.RUnlock()
			group.Add(1)
			return group
		}
	}
	m.mutex.RUnlock()

	// Slow path with write lock
	m.mutex.Lock()
	if m.groups == nil {
		m.groups = make(map[uint32]*WindowGroup)
	} else {
		if group := m.groups[threadID]; group != nil {
			// Another caller raced with our lock and beat us
			m.mutex.Unlock()
			group.Add(1)
			return group
		}
	}

	group := newWindowGroup(threadID, m.removeGroup)
	group.Add(1)
	m.groups[threadID] = group
	m.mutex.Unlock()

	return group
}

// removeGroup is called by window groups to remove themselves from
// the manager.
func (m *windowGroupManager) removeGroup(threadID uint32) {
	m.mutex.Lock()
	delete(m.groups, threadID)
	m.mutex.Unlock()
}

// WindowGroup holds data common to windows that share a thread.
//
// Each WindowGroup keeps track of the number of references to
// the group. When the number of references reaches zero, the
// group is disposed of.
type WindowGroup struct {
	refs            int // Tracks the number of windows that rely on this group
	ignored         int // Tracks the number of refs created by the group itself
	threadID        uint32
	completion      func(uint32) // Used to tell the window group manager to remove this group
	removed         bool         // Has this group been removed from its manager? (used for race detection)
	toolTip         *ToolTip
	activeForm      Form
	oleInit         bool
	accPropServices *win.IAccPropServices

	syncMutex           sync.Mutex
	syncFuncs           []func()                   // Functions queued to run on the group's thread
	layoutResultsByForm map[Form]*formLayoutResult // Layout computations queued for application on the group's thread
}

// newWindowGroup returns a new window group for the given thread ID.
//
// The completion function will be called when the group is disposed of.
func newWindowGroup(threadID uint32, completion func(uint32)) *WindowGroup {
	hr := win.OleInitialize()

	return &WindowGroup{
		threadID:            threadID,
		completion:          completion,
		oleInit:             hr == win.S_OK || hr == win.S_FALSE,
		layoutResultsByForm: make(map[Form]*formLayoutResult),
	}
}

// ThreadID identifies the thread that the group is affiliated with.
func (g *WindowGroup) ThreadID() uint32 {
	return g.threadID
}

// Refs returns the current number of references to the group.
func (g *WindowGroup) Refs() int {
	return g.refs
}

// AccessibilityServices returns an instance of CLSID_AccPropServices class.
func (g *WindowGroup) accessibilityServices() *win.IAccPropServices {
	if g.accPropServices != nil {
		return g.accPropServices
	}

	var accPropServices *win.IAccPropServices
	hr := win.CoCreateInstance(&win.CLSID_AccPropServices, nil, win.CLSCTX_ALL, &win.IID_IAccPropServices, (*unsafe.Pointer)(unsafe.Pointer(&accPropServices)))
	if win.FAILED(hr) {
		return nil
	}

	g.accPropServices = accPropServices
	return accPropServices
}

// accPropIds is a static list of accessibility properties user (may) set for a window
// and we should clear when the window is disposed.
var accPropIds = []win.MSAAPROPID{
	win.PROPID_ACC_DEFAULTACTION,
	win.PROPID_ACC_DESCRIPTION,
	win.PROPID_ACC_HELP,
	win.PROPID_ACC_KEYBOARDSHORTCUT,
	win.PROPID_ACC_NAME,
	win.PROPID_ACC_ROLE,
	win.PROPID_ACC_ROLEMAP,
	win.PROPID_ACC_STATE,
	win.PROPID_ACC_STATEMAP,
	win.PROPID_ACC_VALUEMAP,
}

// accClearHwndProps clears all window properties for Dynamic Annotation to release resources.
func (g *WindowGroup) accClearHwndProps(hwnd win.HWND) {
	if g.accPropServices != nil {
		g.accPropServices.ClearHwndProps(hwnd, win.OBJID_CLIENT, win.CHILDID_SELF, accPropIds)
	}
}

// Add changes the group's reference counter by delta, which may be negative.
//
// If the reference counter becomes zero the group will be disposed of.
//
// If the reference counter goes negative Add will panic.
func (g *WindowGroup) Add(delta int) {
	if g.removed {
		panic("walk: add() called on a WindowGroup that has been removed from its manager")
	}

	g.refs += delta
	if g.refs < 0 {
		panic("walk: negative WindowGroup refs counter")
	}
	if g.refs-g.ignored == 0 {
		g.dispose()
	}
}

// Done decrements the group's reference counter by one.
func (g *WindowGroup) Done() {
	g.Add(-1)
}

// Synchronize adds f to the group's function queue, to be executed
// by the message loop running on the the group's thread.
//
// Synchronize can be called from any thread.
func (g *WindowGroup) Synchronize(f func()) {
	g.syncMutex.Lock()
	defer g.syncMutex.Unlock()
	g.syncFuncs = append(g.syncFuncs, f)
}

// synchronizeLayout causes the given layout computations to be applied
// later by the message loop running on the group's thread.
//
// Any previously queued layout computations for the affected form that
// have not yet been applied will be replaced.
//
// synchronizeLayout can be called from any thread.
func (g *WindowGroup) synchronizeLayout(result *formLayoutResult) {
	g.syncMutex.Lock()
	g.layoutResultsByForm[result.form] = result
	g.syncMutex.Unlock()
}

// RunSynchronized runs all of the function calls queued by Synchronize
// and applies any layout changes queued by synchronizeLayout.
//
// RunSynchronized must be called by the group's thread.
func (g *WindowGroup) RunSynchronized() {
	// Clear the list of callbacks first to avoid deadlock
	// if a callback itself calls Synchronize()...
	g.syncMutex.Lock()
	funcs := g.syncFuncs
	var results []*formLayoutResult
	for _, result := range g.layoutResultsByForm {
		results = append(results, result)
		delete(g.layoutResultsByForm, result.form)
	}
	g.syncFuncs = nil
	g.syncMutex.Unlock()

	for _, result := range results {
		applyLayoutResults(result.results, result.stopwatch)
	}
	for _, f := range funcs {
		f()
	}
}

// ToolTip returns the tool tip control for the group, if one exists.
func (g *WindowGroup) ToolTip() *ToolTip {
	return g.toolTip
}

// CreateToolTip returns a tool tip control for the group.
//
// If a control has not already been prepared for the group one will be
// created.
func (g *WindowGroup) CreateToolTip() (*ToolTip, error) {
	if g.toolTip != nil {
		return g.toolTip, nil
	}

	tt, err := NewToolTip() // This must not call group.ToolTip()
	if err != nil {
		return nil, err
	}
	g.toolTip = tt

	// At this point the ToolTip has already added a reference for itself
	// to the group as part of the ToolTip's InitWindow process. We don't
	// want it to count toward the group's liveness, however, because it
	// would keep the group from cleaning up after itself.
	//
	// To solve this problem we also keep track of the number of
	// references that each group should ignore. The ignored references
	// are subtracted from the total number of references when evaluating
	// liveness. The expectation is that ignored references will be
	// removed as part of the group's disposal process.
	g.ignore(1)

	return tt, nil
}

// ActiveForm returns the currently active form for the group. If no
// form is active it returns nil.
func (g *WindowGroup) ActiveForm() Form {
	return g.activeForm
}

// SetActiveForm updates the currently active form for the group.
func (g *WindowGroup) SetActiveForm(form Form) {
	g.activeForm = form
}

// ignore changes the number of references that the group will ignore.
//
// ignore is used internally by WindowGroup to keep track of the number
// of references created by the group itself. When finished with a group,
// call Done() instead.
func (g *WindowGroup) ignore(delta int) {
	if g.removed {
		panic("walk: ignore() called on a WindowGroup that has been removed from its manager")
	}

	g.ignored += delta
	if g.ignored < 0 {
		panic("walk: negative WindowGroup ignored counter")
	}
	if g.refs-g.ignored == 0 {
		g.dispose()
	}
}

// dispose releases any resources consumed by the group.
func (g *WindowGroup) dispose() {
	if g.accPropServices != nil {
		g.accPropServices.Release()
		g.accPropServices = nil
	}

	if g.oleInit {
		win.OleUninitialize()
		g.oleInit = false
	}

	if g.toolTip != nil {
		g.toolTip.Dispose()
		g.toolTip = nil
	}
	g.removed = true // race detection only
	g.completion(g.threadID)
}

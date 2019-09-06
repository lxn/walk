// Copyright 2019 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"sync"
	"sync/atomic"
)

// The global window group manager instance.
var wgm windowGroupManager

// windowGroupManager manages window groups for each thread with one or
// more windows.
type windowGroupManager struct {
	mutex  sync.RWMutex
	groups map[uint32]*WindowGroup
}

// Group returns a window group for the given thread ID.
//
// The group will have its counter incremented as a result of this call.
// It is the caller's responsibility to call Done when finished with the
// group.
func (m *windowGroupManager) Group(threadID uint32) *WindowGroup {
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
	counter    windowGroupCounter // Accessed atomically, keep at front for proper alignment
	completion func(uint32)       // Called from dispose()
	threadID   uint32

	removed bool // Has the group been removed from its manager? (used for race detection)

	mutex   sync.RWMutex
	toolTip *ToolTip
}

// newWindowGroup returns a new window group for the given thread ID.
//
// The completion function will be called when the group is disposed of.
func newWindowGroup(threadID uint32, completion func(uint32)) *WindowGroup {
	//fmt.Printf("Window Group Created: %d\n", threadID)
	return &WindowGroup{
		threadID:   threadID,
		completion: completion,
	}
}

// ThreadID identifies the thread that the group is affiliated with.
func (g *WindowGroup) ThreadID() uint32 {
	return g.threadID
}

// Refs returns the current number of references to the group.
func (g *WindowGroup) Refs() int {
	refs, _ := g.counter.Value()
	return refs
}

// Add changes the group's reference counter by delta, which may be negative.
//
// If the reference counter becomes zero the group will be disposed of.
//
// If the reference counter goes negative Add will panic.
func (g *WindowGroup) Add(delta int) {
	g.add(delta, 0)
}

// Done decrements the group's reference counter by one.
func (g *WindowGroup) Done() {
	g.Add(-1)
}

// ignore changes the number of references that the group will ignore.
//
// ignore is used internally by WindowGroup to keep track of the number
// of references created by the group itself. When finished with a group,
// call Done() instead.
func (g *WindowGroup) ignore(delta int) {
	g.add(0, delta)
}

func (g *WindowGroup) add(refs, ignored int) {
	// The use of an atomic counter here is theoretically unnecessary because
	// the caller should always be calling this from the same thread. The
	// thread-safe counter is used out of an abundance of caution.

	// Best-effort race detection in case wgm.Group() is called while the
	// group is being disposed.
	if g.removed {
		panic("walk: add() called on a WindowGroup that has been removed from its manager")
	}

	refs, ignored = g.counter.Add(refs, ignored)
	//fmt.Printf("Thread %d: Refs: %d, Ignored: %d\n", g.threadID, refs, ignored)
	if refs < 0 {
		panic("walk: negative WindowGroup refs counter")
	}
	if ignored < 0 {
		panic("walk: negative WindowGroup ignored counter")
	}
	if refs-ignored == 0 {
		g.dispose()
		g.removed = true // race detection only
		g.completion(g.threadID)
	}
}

// ToolTip returns the tool tip control for the group, if one exists.
func (g *WindowGroup) ToolTip() *ToolTip {
	g.mutex.RLock()
	tt := g.toolTip
	g.mutex.RUnlock()
	return tt
}

// CreateToolTip returns a tool tip control for the group.
//
// If a control has not already been prepared for the group one will be
// created.
func (g *WindowGroup) CreateToolTip() (*ToolTip, error) {
	// The use of a mutex here is theoretically unnecessary because the
	// caller should always be calling this from the same thread. The mutex
	// is used out of an abundance of caution and may be removed in the
	// future.

	// Fast path with read lock
	g.mutex.RLock()
	if tt := g.toolTip; tt != nil {
		g.mutex.RUnlock()
		return tt, nil
	}
	g.mutex.RUnlock()

	// Slow path with write lock
	g.mutex.Lock()
	if tt := g.toolTip; tt != nil {
		g.mutex.Unlock()
		return tt, nil
	}

	tt, err := NewToolTip() // This must not call group.ToolTip()
	if err == nil {
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
	}

	g.mutex.Unlock()

	return tt, err
}

// dispose releases any resources consumed by the group.
func (g *WindowGroup) dispose() {
	//fmt.Printf("Window Group Disposed: %d\n", g.threadID)
	if g.toolTip != nil {
		g.toolTip.Dispose()
		g.toolTip = nil
	}
}

// windowGroupCounter is an atomic counter that stores two int32 values
// within a single uint64 state. It stores a pair of integers that are
// accessed atomically.
//
// Care must be taken to ensure the counter has 64-bit alignment even on
// 32-bit systems. This can be accomplished by making it the first member
// of its containing struct.
type windowGroupCounter uint64

// Add changes the counter by the given deltas, which may be negative.
func (wgc *windowGroupCounter) Add(refs, ignored int) (newRefs, newIgnored int) {
	// Constrain individual deltas to 32-bits
	dr := int32(refs)
	di := int32(ignored)

	// Pack into a single uint64 delta
	delta := uint64(uint64(dr)<<32 | uint64(di))

	// Atomic add
	addr := (*uint64)(wgc)
	state := atomic.AddUint64(addr, delta)

	// Unpack state
	return int(int32(state >> 32)), int(int32(state))
}

// Value returns the current value of the counter.
func (wgc *windowGroupCounter) Value() (refs, ignored int) {
	// Atomic load
	addr := (*uint64)(wgc)
	state := atomic.LoadUint64(addr)

	// Unpack state
	return int(int32(state >> 32)), int(int32(state))
}

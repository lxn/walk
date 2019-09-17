// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"runtime"
	"sync"
	"time"

	"github.com/lxn/win"
)

type Settings interface {
	Get(key string) (string, bool)
	Timestamp(key string) (time.Time, bool)
	Put(key, value string) error
	PutExpiring(key, value string) error
	Remove(key string) error
	ExpireDuration() time.Duration
	SetExpireDuration(expireDuration time.Duration)
	Load() error
	Save() error
}

type Persistable interface {
	Persistent() bool
	SetPersistent(value bool)
	SaveState() error
	RestoreState() error
}

type Application struct {
	mutex              sync.RWMutex
	organizationName   string
	productName        string
	settings           Settings
	exiting            bool
	exitCode           int
	panickingPublisher ErrorEventPublisher
}

var appSingleton *Application = new(Application)

func App() *Application {
	return appSingleton
}

func (app *Application) OrganizationName() string {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	return app.organizationName
}

func (app *Application) SetOrganizationName(value string) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.organizationName = value
}

func (app *Application) ProductName() string {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	return app.productName
}

func (app *Application) SetProductName(value string) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.productName = value
}

func (app *Application) Settings() Settings {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	return app.settings
}

func (app *Application) SetSettings(value Settings) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.settings = value
}

func (app *Application) Exit(exitCode int) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.exiting = true
	app.exitCode = exitCode
	win.PostQuitMessage(int32(exitCode))
}

func (app *Application) ExitCode() int {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	return app.exitCode
}

func (app *Application) Panicking() *ErrorEvent {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	return app.panickingPublisher.Event()
}

// ActiveForm returns the currently active form for the caller's thread.
// It returns nil if no form is active or the caller's thread does not
// have any windows associated with it. It should be called from within
// synchronized functions.
func (app *Application) ActiveForm() Form {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	tid := win.GetCurrentThreadId()
	group := wgm.Group(tid)
	if group == nil {
		return nil
	}
	return group.ActiveForm()
}

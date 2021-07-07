// Copyright 2021 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"time"

	"github.com/lxn/win"
)

type TimerID int

type TimerFunc func(wb *WindowBase)

func (wb *WindowBase) AddTimer(d time.Duration, fn TimerFunc) (TimerID, error) {
	wb.timerNextID++
	id := wb.timerNextID

	if wb.timerFuncs == nil {
		wb.timerFuncs = map[TimerID]TimerFunc{}
	}
	wb.timerFuncs[id] = fn

	if win.SetTimer(wb.hWnd, uintptr(id), uint32(d.Milliseconds()), 0) == 0 {
		return 0, lastError("SetTimer")
	}

	return id, nil
}

func (wb *WindowBase) ClearTimer(id TimerID) {
	win.KillTimer(wb.hWnd, uintptr(id))
	delete(wb.timerFuncs, id)
}

func (wb *WindowBase) handleTimer(wParam, lParam uintptr) {
	fn := wb.timerFuncs[TimerID(wParam)]
	if fn != nil {
		fn(wb)
	}
}

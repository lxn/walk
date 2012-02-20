// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
)

import . "github.com/lxn/go-winapi"

type MsgBoxStyle uint

const (
	MsgBoxOK                MsgBoxStyle = MB_OK
	MsgBoxOKCancel          MsgBoxStyle = MB_OKCANCEL
	MsgBoxAbortRetryIgnore  MsgBoxStyle = MB_ABORTRETRYIGNORE
	MsgBoxYesNoCancel       MsgBoxStyle = MB_YESNOCANCEL
	MsgBoxYesNo             MsgBoxStyle = MB_YESNO
	MsgBoxRetryCancel       MsgBoxStyle = MB_RETRYCANCEL
	MsgBoxCancelTryContinue MsgBoxStyle = MB_CANCELTRYCONTINUE
	MsgBoxIconHand          MsgBoxStyle = MB_ICONHAND
	MsgBoxIconQuestion      MsgBoxStyle = MB_ICONQUESTION
	MsgBoxIconExclamation   MsgBoxStyle = MB_ICONEXCLAMATION
	MsgBoxIconAsterisk      MsgBoxStyle = MB_ICONASTERISK
	MsgBoxUserIcon          MsgBoxStyle = MB_USERICON
	MsgBoxIconWarning       MsgBoxStyle = MB_ICONWARNING
	MsgBoxIconError         MsgBoxStyle = MB_ICONERROR
	MsgBoxIconInformation   MsgBoxStyle = MB_ICONINFORMATION
	MsgBoxIconStop          MsgBoxStyle = MB_ICONSTOP
	MsgBoxDefButton1        MsgBoxStyle = MB_DEFBUTTON1
	MsgBoxDefButton2        MsgBoxStyle = MB_DEFBUTTON2
	MsgBoxDefButton3        MsgBoxStyle = MB_DEFBUTTON3
	MsgBoxDefButton4        MsgBoxStyle = MB_DEFBUTTON4
)

func MsgBox(owner RootWidget, title, message string, style MsgBoxStyle) int {
	var ownerHWnd HWND

	if owner != nil {
		ownerHWnd = owner.BaseWidget().hWnd
	}

	return int(MessageBox(
		ownerHWnd,
		syscall.StringToUTF16Ptr(message),
		syscall.StringToUTF16Ptr(title),
		uint32(style)))
}

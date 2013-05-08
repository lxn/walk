// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	. "github.com/lxn/go-winapi"
)

const (
	KeyLButton           = VK_LBUTTON
	KeyRButton           = VK_RBUTTON
	KeyCancel            = VK_CANCEL
	KeyMButton           = VK_MBUTTON
	KeyXButton1          = VK_XBUTTON1
	KeyXButton2          = VK_XBUTTON2
	KeyBack              = VK_BACK
	KeyTab               = VK_TAB
	KeyClear             = VK_CLEAR
	KeyReturn            = VK_RETURN
	KeyShift             = VK_SHIFT
	KeyControl           = VK_CONTROL
	KeyMenu              = VK_MENU
	KeyPause             = VK_PAUSE
	KeyCapital           = VK_CAPITAL
	KeyKana              = VK_KANA
	KeyHangeul           = VK_HANGEUL
	KeyHangul            = VK_HANGUL
	KeyJunja             = VK_JUNJA
	KeyFinal             = VK_FINAL
	KeyHanja             = VK_HANJA
	KeyKanji             = VK_KANJI
	KeyEscape            = VK_ESCAPE
	KeyConvert           = VK_CONVERT
	KeyNonconvert        = VK_NONCONVERT
	KeyAccept            = VK_ACCEPT
	KeyModeChange        = VK_MODECHANGE
	KeySpace             = VK_SPACE
	KeyPrior             = VK_PRIOR
	KeyNext              = VK_NEXT
	KeyEnd               = VK_END
	KeyHome              = VK_HOME
	KeyLeft              = VK_LEFT
	KeyUp                = VK_UP
	KeyRight             = VK_RIGHT
	KeyDown              = VK_DOWN
	KeySelect            = VK_SELECT
	KeyPrint             = VK_PRINT
	KeyExecute           = VK_EXECUTE
	KeySnapshot          = VK_SNAPSHOT
	KeyInsert            = VK_INSERT
	KeyDelete            = VK_DELETE
	KeyHelp              = VK_HELP
	Key0                 = 0x30
	Key1                 = 0x31
	Key2                 = 0x32
	Key3                 = 0x33
	Key4                 = 0x34
	Key5                 = 0x35
	Key6                 = 0x36
	Key7                 = 0x37
	Key8                 = 0x38
	Key9                 = 0x39
	KeyA                 = 0x41
	KeyB                 = 0x42
	KeyC                 = 0x43
	KeyD                 = 0x44
	KeyE                 = 0x45
	KeyF                 = 0x46
	KeyG                 = 0x47
	KeyH                 = 0x48
	KeyI                 = 0x49
	KeyJ                 = 0x4A
	KeyK                 = 0x4B
	KeyL                 = 0x4C
	KeyM                 = 0x4D
	KeyN                 = 0x4E
	KeyO                 = 0x4F
	KeyP                 = 0x50
	KeyQ                 = 0x51
	KeyR                 = 0x52
	KeyS                 = 0x53
	KeyT                 = 0x54
	KeyU                 = 0x55
	KeyV                 = 0x56
	KeyW                 = 0x57
	KeyX                 = 0x58
	KeyY                 = 0x59
	KeyZ                 = 0x5A
	KeyLWin              = VK_LWIN
	KeyRWin              = VK_RWIN
	KeyApps              = VK_APPS
	KeySleep             = VK_SLEEP
	KeyNumpad0           = VK_NUMPAD0
	KeyNumpad1           = VK_NUMPAD1
	KeyNumpad2           = VK_NUMPAD2
	KeyNumpad3           = VK_NUMPAD3
	KeyNumpad4           = VK_NUMPAD4
	KeyNumpad5           = VK_NUMPAD5
	KeyNumpad6           = VK_NUMPAD6
	KeyNumpad7           = VK_NUMPAD7
	KeyNumpad8           = VK_NUMPAD8
	KeyNumpad9           = VK_NUMPAD9
	KeyMultiply          = VK_MULTIPLY
	KeyAdd               = VK_ADD
	KeySeparator         = VK_SEPARATOR
	KeySubtract          = VK_SUBTRACT
	KeyDecimal           = VK_DECIMAL
	KeyDivide            = VK_DIVIDE
	KeyF1                = VK_F1
	KeyF2                = VK_F2
	KeyF3                = VK_F3
	KeyF4                = VK_F4
	KeyF5                = VK_F5
	KeyF6                = VK_F6
	KeyF7                = VK_F7
	KeyF8                = VK_F8
	KeyF9                = VK_F9
	KeyF10               = VK_F10
	KeyF11               = VK_F11
	KeyF12               = VK_F12
	KeyF13               = VK_F13
	KeyF14               = VK_F14
	KeyF15               = VK_F15
	KeyF16               = VK_F16
	KeyF17               = VK_F17
	KeyF18               = VK_F18
	KeyF19               = VK_F19
	KeyF20               = VK_F20
	KeyF21               = VK_F21
	KeyF22               = VK_F22
	KeyF23               = VK_F23
	KeyF24               = VK_F24
	KeyNumlock           = VK_NUMLOCK
	KeyScroll            = VK_SCROLL
	KeyLShift            = VK_LSHIFT
	KeyRShift            = VK_RSHIFT
	KeyLControl          = VK_LCONTROL
	KeyRControl          = VK_RCONTROL
	KeyLMenu             = VK_LMENU
	KeyRMenu             = VK_RMENU
	KeyBrowserBack       = VK_BROWSER_BACK
	KeyBrowserForward    = VK_BROWSER_FORWARD
	KeyBrowserRefresh    = VK_BROWSER_REFRESH
	KeyBrowserStop       = VK_BROWSER_STOP
	KeyBrowserSearch     = VK_BROWSER_SEARCH
	KeyBrowserFavorites  = VK_BROWSER_FAVORITES
	KeyBrowserHome       = VK_BROWSER_HOME
	KeyVolumeMute        = VK_VOLUME_MUTE
	KeyVolumeDown        = VK_VOLUME_DOWN
	KeyVolumeUp          = VK_VOLUME_UP
	KeyMediaNextTrack    = VK_MEDIA_NEXT_TRACK
	KeyMediaPrevTrack    = VK_MEDIA_PREV_TRACK
	KeyMediaStop         = VK_MEDIA_STOP
	KeyMediaPlayPause    = VK_MEDIA_PLAY_PAUSE
	KeyLaunchMail        = VK_LAUNCH_MAIL
	KeyLaunchMediaSelect = VK_LAUNCH_MEDIA_SELECT
	KeyLaunchApp1        = VK_LAUNCH_APP1
	KeyLaunchApp2        = VK_LAUNCH_APP2
	KeyOEM1              = VK_OEM_1
	KeyOEMPlus           = VK_OEM_PLUS
	KeyOEMComma          = VK_OEM_COMMA
	KeyOEMMinus          = VK_OEM_MINUS
	KeyOEMPeriod         = VK_OEM_PERIOD
	KeyOEM2              = VK_OEM_2
	KeyOEM3              = VK_OEM_3
	KeyOEM4              = VK_OEM_4
	KeyOEM5              = VK_OEM_5
	KeyOEM6              = VK_OEM_6
	KeyOEM7              = VK_OEM_7
	KeyOEM8              = VK_OEM_8
	KeyOEM102            = VK_OEM_102
	KeyProcessKey        = VK_PROCESSKEY
	KeyPacket            = VK_PACKET
	KeyAttn              = VK_ATTN
	KeyCRSel             = VK_CRSEL
	KeyEXSel             = VK_EXSEL
	KeyErEOF             = VK_EREOF
	KeyPlay              = VK_PLAY
	KeyZoom              = VK_ZOOM
	KeyNoName            = VK_NONAME
	KeyPA1               = VK_PA1
	KeyOEMClear          = VK_OEM_CLEAR
)

func LogErrors() bool {
	return logErrors
}

func SetLogErrors(v bool) {
	logErrors = v
}

func PanicOnError() bool {
	return panicOnError
}

func SetPanicOnError(v bool) {
	panicOnError = v
}

func TranslationFunc() TranslationFunction {
	return translation
}

func SetTranslationFunc(f TranslationFunction) {
	translation = f
}

type TranslationFunction func(source string, context ...string) string

var translation TranslationFunction

func tr(source string, context ...string) string {
	if translation == nil {
		return source
	}

	return translation(source, context...)
}

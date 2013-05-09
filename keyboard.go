// Copyright 2013 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	. "github.com/lxn/go-winapi"
)

type Key uint16

const (
	KeyLButton           Key = VK_LBUTTON
	KeyRButton               = VK_RBUTTON
	KeyCancel                = VK_CANCEL
	KeyMButton               = VK_MBUTTON
	KeyXButton1              = VK_XBUTTON1
	KeyXButton2              = VK_XBUTTON2
	KeyBack                  = VK_BACK
	KeyTab                   = VK_TAB
	KeyClear                 = VK_CLEAR
	KeyReturn                = VK_RETURN
	KeyShift                 = VK_SHIFT
	KeyControl               = VK_CONTROL
	KeyAlt                   = VK_MENU
	KeyMenu                  = VK_MENU
	KeyPause                 = VK_PAUSE
	KeyCapital               = VK_CAPITAL
	KeyKana                  = VK_KANA
	KeyHangul                = VK_HANGUL
	KeyJunja                 = VK_JUNJA
	KeyFinal                 = VK_FINAL
	KeyHanja                 = VK_HANJA
	KeyKanji                 = VK_KANJI
	KeyEscape                = VK_ESCAPE
	KeyConvert               = VK_CONVERT
	KeyNonconvert            = VK_NONCONVERT
	KeyAccept                = VK_ACCEPT
	KeyModeChange            = VK_MODECHANGE
	KeySpace                 = VK_SPACE
	KeyPrior                 = VK_PRIOR
	KeyNext                  = VK_NEXT
	KeyEnd                   = VK_END
	KeyHome                  = VK_HOME
	KeyLeft                  = VK_LEFT
	KeyUp                    = VK_UP
	KeyRight                 = VK_RIGHT
	KeyDown                  = VK_DOWN
	KeySelect                = VK_SELECT
	KeyPrint                 = VK_PRINT
	KeyExecute               = VK_EXECUTE
	KeySnapshot              = VK_SNAPSHOT
	KeyInsert                = VK_INSERT
	KeyDelete                = VK_DELETE
	KeyHelp                  = VK_HELP
	Key0                     = 0x30
	Key1                     = 0x31
	Key2                     = 0x32
	Key3                     = 0x33
	Key4                     = 0x34
	Key5                     = 0x35
	Key6                     = 0x36
	Key7                     = 0x37
	Key8                     = 0x38
	Key9                     = 0x39
	KeyA                     = 0x41
	KeyB                     = 0x42
	KeyC                     = 0x43
	KeyD                     = 0x44
	KeyE                     = 0x45
	KeyF                     = 0x46
	KeyG                     = 0x47
	KeyH                     = 0x48
	KeyI                     = 0x49
	KeyJ                     = 0x4A
	KeyK                     = 0x4B
	KeyL                     = 0x4C
	KeyM                     = 0x4D
	KeyN                     = 0x4E
	KeyO                     = 0x4F
	KeyP                     = 0x50
	KeyQ                     = 0x51
	KeyR                     = 0x52
	KeyS                     = 0x53
	KeyT                     = 0x54
	KeyU                     = 0x55
	KeyV                     = 0x56
	KeyW                     = 0x57
	KeyX                     = 0x58
	KeyY                     = 0x59
	KeyZ                     = 0x5A
	KeyLWin                  = VK_LWIN
	KeyRWin                  = VK_RWIN
	KeyApps                  = VK_APPS
	KeySleep                 = VK_SLEEP
	KeyNumpad0               = VK_NUMPAD0
	KeyNumpad1               = VK_NUMPAD1
	KeyNumpad2               = VK_NUMPAD2
	KeyNumpad3               = VK_NUMPAD3
	KeyNumpad4               = VK_NUMPAD4
	KeyNumpad5               = VK_NUMPAD5
	KeyNumpad6               = VK_NUMPAD6
	KeyNumpad7               = VK_NUMPAD7
	KeyNumpad8               = VK_NUMPAD8
	KeyNumpad9               = VK_NUMPAD9
	KeyMultiply              = VK_MULTIPLY
	KeyAdd                   = VK_ADD
	KeySeparator             = VK_SEPARATOR
	KeySubtract              = VK_SUBTRACT
	KeyDecimal               = VK_DECIMAL
	KeyDivide                = VK_DIVIDE
	KeyF1                    = VK_F1
	KeyF2                    = VK_F2
	KeyF3                    = VK_F3
	KeyF4                    = VK_F4
	KeyF5                    = VK_F5
	KeyF6                    = VK_F6
	KeyF7                    = VK_F7
	KeyF8                    = VK_F8
	KeyF9                    = VK_F9
	KeyF10                   = VK_F10
	KeyF11                   = VK_F11
	KeyF12                   = VK_F12
	KeyF13                   = VK_F13
	KeyF14                   = VK_F14
	KeyF15                   = VK_F15
	KeyF16                   = VK_F16
	KeyF17                   = VK_F17
	KeyF18                   = VK_F18
	KeyF19                   = VK_F19
	KeyF20                   = VK_F20
	KeyF21                   = VK_F21
	KeyF22                   = VK_F22
	KeyF23                   = VK_F23
	KeyF24                   = VK_F24
	KeyNumlock               = VK_NUMLOCK
	KeyScroll                = VK_SCROLL
	KeyLShift                = VK_LSHIFT
	KeyRShift                = VK_RSHIFT
	KeyLControl              = VK_LCONTROL
	KeyRControl              = VK_RCONTROL
	KeyLMenu                 = VK_LMENU
	KeyRMenu                 = VK_RMENU
	KeyBrowserBack           = VK_BROWSER_BACK
	KeyBrowserForward        = VK_BROWSER_FORWARD
	KeyBrowserRefresh        = VK_BROWSER_REFRESH
	KeyBrowserStop           = VK_BROWSER_STOP
	KeyBrowserSearch         = VK_BROWSER_SEARCH
	KeyBrowserFavorites      = VK_BROWSER_FAVORITES
	KeyBrowserHome           = VK_BROWSER_HOME
	KeyVolumeMute            = VK_VOLUME_MUTE
	KeyVolumeDown            = VK_VOLUME_DOWN
	KeyVolumeUp              = VK_VOLUME_UP
	KeyMediaNextTrack        = VK_MEDIA_NEXT_TRACK
	KeyMediaPrevTrack        = VK_MEDIA_PREV_TRACK
	KeyMediaStop             = VK_MEDIA_STOP
	KeyMediaPlayPause        = VK_MEDIA_PLAY_PAUSE
	KeyLaunchMail            = VK_LAUNCH_MAIL
	KeyLaunchMediaSelect     = VK_LAUNCH_MEDIA_SELECT
	KeyLaunchApp1            = VK_LAUNCH_APP1
	KeyLaunchApp2            = VK_LAUNCH_APP2
	KeyOEM1                  = VK_OEM_1
	KeyOEMPlus               = VK_OEM_PLUS
	KeyOEMComma              = VK_OEM_COMMA
	KeyOEMMinus              = VK_OEM_MINUS
	KeyOEMPeriod             = VK_OEM_PERIOD
	KeyOEM2                  = VK_OEM_2
	KeyOEM3                  = VK_OEM_3
	KeyOEM4                  = VK_OEM_4
	KeyOEM5                  = VK_OEM_5
	KeyOEM6                  = VK_OEM_6
	KeyOEM7                  = VK_OEM_7
	KeyOEM8                  = VK_OEM_8
	KeyOEM102                = VK_OEM_102
	KeyProcessKey            = VK_PROCESSKEY
	KeyPacket                = VK_PACKET
	KeyAttn                  = VK_ATTN
	KeyCRSel                 = VK_CRSEL
	KeyEXSel                 = VK_EXSEL
	KeyErEOF                 = VK_EREOF
	KeyPlay                  = VK_PLAY
	KeyZoom                  = VK_ZOOM
	KeyNoName                = VK_NONAME
	KeyPA1                   = VK_PA1
	KeyOEMClear              = VK_OEM_CLEAR
)

var key2string = map[Key]string{
	KeyLButton:           "LButton",
	KeyRButton:           "RButton",
	KeyCancel:            "Cancel",
	KeyMButton:           "MButton",
	KeyXButton1:          "XButton1",
	KeyXButton2:          "XButton2",
	KeyBack:              "Back",
	KeyTab:               "Tab",
	KeyClear:             "Clear",
	KeyReturn:            "Return",
	KeyShift:             "Shift",
	KeyControl:           "Control",
	KeyAlt:               "Alt / Menu",
	KeyPause:             "Pause",
	KeyCapital:           "Capital",
	KeyKana:              "Kana / Hangul",
	KeyJunja:             "Junja",
	KeyFinal:             "Final",
	KeyHanja:             "Hanja / Kanji",
	KeyEscape:            "Escape",
	KeyConvert:           "Convert",
	KeyNonconvert:        "Nonconvert",
	KeyAccept:            "Accept",
	KeyModeChange:        "ModeChange",
	KeySpace:             "Space",
	KeyPrior:             "Prior",
	KeyNext:              "Next",
	KeyEnd:               "End",
	KeyHome:              "Home",
	KeyLeft:              "Left",
	KeyUp:                "Up",
	KeyRight:             "Right",
	KeyDown:              "Down",
	KeySelect:            "Select",
	KeyPrint:             "Print",
	KeyExecute:           "Execute",
	KeySnapshot:          "Snapshot",
	KeyInsert:            "Insert",
	KeyDelete:            "Delete",
	KeyHelp:              "Help",
	Key0:                 "0",
	Key1:                 "1",
	Key2:                 "2",
	Key3:                 "3",
	Key4:                 "4",
	Key5:                 "5",
	Key6:                 "6",
	Key7:                 "7",
	Key8:                 "8",
	Key9:                 "9",
	KeyA:                 "A",
	KeyB:                 "B",
	KeyC:                 "C",
	KeyD:                 "D",
	KeyE:                 "E",
	KeyF:                 "F",
	KeyG:                 "G",
	KeyH:                 "H",
	KeyI:                 "I",
	KeyJ:                 "J",
	KeyK:                 "K",
	KeyL:                 "L",
	KeyM:                 "M",
	KeyN:                 "N",
	KeyO:                 "O",
	KeyP:                 "P",
	KeyQ:                 "Q",
	KeyR:                 "R",
	KeyS:                 "S",
	KeyT:                 "T",
	KeyU:                 "U",
	KeyV:                 "V",
	KeyW:                 "W",
	KeyX:                 "X",
	KeyY:                 "Y",
	KeyZ:                 "Z",
	KeyLWin:              "LWin",
	KeyRWin:              "RWin",
	KeyApps:              "Apps",
	KeySleep:             "Sleep",
	KeyNumpad0:           "Numpad0",
	KeyNumpad1:           "Numpad1",
	KeyNumpad2:           "Numpad2",
	KeyNumpad3:           "Numpad3",
	KeyNumpad4:           "Numpad4",
	KeyNumpad5:           "Numpad5",
	KeyNumpad6:           "Numpad6",
	KeyNumpad7:           "Numpad7",
	KeyNumpad8:           "Numpad8",
	KeyNumpad9:           "Numpad9",
	KeyMultiply:          "Multiply",
	KeyAdd:               "Add",
	KeySeparator:         "Separator",
	KeySubtract:          "Subtract",
	KeyDecimal:           "Decimal",
	KeyDivide:            "Divide",
	KeyF1:                "F1",
	KeyF2:                "F2",
	KeyF3:                "F3",
	KeyF4:                "F4",
	KeyF5:                "F5",
	KeyF6:                "F6",
	KeyF7:                "F7",
	KeyF8:                "F8",
	KeyF9:                "F9",
	KeyF10:               "F10",
	KeyF11:               "F11",
	KeyF12:               "F12",
	KeyF13:               "F13",
	KeyF14:               "F14",
	KeyF15:               "F15",
	KeyF16:               "F16",
	KeyF17:               "F17",
	KeyF18:               "F18",
	KeyF19:               "F19",
	KeyF20:               "F20",
	KeyF21:               "F21",
	KeyF22:               "F22",
	KeyF23:               "F23",
	KeyF24:               "F24",
	KeyNumlock:           "Numlock",
	KeyScroll:            "Scroll",
	KeyLShift:            "LShift",
	KeyRShift:            "RShift",
	KeyLControl:          "LControl",
	KeyRControl:          "RControl",
	KeyLMenu:             "LMenu",
	KeyRMenu:             "RMenu",
	KeyBrowserBack:       "BrowserBack",
	KeyBrowserForward:    "BrowserForward",
	KeyBrowserRefresh:    "BrowserRefresh",
	KeyBrowserStop:       "BrowserStop",
	KeyBrowserSearch:     "BrowserSearch",
	KeyBrowserFavorites:  "BrowserFavorites",
	KeyBrowserHome:       "BrowserHome",
	KeyVolumeMute:        "VolumeMute",
	KeyVolumeDown:        "VolumeDown",
	KeyVolumeUp:          "VolumeUp",
	KeyMediaNextTrack:    "MediaNextTrack",
	KeyMediaPrevTrack:    "MediaPrevTrack",
	KeyMediaStop:         "MediaStop",
	KeyMediaPlayPause:    "MediaPlayPause",
	KeyLaunchMail:        "LaunchMail",
	KeyLaunchMediaSelect: "LaunchMediaSelect",
	KeyLaunchApp1:        "LaunchApp1",
	KeyLaunchApp2:        "LaunchApp2",
	KeyOEM1:              "OEM1",
	KeyOEMPlus:           "OEMPlus",
	KeyOEMComma:          "OEMComma",
	KeyOEMMinus:          "OEMMinus",
	KeyOEMPeriod:         "OEMPeriod",
	KeyOEM2:              "OEM2",
	KeyOEM3:              "OEM3",
	KeyOEM4:              "OEM4",
	KeyOEM5:              "OEM5",
	KeyOEM6:              "OEM6",
	KeyOEM7:              "OEM7",
	KeyOEM8:              "OEM8",
	KeyOEM102:            "OEM102",
	KeyProcessKey:        "ProcessKey",
	KeyPacket:            "Packet",
	KeyAttn:              "Attn",
	KeyCRSel:             "CRSel",
	KeyEXSel:             "EXSel",
	KeyErEOF:             "ErEOF",
	KeyPlay:              "Play",
	KeyZoom:              "Zoom",
	KeyNoName:            "NoName",
	KeyPA1:               "PA1",
	KeyOEMClear:          "OEMClear",
}

func (k Key) String() string {
	return key2string[k]
}

func AltDown() bool {
	return GetKeyState(KeyAlt)>>15 != 0
}

func ControlDown() bool {
	return GetKeyState(KeyControl)>>15 != 0
}

func ShiftDown() bool {
	return GetKeyState(KeyShift)>>15 != 0
}

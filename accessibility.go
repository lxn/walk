// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import "github.com/lxn/win"

type AccState int32

const (
	AccStateNormal          AccState = win.STATE_SYSTEM_NORMAL
	AccStateUnavailable     AccState = win.STATE_SYSTEM_UNAVAILABLE
	AccStateSelected        AccState = win.STATE_SYSTEM_SELECTED
	AccStateFocused         AccState = win.STATE_SYSTEM_FOCUSED
	AccStatePressed         AccState = win.STATE_SYSTEM_PRESSED
	AccStateChecked         AccState = win.STATE_SYSTEM_CHECKED
	AccStateMixed           AccState = win.STATE_SYSTEM_MIXED
	AccStateIndeterminate   AccState = win.STATE_SYSTEM_INDETERMINATE
	AccStateReadonly        AccState = win.STATE_SYSTEM_READONLY
	AccStateHotTracked      AccState = win.STATE_SYSTEM_HOTTRACKED
	AccStateDefault         AccState = win.STATE_SYSTEM_DEFAULT
	AccStateExpanded        AccState = win.STATE_SYSTEM_EXPANDED
	AccStateCollapsed       AccState = win.STATE_SYSTEM_COLLAPSED
	AccStateBusy            AccState = win.STATE_SYSTEM_BUSY
	AccStateFloating        AccState = win.STATE_SYSTEM_FLOATING
	AccStateMarqueed        AccState = win.STATE_SYSTEM_MARQUEED
	AccStateAnimated        AccState = win.STATE_SYSTEM_ANIMATED
	AccStateInvisible       AccState = win.STATE_SYSTEM_INVISIBLE
	AccStateOffscreen       AccState = win.STATE_SYSTEM_OFFSCREEN
	AccStateSizeable        AccState = win.STATE_SYSTEM_SIZEABLE
	AccStateMoveable        AccState = win.STATE_SYSTEM_MOVEABLE
	AccStateSelfVoicing     AccState = win.STATE_SYSTEM_SELFVOICING
	AccStateFocusable       AccState = win.STATE_SYSTEM_FOCUSABLE
	AccStateSelectable      AccState = win.STATE_SYSTEM_SELECTABLE
	AccStateLinked          AccState = win.STATE_SYSTEM_LINKED
	AccStateTraversed       AccState = win.STATE_SYSTEM_TRAVERSED
	AccStateMultiselectable AccState = win.STATE_SYSTEM_MULTISELECTABLE
	AccStateExtselectable   AccState = win.STATE_SYSTEM_EXTSELECTABLE
	AccStateAlertLow        AccState = win.STATE_SYSTEM_ALERT_LOW
	AccStateAlertMedium     AccState = win.STATE_SYSTEM_ALERT_MEDIUM
	AccStateAlertHigh       AccState = win.STATE_SYSTEM_ALERT_HIGH
	AccStateProtected       AccState = win.STATE_SYSTEM_PROTECTED
	AccStateHasPopup        AccState = win.STATE_SYSTEM_HASPOPUP
	AccStateValid           AccState = win.STATE_SYSTEM_VALID
)

type AccRole int32

const (
	AccRoleTitlebar           AccRole = win.ROLE_SYSTEM_TITLEBAR
	AccRoleMenubar            AccRole = win.ROLE_SYSTEM_MENUBAR
	AccRoleScrollbar          AccRole = win.ROLE_SYSTEM_SCROLLBAR
	AccRoleGrip               AccRole = win.ROLE_SYSTEM_GRIP
	AccRoleSound              AccRole = win.ROLE_SYSTEM_SOUND
	AccRoleCursor             AccRole = win.ROLE_SYSTEM_CURSOR
	AccRoleCaret              AccRole = win.ROLE_SYSTEM_CARET
	AccRoleAlert              AccRole = win.ROLE_SYSTEM_ALERT
	AccRoleWindow             AccRole = win.ROLE_SYSTEM_WINDOW
	AccRoleClient             AccRole = win.ROLE_SYSTEM_CLIENT
	AccRoleMenuPopup          AccRole = win.ROLE_SYSTEM_MENUPOPUP
	AccRoleMenuItem           AccRole = win.ROLE_SYSTEM_MENUITEM
	AccRoleTooltip            AccRole = win.ROLE_SYSTEM_TOOLTIP
	AccRoleApplication        AccRole = win.ROLE_SYSTEM_APPLICATION
	AccRoleDocument           AccRole = win.ROLE_SYSTEM_DOCUMENT
	AccRolePane               AccRole = win.ROLE_SYSTEM_PANE
	AccRoleChart              AccRole = win.ROLE_SYSTEM_CHART
	AccRoleDialog             AccRole = win.ROLE_SYSTEM_DIALOG
	AccRoleBorder             AccRole = win.ROLE_SYSTEM_BORDER
	AccRoleGrouping           AccRole = win.ROLE_SYSTEM_GROUPING
	AccRoleSeparator          AccRole = win.ROLE_SYSTEM_SEPARATOR
	AccRoleToolbar            AccRole = win.ROLE_SYSTEM_TOOLBAR
	AccRoleStatusbar          AccRole = win.ROLE_SYSTEM_STATUSBAR
	AccRoleTable              AccRole = win.ROLE_SYSTEM_TABLE
	AccRoleColumnHeader       AccRole = win.ROLE_SYSTEM_COLUMNHEADER
	AccRoleRowHeader          AccRole = win.ROLE_SYSTEM_ROWHEADER
	AccRoleColumn             AccRole = win.ROLE_SYSTEM_COLUMN
	AccRoleRow                AccRole = win.ROLE_SYSTEM_ROW
	AccRoleCell               AccRole = win.ROLE_SYSTEM_CELL
	AccRoleLink               AccRole = win.ROLE_SYSTEM_LINK
	AccRoleHelpBalloon        AccRole = win.ROLE_SYSTEM_HELPBALLOON
	AccRoleCharacter          AccRole = win.ROLE_SYSTEM_CHARACTER
	AccRoleList               AccRole = win.ROLE_SYSTEM_LIST
	AccRoleListItem           AccRole = win.ROLE_SYSTEM_LISTITEM
	AccRoleOutline            AccRole = win.ROLE_SYSTEM_OUTLINE
	AccRoleOutlineItem        AccRole = win.ROLE_SYSTEM_OUTLINEITEM
	AccRolePagetab            AccRole = win.ROLE_SYSTEM_PAGETAB
	AccRolePropertyPage       AccRole = win.ROLE_SYSTEM_PROPERTYPAGE
	AccRoleIndicator          AccRole = win.ROLE_SYSTEM_INDICATOR
	AccRoleGraphic            AccRole = win.ROLE_SYSTEM_GRAPHIC
	AccRoleStatictext         AccRole = win.ROLE_SYSTEM_STATICTEXT
	AccRoleText               AccRole = win.ROLE_SYSTEM_TEXT
	AccRolePushbutton         AccRole = win.ROLE_SYSTEM_PUSHBUTTON
	AccRoleCheckbutton        AccRole = win.ROLE_SYSTEM_CHECKBUTTON
	AccRoleRadiobutton        AccRole = win.ROLE_SYSTEM_RADIOBUTTON
	AccRoleCombobox           AccRole = win.ROLE_SYSTEM_COMBOBOX
	AccRoleDroplist           AccRole = win.ROLE_SYSTEM_DROPLIST
	AccRoleProgressbar        AccRole = win.ROLE_SYSTEM_PROGRESSBAR
	AccRoleDial               AccRole = win.ROLE_SYSTEM_DIAL
	AccRoleHotkeyfield        AccRole = win.ROLE_SYSTEM_HOTKEYFIELD
	AccRoleSlider             AccRole = win.ROLE_SYSTEM_SLIDER
	AccRoleSpinbutton         AccRole = win.ROLE_SYSTEM_SPINBUTTON
	AccRoleDiagram            AccRole = win.ROLE_SYSTEM_DIAGRAM
	AccRoleAnimation          AccRole = win.ROLE_SYSTEM_ANIMATION
	AccRoleEquation           AccRole = win.ROLE_SYSTEM_EQUATION
	AccRoleButtonDropdown     AccRole = win.ROLE_SYSTEM_BUTTONDROPDOWN
	AccRoleButtonMenu         AccRole = win.ROLE_SYSTEM_BUTTONMENU
	AccRoleButtonDropdownGrid AccRole = win.ROLE_SYSTEM_BUTTONDROPDOWNGRID
	AccRoleWhitespace         AccRole = win.ROLE_SYSTEM_WHITESPACE
	AccRolePageTabList        AccRole = win.ROLE_SYSTEM_PAGETABLIST
	AccRoleClock              AccRole = win.ROLE_SYSTEM_CLOCK
	AccRoleSplitButton        AccRole = win.ROLE_SYSTEM_SPLITBUTTON
	AccRoleIPAddress          AccRole = win.ROLE_SYSTEM_IPADDRESS
	AccRoleOutlineButton      AccRole = win.ROLE_SYSTEM_OUTLINEBUTTON
)

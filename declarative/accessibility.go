// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

// AccState enum defines the state of the window/control
type AccState int32

// Window/control states
const (
	AccStateNormal          = AccState(walk.AccStateNormal)
	AccStateUnavailable     = AccState(walk.AccStateUnavailable)
	AccStateSelected        = AccState(walk.AccStateSelected)
	AccStateFocused         = AccState(walk.AccStateFocused)
	AccStatePressed         = AccState(walk.AccStatePressed)
	AccStateChecked         = AccState(walk.AccStateChecked)
	AccStateMixed           = AccState(walk.AccStateMixed)
	AccStateIndeterminate   = AccState(walk.AccStateIndeterminate)
	AccStateReadonly        = AccState(walk.AccStateReadonly)
	AccStateHotTracked      = AccState(walk.AccStateHotTracked)
	AccStateDefault         = AccState(walk.AccStateDefault)
	AccStateExpanded        = AccState(walk.AccStateExpanded)
	AccStateCollapsed       = AccState(walk.AccStateCollapsed)
	AccStateBusy            = AccState(walk.AccStateBusy)
	AccStateFloating        = AccState(walk.AccStateFloating)
	AccStateMarqueed        = AccState(walk.AccStateMarqueed)
	AccStateAnimated        = AccState(walk.AccStateAnimated)
	AccStateInvisible       = AccState(walk.AccStateInvisible)
	AccStateOffscreen       = AccState(walk.AccStateOffscreen)
	AccStateSizeable        = AccState(walk.AccStateSizeable)
	AccStateMoveable        = AccState(walk.AccStateMoveable)
	AccStateSelfVoicing     = AccState(walk.AccStateSelfVoicing)
	AccStateFocusable       = AccState(walk.AccStateFocusable)
	AccStateSelectable      = AccState(walk.AccStateSelectable)
	AccStateLinked          = AccState(walk.AccStateLinked)
	AccStateTraversed       = AccState(walk.AccStateTraversed)
	AccStateMultiselectable = AccState(walk.AccStateMultiselectable)
	AccStateExtselectable   = AccState(walk.AccStateExtselectable)
	AccStateAlertLow        = AccState(walk.AccStateAlertLow)
	AccStateAlertMedium     = AccState(walk.AccStateAlertMedium)
	AccStateAlertHigh       = AccState(walk.AccStateAlertHigh)
	AccStateProtected       = AccState(walk.AccStateProtected)
	AccStateHasPopup        = AccState(walk.AccStateHasPopup)
	AccStateValid           = AccState(walk.AccStateValid)
)

// AccRole enum defines the role of the window/control in UI.
type AccRole int32

// Window/control system roles
const (
	AccRoleTitlebar           = AccRole(walk.AccRoleTitlebar)
	AccRoleMenubar            = AccRole(walk.AccRoleMenubar)
	AccRoleScrollbar          = AccRole(walk.AccRoleScrollbar)
	AccRoleGrip               = AccRole(walk.AccRoleGrip)
	AccRoleSound              = AccRole(walk.AccRoleSound)
	AccRoleCursor             = AccRole(walk.AccRoleCursor)
	AccRoleCaret              = AccRole(walk.AccRoleCaret)
	AccRoleAlert              = AccRole(walk.AccRoleAlert)
	AccRoleWindow             = AccRole(walk.AccRoleWindow)
	AccRoleClient             = AccRole(walk.AccRoleClient)
	AccRoleMenuPopup          = AccRole(walk.AccRoleMenuPopup)
	AccRoleMenuItem           = AccRole(walk.AccRoleMenuItem)
	AccRoleTooltip            = AccRole(walk.AccRoleTooltip)
	AccRoleApplication        = AccRole(walk.AccRoleApplication)
	AccRoleDocument           = AccRole(walk.AccRoleDocument)
	AccRolePane               = AccRole(walk.AccRolePane)
	AccRoleChart              = AccRole(walk.AccRoleChart)
	AccRoleDialog             = AccRole(walk.AccRoleDialog)
	AccRoleBorder             = AccRole(walk.AccRoleBorder)
	AccRoleGrouping           = AccRole(walk.AccRoleGrouping)
	AccRoleSeparator          = AccRole(walk.AccRoleSeparator)
	AccRoleToolbar            = AccRole(walk.AccRoleToolbar)
	AccRoleStatusbar          = AccRole(walk.AccRoleStatusbar)
	AccRoleTable              = AccRole(walk.AccRoleTable)
	AccRoleColumnHeader       = AccRole(walk.AccRoleColumnHeader)
	AccRoleRowHeader          = AccRole(walk.AccRoleRowHeader)
	AccRoleColumn             = AccRole(walk.AccRoleColumn)
	AccRoleRow                = AccRole(walk.AccRoleRow)
	AccRoleCell               = AccRole(walk.AccRoleCell)
	AccRoleLink               = AccRole(walk.AccRoleLink)
	AccRoleHelpBalloon        = AccRole(walk.AccRoleHelpBalloon)
	AccRoleCharacter          = AccRole(walk.AccRoleCharacter)
	AccRoleList               = AccRole(walk.AccRoleList)
	AccRoleListItem           = AccRole(walk.AccRoleListItem)
	AccRoleOutline            = AccRole(walk.AccRoleOutline)
	AccRoleOutlineItem        = AccRole(walk.AccRoleOutlineItem)
	AccRolePagetab            = AccRole(walk.AccRolePagetab)
	AccRolePropertyPage       = AccRole(walk.AccRolePropertyPage)
	AccRoleIndicator          = AccRole(walk.AccRoleIndicator)
	AccRoleGraphic            = AccRole(walk.AccRoleGraphic)
	AccRoleStatictext         = AccRole(walk.AccRoleStatictext)
	AccRoleText               = AccRole(walk.AccRoleText)
	AccRolePushbutton         = AccRole(walk.AccRolePushbutton)
	AccRoleCheckbutton        = AccRole(walk.AccRoleCheckbutton)
	AccRoleRadiobutton        = AccRole(walk.AccRoleRadiobutton)
	AccRoleCombobox           = AccRole(walk.AccRoleCombobox)
	AccRoleDroplist           = AccRole(walk.AccRoleDroplist)
	AccRoleProgressbar        = AccRole(walk.AccRoleProgressbar)
	AccRoleDial               = AccRole(walk.AccRoleDial)
	AccRoleHotkeyfield        = AccRole(walk.AccRoleHotkeyfield)
	AccRoleSlider             = AccRole(walk.AccRoleSlider)
	AccRoleSpinbutton         = AccRole(walk.AccRoleSpinbutton)
	AccRoleDiagram            = AccRole(walk.AccRoleDiagram)
	AccRoleAnimation          = AccRole(walk.AccRoleAnimation)
	AccRoleEquation           = AccRole(walk.AccRoleEquation)
	AccRoleButtonDropdown     = AccRole(walk.AccRoleButtonDropdown)
	AccRoleButtonMenu         = AccRole(walk.AccRoleButtonMenu)
	AccRoleButtonDropdownGrid = AccRole(walk.AccRoleButtonDropdownGrid)
	AccRoleWhitespace         = AccRole(walk.AccRoleWhitespace)
	AccRolePageTabList        = AccRole(walk.AccRolePageTabList)
	AccRoleClock              = AccRole(walk.AccRoleClock)
	AccRoleSplitButton        = AccRole(walk.AccRoleSplitButton)
	AccRoleIPAddress          = AccRole(walk.AccRoleIPAddress)
	AccRoleOutlineButton      = AccRole(walk.AccRoleOutlineButton)
)

// Accessibility properties
type Accessibility struct {
	Accelerator   string
	DefaultAction string
	Description   string
	Help          string
	Name          string
	Role          AccRole
	RoleMap       string
	State         AccState
	StateMap      string
	ValueMap      string
}

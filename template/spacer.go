package template

import (
	"walk"
)

func Spacer(o walk.Orientation, wf WidgetFlag) GuiTemplate {
	return &SimpleGuiTemplate{wf, func(parent walk.Container) (w walk.Widget, err error) {
		if o == walk.Horizontal {
			w, err = walk.NewHSpacer(parent)
		} else {
			w, err = walk.NewVSpacer(parent)
		}
		return
	}}
}

func HSpacer(wf WidgetFlag) GuiTemplate { return Spacer(walk.Horizontal, wf) }
func VSpacer(wf WidgetFlag) GuiTemplate { return Spacer(walk.Vertical, wf) }

func SpacerFixed(o walk.Orientation, s int, wf WidgetFlag) GuiTemplate {
	return &SimpleGuiTemplate{wf, func(parent walk.Container) (w walk.Widget, err error) {
		if o == walk.Horizontal {
			w, err = walk.NewHSpacerFixed(parent, s)
		} else {
			w, err = walk.NewVSpacerFixed(parent, s)
		}
		return
	}}
}

func HSpacerFixed(s int, wf WidgetFlag) GuiTemplate { return SpacerFixed(walk.Horizontal, s, wf) }
func VSpacerFixed(s int, wf WidgetFlag) GuiTemplate { return SpacerFixed(walk.Vertical, s, wf) }


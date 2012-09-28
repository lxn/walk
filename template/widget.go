package template

import (
	"walk"
)

func Label(s string) GuiTemplate {
	return SimpleGuiTemplateFunc(func(parent walk.Container) error {
		l, err := walk.NewLabel(parent)
		if err == nil {
			l.SetText(s)
		}
		return err
	})
}

func button(
		p **walk.Button,
		title string,
		wf WidgetFlag,
		clickfunc walk.EventHandler,
		nf func(parent walk.Container) (*walk.Button, error),
	) GuiTemplate {
	return &SimpleGuiTemplate{wf, func(parent walk.Container) (walk.Widget, error) {
		button, err := nf(parent)
		if err == nil {
			err = button.SetText(title)
			if err == nil {
				if clickfunc != nil {
					button.Clicked().Attach(clickfunc)
				}
				if p != nil {
					*p = button
				}
			}
		}
		return button, err
	}}
}

func PushButton(p **walk.Button, title string, wf WidgetFlag, clickfunc walk.EventHandler) GuiTemplate {
	return button(p, title, wf, clickfunc, func(parent walk.Container) (b *walk.Button, err error) {
		var bb *walk.PushButton
		bb, err = walk.NewPushButton(parent)
		return &bb.Button, err
	})
}

func ToolButton(p **walk.Button, title string, wf WidgetFlag, clickfunc walk.EventHandler) GuiTemplate {
	return button(p, title, wf, clickfunc, func(parent walk.Container) (b *walk.Button, err error) {
		var bb *walk.ToolButton
		bb, err = walk.NewToolButton(parent)
		return &bb.Button, err
	})
}

func LineEdit(p **walk.LineEdit, value string, wf WidgetFlag) GuiTemplate {
	return &SimpleGuiTemplate{wf, func(parent walk.Container) (walk.Widget, error) {
		lineedit, err := walk.NewLineEdit(parent)
		if err == nil {
			err = lineedit.SetText(value)
			if err == nil && p != nil {
				*p = lineedit
			}
		}
		return lineedit, err
	}}
}

func ProgressBar(p **walk.ProgressBar, min, max int, wf WidgetFlag) GuiTemplate {
	return &SimpleGuiTemplate{wf, func(parent walk.Container) (walk.Widget, error) {
		progressbar, err := walk.NewProgressBar(parent)
		if err == nil {
			progressbar.SetRange(min, max)
			if err == nil && p != nil {
				*p = progressbar
			}
		}
		return progressbar, err
	}}
}

func TextEdit(p **walk.TextEdit, value string, wf WidgetFlag) GuiTemplate {
	return &SimpleGuiTemplate{wf, func(parent walk.Container) (walk.Widget, error) {
		textedit, err := walk.NewTextEdit(parent)
		if err == nil {
			err = textedit.SetText(value)
			if err == nil && p != nil {
				*p = textedit
			}
		}
		return textedit, err
	}}
}


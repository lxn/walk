// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import (
	. "walk/winapi/user32"
)

const splitterHandleWindowClass = `\o/ Walk_SplitterHandle_Class \o/`

var splitterHandleWindowClassRegistered bool

type splitterHandle struct {
	WidgetBase
}

func newSplitterHandle(splitter *Splitter) (*splitterHandle, os.Error) {
	if splitter == nil {
		return nil, newError("splitter cannot be nil")
	}

	ensureRegisteredWindowClass(splitterHandleWindowClass, &splitterHandleWindowClassRegistered)

	sh := &splitterHandle{}
	sh.parent = splitter

	if err := initWidget(
		sh,
		splitter,
		splitterHandleWindowClass,
		WS_CHILD|WS_VISIBLE,
		0); err != nil {
		return nil, err
	}

	return sh, nil
}

func (sh *splitterHandle) LayoutFlags() LayoutFlags {
	splitter := sh.Parent().(*Splitter)
	if splitter.Orientation() == Horizontal {
		return VGrow | VShrink
	}

	return HGrow | HShrink
}

func (sh *splitterHandle) PreferredSize() Size {
	splitter := sh.Parent().(*Splitter)
	handleWidth := splitter.HandleWidth()
	var size Size

	if splitter.Orientation() == Horizontal {
		size.Width = handleWidth
	} else {
		size.Height = handleWidth
	}

	return size
}

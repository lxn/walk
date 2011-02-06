// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

type Page struct {
	info  *PageInfo
	parts []part
}

func (page *Page) addPart(p part) {
	page.parts = append(page.parts, p)
}

func (page *Page) Info() *PageInfo {
	return page.info
}

func (page *Page) Draw(canvas *Canvas) os.Error {
	for _, part := range page.parts {
		if err := part.Draw(canvas); err != nil {
			return err
		}
	}

	return nil
}

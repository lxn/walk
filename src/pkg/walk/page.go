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
	count := len(page.parts)
	if count == cap(page.parts) {
		parts := make([]part, count, count*2)
		copy(parts, page.parts)
		page.parts = parts
	}

	page.parts = page.parts[0 : count+1]
	page.parts[count] = p
}

func (page *Page) Info() *PageInfo {
	return page.info
}

func (page *Page) Draw(surface *Surface) os.Error {
	for _, part := range page.parts {
		err := part.Draw(surface)
		if err != nil {
			return err
		}
	}

	return nil
}

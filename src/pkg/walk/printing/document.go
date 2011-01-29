// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package printing

import (
	"os"
	"syscall"
	"unsafe"
	"utf8"
)

import (
	"walk/drawing"
	. "walk/winapi/gdi32"
)

type part interface {
	Bounds() drawing.Rectangle
	Draw(surface *drawing.Surface) os.Error
}

type item interface {
	PartCount() int
	Part(i int) part
	NextPartMinSize() drawing.Size
	PreferredSize() drawing.Size
	AddNewPart(surface *drawing.Surface, bounds drawing.Rectangle) (part part, more bool, err os.Error)
	Dispose()
}

type Document struct {
	name         string
	items        []item
	pages        []*Page
	nextPageInfo *PageInfo
	insertBounds drawing.Rectangle // Current bounds for pagination
}

func NewDocument(name string) *Document {
	return &Document{
		name:         name,
		nextPageInfo: NewPrinterInfo().pageInfo,
		pages:        make([]*Page, 0, 8),
	}
}

func (doc *Document) Dispose() {
	for _, x := range doc.items {
		x.Dispose()
	}
}

func (doc *Document) addItem(x item) {
	count := len(doc.items)
	if count == cap(doc.items) {
		items := make([]item, count, count*2)
		copy(items, doc.items)
		doc.items = items
	}

	doc.items = doc.items[0 : count+1]
	doc.items[count] = x
}

func (doc *Document) currentPage() *Page {
	if len(doc.pages) == 0 {
		return nil
	}

	return doc.pages[len(doc.pages)-1]
}

func (doc *Document) NextPageInfo() *PageInfo {
	return doc.nextPageInfo
}

func (doc *Document) SetNextPageInfo(value *PageInfo) os.Error {
	if value == nil {
		return newError("value cannot be nil")
	}

	if len(doc.pages) > 0 {
		printerName, err := value.printerInfo.PrinterName()
		if err != nil {
			return err
		}
		prevPrinterName, err := doc.pages[0].info.printerInfo.PrinterName()
		if err != nil {
			return err
		}

		if printerName != prevPrinterName {
			return newError("switching printer not supported")
		}
	}

	doc.nextPageInfo = value

	return nil
}

func (doc *Document) InsertPageBreak() os.Error {
	if page := doc.currentPage(); page != nil && len(page.parts) == 0 {
		return nil
	}

	pageBounds, err := doc.pageBounds()
	if err != nil {
		return err
	}
	doc.insertBounds = pageBounds

	page := &Page{info: doc.nextPageInfo, parts: make([]part, 0, 4)}

	count := len(doc.pages)
	if count == cap(doc.pages) {
		pages := make([]*Page, count, count*2)
		copy(pages, doc.pages)
		doc.pages = pages
	}

	doc.pages = doc.pages[0 : count+1]
	doc.pages[count] = page

	return nil
}

func (doc *Document) PageCount() int {
	return len(doc.pages)
}

func (doc *Document) Page(i int) *Page {
	return doc.pages[i]
}

func (doc *Document) withSurface(f func(surface *drawing.Surface) os.Error) os.Error {
	hdc := doc.nextPageInfo.createDC()
	defer DeleteDC(hdc)

	surface, err := drawing.NewSurfaceFromHDC(hdc)
	if err != nil {
		return err
	}
	defer surface.Dispose()

	return f(surface)
}

func (doc *Document) pageBounds() (bounds drawing.Rectangle, err os.Error) {
	err = doc.withSurface(func(surface *drawing.Surface) os.Error {
		bounds = surface.Bounds()

		return nil
	})

	return
}

func (doc *Document) paginateItem(item item) (err os.Error) {
	err = doc.withSurface(func(surface *drawing.Surface) os.Error {
		pageBounds, err := doc.pageBounds()
		if err != nil {
			return err
		}

		for {
			bounds := doc.insertBounds
			preferredSize := item.PreferredSize()
			if preferredSize.Width > 0 {
				bounds.Width = preferredSize.Width
			}
			if preferredSize.Height > 0 {
				bounds.Height = preferredSize.Height
			}

			part, more, err := item.AddNewPart(surface, bounds)
			if err != nil {
				return err
			}
			if part == nil {
				// No room for next part. Maybe using an empty page?
				nextPartMinSize := item.NextPartMinSize()
				if pageBounds.Width < nextPartMinSize.Width ||
					pageBounds.Height < nextPartMinSize.Height {
					// It does not even fit on an empty page, so we give up.
					return newError("insufficient page size")
				}

				if err := doc.InsertPageBreak(); err != nil {
					return err
				}
			} else {
				page := doc.currentPage()
				page.addPart(part)

				partBounds := part.Bounds()
				doc.insertBounds.Y += partBounds.Height
				doc.insertBounds.Height -= partBounds.Height

				if more {
					if err := doc.InsertPageBreak(); err != nil {
						return err
					}
				} else {
					break
				}
			}
		}

		return nil
	})

	return nil
}

func (doc *Document) AddText(text string, font *drawing.Font, color drawing.Color, preferredSize drawing.Size, format drawing.DrawTextFormat) os.Error {
	item := &simpleTextItem{
		text:          utf8.NewString(text),
		font:          font,
		color:         color,
		preferredSize: preferredSize,
		format:        format,
		parts:         make([]*simpleTextPart, 0, 1),
	}

	return doc.paginateItem(item)
}

func (doc *Document) Print() (err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	hdc := doc.pages[0].info.createDC()
	defer DeleteDC(hdc)

	var di DOCINFO
	di.CbSize = unsafe.Sizeof(di)
	di.LpszDocName = syscall.StringToUTF16Ptr(doc.name)

	if StartDoc(hdc, &di) <= 0 {
		return newError("StartDoc failed")
	}
	defer func() {
		if EndDoc(hdc) <= 0 {
			err = newError("EndDoc failed")
		}
	}()

	surface, err := drawing.NewSurfaceFromHDC(hdc)
	if err != nil {
		return err
	}
	defer surface.Dispose()

	for i, page := range doc.pages {
		if i > 0 {
			// TODO: Only reset dc if required.
			page.info.resetDC(hdc)
		}

		if StartPage(hdc) <= 0 {
			return newError("StartPage failed")
		}

		err = page.Draw(surface)
		if err != nil {
			return err
		}

		if EndPage(hdc) <= 0 {
			return newError("EndPage failed")
		}
	}

	return nil
}

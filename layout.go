// Copyright 2019 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"sync"

	"github.com/lxn/win"
)

func createLayoutItemForWidget(widget Widget) LayoutItem {
	ctx := newLayoutContext(widget.Handle())

	return createLayoutItemForWidgetWithContext(widget, ctx)
}

func createLayoutItemForWidgetWithContext(widget Widget, ctx *LayoutContext) LayoutItem {
	var item LayoutItem

	if container, ok := widget.(Container); ok {
		if container.Layout() == nil {
			return nil
		}

		item = CreateLayoutItemsForContainerWithContext(container, ctx)
	} else {
		item = widget.CreateLayoutItem(ctx)
	}

	lib := item.AsLayoutItemBase()
	lib.ctx = ctx
	lib.handle = widget.Handle()
	lib.visible = widget.AsWidgetBase().visible
	lib.geometry = widget.AsWidgetBase().geometry
	lib.geometry.Alignment = widget.Alignment()
	lib.geometry.MinSize = widget.MinSizePixels()
	lib.geometry.MaxSize = widget.MaxSizePixels()
	lib.geometry.ConsumingSpaceWhenInvisible = widget.AlwaysConsumeSpace()

	return item
}

func CreateLayoutItemsForContainer(container Container) ContainerLayoutItem {
	ctx := newLayoutContext(container.Handle())

	return CreateLayoutItemsForContainerWithContext(container, ctx)
}

func CreateLayoutItemsForContainerWithContext(container Container, ctx *LayoutContext) ContainerLayoutItem {
	var containerItem ContainerLayoutItem
	var clib *ContainerLayoutItemBase

	layout := container.Layout()
	if layout == nil || container.Children().Len() == 0 {
		layout = NewHBoxLayout()
		layout.SetMargins(Margins{})
	}

	if widget, ok := container.(Widget); ok {
		containerItem = widget.CreateLayoutItem(ctx).(ContainerLayoutItem)
	} else {
		containerItem = layout.CreateLayoutItem(ctx)
	}

	clib = containerItem.AsContainerLayoutItemBase()
	clib.ctx = ctx
	clib.handle = container.Handle()
	cb := container.AsContainerBase()
	clib.visible = cb.visible
	clib.geometry = cb.geometry
	clib.geometry.ConsumingSpaceWhenInvisible = cb.AlwaysConsumeSpace()

	if lb := layout.asLayoutBase(); lb != nil {
		clib.alignment = lb.alignment
		clib.margins96dpi = lb.margins96dpi
		clib.spacing96dpi = lb.spacing96dpi
	}

	if len(clib.children) == 0 {
		children := container.Children()
		count := children.Len()

		for i := 0; i < count; i++ {
			item := createLayoutItemForWidgetWithContext(children.At(i), ctx)
			if item != nil {
				lib := item.AsLayoutItemBase()
				lib.ctx = ctx
				lib.parent = containerItem

				clib.children = append(clib.children, item)
			}
		}
	}

	return containerItem
}

func startLayoutPerformer(form Form) (performLayout chan ContainerLayoutItem, layoutResults chan []LayoutResult, inSizeLoop chan bool, updateStopwatch chan *stopwatch, quit chan struct{}) {
	performLayout = make(chan ContainerLayoutItem)
	layoutResults = make(chan []LayoutResult)
	inSizeLoop = make(chan bool)
	updateStopwatch = make(chan *stopwatch)
	quit = make(chan struct{})

	var stopwatch *stopwatch

	go func() {
		sizing := false
		busy := false
		var cancel chan struct{}
		done := make(chan []LayoutResult)

		for {
			select {
			case root := <-performLayout:
				if busy {
					close(cancel)
				}

				busy = true
				cancel = make(chan struct{})

				go layoutTree(root, root.Geometry().ClientSize, cancel, done, stopwatch)

			case results := <-done:
				busy = false
				if cancel != nil {
					close(cancel)
					cancel = nil
				}

				if sizing {
					layoutResults <- results
				} else {
					form.AsFormBase().synchronizeLayout(&formLayoutResult{form, stopwatch, results})
				}

			case sizing = <-inSizeLoop:

			case stopwatch = <-updateStopwatch:

			case <-quit:
				close(performLayout)
				close(layoutResults)
				close(inSizeLoop)
				close(updateStopwatch)
				if cancel != nil {
					close(cancel)
				}
				close(done)
				close(quit)
				return
			}
		}
	}()

	return
}

// layoutTree lays out tree. size parameter is in native pixels.
func layoutTree(root ContainerLayoutItem, size Size, cancel chan struct{}, done chan []LayoutResult, stopwatch *stopwatch) {
	const minSizeCacheSubject = "layoutTree - populating min size cache"

	if stopwatch != nil {
		stopwatch.Start(minSizeCacheSubject)
	}

	// Populate some caches now, so we later need only read access to them from multiple goroutines.
	ctx := root.Context()

	populateContextForItem := func(item LayoutItem) {
		ctx.layoutItem2MinSizeEffective[item] = minSizeEffective(item)
	}

	var populateContextForContainer func(container ContainerLayoutItem)
	populateContextForContainer = func(container ContainerLayoutItem) {
		for _, child := range container.AsContainerLayoutItemBase().children {
			if cli, ok := child.(ContainerLayoutItem); ok {
				populateContextForContainer(cli)
			} else {
				populateContextForItem(child)
			}
		}

		populateContextForItem(container)
	}

	populateContextForContainer(root)

	if stopwatch != nil {
		stopwatch.Stop(minSizeCacheSubject)
	}

	const layoutSubject = "layoutTree - computing layout"

	if stopwatch != nil {
		stopwatch.Start(layoutSubject)
	}

	results := make(chan LayoutResult)
	finished := make(chan struct{})

	go func() {
		defer func() {
			close(results)
			close(finished)
		}()

		var wg sync.WaitGroup

		var layoutSubtree func(container ContainerLayoutItem, size Size)
		layoutSubtree = func(container ContainerLayoutItem, size Size) {
			wg.Add(1)

			go func() {
				defer wg.Done()

				clib := container.AsContainerLayoutItemBase()

				clib.geometry.ClientSize = size

				items := container.PerformLayout()

				select {
				case <-cancel:
					return

				case results <- LayoutResult{container, items}:
				}

				for _, item := range items {
					select {
					case <-cancel:
						return

					default:
					}

					item.Item.Geometry().Size = item.Bounds.Size()

					if childContainer, ok := item.Item.(ContainerLayoutItem); ok {
						layoutSubtree(childContainer, item.Bounds.Size())
					}
				}
			}()
		}

		layoutSubtree(root, size)

		wg.Wait()

		select {
		case <-cancel:
			return

		case finished <- struct{}{}:
		}
	}()

	var layoutResults []LayoutResult

	for {
		select {
		case result := <-results:
			layoutResults = append(layoutResults, result)

		case <-finished:
			if stopwatch != nil {
				stopwatch.Stop(layoutSubject)
			}

			done <- layoutResults
			return

		case <-cancel:
			if stopwatch != nil {
				stopwatch.Cancel(layoutSubject)
			}
			return
		}
	}
}

func applyLayoutResults(results []LayoutResult, stopwatch *stopwatch) error {
	if stopwatch != nil {
		const subject = "applyLayoutResults"
		stopwatch.Start(subject)
		defer stopwatch.Stop(subject)
	}

	var form Form

	for _, result := range results {
		if len(result.items) == 0 {
			continue
		}

		hdwp := win.BeginDeferWindowPos(int32(len(result.items)))
		if hdwp == 0 {
			return lastError("BeginDeferWindowPos")
		}

		var maybeInvalidate bool
		if wnd := windowFromHandle(result.container.Handle()); wnd != nil {
			if ctr, ok := wnd.(Container); ok {
				if cb := ctr.AsContainerBase(); cb != nil {
					maybeInvalidate = cb.hasComplexBackground()
				}
			}
		}

		for _, ri := range result.items {
			if ri.Item.Handle() != 0 {
				window := windowFromHandle(ri.Item.Handle())
				if window == nil {
					continue
				}

				if form == nil {
					if form = window.Form(); form != nil {
						defer func() {
							if focusedWindow := windowFromHandle(win.GetFocus()); focusedWindow == nil || focusedWindow == form || focusedWindow.Form() != form {
								form.AsFormBase().clientComposite.focusFirstCandidateDescendant()
							}
						}()
					}
				}

				widget := window.(Widget)

				oldBounds := widget.BoundsPixels()

				if ri.Bounds == oldBounds {
					continue
				}

				if ri.Bounds.X == oldBounds.X && ri.Bounds.Y == oldBounds.Y && ri.Bounds.Width == oldBounds.Width {
					if _, ok := widget.(*ComboBox); ok {
						if ri.Bounds.Height == oldBounds.Height+1 {
							continue
						}
					} else if ri.Bounds.Height == oldBounds.Height {
						continue
					}
				}

				if maybeInvalidate {
					if ri.Bounds.Width == oldBounds.Width && ri.Bounds.Height == oldBounds.Height && (ri.Bounds.X != oldBounds.X || ri.Bounds.Y != oldBounds.Y) {
						widget.Invalidate()
					}
				}

				if hdwp = win.DeferWindowPos(
					hdwp,
					ri.Item.Handle(),
					0,
					int32(ri.Bounds.X),
					int32(ri.Bounds.Y),
					int32(ri.Bounds.Width),
					int32(ri.Bounds.Height),
					win.SWP_NOACTIVATE|win.SWP_NOOWNERZORDER|win.SWP_NOZORDER); hdwp == 0 {

					return lastError("DeferWindowPos")
				}

				if widget.GraphicsEffects().Len() == 0 {
					continue
				}

				widget.AsWidgetBase().invalidateBorderInParent()
			}
		}

		if !win.EndDeferWindowPos(hdwp) {
			return lastError("EndDeferWindowPos")
		}
	}

	return nil
}

// Margins define margins in 1/96" units or native pixels.
type Margins struct {
	HNear, VNear, HFar, VFar int
}

func (m Margins) isZero() bool {
	return m.HNear == 0 && m.HFar == 0 && m.VNear == 0 && m.VFar == 0
}

type Layout interface {
	Container() Container
	SetContainer(value Container)
	Margins() Margins
	SetMargins(value Margins) error
	Spacing() int
	SetSpacing(value int) error
	CreateLayoutItem(ctx *LayoutContext) ContainerLayoutItem
	asLayoutBase() *LayoutBase
}

type LayoutBase struct {
	layout       Layout
	container    Container
	margins96dpi Margins
	margins      Margins // in native pixels
	spacing96dpi int
	spacing      int // in native pixels
	alignment    Alignment2D
	resetNeeded  bool
	dirty        bool
}

func (l *LayoutBase) asLayoutBase() *LayoutBase {
	return l
}

func (l *LayoutBase) Container() Container {
	return l.container
}

func (l *LayoutBase) SetContainer(value Container) {
	if value == l.container {
		return
	}

	if l.container != nil {
		l.container.SetLayout(nil)
	}

	l.container = value

	if value != nil && value.Layout() != l.layout {
		value.SetLayout(l.layout)
	}

	l.updateMargins()
	l.updateSpacing()

	if l.container != nil {
		l.container.RequestLayout()
	}
}

func (l *LayoutBase) Margins() Margins {
	return l.margins96dpi
}

func (l *LayoutBase) SetMargins(value Margins) error {
	if value == l.margins96dpi {
		return nil
	}

	if value.HNear < 0 || value.VNear < 0 || value.HFar < 0 || value.VFar < 0 {
		return newError("margins must be positive")
	}

	l.margins96dpi = value

	l.updateMargins()

	if l.container != nil {
		l.container.RequestLayout()
	}

	return nil
}

func (l *LayoutBase) Spacing() int {
	return l.spacing96dpi
}

func (l *LayoutBase) SetSpacing(value int) error {
	if value == l.spacing96dpi {
		return nil
	}

	if value < 0 {
		return newError("spacing cannot be negative")
	}

	l.spacing96dpi = value

	l.updateSpacing()

	if l.container != nil {
		l.container.RequestLayout()
	}

	return nil
}

func (l *LayoutBase) updateMargins() {
	if l.container != nil {
		l.margins = MarginsFrom96DPI(l.margins96dpi, l.container.AsWindowBase().DPI())
	}
}

func (l *LayoutBase) updateSpacing() {
	if l.container != nil {
		l.spacing = IntFrom96DPI(l.spacing96dpi, l.container.AsWindowBase().DPI())
	}
}

func (l *LayoutBase) Alignment() Alignment2D {
	return l.alignment
}

func (l *LayoutBase) SetAlignment(alignment Alignment2D) error {
	if alignment != l.alignment {
		if alignment < AlignHVDefault || alignment > AlignHFarVFar {
			return newError("invalid Alignment value")
		}

		l.alignment = alignment

		if l.container != nil {
			l.container.RequestLayout()
		}
	}

	return nil
}

type IdealSizer interface {
	// IdealSize returns ideal window size in native pixels.
	IdealSize() Size
}

type MinSizer interface {
	// MinSize returns minimum window size in native pixels.
	MinSize() Size
}

type MinSizeForSizer interface {
	// MinSize returns minimum window size for given size. Both sizes are in native pixels.
	MinSizeForSize(size Size) Size
}

type HeightForWidther interface {
	HasHeightForWidth() bool

	// HeightForWidth returns appropriate height if element has given width. width parameter and
	// return value are in native pixels.
	HeightForWidth(width int) int
}

type LayoutContext struct {
	layoutItem2MinSizeEffective map[LayoutItem]Size // in native pixels
	dpi                         int
}

func (ctx *LayoutContext) DPI() int {
	return ctx.dpi
}

func newLayoutContext(handle win.HWND) *LayoutContext {
	return &LayoutContext{
		layoutItem2MinSizeEffective: make(map[LayoutItem]Size),
		dpi:                         int(win.GetDpiForWindow(handle)),
	}
}

type LayoutItem interface {
	AsLayoutItemBase() *LayoutItemBase
	Context() *LayoutContext
	Handle() win.HWND
	Geometry() *Geometry
	Parent() ContainerLayoutItem
	Visible() bool
	LayoutFlags() LayoutFlags
}

type ContainerLayoutItem interface {
	LayoutItem
	MinSizer
	MinSizeForSizer
	HeightForWidther
	AsContainerLayoutItemBase() *ContainerLayoutItemBase

	// MinSizeEffectiveForChild returns minimum effective size for a child in native pixels.
	MinSizeEffectiveForChild(child LayoutItem) Size

	PerformLayout() []LayoutResultItem
	Children() []LayoutItem
	containsHandle(handle win.HWND) bool
}

type LayoutItemBase struct {
	ctx      *LayoutContext
	handle   win.HWND
	geometry Geometry
	parent   ContainerLayoutItem
	visible  bool
}

func (lib *LayoutItemBase) AsLayoutItemBase() *LayoutItemBase {
	return lib
}

func (lib *LayoutItemBase) Context() *LayoutContext {
	return lib.ctx
}

func (lib *LayoutItemBase) Handle() win.HWND {
	return lib.handle
}

func (lib *LayoutItemBase) Geometry() *Geometry {
	return &lib.geometry
}

func (lib *LayoutItemBase) Parent() ContainerLayoutItem {
	return lib.parent
}

func (lib *LayoutItemBase) Visible() bool {
	return lib.visible
}

type ContainerLayoutItemBase struct {
	LayoutItemBase
	children     []LayoutItem
	margins96dpi Margins
	spacing96dpi int
	alignment    Alignment2D
}

func (clib *ContainerLayoutItemBase) AsContainerLayoutItemBase() *ContainerLayoutItemBase {
	return clib
}

var clibMinSizeEffectiveForChildMutex sync.Mutex

func (clib *ContainerLayoutItemBase) MinSizeEffectiveForChild(child LayoutItem) Size {
	// NOTE: This map is pre-populated in startLayoutTree before performing layout.
	// For other usages it is not pre-populated and we assume this method will then
	// be called from the main goroutine exclusively.
	// If we want to do concurrent size measurement, we will need to pre-populate also.

	// FIXME: There seems to be a bug in pre-population, so we use a mutex for now.

	clibMinSizeEffectiveForChildMutex.Lock()

	if clib.ctx != nil {
		if size, ok := clib.ctx.layoutItem2MinSizeEffective[child]; ok {
			clibMinSizeEffectiveForChildMutex.Unlock()
			return size
		}
	}

	if clib.ctx == nil {
		if clib.parent == nil {
			clib.ctx = newLayoutContext(clib.Handle())
		} else {
			clib.ctx = clib.parent.Context()
		}
	}

	child.AsLayoutItemBase().ctx = clib.ctx

	clibMinSizeEffectiveForChildMutex.Unlock()

	size := minSizeEffective(child)

	clibMinSizeEffectiveForChildMutex.Lock()

	if clib.ctx != nil {
		clib.ctx.layoutItem2MinSizeEffective[child] = size
	}

	clibMinSizeEffectiveForChildMutex.Unlock()

	return size
}

func (clib *ContainerLayoutItemBase) Children() []LayoutItem {
	return clib.children
}

func (clib *ContainerLayoutItemBase) SetChildren(children []LayoutItem) {
	clib.children = children
}

func (clib *ContainerLayoutItemBase) containsHandle(handle win.HWND) bool {
	for _, item := range clib.children {
		if item.Handle() == handle {
			return true
		}
	}

	return false
}

func (clib *ContainerLayoutItemBase) HasHeightForWidth() bool {
	for _, child := range clib.children {
		if hfw, ok := child.(HeightForWidther); ok && hfw.HasHeightForWidth() {
			return true
		}
	}

	return false
}

type greedyLayoutItem struct {
	LayoutItemBase
}

func NewGreedyLayoutItem() LayoutItem {
	return new(greedyLayoutItem)
}

func (*greedyLayoutItem) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | GrowableHorz | GreedyHorz | ShrinkableVert | GrowableVert | GreedyVert
}

func (li *greedyLayoutItem) IdealSize() Size {
	return SizeFrom96DPI(Size{100, 100}, li.ctx.dpi)
}

func (li *greedyLayoutItem) MinSize() Size {
	return SizeFrom96DPI(Size{50, 50}, li.ctx.dpi)
}

type Geometry struct {
	Alignment                   Alignment2D
	MinSize                     Size // in native pixels
	MaxSize                     Size // in native pixels
	IdealSize                   Size // in native pixels
	Size                        Size // in native pixels
	ClientSize                  Size // in native pixels
	ConsumingSpaceWhenInvisible bool
}

type formLayoutResult struct {
	form      Form
	stopwatch *stopwatch
	results   []LayoutResult
}

type LayoutResult struct {
	container ContainerLayoutItem
	items     []LayoutResultItem
}

type LayoutResultItem struct {
	Item   LayoutItem
	Bounds Rectangle // in native pixels
}

func shouldLayoutItem(item LayoutItem) bool {
	if item == nil {
		return false
	}

	_, isSpacer := item.(*spacerLayoutItem)

	return isSpacer || item.Visible() || item.Geometry().ConsumingSpaceWhenInvisible
}

func itemsToLayout(allItems []LayoutItem) []LayoutItem {
	filteredItems := make([]LayoutItem, 0, len(allItems))

	for i := 0; i < cap(filteredItems); i++ {
		item := allItems[i]

		if !shouldLayoutItem(item) {
			continue
		}

		var idealSize Size
		if hfw, ok := item.(HeightForWidther); !ok || !hfw.HasHeightForWidth() {
			if is, ok := item.(IdealSizer); ok {
				idealSize = is.IdealSize()
			}
		}
		if idealSize.Width == 0 && idealSize.Height == 0 && item.LayoutFlags() == 0 {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems
}

func anyVisibleItemInHierarchy(item LayoutItem) bool {
	if item == nil || !item.Visible() {
		return false
	}

	if cli, ok := item.(ContainerLayoutItem); ok {
		for _, child := range cli.AsContainerLayoutItemBase().children {
			if anyVisibleItemInHierarchy(child) {
				return true
			}
		}
	} else if _, ok := item.(*spacerLayoutItem); !ok {
		return true
	}

	return false
}

// minSizeEffective returns minimum effective size in native pixels
func minSizeEffective(item LayoutItem) Size {
	geometry := item.Geometry()

	var s Size
	if msh, ok := item.(MinSizer); ok {
		s = msh.MinSize()
	} else if is, ok := item.(IdealSizer); ok {
		s = is.IdealSize()
	}

	size := maxSize(geometry.MinSize, s)

	max := geometry.MaxSize
	if max.Width > 0 && size.Width > max.Width {
		size.Width = max.Width
	}
	if max.Height > 0 && size.Height > max.Height {
		size.Height = max.Height
	}

	return size
}

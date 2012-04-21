package walk
import (
	"unsafe"
	"syscall"
)
import . "github.com/lxn/go-winapi"


const (
	LB_ADDSTRING = 0x0180

	LBS_NOTIFY = 0x0001
	LBS_SORT   = 0x0002
	
	LBS_STANDARD = LBS_NOTIFY | LBS_SORT | WS_VSCROLL | WS_BORDER

	LB_GETCURSEL = 0x188
	LB_GETTEXT = 0x0189
	LB_GETCOUNT = 0x18B
	LB_GETTEXTLEN = 0x018A
	LBN_SELCHANGE = 1
	LBN_DBLCLK = 2


)

type ListBox struct{
	WidgetBase
	Items                         []string
	maxItemTextWidth              int
	selectedIndexChangedPublisher EventPublisher
	dbClickedPublisher            EventPublisher
}

func NewListBox(parent Container)(*ListBox, error){
	lb := &ListBox{}
	err := initChildWidget(
		lb,
		parent,
		"LISTBOX",
		WS_TABSTOP | WS_VISIBLE | LBS_STANDARD,
		0)
	if err != nil{
		return nil, err
	}
	return lb, nil
}

func (*ListBox) origWndProcPtr() uintptr {
	return checkBoxOrigWndProcPtr
}

func (*ListBox) setOrigWndProcPtr(ptr uintptr) {
	checkBoxOrigWndProcPtr = ptr
}

func (*ListBox) LayoutFlags() LayoutFlags {
	return GrowableHorz | GrowableVert
}

func (this *ListBox) SetItems(items []string){
	this.Items = items
	//Should remove the original content?
	for _, item := range items{
		this.AddItem(item)
	}
}

func (this *ListBox) AddItem(item string){
	this.Items = append(this.Items, item)
	SendMessage (this.hWnd, LB_ADDSTRING, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(item))))
}

func (this *ListBox) calculateMaxItemTextWidth() int {
	hdc := GetDC(this.hWnd)
	if hdc == 0 {
		newError("GetDC failed")
		return -1
	}
	defer ReleaseDC(this.hWnd, hdc)

	hFontOld := SelectObject(hdc, HGDIOBJ(this.Font().handleForDPI(0)))
	defer SelectObject(hdc, hFontOld)

	var maxWidth int

	for _, item := range this.Items {
		var s SIZE
		str := syscall.StringToUTF16(item)

		if !GetTextExtentPoint32(hdc, &str[0], int32(len(str)-1), &s) {
			newError("GetTextExtentPoint32 failed")
			return -1
		}

		maxWidth = maxi(maxWidth, int(s.CX))
	}

	return maxWidth
}


func (this *ListBox) SizeHint() Size {

	defaultSize := this.dialogBaseUnitsToPixels(Size{50, 12})

	if this.Items != nil && this.maxItemTextWidth <= 0 {
		this.maxItemTextWidth = this.calculateMaxItemTextWidth()
	}

	// FIXME: Use GetThemePartSize instead of guessing
	w := maxi(defaultSize.Width, this.maxItemTextWidth+24)
	h := defaultSize.Height + 1

	return Size{w, h}	

}

func (this *ListBox) SelectedIndex() int{
	return int(SendMessage (this.hWnd, LB_GETCURSEL, 0, 0))
}

func (this *ListBox) SelectedItem() string{
	index := this.SelectedIndex()
	length := int(SendMessage(this.hWnd, LB_GETTEXTLEN, uintptr(index), 0)) + 1
	buffer := make([]uint16, length +1)
	SendMessage(this.hWnd, LB_GETTEXT, uintptr(index), uintptr(unsafe.Pointer(&buffer[0])))
	return syscall.UTF16ToString(buffer)
}

func (this *ListBox) SelectedIndexChanged() *Event{
	return this.selectedIndexChangedPublisher.Event()
}

func (this *ListBox) DBClicked() *Event{
	return this.dbClickedPublisher.Event()
}

func (this *ListBox) wndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint32(wParam)) {
		case LBN_SELCHANGE:
			this.selectedIndexChangedPublisher.Publish()
		case LBN_DBLCLK:
			this.dbClickedPublisher.Publish()
		}
	}

	return this.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}
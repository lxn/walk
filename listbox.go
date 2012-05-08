package walk
import (
	"unsafe"
	"syscall"
	"errors"
)
import . "github.com/lxn/go-winapi"

const (
	LB_ERR       = -1
	LB_ERRSPACE  = -2
)

type ListBox struct{
	WidgetBase
	maxItemTextWidth              int
	selectedIndexChangedPublisher EventPublisher
	dbClickedPublisher            EventPublisher
}

func NewListBox(parent Container)(*ListBox, error){
	//TODO: move to go-winapi/listbox
	//LBS_STANDARD := LBS_NOTIFY | LBS_SORT | WS_VSCROLL | WS_BORDER

	lb := &ListBox{}
	err := initChildWidget(
		lb,
		parent,
		"LISTBOX",
		WS_TABSTOP | WS_VISIBLE | LBS_NOTIFY | LBS_SORT | WS_VSCROLL | WS_BORDER,
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

func (this *ListBox) AddString(item string){
	SendMessage (this.hWnd, LB_ADDSTRING, 0, 
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(item))))
}

//If this parameter is -1, the string is added to the end of the list.
func (this *ListBox) InsertString(index int, item string) error{
	if index < -1{
		return errors.New("Invalid index")
	}
	
	ret := int(SendMessage(this.hWnd, LB_INSERTSTRING, uintptr(index), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(item)))))
	if ret == LB_ERRSPACE || ret == LB_ERR {
		return errors.New("LB_ERR or LB_ERRSPACE")
	}
	return nil
}

func (this *ListBox) DeleteString(index uint) error{
	ret := int(SendMessage(this.hWnd, LB_DELETESTRING, uintptr(index), 0))
	if ret == LB_ERR {
		return errors.New("LB_ERR")
	}
	return nil
}

func (this *ListBox) GetString(index uint) string{
	len := int(SendMessage(this.hWnd, LB_GETTEXTLEN, uintptr(index), 0))
	if len == LB_ERR{
		return ""
	}

	buf := make([]uint16, len + 1)
	_ = SendMessage(this.hWnd, LB_GETTEXT, uintptr(index), uintptr(unsafe.Pointer(&buf[0])))
	
	if len == LB_ERR{
		return ""
	}
	return syscall.UTF16ToString(buf)
}
	

func (this *ListBox) ResetContent(){
	SendMessage(this.hWnd, LB_RESETCONTENT, 0, 0)
}

//The return value is the number of items in the list box, 
//or LB_ERR (-1) if an error occurs.
func (this *ListBox) GetCount() (uint, error){
    retPtr := SendMessage(this.hWnd, LB_GETCOUNT, 0, 0)
	ret := int(retPtr)
	if ret == LB_ERR{
		return 0, errors.New("LB_ERR")
	}
	return uint(ret), nil
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

	count, _ := this.GetCount()
	var i uint
	for i = 0; i < count ; i++{
		item  := this.GetString(i)
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

	if this.maxItemTextWidth <= 0 {
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
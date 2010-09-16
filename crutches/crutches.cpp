#define _WIN32_WINNT 0x501 // XP
#define _WIN32_IE 0x501

#include "crutches.h"

#include <queue>

#include <commctrl.h>

using namespace std;

#define WM_RESIZE_KEY 0
#define WM_COMMAND_KEY 1
#define WM_CONTEXTMENU_KEY 2
#define WM_ITEMCHANGED_KEY 3
#define WM_ITEMACTIVATE_KEY 4
#define WM_CLOSE_KEY 5

UINT _resizeMsgId;
UINT _commandMsgId;
UINT _contextMenuMsgId;
UINT _itemChangedMsgId;
UINT _itemActivateMsgId;
UINT _closeMsgId;

queue<Message> msgQueue;


UINT WINAPI GetRegisteredMessageId(UINT key)
{
    switch (key) {
    case WM_RESIZE_KEY:
        if (0 == _resizeMsgId) {
            _resizeMsgId = RegisterWindowMessageW(L"resize_0b0f95e6-7ef7-4767-b484-940e7a3cf4f1");
        }
        return _resizeMsgId;

    case WM_COMMAND_KEY:
        if (0 == _commandMsgId) {
            _commandMsgId = RegisterWindowMessageW(L"command_442946bf-f806-434b-baa3-98439930eecd");
        }
        return _commandMsgId;

    case WM_CONTEXTMENU_KEY:
        if (0 == _contextMenuMsgId) {
            _contextMenuMsgId = RegisterWindowMessageW(L"contextmenu_50fe6189-a94b-4826-8dc3-48e179f89ffc");
        }
        return _contextMenuMsgId;

    case WM_ITEMCHANGED_KEY:
        if (0 == _itemChangedMsgId) {
            _itemChangedMsgId = RegisterWindowMessageW(L"itemchanged_b453604c-2195-4df2-8d9d-e6486d9bf73a");
        }
        return _itemChangedMsgId;

    case WM_ITEMACTIVATE_KEY:
        if (0 == _itemActivateMsgId) {
            _itemActivateMsgId = RegisterWindowMessageW(L"itemactivate_1ff744a7-ff21-464c-b0a7-786843b75976");
        }
        return _itemActivateMsgId;

    case WM_CLOSE_KEY:
        if (0 == _closeMsgId) {
            _closeMsgId = RegisterWindowMessageW(L"close_cd8f8d08-cdb1-42c0-9d8c-bcf63f4114d7");
        }
        return _closeMsgId;
    }

    return 0;
}

int WINAPI GetCustomMessage(Message* msg)
{
    if (0 == msg) {
        SetLastError(ERROR_INVALID_PARAMETER);
        return -1;
    }

    if (msgQueue.empty()) {
        return 0;
    }

    (*msg) = msgQueue.back();
    msgQueue.pop();

    return msgQueue.size() + 1;
}

bool EnqueueMessage(const Message &msg)
{
    if (msgQueue.size() >= 1000) {
        return false;
    }

    msgQueue.push(msg);

    return true;
}

LRESULT CALLBACK ContainerWndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam)
{
    switch(msg)
    {
    case WM_CLOSE:
        PostMessageW(hwnd, _closeMsgId, wParam, lParam);
        break;

    case WM_CONTEXTMENU:
        PostMessageW(hwnd, _contextMenuMsgId, wParam, lParam);

    case WM_COMMAND:        
        PostMessageW(hwnd, _commandMsgId, wParam, lParam);
        break;

    case WM_NOTIFY:
        switch (((NMHDR*)(lParam))->code) {
        case LVN_ITEMCHANGED: {
                NMLISTVIEW* nmlv = (NMLISTVIEW*)lParam;

                Message m;
                m.hwnd = nmlv->hdr.hwndFrom;
                m.msg = _itemChangedMsgId;
                m.wParam = 0;
                m.lParam = nmlv->iItem;

                EnqueueMessage(m);
                break;
            }

        case LVN_ITEMACTIVATE: {
                NMITEMACTIVATE* nmia = (NMITEMACTIVATE*)lParam;

                Message m;
                m.hwnd = nmia->hdr.hwndFrom;
                m.msg = _itemActivateMsgId;
                m.wParam = 0;
                m.lParam = nmia->iItem;

                EnqueueMessage(m);
                break;
            }
        }

        break;

    case WM_SIZE:
    case WM_SIZING: {
            LRESULT ret = DefDlgProcW(hwnd, msg, wParam, lParam);
            PostMessageW(hwnd, _resizeMsgId, 0, 0);
            return ret;
        }

    default:
        return DefDlgProcW(hwnd, msg, wParam, lParam);
    }
    return 0;
}

ATOM WINAPI RegisterWindowClass(HINSTANCE hInst)
{
    WNDCLASSEX wcButton;
    WNDCLASSEX wc;

    wc.hCursor       = LoadCursor(hInst, IDC_ARROW);

    if (GetClassInfoExW(hInst, L"BUTTON", &wcButton) != FALSE) {
        wc.hCursor = wcButton.hCursor;
    }

    wc.cbSize        = sizeof(WNDCLASSEX);
    wc.style         = 0;
    wc.lpfnWndProc   = ContainerWndProc;
    wc.cbClsExtra    = 0;
    wc.cbWndExtra    = DLGWINDOWEXTRA;
    wc.hInstance     = hInst;
    wc.hIcon         = LoadIconW(hInst, IDI_APPLICATION);
//    wc.hCursor       = LoadCursor(hInst, IDC_ARROW);
    wc.hbrBackground = (HBRUSH)(COLOR_3DFACE+1);
    wc.lpszMenuName  = 0;
    wc.lpszClassName = L"Container_WindowClass";
    wc.hIconSm       = LoadIcon(hInst, IDI_APPLICATION);

    return RegisterClassExW(&wc);
}

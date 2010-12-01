// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "crutch.h"

uintptr crutches·stdcall_return(uintptr retval, void* addr_of_caller_first_arg, uint32 no_of_caller_args);

void *DefDlgProcW;
void *GetClassInfoExW;
void *LoadCursorW;
void *LoadIconW;
void *PostMessageW;
void *RegisterClassExW;
void *RegisterWindowMessageW;

void
·initcrutch(void) {
    // TODO: make sure that 'runtime·osinit()' was already called,
    // or use their approach with 'get_proc_addr2()'
    DefDlgProcW = runtime·get_proc_addr("user32.dll", "DefDlgProcW");
    GetClassInfoExW = runtime·get_proc_addr("user32.dll", "GetClassInfoExW");
    LoadCursorW = runtime·get_proc_addr("user32.dll", "LoadCursorW");
    LoadIconW = runtime·get_proc_addr("user32.dll", "LoadIconW");
    PostMessageW = runtime·get_proc_addr("user32.dll", "PostMessageW");
    RegisterClassExW = runtime·get_proc_addr("user32.dll", "RegisterClassExW");
    RegisterWindowMessageW = runtime·get_proc_addr("user32.dll", "RegisterWindowMessageW");
}

enum {
    WM_RESIZE_KEY = 0,
    WM_COMMAND_KEY = 1,
    WM_CONTEXTMENU_KEY = 2,
    WM_ITEMCHANGED_KEY = 3,
    WM_ITEMACTIVATE_KEY = 4,
    WM_CLOSE_KEY = 5,
    NUM_WM_KEYS = 6
};

uint32 msgIds[NUM_WM_KEYS];
byte* msgClsids[NUM_WM_KEYS] = {
    (byte*)L"resize_0b0f95e6-7ef7-4767-b484-940e7a3cf4f1",
    (byte*)L"command_442946bf-f806-434b-baa3-98439930eecd",
    (byte*)L"contextmenu_50fe6189-a94b-4826-8dc3-48e179f89ffc",
    (byte*)L"itemchanged_b453604c-2195-4df2-8d9d-e6486d9bf73a",
    (byte*)L"itemactivate_1ff744a7-ff21-464c-b0a7-786843b75976",
    (byte*)L"close_cd8f8d08-cdb1-42c0-9d8c-bcf63f4114d7",
};

typedef struct NMHDR NMHDR;
struct NMHDR {
  HANDLE   hwndFrom;
  uint32   idFrom;
  uint32   code;
};

typedef struct POINT POINT;
struct POINT {
  uint32  x;
  uint32  y;
};

typedef struct NMLISTVIEW NMLISTVIEW;
struct NMLISTVIEW {
  NMHDR    hdr;
  int32    iItem;
  int32    iSubItem;
  uint32   uNewState;
  uint32   uOldState;
  uint32   uChanged;
  POINT    ptAction;
  uint32   lParam;
};

typedef struct NMITEMACTIVATE NMITEMACTIVATE;
struct NMITEMACTIVATE {
  NMHDR    hdr;
  int32    iItem;
  int32    iSubItem;
  uint32   uNewState;
  uint32   uOldState;
  uint32   uChanged;
  POINT    ptAction;
  uint32   lParam;
  uint32   uKeyFlags;
};

#define WM_CLOSE                        0x0010
#define WM_CONTEXTMENU                  0x007B
#define WM_COMMAND                      0x0111
#define WM_NOTIFY                       0x004e
#define WM_SIZE                         0x0005
#define WM_SIZING                       0x0214

#define LVN_ITEMCHANGED         ((uint32)-101U)
#define LVN_ITEMACTIVATE        ((uint32)-114U)


#pragma textflag 7
static uint32
internalContainerWndProc(HANDLE hwnd, uint32 uMsg, uint32 wParam, uint32 lParam) {
    switch(uMsg)
    {
    case WM_CLOSE:
        crutches·wildcall(PostMessageW, 4, hwnd, msgIds[WM_CLOSE_KEY], wParam, lParam);
        break;

    case WM_CONTEXTMENU:
        crutches·wildcall(PostMessageW, 4, hwnd, msgIds[WM_CONTEXTMENU_KEY], wParam, lParam);
        // FIXME: is the lack of "break;" here intentional?

    case WM_COMMAND:
        crutches·wildcall(PostMessageW, 4, hwnd, msgIds[WM_COMMAND_KEY], wParam, lParam);
        break;

    case WM_NOTIFY:
        switch (((NMHDR*)(lParam))->code) {
        case LVN_ITEMCHANGED: {
                NMLISTVIEW* nmlv = (NMLISTVIEW*)lParam;

                Message m;
                m.hwnd = nmlv->hdr.hwndFrom;
                m.msg = msgIds[WM_ITEMCHANGED_KEY];
                m.wParam = 0;
                m.lParam = nmlv->iItem;

                crutches·nosplit_enqueue(&queue, &m);
                break;
            }

        case LVN_ITEMACTIVATE: {
                NMITEMACTIVATE* nmia = (NMITEMACTIVATE*)lParam;

                Message m;
                m.hwnd = nmia->hdr.hwndFrom;
                m.msg = msgIds[WM_ITEMACTIVATE_KEY];
                m.wParam = 0;
                m.lParam = nmia->iItem;

                crutches·nosplit_enqueue(&queue, &m);
                break;
            }
        }

        break;

    case WM_SIZE:
    case WM_SIZING: {
            uint32 ret = (uint32)crutches·wildcall(DefDlgProcW, 4, hwnd, uMsg, wParam, lParam);
            crutches·wildcall(PostMessageW, 4, hwnd, msgIds[WM_RESIZE_KEY], 0, 0);
            return ret;
        }

    default:
        return (uint32)crutches·wildcall(DefDlgProcW, 4, hwnd, uMsg, wParam, lParam);
    }
    return 0;
}

#pragma textflag 7
uint32
crutches·containerWndProc(HANDLE hwnd, uint32 uMsg, uint32 wParam, uint32 lParam) {
    return (uint32)crutches·stdcall_return(
        internalContainerWndProc(hwnd, uMsg, wParam, lParam),
        &hwnd,
        4);
}

typedef struct WNDCLASSEX WNDCLASSEX;
struct WNDCLASSEX {
  uint32    cbSize;
  uint32    style;
  void*     lpfnWndProc;
  int32     cbClsExtra;
  int32     cbWndExtra;
  HANDLE    hInstance;
  HANDLE    hIcon;
  HANDLE    hCursor;
  HANDLE    hbrBackground;
  byte*     lpszMenuName;
  byte*     lpszClassName;
  HANDLE    hIconSm;
};

static uintptr
Syscall3(void* func, uint32 arg0, uint32 arg1, uint32 arg2) {
    StdcallParams p;
    p.fn = func;
    p.args[0] = arg0;
    p.args[1] = arg1;
    p.args[2] = arg2;
    p.n = 3;
    runtime·syscall(&p);
    return p.r;
}

void
·registerWindowClass(uintptr hInst, uintptr r1) {
    WNDCLASSEX wcButton;
    WNDCLASSEX wc;

    wc.hCursor       = Syscall3(LoadCursorW, hInst, 32512, 0); //IDC_ARROW

    if(Syscall3(GetClassInfoExW, hInst, (uint32)L"BUTTON", (uint32)&wcButton)) {
        wc.hCursor = wcButton.hCursor;
    }

    wc.cbSize        = sizeof(WNDCLASSEX);
    wc.style         = 0;
    wc.lpfnWndProc   = (void*)crutches·containerWndProc;
    wc.cbClsExtra    = 0;
    wc.cbWndExtra    = 30; //DLGWINDOWEXTRA;
    wc.hInstance     = hInst;
    wc.hIcon         = Syscall3(LoadIconW, hInst, 32512, 0); //IDI_APPLICATION
//    wc.hCursor       = LoadCursor(hInst, IDC_ARROW);
    wc.hbrBackground = 15+1; //COLOR_3DFACE+1
    wc.lpszMenuName  = 0;
    wc.lpszClassName = (byte*)L"Container_WindowClass";
    wc.hIconSm       = Syscall3(LoadIconW, hInst, 32512, 0); //IDI_APPLICATION

    r1 = Syscall3(RegisterClassExW, (uint32)&wc, 0, 0);
    FLUSH(&r1);
}

void
·getCustomMessage(uintptr msgPointer, uintptr r1) {
    r1 = crutches·cansplit_dequeue(&queue, (Message*)msgPointer);
    FLUSH(&r1);
}

void
·getRegisteredMessageId(uint32 key, uint32 r1)
{
    if( key >= NUM_WM_KEYS ) {
        r1 = 0;
        FLUSH(&r1);
        return;
    }

    if( 0 == msgIds[key] ) {
        msgIds[key] = Syscall3(
            RegisterWindowMessageW,
            (uintptr)msgClsids[key],
            0,
            0);
    }

    r1 = msgIds[key];
    FLUSH(&r1);
    return;
}

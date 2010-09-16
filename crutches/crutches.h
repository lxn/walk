// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef CRUTCHES_H
#define CRUTCHES_H

#include <windows.h>

struct Message {
    HWND hwnd;
    UINT msg;
    WPARAM wParam;
    LPARAM lParam;
};

extern "C"
{
    ATOM WINAPI RegisterWindowClass(HINSTANCE hInst);

    UINT WINAPI GetRegisteredMessageId(UINT key);

    int WINAPI GetCustomMessage(Message* msg);
}

#endif // CRUTCHES_H

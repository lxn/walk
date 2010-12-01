// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "os.h"

void* crutches·wildcall(void* fn, int32 count, ...);

// some WinApi struct
// internal structure doesn't matter for us
// 15*uint32 should be enough based on quick overview of the fields
typedef byte PCRITICAL_SECTION[4*15];
typedef PCRITICAL_SECTION CriticalSection;

typedef uintptr HANDLE;

typedef struct Message Message;
struct Message {
	HANDLE hwnd;
	uint32 msg;
	uintptr wParam;
	uintptr lParam;
};

typedef struct Queue Queue;
struct Queue {
    Message *data;
    int32 capacity;
    int32 head, len;
    CriticalSection lock;
    uintptr hasdataEvent;
};

extern Queue queue;

void crutches·nosplit_enqueue(Queue* q, Message* msg);
int32 crutches·nosplit_dequeue(Queue* q, Message* msg);

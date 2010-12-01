// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "crutch.h"

extern void* crutches·wildcall_raw(void* fn, int32 count, uintptr* args);

#pragma textflag 7
void*
crutches·wildcall(void* fn, int32 count, ...) {
    return crutches·wildcall_raw(fn, count, (uintptr*)(&count+1));
}

void* InitializeCriticalSectionAndSpinCount;
void* EnterCriticalSection;
void* LeaveCriticalSection;
//void* DeleteCriticalSection;
void* CreateEvent;
void* WaitForSingleObject;
void* SetEvent;

void
newqueue(Queue* q, Message* buffer, int32 capacity) {
    q->data = buffer;
    q->capacity = capacity;

    StdcallParams p;
    p.fn = InitializeCriticalSectionAndSpinCount;
    p.args[0] = (uintptr)q->lock;
    p.args[1] = 0;
    p.args[2] = 0;
    p.n = 3;
    runtime·syscall(&p);
    // TODO: handle error in InitializeCriticalSection somehow... (panic?)

    p.fn = CreateEvent;
    p.args[0] = 0;
    p.args[1] = 0;
    p.args[2] = 0;
    p.args[3] = 0; // note: zeroes up to here are important
    p.args[4] = 0;
    p.args[5] = 0;
    p.n = 6;
    runtime·syscall(&p);
    q->hasdataEvent = p.r;
    //q->hasdataEvent = CreateEvent(0, 0, 0, 0);
    // TODO: handle error in CreateEvent (panic?)

    q->head = q->len = 0;
}

#pragma textflag 7
void
crutches·nosplit_enqueue(Queue* q, Message* msg) {
    crutches·wildcall(EnterCriticalSection, 1, q->lock);
    if(q->len < q->capacity) {
        int32 qtail = (q->head+q->len) % q->capacity;
        q->data[qtail] = *msg;
        q->len++;
    }
    // TODO: remove the event if the evented approach is not needed
    crutches·wildcall(SetEvent, 1, q->hasdataEvent);
    crutches·wildcall(LeaveCriticalSection, 1, q->lock);
    // FIXME: handle full queue
    // IMPLEM. NOTE: take care not to "empty" the list by overflowing
}

#pragma textflag 7
int32
crutches·nosplit_dequeue(Queue* q, Message* msg) {
    if(0 == msg) {
        // TODO: SetLastError(ERROR_INVALID_PARAMETER);
        return -1;
    }

    crutches·wildcall(EnterCriticalSection, 1, q->lock);
    // TODO: error if qlen < 0
    int32 result = q->len;
    if(q->len > 0) {
        *msg = q->data[q->head];
        q->len--;
        q->head++;
        q->head = q->head % q->capacity;
    }
    crutches·wildcall(LeaveCriticalSection, 1, q->lock);
    return result;
/*
    crutches·wildcall(EnterCriticalSection, 1, q->lock);
    for(;;) {
        if(q->len > 0) {
            int32 msg = qdata[q->head];
            q->len--;
            q->head++;
            q->head = q->head % q->capacity;
            crutches·wildcall(LeaveCriticalSection, 1, q->lock);
            return msg;
        } else {
            crutches·wildcall(LeaveCriticalSection, 1, q->lock);
            crutches·wildcall(WaitForSingleObject, 2, q->hasdataEvent, 0xffffffff);
            crutches·wildcall(EnterCriticalSection, 1, q->lock);
        }
    }
*/
}

// TODO: build using some dynamic malloc to reduce executable size [?]
#define qqmaxsize 500
Message qqdata[qqmaxsize];
Queue queue;

void
·initqueue(void) {
    // TODO: make sure that 'runtime·osinit()' was already called,
    // or use their approach with 'get_proc_addr2()'
    InitializeCriticalSectionAndSpinCount = runtime·get_proc_addr("kernel32.dll", "InitializeCriticalSectionAndSpinCount");
    EnterCriticalSection = runtime·get_proc_addr("kernel32.dll", "EnterCriticalSection");
    LeaveCriticalSection = runtime·get_proc_addr("kernel32.dll", "LeaveCriticalSection");
    CreateEvent = runtime·get_proc_addr("kernel32.dll", "CreateEventA");
    WaitForSingleObject = runtime·get_proc_addr("kernel32.dll", "WaitForSingleObject");
    SetEvent = runtime·get_proc_addr("kernel32.dll", "SetEvent");

    newqueue(&queue, qqdata, qqmaxsize);
}


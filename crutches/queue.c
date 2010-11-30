// by Mateusz Czapliński
// <czapkofan@gmail.com>

#include "crutch.h"

extern void* crutches·wildcall_raw(void* fn, int32 count, uintptr* args);

#pragma textflag 7
void*
crutches·wildcall(void* fn, int32 count, ...) {
    return crutches·wildcall_raw(fn, count, (uintptr*)(&count+1));
}

// some WinApi struct
// internal structure doesn't matter for us
// 15*uint32 should be enough based on quick overview of the fields
typedef byte PCRITICAL_SECTION[4*15];
typedef PCRITICAL_SECTION CriticalSection;

void* InitializeCriticalSectionAndSpinCount;
void* EnterCriticalSection;
void* LeaveCriticalSection;
//void* DeleteCriticalSection;
void* CreateEvent;
void* WaitForSingleObject;
void* SetEvent;

// TODO: build using some dynamic malloc to reduce executable size [?]
#define qmaxsize 100
int32 qdata[qmaxsize];
int32 qhead, qlen;
CriticalSection qlock;
uintptr qhasdataevent;

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

    StdcallParams p;
    p.fn = InitializeCriticalSectionAndSpinCount;
    p.args[0] = (uintptr)qlock;
    p.args[1] = 0;
    p.args[2] = 0;
    p.n = 3;
    runtime·syscall(&p);
    // TODO: handle error in InitializeCriticalSection somehow... (panic?)
    
    p.fn = CreateEvent;
    p.args[0] = 0;
    p.args[1] = 0;
    p.args[2] = 0;
    p.args[3] = 0;
    p.args[4] = 0;
    p.args[5] = 0;
    p.n = 6;
    runtime·syscall(&p);
    qhasdataevent = p.r;
    //qhasdataevent = CreateEvent(0, 0, 0, 0);
    // TODO: handle error in CreateEvent (panic?)

    qhead = qlen = 0;
    
    // DEBUG/TEST
    qhead = 98;
    qdata[98] = 11;
    qdata[99] = 10;
    qdata[0] = 9;
    qdata[1] = 8;
    qlen = 4;
    
    p.fn = SetEvent;
    p.args[0] = qhasdataevent;
    p.args[1] = 0;
    p.args[2] = 0;
    p.n = 3;
    runtime·syscall(&p);
}

#pragma textflag 7
void
crutches·nosplit_enqueue(int32 msg) {
    crutches·wildcall(EnterCriticalSection, 1, qlock);
    if(qlen < qmaxsize) {
        int32 qtail = (qhead+qlen) % qmaxsize;
        qdata[qtail] = msg;
        qlen++;
    }
    crutches·wildcall(SetEvent, 1, qhasdataevent);
    crutches·wildcall(LeaveCriticalSection, 1, qlock);
    // IMPLEM. NOTE: take care not to "empty" the list by overflowing
    // FIXME: handle full queue
}

#pragma textflag 7
int32
crutches·nosplit_dequeue(void) {
    crutches·wildcall(EnterCriticalSection, 1, qlock);
    for(;;) {
        if(qlen > 0) {
            int32 msg = qdata[qhead];
            qlen--;
            qhead++;
            qhead = qhead % qmaxsize;
            crutches·wildcall(LeaveCriticalSection, 1, qlock);
            return msg;
        } else {
            crutches·wildcall(LeaveCriticalSection, 1, qlock);
            crutches·wildcall(WaitForSingleObject, 2, qhasdataevent, 0xffffffff);
            crutches·wildcall(EnterCriticalSection, 1, qlock);
        }
    }
}


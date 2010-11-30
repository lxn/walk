// by Mateusz Czapliński
// <czapkofan@gmail.com>

#include "crutch.h"

uintptr crutches·stdcall_return(uintptr retval, void* addr_of_caller_first_arg, uint32 no_of_caller_args);

#pragma textflag 7
int32
crutches·foobar(int32 arg) {
    //crutches·wildcall(crutches·nosplit_enqueue, 1, arg+100);
    crutches·nosplit_enqueue(arg+100);
    return (int32)crutches·stdcall_return(arg+100, &arg, 1);
}

void
crutches·Callme(uintptr procaddr, uintptr ms, uintptr times, uintptr r1) {
	StdcallParams p;
	p.fn = (void*)procaddr;
	p.args[0] = ms;
	p.args[1] = times;
	p.args[2] = (uintptr)crutches·foobar; // 0;
	p.args[3] = 0;
	p.args[4] = 0;
	p.args[5] = 0;
	p.n = 6;
	runtime·syscall(&p);
	r1 = p.r;
	//r2 = 0;
	//err = p.err;
	FLUSH(&r1);
	//FLUSH(&r2);
	//FLUSH(&err);
}

void
crutches·WaitForMessage(uintptr r1) {
	StdcallParams p;
	p.fn = (void*)crutches·nosplit_dequeue;
	p.args[0] = 0;
	p.args[1] = 0;
	p.args[2] = 0;
	p.n = 3;
	runtime·syscall(&p);
	r1 = p.r;
	//r2 = 0;
	//err = p.err;
	FLUSH(&r1);
	//FLUSH(&r2);
	//FLUSH(&err);
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

void *RegisterWindowMessageW;

void
·initcrutch(void) {
    // TODO: make sure that 'runtime·osinit()' was already called,
    // or use their approach with 'get_proc_addr2()'
    RegisterWindowMessageW = runtime·get_proc_addr("user32.dll", "RegisterWindowMessageW");
    
}

void
·getRegisteredMessage3(uint32 key, uint32 r1)
{
    if( key >= NUM_WM_KEYS ) {
        r1 = 0;
        FLUSH(&r1);
        return;
    }
        
    if( 0 == msgIds[key] ) {
        StdcallParams p;
        p.fn = RegisterWindowMessageW;
        p.args[0] = (uintptr)msgClsids[key];
        p.args[1] = 0;
        p.args[2] = 0;
        p.n = 3;
        runtime·syscall(&p);
        msgIds[key] = p.r;
    }

    r1 = msgIds[key];
    FLUSH(&r1);
    return;
}

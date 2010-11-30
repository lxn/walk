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

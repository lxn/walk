// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// 'wildcall()' - call into non-8c stdcall code, trying to restore some
// invariants on return ('cld' seems to be enough as of now?)
// based on stdcall_raw()

// notes:
// cld = si, di ++
// movs = ds:si -> es:di
// sal = "shift arithmetic left"

// NOTE: we advertise as a stdcall callback function, so we must
//  behave: ax,cx,dx we can play with, but the others we must
//  leave unchanged. See: http://en.wikipedia.org/wiki/X86_calling_conventions
// TODO: what's the $0/$4 by the fun.? which one do we need?
// void* crutches·wildcall_raw(void* fn, int32 count, uintptr* args)
TEXT crutches·wildcall_raw(SB),7,$0
	PUSHL	BP
	PUSHL	SI
	PUSHL	DI
	PUSHL	BX
	// Copy arguments from stack.
	MOVL	fn+0(FP), AX
	MOVL	count+4(FP), CX		// words
	MOVL	args+8(FP), BP

	// Extract 'args' contents
	SUBL	$(10*4), SP		// padding [?]
	MOVL	CX, BX
	SALL	$2, BX
	SUBL	BX, SP			// room for args
	MOVL	SP, DI
	MOVL	BP, SI
	CLD
	REP; MOVSL

	// Call stdcall function.
	CALL	AX
	
	ADDL	$(10*4), SP		// restore SP?
	POPL	BX
	POPL	DI
	POPL	SI
	POPL	BP
	CLD						// do 8c good?
	RET

// Ugly, ugly hack. We'd better do it like in:
// http://codereview.appspot.com/1696051/diff/21001/src/pkg/runtime/windows/386/sys.s
// http://codereview.appspot.com/1696051/diff/21001/src/pkg/runtime/windows/thread.c
// But now, what is here might work sooner for me.
//
// We must copy the "return to" address (of the Function 
// which called us) in place of the "first arg" (of the Function),
// then reduce stack to there and RET.
// uintptr crutches·stdcall_return(uintptr retval, void* addr_of_caller_first_arg, uint32 no_of_caller_args);
TEXT crutches·stdcall_return(SB),7,$0
	MOVL	retval+0(FP), AX
	MOVL	first+4(FP), DX

	MOVL	num+8(FP), CX
	SUBL	$1, CX
	SALL	$2, CX
	ADDL	DX, CX		// CX := &last_arg

	SUBL	$4, DX		// points to the "return to" address
	MOVL	DX, SP
	MOVL	0(SP), DX	// get the "return to" address
	MOVL	CX, SP
	MOVL	DX, 0(SP)	// dump the "return to" address

	RET					// yeehaw! win or die!


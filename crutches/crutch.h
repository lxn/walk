// by Mateusz Czapliński
// <czapkofan@gmail.com>

#include "runtime.h"
#include "os.h"

void* crutches·wildcall(void* fn, int32 count, ...);

typedef uintptr HANDLE;

typedef struct Message Message;
struct Message {
	HANDLE hwnd;
	uint32 msg;
	uintptr wParam;
	uintptr lParam;
};

void crutches·nosplit_enqueue(Message* msg);
int32 crutches·nosplit_dequeue(Message* msg);

//void crutches·nosplit_enqueue(int32 msg);
//int32 crutches·nosplit_dequeue(void);

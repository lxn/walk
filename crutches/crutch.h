// by Mateusz Czapliński
// <czapkofan@gmail.com>

#include "runtime.h"
#include "os.h"

void crutches·nosplit_enqueue(int32 msg);
int32 crutches·nosplit_dequeue(void);

void* crutches·wildcall(void* fn, int32 count, ...);

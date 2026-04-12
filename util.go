package sqlite

/*
#include <stdlib.h>
#include "./cUtil.h"
*/
import "C"
import "unsafe"

func cUcharStrToCharStr(ucharStr *C.uchar) string {
	cStr := C.ucharStrToCharStr((*C.uchar)(unsafe.Pointer(ucharStr)))

	str := C.GoString(cStr)

	C.free(unsafe.Pointer(cStr))
	return str
}

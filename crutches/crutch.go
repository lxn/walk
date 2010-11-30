// by Mateusz Czapliński
// <czapkofan@gmail.com>

package crutches

func Callme(procaddr, ms, times uintptr) uintptr
func WaitForMessage() uintptr
func initqueue()
func initcrutch()

func init() {
	initqueue()
}

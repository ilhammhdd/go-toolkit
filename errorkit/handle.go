package errorkit

import (
	"fmt"
	"runtime"
)

func ErrorHandled(err error) (ok bool) {
	if err != nil {
		stack := make([]byte, 1024*8)
		stack = stack[:runtime.Stack(stack, false)]
		fmt.Printf("ERROR: %s\n%s\n", err.Error(), stack)
		return true
	}

	return false
}

func PanicThenHandle(err error) {
	defer func() {
		stack := make([]byte, 1024*8)
		stack = stack[:runtime.Stack(stack, false)]
		fmt.Printf("INVOKE PANIC: %s\n%s\n", recover(), stack)
	}()

	if err != nil {
		panic(err)
	}
}

func HandlePanic() {
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 1024*8)
			stack = stack[:runtime.Stack(stack, false)]
			fmt.Printf("PANIC: %s\n%s\n", r, stack)
		}
	}()
}

package errorkit

import (
	"log"
	"runtime"
	"time"
)

func ErrorHandled(err error, stackSize uint32) (ok bool) {
	if err != nil {
		stack := make([]byte, stackSize)
		stack = stack[:runtime.Stack(stack, false)]
		log.Printf("\n%s: %s\nstack: %s", time.Now().UTC().Format(time.RFC3339Nano), err.Error(), stack)
		handlePanic()
		return true
	}

	return false
}

func PanicThenHandle(err error, stackSize uint32) {
	defer func() {
		stack := make([]byte, stackSize)
		stack = stack[:runtime.Stack(stack, false)]
		log.Printf("\n%s: INVOKE PANIC: %s\nstack: %s", time.Now().UTC().Format(time.RFC3339Nano), recover(), stack)
	}()

	if err != nil {
		panic(err)
	}
}

func handlePanic() {
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 1024*8)
			stack = stack[:runtime.Stack(stack, false)]
			log.Printf("\n%s: PANIC: %s\nstack: %s", time.Now().UTC().Format(time.RFC3339Nano), r, stack)
		}
	}()
}

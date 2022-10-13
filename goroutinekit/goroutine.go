package safekit

import (
	"log"
	"runtime"
)

func Do(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]
				log.Printf("PANIC: %s\n%s\n", r, stack)
			}
		}()
		fn()
	}()
}

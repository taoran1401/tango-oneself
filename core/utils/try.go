package utils

import (
	"fmt"
	"runtime"
)

func try(entry func()) {
	defer func() {
		err := recover()
		switch err.(type) {
		case runtime.Error:
			fmt.Println(err)
			break
		default:
			fmt.Println(err)
			break
		}
	}()
	entry()
}

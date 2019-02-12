package mh

import (
	"fmt"
	"runtime"
)

func stderr(err error, levelsabove int) {
	if err == nil {
		return
	}
	_, ffile, fline, fok := runtime.Caller(levelsabove)
	if fok {
		fmt.Println("ERR:", ffile, "@", fline, err.Error())
	}
}

func stdout(v ...interface{}) {
	fmt.Println(v...)
}

package util

import (
	"fmt"
	"log"
	"os"
)

func FatalIfErr(e error, msg string, arg interface{}) {
	if e != nil {
		log.Fatalf("Error: %v Msg: %v Arg: %v", e, msg, arg)
	}
}

func PrintIfErr(e error, msg string, arg interface{}) {
	if e != nil {
		log.Printf("Error: %v Msg: %v Arg: %v", e, msg, arg)
	}
}

func Exit(code int, format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(code)
}

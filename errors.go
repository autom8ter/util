package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

type ErrorCfg struct {
	Message string                 `json:"message"`
	Err     string                 `json:"error"`
	Config  map[string]interface{} `json:"config"`
}

func NewErrCfg(msg string, e error) *ErrorCfg {

	err := &ErrorCfg{
		Message: msg,

		Config: viper.AllSettings(),
	}
	if e == nil {
		err.Err = ""
	} else {
		err.Err = e.Error()
	}
	return err
}

func NewErr(msg string) error {
	return errors.New(msg)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (e *ErrorCfg) String() string {
	output, _ := json.MarshalIndent(e, "", "  ")
	return fmt.Sprintf("%s", output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (e *ErrorCfg) Error() string {
	output, _ := json.MarshalIndent(e, "", "  ")
	return fmt.Sprintf("%s", output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (e *ErrorCfg) FailIfErr() {
	if e.Err != "" {
		logrus.Fatal(e.String())
	}
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (e *ErrorCfg) WarnIfErr() {
	if e.Err != "" {
		logrus.Warn(e.String())
	}
}

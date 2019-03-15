package util

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"reflect"
)

type Struct struct{}

func NewStruct() *Struct {
	return &Struct{}
}

func newStruct(obj interface{}) *structs.Struct {
	return structs.New(obj)
}

func (s *Struct) StructValues(obj interface{}) []interface{} {
	return newStruct(obj).Values()
}

func (s *Struct) StructMap(obj interface{}) map[string]interface{} {
	return newStruct(obj).Map()
}

func (s *Struct) StructFields(obj interface{}) []*structs.Field {
	return newStruct(obj).Fields()
}

func (s *Struct) StructField(obj interface{}, name string) *structs.Field {
	return newStruct(obj).Field(name)
}

func (s *Struct) StructFieldOK(obj interface{}, name string) (*structs.Field, bool) {
	return newStruct(obj).FieldOk(name)
}

func (s *Struct) StructName(obj interface{}, name string) string {
	return newStruct(obj).Name()
}

func (s *Struct) StructHasEmptyFields(obj interface{}, name string) bool {
	return newStruct(obj).HasZero()
}

func (s *Struct) StructUninitialized(obj interface{}, name string) bool {
	return newStruct(obj).IsZero()
}

func (s *Struct) Flagify(obj interface{}) {
	fields := s.StructFields(obj)
	for _, f := range fields {
		if f.IsZero() {
			switch f.Kind() {
			case reflect.String:
				err := f.Set(pflag.String(f.Name(), "", fmt.Sprintf("name: %s kind: %s", f.Name(), f.Kind())))
				logrus.Warnln("failed to reflect flags to struct", err.Error())
			case reflect.Bool:
				err := f.Set(pflag.Bool(f.Name(), false, fmt.Sprintf("name: %s kind: %s", f.Name(), f.Kind())))
				logrus.Warnln("failed to reflect flags to struct", err.Error())
			case reflect.Int:
				err := f.Set(pflag.Int(f.Name(), 0, fmt.Sprintf("name: %s kind: %s", f.Name(), f.Kind())))
				logrus.Warnln("failed to reflect flags to struct", err.Error())
			case reflect.Slice:
				err := f.Set(pflag.StringSlice(f.Name(), []string{}, fmt.Sprintf("name: %s kind: %s", f.Name(), f.Kind())))
				logrus.Warnln("failed to reflect flags to struct", err.Error())
			case reflect.Map:
				err := f.Set(pflag.StringToString(f.Name(), make(map[string]string), fmt.Sprintf("name: %s kind: %s", f.Name(), f.Kind())))
				logrus.Warnln("failed to reflect flags to struct", err.Error())
			}
		}
	}
}

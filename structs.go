package util

import "github.com/fatih/structs"

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

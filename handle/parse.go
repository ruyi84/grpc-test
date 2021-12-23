package handle

import (
	jsoniter "github.com/json-iterator/go"
	"grpc-test/parse"
	"os"
)

type Value struct {
	Name      string
	ValueType parse.FieldType
	Default   interface{}
}

type Test struct {
	MethodName string
	Value      []Value
}

func ParseParameter(parseFile string) []Test {
	paramBytes, err := os.ReadFile(parseFile)
	if err != nil {
		panic(err)
	}

	var paramList []Test

	err = jsoniter.Unmarshal(paramBytes, &paramList)
	if err != nil {
		panic(err)
	}

	return paramList
}

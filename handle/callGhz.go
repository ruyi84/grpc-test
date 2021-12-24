package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bojand/ghz/runner"
	jsoniter "github.com/json-iterator/go"
	"grpc-test/module"
	"io/ioutil"
)

type StdinJson struct {
	MethodName string
	Value      map[string]interface{}
}

// callGhz
func CallGHZ(protoFile, paramFile, callName, host string) {
	testList := ParseParameter(paramFile)

	for _, test := range testList {
		input := StdinJson{
			MethodName: test.MethodName,
		}

		input.Value = make(map[string]interface{})
		for _, value := range test.Value {
			input.Value[value.Name] = value.Default
		}

		ghz(protoFile, callName, test.MethodName, host, input)
	}
}

func ghz(protoFile, callName, methodName, host string, jsonInfo StdinJson) {

	reqInfo, err := jsoniter.MarshalToString(jsonInfo.Value)
	if err != nil {
		panic(err)
	}

	callName = fmt.Sprintf("%s.%s", callName, methodName)
	if module.Token != "" {
		runner.WithMetadataFromJSON(fmt.Sprintf(`{"authorization":"Bearer %s"}`, module.Token))
	} else {

	}
	report, err := runner.Run(
		callName,
		host,
		runner.WithProtoFile(protoFile, []string{}),
		runner.WithDataFromJSON(reqInfo),
		runner.WithInsecure(true),
	)

	if err != nil {
		panic(err)
	}

	marshalJSON, err := report.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var result bytes.Buffer
	err = json.Indent(&result, marshalJSON, "", " ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(module.OutPath, result.Bytes(), 0666)
	if err != nil {
		panic(err)
	}
}

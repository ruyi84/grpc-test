package handle

import (
	"fmt"
	"github.com/bojand/ghz/runner"
	jsoniter "github.com/json-iterator/go"
)

type StdinJson struct {
	MethodNafme string
	Value       map[string]interface{}
}

// callGhz
func CallGHZ(protoFile, paramFile, callName, host string) {
	testList := ParseParameter(paramFile)

	for _, test := range testList {
		input := StdinJson{
			MethodNafme: test.MethodName,
		}
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
	report, err := runner.Run(
		callName,
		host,
		runner.WithProtoFile(protoFile, []string{}),
		runner.WithDataFromJSON(reqInfo),
		runner.WithInsecure(true),
		//runner.WithMetadataFromJSON(`{"authorization":"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzg0NTA0OTUsInByaXYiOjEsInJvbGUiOiJtZ3IiLCJ1c2VyIjoibGluIn0.mPICgvgBxFnLsQHJlULWu31ckBkuIY559PHEJobsqYU"}`),
	)

	if err != nil {
		panic(err)
	}

	json, err := report.MarshalJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json))

}

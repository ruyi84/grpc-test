package main

import (
	"flag"
	"fmt"
	"github.com/bojand/ghz/runner"
	"github.com/fullstorydev/grpcurl"
	"github.com/json-iterator/go"
	"grpc-test/parse"
	"io/ioutil"
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

type StdinJson struct {
	MethodNafme string
	Value       map[string]interface{}
}

var operation string // 操作类型：输出解析文件，交互式，参数命令式
var intPath string
var outPath string

func init() {
	flag.StringVar(&operation, "operation", "", "Input your type of operation.")
	flag.StringVar(&intPath, "intPath", "", "Input your type of input path.")
	flag.StringVar(&outPath, "outPath", "", "Input your type of output path.")

}

func main() {

	flag.Parse()

	switch operation {
	case "parseProto":
	case "ghz":
	}

	// 解析proto，获取对应的接口列表和参数列表
	parseJson(intPath, outPath)

}

func ghz(protoFile, intFile, callName, host string) {

	readFile, err := ioutil.ReadFile(intFile)
	if err != nil {
		panic(err)
	}

	var stdinList []StdinJson
	err = jsoniter.Unmarshal(readFile, &stdinList)
	if err != nil {
		panic(err)
	}

	for _, jsonInfo := range stdinList {
		reqInfo, err := jsoniter.MarshalToString(jsonInfo.Value)
		if err != nil {
			panic(err)
		}

		callName := fmt.Sprintf("%s.%s", callName, jsonInfo.MethodNafme)
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

}

func parseJson(filePath, outPath string) {
	descSource, err := grpcurl.DescriptorSourceFromProtoFiles([]string{}, filePath)
	if err != nil {
		panic(err)
	}

	methods, err := parse.GetMethods(descSource, nil)
	if err != nil {
		panic(err)
	}

	var list []Test
	var stdinList []StdinJson
	for _, method := range methods {
		r, err := parse.GatherMetadataForMethod(method)
		if err != nil {
			panic(err)
		}
		var valueList []Value
		valueMap := make(map[string]interface{})
		for _, v := range r.MessageTypes {
			for _, def := range v {
				valueList = append(valueList, Value{
					Name:      def.Name,
					ValueType: def.Type,
					Default:   def.DefaultVal,
				})
				valueMap[def.Name] = def.DefaultVal
			}

		}
		list = append(list, Test{
			MethodName: method.GetName(),
			Value:      valueList,
		})

		stdinList = append(stdinList, StdinJson{
			MethodNafme: method.GetName(),
			Value:       valueMap,
		})
	}

	indent, err := jsoniter.MarshalIndent(stdinList, "", "	")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(outPath, indent, 0644)
	if err != nil {
		panic(err)
	}

	return
}

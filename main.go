package main

import (
	"flag"
	"github.com/fullstorydev/grpcurl"
	"github.com/json-iterator/go"
	"grpc-test/handle"
	"grpc-test/parse"
	"os"
)

var operation string // 操作类型：输出解析文件，交互式，参数命令式
var intPath string
var outPath string
var protoFile string
var serviceName string
var host string

func init() {
	flag.StringVar(&operation, "operation", "", "Input your type of operation.")
	flag.StringVar(&serviceName, "serviceName", "", "Input your serviceName.")
	flag.StringVar(&host, "host", "", "Input your host.")
	flag.StringVar(&protoFile, "protoFile", "", "Input your type of input path.")
	flag.StringVar(&intPath, "intPath", "", "Input your protoFile.")
	flag.StringVar(&outPath, "outPath", "", "Input your type of output path.")

}

func main() {

	flag.Parse()

	switch operation {
	case "parseProto":
		// 解析proto，获取对应的接口列表和参数列表
		parseJson(intPath, outPath)
	case "ghz": // 倒入protoFile和对应的parseFile，然后遍历执行所有的请求，并且输出对应的测试结果
		handle.CallGHZ(protoFile, intPath, serviceName, host)
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

	var list []handle.Test
	var stdinList []handle.StdinJson
	for _, method := range methods {
		r, err := parse.GatherMetadataForMethod(method)
		if err != nil {
			panic(err)
		}
		var valueList []handle.Value
		valueMap := make(map[string]interface{})
		for _, v := range r.MessageTypes {
			for _, def := range v {
				valueList = append(valueList, handle.Value{
					Name:      def.Name,
					ValueType: def.Type,
					Default:   def.DefaultVal,
				})
				valueMap[def.Name] = def.DefaultVal
			}

		}
		list = append(list, handle.Test{
			MethodName: method.GetName(),
			Value:      valueList,
		})

		stdinList = append(stdinList, handle.StdinJson{
			MethodNafme: method.GetName(),
			Value:       valueMap,
		})
	}

	indent, err := jsoniter.MarshalIndent(stdinList, "", "    ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(outPath, indent, 0644)
	if err != nil {
		panic(err)
	}

	return
}

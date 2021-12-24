package handle

import (
	"github.com/fullstorydev/grpcurl"
	jsoniter "github.com/json-iterator/go"
	"grpc-test/parse"
	"io/ioutil"
	"strings"
)

func ParseJson(filePath, outPath string) {
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
		for k, v := range r.MessageTypes {

			// 用户去除嵌套参数
			split := strings.Split(k, ".")
			if len(split) > 2 {
				continue
			}
			// 仅保留请求参数本身，忽略嵌套结构
			if !strings.Contains(k, "Req") {
				continue
			}

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
			MethodName: method.GetName(),
			Value:      valueMap,
		})
	}

	indent, err := jsoniter.MarshalIndent(list, "", "    ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(outPath, indent, 0644)
	if err != nil {
		panic(err)
	}

	return
}

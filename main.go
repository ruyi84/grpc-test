package main

import (
	"grpc-test/handle"
	"grpc-test/module"
)

func main() {

	module.InitParm()

	switch module.Operation {
	case "parseProto":
		// 解析proto，获取对应的接口列表和参数列表
		handle.ParseJson(module.IntPath, module.OutPath)
	case "ghz": // 倒入protoFile和对应的parseFile，然后遍历执行所有的请求，并且输出对应的测试结果
		handle.CallGHZ(module.ProtoFile, module.IntPath, module.ServiceName, module.Host)
	}
}

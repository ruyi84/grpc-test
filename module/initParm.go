package module

import "flag"

var Operation string // 操作类型：输出解析文件，交互式，参数命令式
var IntPath string
var OutPath string
var ProtoFile string
var ServiceName string
var Host string
var Token string

func InitParm() {
	flag.StringVar(&Operation, "operation", "", "Input your type of Operation.")
	flag.StringVar(&ServiceName, "serviceName", "", "Input your ServiceName.")
	flag.StringVar(&Host, "host", "", "Input your Host.")
	flag.StringVar(&ProtoFile, "protoFile", "", "Input your type of input path.")
	flag.StringVar(&IntPath, "intPath", "", "Input your ProtoFile.")
	flag.StringVar(&OutPath, "outPath", "", "Input your type of output path.")
	flag.StringVar(&Token, "token", "", "Input your Token.")

	flag.Parse()
}

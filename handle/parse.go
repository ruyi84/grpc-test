package handle

import (
	"bytes"
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
	"grpc-test/parse"
	"io/ioutil"
)

//var str = `[{"MethodName":"Login","Value":{"username":"","password":"","totpCode":""}},{"MethodName":"Logout","Value":{"username":""}},{"MethodName":"ListUser","Value":{"sort":0,"pageIdx":0,"search":"","isSelf":false,"filterPassStrength":null,"filterStartPassTm":"","filterEndPassTm":"","pageSize":0,"filterRole":null,"filterMfaStatus":null,"filterStartCreateTm":"","filterEndCreateTm":""}},{"MethodName":"AddUser","Value":{"username":"","role":"","realName":"","department":"","password":"","mobile":"","email":"","remark":"","address":"","post":""}},{"MethodName":"UpdateUser","Value":{"mobile":"","email":"","remark":"","department":"","post":"","username":"","role":"","address":"","realName":""}},{"MethodName":"UpdateUserPassword","Value":{"oldPassword":"","newPassword":"","username":""}},{"MethodName":"DeleteUser","Value":{"username":""}},{"MethodName":"UserExists","Value":{"username":""}},{"MethodName":"CheckMfa","Value":{"password":"","username":""}},{"MethodName":"EnableMfa","Value":{"username":"","password":"","secret":"","mfaCode":""}},{"MethodName":"DisableMfa","Value":{"username":""}},{"MethodName":"UpdateAvatar","Value":{"userId":0,"file":""}},{"MethodName":"ResetPassword","Value":{"username":"","newPassword":""}},{"MethodName":"GetPwdUpdateTm","Value":{"userName":""}},{"MethodName":"ListAgent","Value":{"pageIdx":0,"IP":"","status":null,"domain":null,"tmSort":0,"pageSize":0,"dcHostname":"","packetStatus":null,"keyword":""}},{"MethodName":"UpdateAgent","Value":{"key":"","value":"","ID":"","bindNetIface":null,"perfLimit":{},"remark":""}},{"MethodName":"CmdAgent","Value":{"ID":"","cmd":""}},{"MethodName":"DownloadAgent","Value":{"Type":""}},{"MethodName":"DownCertificate","Value":{"domain":""}},{"MethodName":"DeleteAgent","Value":{"agentId":""}},{"MethodName":"ListGateway","Value":{"filterStatus":null,"keyWord":"","pageIdx":0,"pageSize":0}},{"MethodName":"ListWecBeat","Value":{"dcHostname":"","status":null,"tmSort":0,"pageIdx":0,"pageSize":0}},{"MethodName":"UpdateAgentVersion","Value":{"agentId":"","version":""}},{"MethodName":"DeleteWecBeat","Value":{"beatId":""}},{"MethodName":"ListDomainNameFromAgent","Value":{}},{"MethodName":"GetDCList","Value":{}},{"MethodName":"AddWecConf","Value":{"fqdn":"","ip":""}},{"MethodName":"TestWecConf","Value":{"fqdn":"","ip":""}},{"MethodName":"WecBeatInfo","Value":{"id":""}},{"MethodName":"ListWecBeatEventInfo","Value":{"wecBeatID":"","orderEventID":0,"orderCount":0,"lastEventTmSort":0,"pageIdx":0,"pageSize":0}},{"MethodName":"ListDomain","Value":{"pageIdx":0,"pageSize":0,"filterDomain":"","filterStatus":"","filterKeyword":""}},{"MethodName":"AddDomain","Value":{"ldapAddr":"","username":"","password":"","DNS":""}},{"MethodName":"TestDomain","Value":{"password":"","DNS":"","ldapAddr":"","username":""}},{"MethodName":"UpdateDomain","Value":{"ID":"","ldapAddr":"","username":"","password":"","DNS":""}},{"MethodName":"DeleteDomain","Value":{"ID":"","Name":""}},{"MethodName":"GetDomainObject","Value":{"domain":"","name":""}},{"MethodName":"GetDomainObjectInfo","Value":{"name":"","domain":""}},{"MethodName":"UpdateDomainData","Value":{"domainID":""}},{"MethodName":"GetDomainInfo","Value":{"domainID":"","PageSize":0,"PageIdx":0,"Sort":0,"Keyword":""}},{"MethodName":"SetMsRCP","Value":{"domainID":"","dcHostName":"","isPullLog":false}},{"MethodName":"ListDomainNameForEventList","Value":{}},{"MethodName":"ListDomainName","Value":{}},{"MethodName":"AddDomainEntry","Value":{"domain":"","entryType":"","entryCnt":""}},{"MethodName":"DeleteDomainEntry","Value":{"IDList":null}},{"MethodName":"ListDomainEntry","Value":{"entryType":"","timeSort":0,"pageIdx":0,"pageSize":0,"domain":null,"startTm":"","endTm":"","search":""}},{"MethodName":"UpdateKerberosConf","Value":{"showType":"","filedSuffix":"","ID":"","list":null,"cntName":"","desc":"","valueList":null}},{"MethodName":"ListKerberosConf","Value":{"pageIdx":0,"pageSize":0,"domain":null,"ruleName":null,"ruleType":null,"sortTm":0}},{"MethodName":"GetKerberosConf","Value":{"id":""}},{"MethodName":"AddRuleWhitelist","Value":{"ruleID":"","remark":"","domain":"","matchType":0,"ruleCnt":null,"eventID":"","name":"","match":"","value":""}},{"MethodName":"DeleteRuleWhitelist","Value":{"IDList":null}},{"MethodName":"UpdateRuleWhitelist","Value":{"name":"","match":"","value":"","ID":"","matchType":0,"ruleCnt":null,"remark":""}},{"MethodName":"ListRuleWhitelist","Value":{"startTm":"","endTm":"","sortTm":0,"pageIdx":0,"pageSize":0,"domain":null,"ruleName":null,"ruleType":null}},{"MethodName":"GetRuleWhitelistInfo","Value":{"ID":""}},{"MethodName":"ListWhiteField","Value":{"ruleID":""}},{"MethodName":"GetWhiteFieldValue","Value":{"domain":"","filedName":"","search":""}},{"MethodName":"GetSystemInfo","Value":{}},{"MethodName":"DownloadSystemLog","Value":{"systemType":null}},{"MethodName":"GetSystemLog","Value":{"pageIdx":0,"pageSize":0}},{"MethodName":"UpdateReboot","Value":{"type":""}},{"MethodName":"GetLicence","Value":{}},{"MethodName":"UpdateLicence","Value":{"code":""}},{"MethodName":"SelectOptionAuditLogEvent","Value":{}},{"MethodName":"ListAuditLog","Value":{"pageIdx":0,"pageSize":0,"startTm":"","endTm":"","keyword":"","filterEvent":null,"createSort":0}},{"MethodName":"DeleteAuditLog","Value":{"id":null}},{"MethodName":"ExportAuditLog","Value":{"startTm":"","endTm":""}},{"MethodName":"GetTaskState","Value":{"taskUUID":""}},{"MethodName":"UpdateSystemIcon","Value":{"file":""}},{"MethodName":"GetSystemIcon","Value":{}},{"MethodName":"NetworkDiag","Value":{"type":"DiagDefault","target":""}},{"MethodName":"ListSystemStatus","Value":{"type":""}},{"MethodName":"SetSystemTime","Value":{"dateOrIp":""}},{"MethodName":"GetSystemPing","Value":{}},{"MethodName":"SetStorage","Value":{"storageDataAlarmValue":0,"storageLogValue":0,"storageAutoClearValue":0}},{"MethodName":"GetStorage","Value":{}},{"MethodName":"ListExportTask","Value":{"pageSize":0,"modelType":null,"status":null,"startTm":"","endTm":"","sortTm":0,"pageIdx":0}},{"MethodName":"ListThreatEvent","Value":{"startTm":"","advancedSearch":null,"name":"","value":null,"searchType":0,"sortTm":0,"type":"","pageIdx":0,"pageSize":0,"threatID":null,"threatLevel":null,"endTm":""}},{"MethodName":"UpdateThreatEvent","Value":{"status":0,"ID":""}},{"MethodName":"ExportThreatEvent","Value":{"startTm":"","endTm":""}},{"MethodName":"ListThreatActivity","Value":{"pageIdx":0,"pageSize":0,"eventID":"","orderCreateTm":0,"uniqueID":""}},{"MethodName":"ListThreatRawLog","Value":{"pageIdx":0,"pageSize":0,"ID":"","eventID":""}},{"MethodName":"GetRuleInfo","Value":{"ruleId":""}},{"MethodName":"GetDCNameList","Value":{}},{"MethodName":"GetTarget","Value":{}},{"MethodName":"GetDomainFromAlert","Value":{}},{"MethodName":"ListThreatEventSearch","Value":{"search":"","searchType":""}},{"MethodName":"ListRuleTypes","Value":{"ruleID":""}},{"MethodName":"StateAlertEventByRule","Value":{}},{"MethodName":"GetThreatEventByUniqueID","Value":{"uniqueID":""}},{"MethodName":"StatsAlertActivity","Value":{"domain":""}},{"MethodName":"StatsRiskAssets","Value":{"domain":""}},{"MethodName":"StatsAlertEvents","Value":{"domain":""}},{"MethodName":"StatsScanEvents","Value":{"domain":""}},{"MethodName":"StatsAssets","Value":{"domain":""}},{"MethodName":"AlarmAnalysis","Value":{"interval":"","startTm":"","endTm":"","template":"","domain":""}},{"MethodName":"RiskTrend","Value":{"interval":"","startTm":"","endTm":"","template":"","domain":""}},{"MethodName":"ListStatsAlertName","Value":{"domain":""}},{"MethodName":"ListStatsAlertType","Value":{"domain":""}},{"MethodName":"ScanInspection","Value":{"type":"","domains":null}},{"MethodName":"StopScan","Value":{"taskUUID":""}},{"MethodName":"ExportScanEvent","Value":{"domainName":null,"startTm":"","endTm":""}},{"MethodName":"GetScanTaskState","Value":{"label":"","taskUUID":""}},{"MethodName":"GetScanScore","Value":{"domain":""}},{"MethodName":"GetScanRule","Value":{"label":""}},{"MethodName":"LastScanInfo","Value":{}},{"MethodName":"ListOnlineDomain","Value":{}},{"MethodName":"ListDomainByScanEvent","Value":{}},{"MethodName":"ScanLeakEvent","Value":{"scanType":"","plugins":null,"domains":null,"fqdn":""}},{"MethodName":"GetScanLeakEventStatus","Value":{"taskID":""}},{"MethodName":"StatsLeakEvent","Value":{}},{"MethodName":"ListScanPluginType","Value":{}},{"MethodName":"ExportLeakEvent","Value":{"domainName":null,"startTm":"","endTm":""}},{"MethodName":"ProtectInfo","Value":{"tmgType":"","domain":""}},{"MethodName":"IsScanning","Value":{"tmgType":""}},{"MethodName":"BaselineScanStatus","Value":{}},{"MethodName":"LeakScanStatus","Value":{}},{"MethodName":"LeakScanLoopholeStatistics","Value":{"domain":""}},{"MethodName":"BaselineNewResult","Value":{"domain":""}},{"MethodName":"OverviewLeakEventList","Value":{"orderUpdateTm":0,"domain":""}},{"MethodName":"EventList","Value":{"points":0,"lastScanTime":0,"pageIdx":0,"pageSize":0,"categorys":null,"domains":null,"levels":null,"eventName":""}},{"MethodName":"EventDetails","Value":{"ruleId":""}},{"MethodName":"GetInstanceList","Value":{"pageIdx":0,"pageSize":0,"eventID":""}},{"MethodName":"ListLeakEvent","Value":{"pageSize":0,"startCreateTm":"","startUpdateTm":"","endUpdateTm":"","levels":null,"statuses":null,"search":"","pageIdx":0,"orderCreateTm":0,"orderLastCheckTm":0,"domains":null,"leakTypes":null,"endCreateTm":""}},{"MethodName":"ListTaskManagerGroup","Value":{"pageIdx":0,"pageSize":0,"triggers":null,"tmgTypes":null,"endTmEnd":"","orderCreateTm":0,"statuses":null,"endTmStart":"","orderEndTm":0}},{"MethodName":"DetailTaskManagerGroup","Value":{"orderUpdateTm":0,"orderPoints":0,"id":"","pageIdx":0,"pageSize":0,"orderLastCheckTm":0,"orderCreateTm":0}},{"MethodName":"DeleteTaskManagerGroup","Value":{"ids":null}},{"MethodName":"SetCronTask","Value":{"id":"","enable":false,"cycle":"every_day","Rate":0}},{"MethodName":"ListCronTask","Value":{"pageIdx":0,"pageSize":0}},{"MethodName":"ListScanPlugin","Value":{"pageIdx":0,"pageSize":0,"enable":""}},{"MethodName":"UpdateScanPluginEnable","Value":{"id":0,"enable":0}},{"MethodName":"UpdateScanPluginMetaData","Value":{"id":0,"metaData":{},"key":"","value":""}},{"MethodName":"ListNotify","Value":{"endCreateTm":"","createTmSort":0,"notifyType":"","startCreateTm":"","pageSize":0,"notifySubType":"","keyword":"","notifyStatus":"","pageIdx":0}},{"MethodName":"UpdateNotify","Value":{"IDList":null}},{"MethodName":"StatsNotify","Value":{}},{"MethodName":"ListNotifyConf","Value":{"notifyType":null,"target":"","enable":null,"sortTime":0,"pageIdx":0,"pageSize":0,"moduleType":null}},{"MethodName":"UpdateNotifyConf","Value":{"senderName":"","id":"","sender":"","server":"","port":"","alertInterval":0,"remark":"","target":"","senderIdentity":"","level":null}},{"MethodName":"UpdateNotifyConfEnable","Value":{"id":"","enable":""}},{"MethodName":"GetNotifyConfInfo","Value":{"id":""}},{"MethodName":"ListNotifyTarget","Value":{}},{"MethodName":"TestEmailSend","Value":{"config":{},"key":"","value":""}},{"MethodName":"SelectOptionNotify","Value":{}},{"MethodName":"GenerateEventReport","Value":{"name":"","startTm":"","endTm":""}},{"MethodName":"ListEventReport","Value":{"pageSize":0,"pageIdx":0}},{"MethodName":"StatusEventReport","Value":{"taskID":""}},{"MethodName":"DownloadEventReport","Value":{"ID":"","type":""}},{"MethodName":"DeleteEventReport","Value":{"ReportID":""}},{"MethodName":"ListAssetsUser","Value":{"pageIdx":0,"pageSize":0,"domainName":""}},{"MethodName":"ListAssetsComputer","Value":{"pageIdx":0,"pageSize":0,"domainName":""}},{"MethodName":"ListAssetsGroup","Value":{"pageIdx":0,"pageSize":0,"domainName":""}},{"MethodName":"GetAssetsDetailsByAlert","Value":{"domain":"","type":"","name":""}},{"MethodName":"ListGroupByAssets","Value":{"name":"","assetsType":"","domain":""}},{"MethodName":"ListAssetsUserByGroup","Value":{"pageIdx":0,"pageSize":0,"domain":"","groupName":""}},{"MethodName":"GetAssetsEntry","Value":{"domain":"","type":"","search":""}},{"MethodName":"GetAssetsLabel","Value":{"domain":"","assetsType":"","name":""}},{"MethodName":"GetAssetsActivities","Value":{"domain":"","type":"","search":""}},{"MethodName":"StatsAssetsActivitiesLevel","Value":{"domain":"","search":""}},{"MethodName":"GetAssetsSensitiveGroupLabelInfo","Value":{"domain":"","type":"","name":""}},{"MethodName":"ListUsersSensitiveGroup","Value":{"domain":"","name":""}},{"MethodName":"ListAttackPath","Value":{"startTm":"","name":"","nameType":"","pageIdx":0,"domain":"","ruleIDList":null,"endTm":"","level":"","pageSize":0}},{"MethodName":"ExportAttackPath","Value":{"startTm":"","endTm":"","name":"","nameType":"","ruleList":null,"domain":""}},{"MethodName":"ListSearchLogEvent","Value":{"pageSize":0,"domain":"","logType":"","startTm":"","endTm":"","search":"","pageIdx":0}},{"MethodName":"GetSearchLogField","Value":{"domain":"","logType":""}},{"MethodName":"GetSearchFieldInfo","Value":{"domain":"","fieldType":"","filedName":""}},{"MethodName":"GetSearchChartData","Value":{"domain":"","type":"minutes"}},{"MethodName":"AddSearchTemplate","Value":{"title":"","content":"","logType":""}},{"MethodName":"ListSearchTemplate","Value":{"PageSize":0,"logType":"","PageIdx":0}},{"MethodName":"DeleteSearchTemplate","Value":{"id":""}},{"MethodName":"ListSearchTableColumns","Value":{"logType":""}}]`

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
	paramBytes, err := ioutil.ReadFile(parseFile)
	if err != nil {
		panic(err)
	}

	var paramList []Test

	var bytes bytes.Buffer
	err = json.Indent(&bytes, paramBytes, " ", "")
	if err != nil {
		panic(err)
	}

	err = jsoniter.Unmarshal(bytes.Bytes(), &paramList)
	if err != nil {
		panic(err)
	}

	return paramList
}

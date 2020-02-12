package constant

const (
	CommitLogProperty_AppId  = "appID"
	CommitLogProperty_Status = "statusCode"
	CommitLogProperty_Source = "source"
)

//log source
const (
	CommitLogSource_EdgeWize = "EdgeWize"
)

const (
	//Monitor服务启动
	LogStatus_MonitorStart = 920
	//Monitor服务停止
	LogStatus_MonitorStop = 921
	//Monitor服务运行失败
	LogStatus_MonitorFailed = 922
	//采集数据成功
	LogStatus_ColllectSuccess = 923
	//采集数据失败
	LogStatus_CollectFailed = 924
	//上报数据成功
	LogStatus_PostSuccess = 925
	//上报数据失败
	LogStatus_PostFailed = 926
)

package Errmsg
const (
	SUCCESS 			= 200
	ERROR   			= 500
	Error_QueryInfo 	= 1001 // 接口查询参数出错
	Error_NotFound		= 1002 // 未找到符合条件结果
	Error_InsertInfo	= 1003 // 新增数据失败
	ERROR_InitError		= 2001 // 数据初始化失败
)

var codemsg = map[int]string{
	SUCCESS: 				"OK",
	ERROR:   				"FAIL",
	Error_QueryInfo:		"QueryInfo Error",
	Error_NotFound:			"No Result Found",
	Error_InsertInfo: 		"InsertInfo Error",
	ERROR_InitError:		"InitData Error",
}

func GetErrMsg(code int) string {
	return codemsg[code]
}

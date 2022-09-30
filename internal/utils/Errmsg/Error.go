package Errmsg

const (
	SUCCESS = 200
	ERROR   = 500
	// aliyunSTS 模块报错
	ErrorCreateSTSClient = 101
	ErrorAssumeRole      = 102
	// aliyun ResourceManager 模块报错
	ErrorCreateRMClient       = 201
	ErrorListAliCloudAccounts = 202
	// aliyun RDS模块报错
	ErrorCreateRDSClient        = 301
	ErrorDescribeRDSInstances   = 302
	ErrorAddInventory           = 303
	DescribeDBInstanceAttribute = 304
	ErrorDescribeRDSAccount     = 305
)

var codemsg = map[int]string{
	SUCCESS:                     "OK",
	ERROR:                       "FAIL",
	ErrorCreateSTSClient:        "createSTSClient fail",
	ErrorAssumeRole:             "AssumeRole fail",
	ErrorCreateRMClient:         "ErrorCreateRMClient fail",
	ErrorListAliCloudAccounts:   "ErrorListAliCloudAccounts fail",
	ErrorCreateRDSClient:        "ErrorCreateRDSClient fail",
	ErrorDescribeRDSInstances:   "ErrorDescribeRDSInstances fail",
	ErrorAddInventory:           "ErrorAddInventory fail",
	DescribeDBInstanceAttribute: "DescribeDBInstanceAttribute fail",
	ErrorDescribeRDSAccount:     "ErrorDescribeRDSAccount fail",
}

func GetErrMsg(code int) string {
	return codemsg[code]
}

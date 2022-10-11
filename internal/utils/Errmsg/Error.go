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
	ErrorDescribeDatabases      = 306
	ErrorDescribeInstanceSSL    = 307
	ErrorLockAccount            = 308
	// aliyun UpdateDBDetails 模块报错
	ErrorCheckInventoryExist          = 401
	ErrorCheckAccountExist            = 402
	ErrorCreateRDSAccount             = 403
	ErrorUnlockRDSAccount             = 404
	ErrorUpdatePgsqlDetailsConnection = 405
)

var CodeMsg = map[int]string{
	SUCCESS:                           "OK",
	ERROR:                             "FAIL",
	ErrorCreateSTSClient:              "createSTSClient fail",
	ErrorAssumeRole:                   "AssumeRole fail",
	ErrorCreateRMClient:               "CreateRMClient fail",
	ErrorListAliCloudAccounts:         "ListAliCloudAccounts fail",
	ErrorCreateRDSClient:              "CreateRDSClient fail",
	ErrorDescribeRDSInstances:         "DescribeRDSInstances fail",
	ErrorDescribeInstanceSSL:          "ErrorDescribeInstanceSSL",
	ErrorAddInventory:                 "AddInventory fail",
	DescribeDBInstanceAttribute:       "DescribeDBInstanceAttribute fail",
	ErrorDescribeRDSAccount:           "DescribeRDSAccount fail",
	ErrorDescribeDatabases:            "ErrorDescribeDatabases Fail",
	ErrorCheckInventoryExist:          "Inventory Not Exist",
	ErrorUnlockRDSAccount:             "UnlockAccount Fail",
	ErrorUpdatePgsqlDetailsConnection: "Error Connect the RemoteDB",
	ErrorLockAccount:                  "Error Lock Account",
}

func GetErrMsg(code int) string {
	return CodeMsg[code]
}

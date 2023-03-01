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
	ErrorCreateRDSClient                  = 301
	ErrorDescribeRDSInstances             = 302
	ErrorAddInventory                     = 303
	DescribeDBInstanceAttribute           = 304
	ErrorDescribeRDSAccount               = 305
	ErrorDescribeDatabases                = 306
	ErrorDescribeInstanceSSL              = 307
	ErrorDeleteAccount                    = 308
	ErrorGrantAccountPrivilege            = 309
	ErrorDescribeDBInstanceNetInfo        = 310
	ErrorAllocateInstancePublicConnection = 311
	ErrorModifySecurityIps                = 312
	ErrorReleaseInstancePublicConnection  = 313
	ErrorDescribeDBInstanceIPArrayList    = 314
	// aliyun UpdateDBDetails 模块报错
	ErrorCheckInventoryExist          = 401
	ErrorCheckAccountExist            = 402
	ErrorCreateRDSAccount             = 403
	ErrorUnlockRDSAccount             = 404
	ErrorUpdatePgsqlDetailsConnection = 405
	ErrorUpdateMysqlDetailsConnection = 406
	// Category 模块报错
	// API 相关
	ErrorQueryInput = 2001
	ErrorNotFound   = 2002
	ErrorCreate     = 2003
)

var CodeMsg = map[int]string{
	SUCCESS:                               "OK",
	ERROR:                                 "FAIL",
	ErrorCreateSTSClient:                  "createSTSClient fail",
	ErrorAssumeRole:                       "AssumeRole fail",
	ErrorCreateRMClient:                   "CreateRMClient fail",
	ErrorListAliCloudAccounts:             "ListAliCloudAccounts fail",
	ErrorCreateRDSClient:                  "CreateRDSClient fail",
	ErrorDescribeRDSInstances:             "DescribeRDSInstances fail",
	ErrorDescribeInstanceSSL:              "ErrorDescribeInstanceSSL",
	ErrorAddInventory:                     "AddInventory fail",
	DescribeDBInstanceAttribute:           "DescribeDBInstanceAttribute fail",
	ErrorDescribeRDSAccount:               "DescribeRDSAccount fail",
	ErrorDescribeDatabases:                "ErrorDescribeDatabases Fail",
	ErrorCheckInventoryExist:              "Inventory Not Exist",
	ErrorUnlockRDSAccount:                 "UnlockAccount Fail",
	ErrorUpdatePgsqlDetailsConnection:     "Error Connect the RemoteDB",
	ErrorDeleteAccount:                    "Error Delete Account",
	ErrorUpdateMysqlDetailsConnection:     "Error Connect the RemoteDB",
	ErrorGrantAccountPrivilege:            "Error Grant Account Privilege",
	ErrorDescribeDBInstanceNetInfo:        "Error Describe DBInstance NetInfo",
	ErrorAllocateInstancePublicConnection: "Error Allocate Instance Public Connection",
	ErrorModifySecurityIps:                "Error Modify Security Ips",
	ErrorReleaseInstancePublicConnection:  "Error Release Instance Public Connection",
	ErrorQueryInput:                       "Error Query Input",
	ErrorNotFound:                         "Error Data Not Found",
	ErrorDescribeDBInstanceIPArrayList:    "Error Describe DBInstance IPArrayList",
	ErrorCreate:                           "Error in Creation",
}

func GetErrMsg(code int) string {
	return CodeMsg[code]
}

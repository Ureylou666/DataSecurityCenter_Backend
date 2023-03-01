package aliyunSDK

/*
aliyun RDS 相关sdk接口
- 创建客户端
- 查询RDS实例列表
- 查询RDS实例账号信息
- 创建RDS实例账号
- 查询RDS实例下数据库信息
- 授权RDS实例账号权限
*/

import (
	"Backend/internal/utils"
	"Backend/internal/utils/Errmsg"
	"Backend/internal/utils/setting"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	rds20140815 "github.com/alibabacloud-go/rds-20140815/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

// CreateRDSClient 创建RDS连接客户端
func CreateRDSClient(accountId string) (ErrCode int, ErrMessage error, client *rds20140815.Client) {
	// 获取临时STS token
	previousCode, previousMsg, AccessKeyId, AccessKeySecret, SecurityToken := AssumeRole(accountId)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg, nil
	}
	// 成功获取临时AK SK
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: tea.String(AccessKeyId),
		// 您的 AccessKey Secret
		AccessKeySecret: tea.String(AccessKeySecret),
		// 您的 AccessToken
		SecurityToken: tea.String(SecurityToken),
	}
	// 访问的域名
	config.Endpoint = tea.String("rds.aliyuncs.com")
	client = &rds20140815.Client{}
	client, ErrMessage = rds20140815.NewClient(config)
	// 错误判断
	if ErrMessage != nil {
		return Errmsg.ErrorCreateRDSClient, ErrMessage, nil
	} else {
		return Errmsg.SUCCESS, nil, client
	}
}

// DescribeRDSInstances 查询项目组RDS实例列表
func DescribeRDSInstances(client *rds20140815.Client) (ErrCode int, ErrMessage error, RDSInstanceList []*rds20140815.DescribeDBInstancesResponseBodyItemsDBInstance) {
	describeDBInstancesRequest := &rds20140815.DescribeDBInstancesRequest{
		RegionId:   tea.String("cn-shanghai"), // 目前默认都在cn-shanghai
		PageSize:   tea.Int32(100),
		PageNumber: tea.Int32(1),
	}
	describeDBInstancesResponse, ErrMessage := client.DescribeDBInstances(describeDBInstancesRequest)
	if ErrMessage != nil {
		return Errmsg.ErrorDescribeRDSInstances, ErrMessage, nil
	}
	return Errmsg.SUCCESS, nil, describeDBInstancesResponse.Body.Items.DBInstance
}

// DescribeDBInstanceAttribute 查询RDS实例的详细信息
func DescribeDBInstanceAttribute(InstanceId string, client *rds20140815.Client) (ErrCode int, ErrMessage error, DBInstanceDetail *rds20140815.DescribeDBInstanceAttributeResponseBodyItemsDBInstanceAttribute) {
	describeDBInstanceAttributeRequest := &rds20140815.DescribeDBInstanceAttributeRequest{
		DBInstanceId: tea.String(InstanceId),
	}
	runtime := &util.RuntimeOptions{}
	describeDBInstanceAttributeResponse, ErrMessage := client.DescribeDBInstanceAttributeWithOptions(describeDBInstanceAttributeRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.DescribeDBInstanceAttribute, ErrMessage, nil
	}
	// 默认一个.
	return Errmsg.SUCCESS, nil, describeDBInstanceAttributeResponse.Body.Items.DBInstanceAttribute[0]
}

// DescribeRDSAccount 查询实例的账号信息
func DescribeRDSAccount(InstanceId string, client *rds20140815.Client) (ErrCode int, ErrMessage error, RDSAccountList []*rds20140815.DescribeAccountsResponseBodyAccountsDBInstanceAccount) {
	describeAccountsRequest := &rds20140815.DescribeAccountsRequest{
		DBInstanceId: tea.String(InstanceId),
		PageSize:     tea.Int32(200),
		PageNumber:   tea.Int32(1),
	}
	runtime := &util.RuntimeOptions{}
	describeAccountsResponse, ErrMessage := client.DescribeAccountsWithOptions(describeAccountsRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorDescribeRDSAccount, ErrMessage, nil
	}
	return Errmsg.SUCCESS, nil, describeAccountsResponse.Body.Accounts.DBInstanceAccount
}

// CreateRDSAccount 创建管理数据库的账号
func CreateRDSAccount(InstanceId string, DBEngine string, TempPasswd string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	setting.LoadAuditAccount()
	var accountType string
	// 不同数据库类型需要创建不同账号 PgSQL 直接设定为Super Mysql设置为normal
	switch DBEngine {
	case "PostgreSQL":
		{
			accountType = "Super"
		}
	case "MySQL":
		{
			accountType = "Normal"
		}
	case "SQLServer":
		{
			accountType = "Normal"
		}
	}
	// 调用sdk创建用户
	createAccountRequest := &rds20140815.CreateAccountRequest{
		DBInstanceId:       tea.String(InstanceId),
		AccountName:        tea.String(setting.AccountName),
		AccountDescription: tea.String(setting.AccountDescription),
		AccountType:        tea.String(accountType),
		AccountPassword:    tea.String(TempPasswd),
	}
	runtime := &util.RuntimeOptions{
		// 超时设置，该产品部分接口调用比较慢，请您适当调整超时时间。
		ReadTimeout:    tea.Int(50000),
		ConnectTimeout: tea.Int(50000),
	}
	_, err := client.CreateAccountWithOptions(createAccountRequest, runtime)
	if err != nil {
		return Errmsg.ErrorCreateRDSAccount, err
	}
	return Errmsg.SUCCESS, nil
}

// GrantAccountPrivilege 授权账号访问数据库
func GrantAccountPrivilege(InstanceId string, DBName string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	setting.LoadAuditAccount()
	grantAccountPrivilegeRequest := &rds20140815.GrantAccountPrivilegeRequest{
		DBInstanceId:     tea.String(InstanceId),
		AccountName:      tea.String(setting.AccountName),
		DBName:           tea.String(DBName),
		AccountPrivilege: tea.String("ReadOnly"),
	}
	runtime := &util.RuntimeOptions{
		// 超时设置，该产品部分接口调用比较慢，请您适当调整超时时间。
		ReadTimeout:    tea.Int(50000),
		ConnectTimeout: tea.Int(50000),
	}
	_, ErrMessage = client.GrantAccountPrivilegeWithOptions(grantAccountPrivilegeRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorGrantAccountPrivilege, ErrMessage
	}
	return Errmsg.SUCCESS, nil
}

// UnlockAccount 解锁RDS PostgresSQL实例的账号
func UnlockAccount(InstanceId string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	unlockAccountRequest := &rds20140815.UnlockAccountRequest{
		DBInstanceId: tea.String(InstanceId),
		AccountName:  tea.String("cnisdp"),
	}
	runtime := &util.RuntimeOptions{
		// 超时设置，该产品部分接口调用比较慢，请您适当调整超时时间。
		ReadTimeout:    tea.Int(50000),
		ConnectTimeout: tea.Int(50000),
	}
	_, err := client.UnlockAccountWithOptions(unlockAccountRequest, runtime)
	if err != nil {
		return Errmsg.ErrorUnlockRDSAccount, err
	}
	return Errmsg.SUCCESS, nil
}

// DeleteAccount 删除审计数据库用户
func DeleteAccount(InstanceId string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	setting.LoadAuditAccount()
	deleteAccountRequest := &rds20140815.DeleteAccountRequest{
		DBInstanceId: tea.String(InstanceId),
		AccountName:  tea.String("cnisdp"),
	}
	runtime := &util.RuntimeOptions{
		// 超时设置，该产品部分接口调用比较慢，请您适当调整超时时间。
		ReadTimeout:    tea.Int(50000),
		ConnectTimeout: tea.Int(50000),
	}
	_, ErrMessage = client.DeleteAccountWithOptions(deleteAccountRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorDeleteAccount, ErrMessage
	}
	return Errmsg.SUCCESS, nil
}

// DescribeDatabases 查询RDS实例下的数据库信息
func DescribeDatabases(InstanceId string, client *rds20140815.Client) (ErrCode int, ErrMessage error, DatabaseList []*rds20140815.DescribeDatabasesResponseBodyDatabasesDatabase) {
	describeDatabasesRequest := &rds20140815.DescribeDatabasesRequest{
		DBInstanceId: tea.String(InstanceId),
		PageSize:     tea.Int32(100),
		PageNumber:   tea.Int32(1),
	}
	runtime := &util.RuntimeOptions{}
	describeDatabasesResponse, ErrMessage := client.DescribeDatabasesWithOptions(describeDatabasesRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorDescribeDatabases, ErrMessage, nil
	}
	return Errmsg.SUCCESS, nil, describeDatabasesResponse.Body.Databases.Database
}

// DescribeDBInstranceSSL 查询实例SSL设置
func DescribeDBInstranceSSL(InstanceId string, client *rds20140815.Client) (ErrCode int, ErrMessage error, SSLEnabled bool) {
	describeDBInstanceSSLRequest := &rds20140815.DescribeDBInstanceSSLRequest{
		DBInstanceId: tea.String(InstanceId),
	}
	runtime := &util.RuntimeOptions{}
	describeDBInstanceSSLResponse, ErrMessage := client.DescribeDBInstanceSSLWithOptions(describeDBInstanceSSLRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorDescribeInstanceSSL, ErrMessage, false
	}
	if *describeDBInstanceSSLResponse.Body.SSLEnabled == "on" || *describeDBInstanceSSLResponse.Body.SSLEnabled == "Yes" {
		return Errmsg.SUCCESS, nil, true
	}
	if *describeDBInstanceSSLResponse.Body.SSLEnabled == "off" || *describeDBInstanceSSLResponse.Body.SSLEnabled == "No" {
		return Errmsg.SUCCESS, nil, false
	}
	// 默认为false吧
	return Errmsg.SUCCESS, nil, false
}

// DescribeDBInstanceNetInfo 查询实例的所有连接地址信息
func DescribeDBInstanceNetInfo(InstanceId string, client *rds20140815.Client) (ErrCode int, ErrMessage error, PublicConnectionString string, IPAddress string) {
	describeDBInstanceNetInfoRequest := &rds20140815.DescribeDBInstanceNetInfoRequest{
		DBInstanceId: tea.String(InstanceId),
	}
	runtime := &util.RuntimeOptions{}
	describeDBInstanceNetInfoResponse, ErrMessage := client.DescribeDBInstanceNetInfoWithOptions(describeDBInstanceNetInfoRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorDescribeDBInstanceNetInfo, ErrMessage, "", ""
	}
	DBInstanceNetInfo := describeDBInstanceNetInfoResponse.Body.DBInstanceNetInfos.DBInstanceNetInfo
	for i := 0; i < len(DBInstanceNetInfo); i++ {
		if *DBInstanceNetInfo[i].IPType == "Public" {
			return Errmsg.SUCCESS, nil, *DBInstanceNetInfo[i].ConnectionString, *DBInstanceNetInfo[i].IPAddress
		}
	}
	return Errmsg.SUCCESS, nil, "", ""
}

// AllocateInstancePublicConnection 申请实例的外网地址
func AllocateInstancePublicConnection(InstanceId string, DBPort string, client *rds20140815.Client) (ErrCode int, ErrMessage error, PublicConnectionString string) {
	allocateInstancePublicConnectionRequest := &rds20140815.AllocateInstancePublicConnectionRequest{
		DBInstanceId:           tea.String(InstanceId),
		ConnectionStringPrefix: tea.String(utils.GenerateString()),
		Port:                   tea.String(DBPort),
	}
	runtime := &util.RuntimeOptions{
		// 超时设置，该产品部分接口调用比较慢，请您适当调整超时时间。
		ReadTimeout:    tea.Int(50000),
		ConnectTimeout: tea.Int(50000),
	}
	allocateInstancePublicConnectionResponse, ErrMessage := client.AllocateInstancePublicConnectionWithOptions(allocateInstancePublicConnectionRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorAllocateInstancePublicConnection, ErrMessage, ""
	}
	return Errmsg.SUCCESS, nil, *allocateInstancePublicConnectionResponse.Body.ConnectionString
}

// ModifySecurityIps 修改RDS实例IP白名单
func ModifySecurityIps(InstanceId string, ModifyMode string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	setting.LoadServer()
	modifySecurityIpsRequest := &rds20140815.ModifySecurityIpsRequest{
		ModifyMode:            tea.String(ModifyMode),
		DBInstanceId:          tea.String(InstanceId),
		SecurityIps:           tea.String(setting.IPAddr),
		DBInstanceIPArrayName: tea.String("cnisdp"),
	}
	runtime := &util.RuntimeOptions{
		// 超时设置，该产品部分接口调用比较慢，请您适当调整超时时间。
		ReadTimeout:    tea.Int(50000),
		ConnectTimeout: tea.Int(50000),
	}
	_, ErrMessage = client.ModifySecurityIpsWithOptions(modifySecurityIpsRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorModifySecurityIps, ErrMessage
	}
	return Errmsg.SUCCESS, nil
}

// DescribeDBInstanceIPArrayList 查询RDS实例IP白名单
func DescribeDBInstanceIPArrayList(InstanceId string, client *rds20140815.Client) (ErrCode int, ErrMessage error, WhitelistArray []*rds20140815.DescribeDBInstanceIPArrayListResponseBodyItemsDBInstanceIPArray) {
	describeDBInstanceIPArrayListRequest := &rds20140815.DescribeDBInstanceIPArrayListRequest{
		DBInstanceId: tea.String(InstanceId),
	}
	runtime := &util.RuntimeOptions{
		// 超时设置，该产品部分接口调用比较慢，请您适当调整超时时间。
		ReadTimeout:    tea.Int(50000),
		ConnectTimeout: tea.Int(50000),
	}
	// 复制代码运行请自行打印 API 的返回值
	response, ErrMessage := client.DescribeDBInstanceIPArrayListWithOptions(describeDBInstanceIPArrayListRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorDescribeDBInstanceIPArrayList, ErrMessage, nil
	}
	WhitelistArray = response.Body.Items.DBInstanceIPArray
	return Errmsg.SUCCESS, nil, WhitelistArray
}

// ReleaseInstancePublicConnection 释放实例的外网连接地址
func ReleaseInstancePublicConnection(InstanceId string, PublicConnectionString string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	releaseInstancePublicConnectionRequest := &rds20140815.ReleaseInstancePublicConnectionRequest{
		DBInstanceId:            tea.String(InstanceId),
		CurrentConnectionString: tea.String(PublicConnectionString),
	}
	runtime := &util.RuntimeOptions{
		// 超时设置，该产品部分接口调用比较慢，请您适当调整超时时间。
		ReadTimeout:    tea.Int(50000),
		ConnectTimeout: tea.Int(50000),
	}
	_, ErrMessage = client.ReleaseInstancePublicConnectionWithOptions(releaseInstancePublicConnectionRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorReleaseInstancePublicConnection, ErrMessage
	}
	return Errmsg.SUCCESS, nil
}

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
func CreateRDSAccount(DBInstanceId string, client *rds20140815.Client) {
	setting.LoadAuditAccount()
	// 调用sdk创建用户
	createAccountRequest := &rds20140815.CreateAccountRequest{
		DBInstanceId:       tea.String(DBInstanceId),
		AccountName:        tea.String(setting.AccountName),
		AccountDescription: tea.String(setting.AccountDescription),
		AccountType:        tea.String(setting.AccountType),
		AccountPassword:    tea.String(setting.AccountPassword),
	}
	runtime := &util.RuntimeOptions{
		// 超时设置，该产品部分接口调用比较慢，请您适当调整超时时间。
		ReadTimeout:    tea.Int(50000),
		ConnectTimeout: tea.Int(50000),
	}
	_, _ = client.CreateAccountWithOptions(createAccountRequest, runtime)
}

// UnlockAccount 解锁RDS PostgresSQL实例的账号
func UnlockAccount(DBInstanceId string, client *rds20140815.Client) {
	unlockAccountRequest := &rds20140815.UnlockAccountRequest{
		DBInstanceId: tea.String(DBInstanceId),
		AccountName:  tea.String("cnisdp"),
	}
	runtime := &util.RuntimeOptions{
		// 超时设置，该产品部分接口调用比较慢，请您适当调整超时时间。
		ReadTimeout:    tea.Int(50000),
		ConnectTimeout: tea.Int(50000),
	}
	_, _ = client.UnlockAccountWithOptions(unlockAccountRequest, runtime)
}

// DescribeDatabases 查询RDS实例下的数据库信息
func DescribeDatabases(DBInstanceId string, client *rds20140815.Client) []*rds20140815.DescribeDatabasesResponseBodyDatabasesDatabase {
	describeDatabasesRequest := &rds20140815.DescribeDatabasesRequest{
		DBInstanceId: tea.String(DBInstanceId),
		PageSize:     tea.Int32(100),
		PageNumber:   tea.Int32(1),
	}
	runtime := &util.RuntimeOptions{}
	describeDatabasesResponse, _ := client.DescribeDatabasesWithOptions(describeDatabasesRequest, runtime)
	return describeDatabasesResponse.Body.Databases.Database
}

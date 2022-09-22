package aliyunRDS

/*
InitConnect：
1. 获取RDS实例清单
2. 获取RDS实例下数据库db清单
3. 授权ISDP审计账号权限（判断、新增、授权）
*/

import (
	"Backend/internal/application/setting"
	"Backend/internal/model"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	rds20140815 "github.com/alibabacloud-go/rds-20140815/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func InitRDSData(Group string) int {
	// 清空数据库中该项目组数据
	//model.RestoreData(Group)
	// 获取key值
	AccessKey := viper.GetString(Group + ".AccessKey")
	AccessSecret := viper.GetString(Group + ".AccessSecret")
	//打开client
	Client, _err := CreateClient(tea.String(AccessKey), tea.String(AccessSecret))
	//初始化RDS实例列表
	initDescribeRDSInstances(Group, Client)
	if _err != nil {
		return 500
	} else {
		return 200
	}
}

// CreateClient 创建连接客户端
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *rds20140815.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("rds.aliyuncs.com")
	_result = &rds20140815.Client{}
	_result, _err = rds20140815.NewClient(config)
	return _result, _err
}

// 查询项目组RDS实例列表
func initDescribeRDSInstances(Group string, client *rds20140815.Client) {
	var input model.DataInventory // input 表示录入到数据库中的值
	describeDBInstancesRequest := &rds20140815.DescribeDBInstancesRequest{
		RegionId:   tea.String("cn-shanghai"), // 目前默认都在cn-shanghai
		PageSize:   tea.Int32(100),
		PageNumber: tea.Int32(1),
	}
	rawresponse, _err := client.DescribeDBInstances(describeDBInstancesRequest)
	if _err != nil {
		return
	}
	//rawInstancelist指aliyunSDK返回的原始数据
	rawInstancelist := rawresponse.Body.Items.DBInstance
	// 遍历每一个 RDS实例
	for i := 0; i < len(rawInstancelist); i++ {
		// 判断 rds实例是否已经在数据库中
		if !model.CheckInventoryExist(*rawInstancelist[i].DBInstanceId) {
			input.UUID = uuid.New().String()
			input.CreationTime = *rawInstancelist[i].CreateTime
			input.GroupName = Group
			input.RDSInstanceID = *rawInstancelist[i].DBInstanceId
			if rawInstancelist[i].DBInstanceDescription != nil {
				input.RDSInstanceDescription = *rawInstancelist[i].DBInstanceDescription
			} else {
				input.RDSInstanceDescription = ""
			}
			input.RDSInstanceNetType = *rawInstancelist[i].DBInstanceNetType
			input.RDSInstanceType = *rawInstancelist[i].DBInstanceType
			input.RDSInstanceNetworkType = *rawInstancelist[i].InstanceNetworkType
			input.RDSEngine = *rawInstancelist[i].Engine
			input.RDSEngineVersion = *rawInstancelist[i].EngineVersion
			input.RDSInstanceStatus = *rawInstancelist[i].DBInstanceStatus
			input.RDSConnectionString = *rawInstancelist[i].ConnectionString
			input.RegionId = *rawInstancelist[i].RegionId
			model.AddInventory(&input)
		}
		// 初始化RDS实例下审计账号
		initAuditAccount(*rawInstancelist[i].DBInstanceId, client)
		// 初始化RDS实例下数据库
		initDescribeDatabases(*rawInstancelist[i].DBInstanceId, client)
	}
}

/*
	初始化 isdp审计账号
		- 遍历各个RDS实例下账号
		- 判断是否已有ISDP审计账号
		- 启用 / 创建 ISDP 账号
		- 获取instance下DB清单
		- 授权 DB
		- 锁定 ISDP账号
*/
func initAuditAccount(InstanceName string, client *rds20140815.Client) {
	var input model.DatabaseAccount // input 表示录入到数据库中的值
	describeAccountsRequest := &rds20140815.DescribeAccountsRequest{
		DBInstanceId: tea.String(InstanceName),
		PageSize:     tea.Int32(200),
		PageNumber:   tea.Int32(1),
	}
	runtime := &util.RuntimeOptions{}
	Accountlist, _err := client.DescribeAccountsWithOptions(describeAccountsRequest, runtime)
	if _err != nil {
		return
	}
	//rawAccountlist指aliyunSDK返回的原始数据
	rawAccountlist := Accountlist.Body.Accounts.DBInstanceAccount
	// 更新数据库账户列表
	for i := 0; i < len(rawAccountlist); i++ {
		if !model.CheckAccountExist(InstanceName, *rawAccountlist[i].AccountName) {
			input.UUID = uuid.New().String()
			if rawAccountlist[i].AccountDescription != nil {
				input.AccountDescription = *rawAccountlist[i].AccountDescription
			} else {
				input.AccountDescription = ""
			}
			input.AccountStatus = *rawAccountlist[i].AccountStatus
			input.DBInstanceId = *rawAccountlist[i].DBInstanceId
			input.AccountType = *rawAccountlist[i].AccountType
			input.AccountName = *rawAccountlist[i].AccountName
			model.AddDatabaseAccount(&input)
		}
	}
	// 判断是否已开通审计账户，未开通则创建cnisdp账号
	if !model.CheckAccountExist(InstanceName, "cnisdp") {
		setting.LoadAuditAccount()
		// 调用sdk创建用户
		createAccountRequest := &rds20140815.CreateAccountRequest{
			DBInstanceId:       tea.String(InstanceName),
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
		_, _err = client.CreateAccountWithOptions(createAccountRequest, runtime)
		// 添加审计用户到账号列表中
		input.UUID = uuid.New().String()
		input.AccountDescription = setting.AccountDescription
		input.AccountStatus = "Available"
		input.DBInstanceId = InstanceName
		input.AccountType = setting.AccountType
		input.AccountName = setting.AccountName
		model.AddDatabaseAccount(&input)
	}
	// 解锁cnisdp账号
	unlockAccountRequest := &rds20140815.UnlockAccountRequest{
		DBInstanceId: tea.String(InstanceName),
		AccountName:  tea.String("cnisdp"),
	}
	runtime = &util.RuntimeOptions{
		// 超时设置，该产品部分接口调用比较慢，请您适当调整超时时间。
		ReadTimeout:    tea.Int(50000),
		ConnectTimeout: tea.Int(50000),
	}
	_, _err = client.UnlockAccountWithOptions(unlockAccountRequest, runtime)
}

func initDescribeDatabases(InstanceName string, client *rds20140815.Client) {
	var inputDB model.DataDatabase //input表示输入到系统中的数据
	var inputAccountPrivilege model.AccountPrivilege
	describeDatabasesRequest := &rds20140815.DescribeDatabasesRequest{
		DBInstanceId: tea.String(InstanceName),
		PageSize:     tea.Int32(100),
		PageNumber:   tea.Int32(1),
	}
	runtime := &util.RuntimeOptions{}
	rawresponse, _err := client.DescribeDatabasesWithOptions(describeDatabasesRequest, runtime)
	if _err != nil {
		return
	}
	dblist := rawresponse.Body.Databases.Database
	for i := 0; i < len(dblist); i++ {
		inputDB.UUID = uuid.New().String()
		inputDB.RDSInstanceID = *dblist[i].DBInstanceId
		if dblist[i].DBDescription != nil {
			inputDB.DatabaseDescription = *dblist[i].DBDescription
		} else {
			inputDB.DatabaseDescription = ""
		}
		inputDB.DatabaseName = *dblist[i].DBName
		inputDB.DatabaseStatus = *dblist[i].DBStatus
		inputDB.DatabaseEngine = *dblist[i].Engine
		model.AddDatabase(&inputDB)
		// 输入数据库权限清单
		accountlist := dblist[i].Accounts.AccountPrivilegeInfo
		for j := 0; j < len(accountlist); j++ {
			inputAccountPrivilege.UUID = uuid.New().String()
			inputAccountPrivilege.RDSInstanceID = *dblist[i].DBInstanceId
			inputAccountPrivilege.DatabaseName = *dblist[i].DBName
			inputAccountPrivilege.AccountName = *accountlist[j].Account
			inputAccountPrivilege.Privilege = *accountlist[j].AccountPrivilege
			model.AddAccountPrivilege(&inputAccountPrivilege)
		}
	}
}

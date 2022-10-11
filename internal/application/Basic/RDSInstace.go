package Basic

import (
	"Backend/internal/application/aliyunSDK"
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	rds20140815 "github.com/alibabacloud-go/rds-20140815/v2/client"
	"github.com/google/uuid"
)

/*
初始化
- 初始化cloud账号列表
- 初始化rds列表
- 初始化db列表
*/

// UpdateRDS 当accountId是""时，初始化所有，也可以对单个accountId进行更新
func UpdateRDS(CloudAccountID string) (ErrCode int, ErrMessage error) {
	// 若不指定对应accountID，则获取全量CloudAccountList
	previousCode, previousMsg, cloudAccount := local.GetCloudAccount(CloudAccountID)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	// 删除数据库中CloudAccountID下所有 RDS 实例
	if CloudAccountID != "" {
		local.DeleteCloudInventory(CloudAccountID)
	}
	// 遍历 CloudAccountList 中RDS实例 进行对应更新
	for i := 0; i < len(cloudAccount); i++ {
		// 创建RDSClient
		previousCode, previousMsg, client := aliyunSDK.CreateRDSClient(cloudAccount[i].AccountId)
		if previousCode != Errmsg.SUCCESS {
			return previousCode, previousMsg
		}
		// 更新该云账号下RDS 实例Instance List信息
		previousCode, previousMsg = UpdateRDSInstanceList(cloudAccount[i].AccountId, client)
		if previousCode != Errmsg.SUCCESS {
			return previousCode, previousMsg
		}
	}
	return Errmsg.SUCCESS, nil
}

// UpdateRDSInstanceList 更新该云账号下RDS实例详情、实例账号信息
func UpdateRDSInstanceList(CloudAccountID string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	var RDSInstanceList []*rds20140815.DescribeDBInstancesResponseBodyItemsDBInstance
	previousCode, previousMsg, RDSInstanceList := aliyunSDK.DescribeRDSInstances(client)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	// 更新云账号下rds数量
	local.UpdateCloudAccountRDS(CloudAccountID, len(RDSInstanceList))
	// 若cloudAccount下存在rds实例，则将获取到的实例录入到数据库，新实例数据则为插入，原有实例则为update
	if len(RDSInstanceList) > 0 {
		for i := 0; i < len(RDSInstanceList); i++ {
			// 更新云账号下RDS 实例详情
			previousCode, previousMsg = UpdateRDSInstanceDetails(CloudAccountID, *RDSInstanceList[i].DBInstanceId, client)
			if previousCode != Errmsg.SUCCESS {
				return previousCode, previousMsg
			}
			// 更新云账号下RDS账号信息
			previousCode, previousMsg = UpdateRDSAccount(*RDSInstanceList[i].DBInstanceId, client)
			if previousCode != Errmsg.SUCCESS {
				return previousCode, previousMsg
			}
			// 更新该云账号下RDS DBList 信息
			previousCode, previousMsg = UpdateDatabaseList(*RDSInstanceList[i].DBInstanceId, client)
			if previousCode != Errmsg.SUCCESS {
				return previousCode, previousMsg
			}
		}
	}
	return Errmsg.SUCCESS, nil
}

// UpdateRDSInstanceDetails 更新RDS实例详情
func UpdateRDSInstanceDetails(CloudAccountID string, InstanceId string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	var InstanceDetails *rds20140815.DescribeDBInstanceAttributeResponseBodyItemsDBInstanceAttribute
	previousCode, previousMsg, InstanceDetails := aliyunSDK.DescribeDBInstanceAttribute(InstanceId, client)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	var input local.DataInventory
	// 录入前先进行格式化数据
	input.UUID = uuid.New().String()
	input.CreationTime = *InstanceDetails.CreationTime
	input.CloudAccountID = CloudAccountID
	input.RDSInstanceID = *InstanceDetails.DBInstanceId
	if InstanceDetails.DBInstanceDescription != nil {
		input.RDSInstanceDescription = *InstanceDetails.DBInstanceDescription
	} else {
		input.RDSInstanceDescription = ""
	}
	input.RDSInstanceNetType = *InstanceDetails.DBInstanceNetType
	if len(InstanceDetails.ReadOnlyDBInstanceIds.ReadOnlyDBInstanceId) > 0 {
		input.RDSReadOnlyDBInstanceId = *InstanceDetails.ReadOnlyDBInstanceIds.ReadOnlyDBInstanceId[0].DBInstanceId
	}
	input.RDSInstanceType = *InstanceDetails.DBInstanceType
	input.RDSInstanceNetworkType = *InstanceDetails.InstanceNetworkType
	input.RDSEngine = *InstanceDetails.Engine
	input.RDSEngineVersion = *InstanceDetails.EngineVersion
	input.RDSInstanceStatus = *InstanceDetails.DBInstanceStatus
	input.RDSConnectionString = *InstanceDetails.ConnectionString
	input.RDSConnectionPort = *InstanceDetails.Port
	input.RegionId = *InstanceDetails.RegionId
	// 获取SSL是否开启
	previousCode, previousMsg, input.SSLEnabled = aliyunSDK.DescribeDBInstranceSSL(InstanceId, client)
	// 更新Inventory
	local.DeleteInventory(input.RDSInstanceID)
	previousCode, previousMsg = local.AddInventory(&input)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	return Errmsg.SUCCESS, nil
}

// UpdateRDSAccount 更新RDS实例账号信息
func UpdateRDSAccount(InstanceID string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	var RDSAccountList []*rds20140815.DescribeAccountsResponseBodyAccountsDBInstanceAccount
	var input local.RDSAccount
	previousCode, previousMsg, RDSAccountList := aliyunSDK.DescribeRDSAccount(InstanceID, client)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	for i := 0; i < len(RDSAccountList); i++ {
		// 先删除 后更新
		local.DeleteAccount(InstanceID, *RDSAccountList[i].AccountName)
		if RDSAccountList[i].AccountDescription != nil {
			input.AccountDescription = *RDSAccountList[i].AccountDescription
		} else {
			input.AccountDescription = ""
		}
		input.AccountStatus = *RDSAccountList[i].AccountStatus
		input.RDSInstanceID = *RDSAccountList[i].DBInstanceId
		input.AccountType = *RDSAccountList[i].AccountType
		input.AccountName = *RDSAccountList[i].AccountName
		local.AddDatabaseAccount(&input)
	}
	return Errmsg.SUCCESS, nil
}

// UpdateDatabaseList 更新RDS实例下Database 列表
func UpdateDatabaseList(InstanceID string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	var DatabaseList []*rds20140815.DescribeDatabasesResponseBodyDatabasesDatabase
	var input local.DataDatabase
	previousCode, previousMsg, DatabaseList := aliyunSDK.DescribeDatabases(InstanceID, client)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	// 先删除 后更新
	local.DeleteDatabase(InstanceID)
	for i := 0; i < len(DatabaseList); i++ {
		input.UUID = uuid.New().String()
		input.InstanceID = *DatabaseList[i].DBInstanceId
		if DatabaseList[i].DBDescription != nil {
			input.Description = *DatabaseList[i].DBDescription
		}
		input.Status = *DatabaseList[i].DBStatus
		input.DatabaseName = *DatabaseList[i].DBName
		input.Engine = *DatabaseList[i].Engine
		local.AddDatabase(&input)
	}
	return Errmsg.SUCCESS, nil
}

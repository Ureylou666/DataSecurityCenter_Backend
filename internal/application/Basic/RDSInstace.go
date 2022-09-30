package Basic

import (
	"Backend/internal/application/aliyunSDK"
	"Backend/internal/model"
	"Backend/internal/utils/Errmsg"
	"fmt"
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
	previousCode, previousMsg, cloudAccount := model.GetCloudAccount(CloudAccountID)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	// 遍历 CloudAccountList 中RDS实例 进行对应更新
	for i := 0; i < len(cloudAccount); i++ {
		// 创建RDSClient
		previousCode, previousMsg, client := aliyunSDK.CreateRDSClient(cloudAccount[i].AccountId)
		if previousCode == Errmsg.SUCCESS {
			// 更新该云账号下RDS实例信息
			_, _ = UpdateRDSInstanceList(cloudAccount[i].AccountId, client)
		} else {
			fmt.Println(previousCode, ": ", previousMsg)
		}
	}
	return Errmsg.SUCCESS, nil
}

// UpdateRDSInstanceList 更新该云账号下RDS实例信息
func UpdateRDSInstanceList(CloudAccountID string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	var RDSInstanceList []*rds20140815.DescribeDBInstancesResponseBodyItemsDBInstance
	previousCode, previousMsg, RDSInstanceList := aliyunSDK.DescribeRDSInstances(client)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	// 更新云账号下rds数量
	model.UpdateCloudAccountRDS(CloudAccountID, len(RDSInstanceList))
	// 若cloudAccount下存在rds实例，则将获取到的实例录入到数据库，新实例数据则为插入，原有实例则为update
	if len(RDSInstanceList) > 0 {
		for i := 0; i < len(RDSInstanceList); i++ {
			previousCode, previousMsg = UpdateRDSInstanceDetails(CloudAccountID, *RDSInstanceList[i].DBInstanceId, client)
			if previousCode != Errmsg.SUCCESS {
				return previousCode, previousMsg
			}
			// 更新云账号下RDS账号信息
			previousCode, previousMsg = UpdateRDSAccount(*RDSInstanceList[i].DBInstanceId, client)
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
	var input model.DataInventory
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
	// 更新Inventory
	model.DeleteInventory(input.RDSInstanceID)
	previousCode, previousMsg = model.AddInventory(&input)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	return Errmsg.SUCCESS, nil
}

// UpdateRDSAccount 更新RDS实例账号信息
func UpdateRDSAccount(InstanceID string, client *rds20140815.Client) (ErrCode int, ErrMessage error) {
	var RDSAccountList []*rds20140815.DescribeAccountsResponseBodyAccountsDBInstanceAccount
	var input model.DatabaseAccount
	previousCode, previousMsg, RDSAccountList := aliyunSDK.DescribeRDSAccount(InstanceID, client)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	for i := 0; i < len(RDSAccountList); i++ {
		// 先删除 后更新
		model.DeleteAccount(InstanceID, *RDSAccountList[i].AccountName)
		input.UUID = uuid.New().String()
		if RDSAccountList[i].AccountDescription != nil {
			input.AccountDescription = *RDSAccountList[i].AccountDescription
		} else {
			input.AccountDescription = ""
		}
		input.AccountStatus = *RDSAccountList[i].AccountStatus
		input.DBInstanceId = *RDSAccountList[i].DBInstanceId
		input.AccountType = *RDSAccountList[i].AccountType
		input.AccountName = *RDSAccountList[i].AccountName
		model.AddDatabaseAccount(&input)
	}
	return Errmsg.SUCCESS, nil
}

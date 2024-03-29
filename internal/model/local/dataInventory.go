package local

import (
	"Backend/internal/utils/Errmsg"
)

// DataInventory 定义数据资产数据结构，与aliyun sdk中定义一致

// RDS实例详情
type DataInventory struct {
	UUID                    string `gorm:"primaryKey" json:"UUID"`
	CreationTime            string `gorm:"type:varchar(50)" json:"CreationTime"`            //创建该数据资产实例的时间。使用时间戳表示，单位：毫秒。
	CloudAccountID          string `gorm:"type:varchar(50)" json:"CloudAccountID"`          //数据资产实例所属部门的名称。
	RDSInstanceID           string `gorm:"type:varchar(50)" json:"RDSInstanceID"`           //RDS实例ID
	RDSInstanceDescription  string `gorm:"type:varchar(200)" json:"RDSInstanceDescription"` //数据资产实例的描述信息。
	RDSInstanceNetType      string `gorm:"type:varchar(50)" json:"RDSInstanceNetType"`      //内网访问OR外网访问
	RDSInstanceType         string `gorm:"type:varchar(50)" json:"RDSInstanceType"`         //实例类型，取值： Primary：主实例、Readonly：只读实例、Guard：灾备实例、Temp：临时实例
	RDSInstanceNetworkType  string `gorm:"type:varchar(50)" json:"RDSInstanceNetworkType"`  //实例的网络类型，取值：	VPC：专有网络下的实例、Classic：经典网络下的实例
	RDSReadOnlyDBInstanceId string `gorm:"type:varchar(50)" json:"ReadOnlyDBInstanceId"`    //只读实例ID,取第一个即可
	RDSEngine               string `gorm:"type:varchar(50)" json:"RDSEngine"`               //数据库类型，取值：MySQL、SQLServer、PostgreSQL、MariaDB
	RDSEngineVersion        string `gorm:"type:varchar(50)" json:"RDSEngineVersion"`        //数据库版本
	RDSInstanceStatus       string `gorm:"type:varchar(50)" json:"RDSInstanceStatus"`       //实例状态
	RDSConnectionString     string `gorm:"type:varchar(100)" json:"RDSConnectionString"`    //连接字符串
	SSLEnabled              bool   `gorm:"type:boolean" json:"SSLEnabled"`                  // 是否支持ssl
	RDSConnectionPort       string `gorm:"type:varchar(10)" json:"RDSConnectionPort"`       //连接端口
	RegionId                string `gorm:"type:varchar(50)" json:"RegionId"`                //地区
	DepartName              string `gorm:"type:varchar(50)" json:"DepartName"`              //主体
	DBCount                 int    `gorm:"type:int" json:"DBCount"`                         //实例下数据库数量
	SensitiveCount          int    `gorm:"type:int" json:"SensitiveCount"`                  //实例下敏感数据库数量
}

// RDSQueryInfo 定义api接口查询参数
type RDSQueryInfo struct {
	AccountID string
	PageNum   int
	PageSize  int
}

// CheckInventoryExist 判断数据库中是否存在相关资产
func CheckInventoryExist(InstanceID string) bool {
	var Num int64
	db.Where("rds_instance_id = ?", InstanceID).Find(&DataInventory{}).Count(&Num)
	if Num > 0 {
		return true
	} else {
		return false
	}
}

// AddInventory 新增aliyun数据实例资产
func AddInventory(data *DataInventory) (ErrCode int, ErrMessage error) {
	ErrMessage = db.Create(&data).Error
	if ErrMessage != nil {
		return Errmsg.ErrorAddInventory, ErrMessage
	}
	return Errmsg.SUCCESS, nil
}

// DeleteInventory 通过RDS InstanceID 删除资产
func DeleteInventory(InstanceID string) {
	db.Where("rds_instance_id = ?", InstanceID).Delete(&DataInventory{})
}

// DeleteCloudInventory 通过CloudAccountID 删除资产
func DeleteCloudInventory(CloudAccountID string) {
	db.Where("cloud_account_id = ?", CloudAccountID).Delete(&DataInventory{})
}

// GetAccountID 通过InstanceID获取
func GetAccountID(InstanceID string) (CloudAccountID string) {
	var result DataInventory
	db.Where("rds_instance_id = ?", InstanceID).First(&result)
	return result.CloudAccountID
}

// GetConnectString 通过InstanceID获取ConnectString相关
func GetConnectString(InstanceID string) (RdsDetails DataInventory) {
	db.Where("rds_instance_id = ?", InstanceID).First(&RdsDetails)
	return RdsDetails
}

// GetRDSInventory 获取aliyun RDS资产列表 并进行分页展示
func GetRDSInventory(query RDSQueryInfo) ([]DataInventory, int64, int64) {
	var result []DataInventory
	var resTotal, inventoryTotal int64
	// 获取数据库中Instance总数
	// cloudaccountid 精确搜索
	// 分页处理
	if query.AccountID == "" {
		db.Find(&result).Count(&inventoryTotal)
		if query.PageNum == 0 || query.PageSize == 0 {
			db.Limit(-1).Find(&result)
			resTotal = int64(len(result))
		} else {
			db.Limit(query.PageSize).Offset((query.PageNum - 1) * query.PageSize).Find(&result)
			resTotal = int64(len(result))
		}
		return result, resTotal, inventoryTotal
	}
	db.Where("cloud_account_id = ?", query.AccountID).Find(&result).Count(&inventoryTotal)
	if query.PageNum == 0 || query.PageSize == 0 {
		db.Where("cloud_account_id = ?", query.AccountID).Limit(-1).Find(&result)
		resTotal = int64(len(result))
	} else {
		db.Where("cloud_account_id = ?", query.AccountID).Limit(query.PageSize).Offset((query.PageNum - 1) * query.PageSize).Find(&result)
		resTotal = int64(len(result))
	}
	return result, resTotal, inventoryTotal
}

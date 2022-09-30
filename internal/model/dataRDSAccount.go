package model

import (
	"Backend/internal/utils/Errmsg"
)

type DatabaseAccount struct {
	UUID               string `gorm:"primaryKey" json:"UUID"`
	RDSInstanceID      string `gorm:"type:varchar(50)" json:"RDSInstanceID"`       //RDS实例ID
	AccountDescription string `gorm:"type:varchar(200)" json:"AccountDescription"` //账号描述信息。
	AccountStatus      string `gorm:"type:varchar(200)" json:"AccountStatus"`      //账号状态
	DBInstanceId       string `gorm:"type:varchar(200)" json:"DBInstanceId"`       //实例名
	AccountType        string `gorm:"type:varchar(200)" json:"AccountType"`        //账号类型 是否管理员
	AccountName        string `gorm:"type:varchar(200)" json:"AccountName"`        // 账号名
}

// AddDatabaseAccount 新增aliyun数据实例资产
func AddDatabaseAccount(data *DatabaseAccount) int {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
}

func CheckAccountExist(InstanceID string, AccountName string) bool {
	var Num int64
	db.Where("db_instance_id = ? AND account_name = ?", InstanceID, AccountName).Find(&DatabaseAccount{}).Count(&Num)
	if Num > 0 {
		return true
	} else {
		return false
	}
}

func DeleteAccount(InstanceID string, AccountName string) {
	db.Where("db_instance_id = ? AND account_name = ?", InstanceID, AccountName).Delete(&DatabaseAccount{})
}

package model

import "Backend/internal/Errmsg"

type AccountPrivilege struct {
	UUID          string `gorm:"primaryKey" json:"UUID"`
	RDSInstanceID string `gorm:"type:varchar(50)" json:"RDSInstanceID"` //数据库所属实例id
	DatabaseName  string `gorm:"type:varchar(200)" json:"DatabaseName"` //数据资产实例的名称。
	AccountName   string `gorm:"type:varchar(200)" json:"AccountName"`  //数据库账号名
	Privilege     string `gorm:"type:varchar(50)" json:"Privilege"`     //数据库账号权限
}

// AddAccountPrivilege 新增数据库账号权限
func AddAccountPrivilege(data *AccountPrivilege) int {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
}

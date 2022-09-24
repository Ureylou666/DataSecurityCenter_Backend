package model

import "Backend/internal/Errmsg"

type CloudAccount struct {
	UUID        string `gorm:"primaryKey" json:"UUID"`
	Type        string `gorm:"type:varchar(200)" json:"Type"`        //成员类型。取值：CloudAccount：云账号。ResourceAccount：资源账号
	DisplayName string `gorm:"type:varchar(200)" json:"DisplayName"` //成员名称。
	JoinTime    string `gorm:"type:varchar(200)" json:"JoinTime"`    //成员加入时间（UTC时间）
	AccountId   string `gorm:"type:varchar(200)" json:"AccountId"`   //成员账号ID。
	ModifyTime  string `gorm:"type:varchar(200)" json:"ModifyTime"`  //成员修改时间（UTC时间）。
}

// AddData 新增aliyun数据实例资产
func AddCloudAccount(data *CloudAccount) int {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
}

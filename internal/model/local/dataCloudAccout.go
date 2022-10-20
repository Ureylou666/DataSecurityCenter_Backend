package local

import (
	"Backend/internal/utils/Errmsg"
)

type CloudAccount struct {
	UUID        string `gorm:"primaryKey" json:"UUID"`
	GroupName   string `gorm:"type:varchar(50)" json:"GroupName,omitempty"`
	Type        string `gorm:"type:varchar(200)" json:"Type"`        //成员类型。取值：CloudAccount：云账号。ResourceAccount：资源账号
	DisplayName string `gorm:"type:varchar(200)" json:"DisplayName"` //成员名称。
	JoinTime    string `gorm:"type:varchar(200)" json:"JoinTime"`    //成员加入时间（UTC时间）
	AccountId   string `gorm:"type:varchar(200)" json:"AccountId"`   //成员账号ID。
	RDSCount    int    `gorm:"type:int" json:"RDSCount"`             //成员RDS数量
	ModifyTime  string `gorm:"type:varchar(200)" json:"ModifyTime"`  //成员修改时间（UTC时间）。
}

type CloudAccountQueryInfo struct {
	GroupName string
	PageNum   int
	PageSize  int
}

// AddCloudAccount 新增 aliyun Cloud Account
func AddCloudAccount(data *CloudAccount) (ErrCode int, ErrMessage error) {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR, err
	}
	return Errmsg.SUCCESS, nil
}

// UpdateCloudAccountRDS 更新云账号下RDS数量
func UpdateCloudAccountRDS(CloudAccountID string, RDSCount int) {
	db.Model(&CloudAccount{}).Where("account_id = ?", CloudAccountID).Update("RDSCount", RDSCount)
}

// GetCloudAccount 查询返回 aliyun Cloud Account
func GetCloudAccount(CloudAccountID string) (ErrCode int, ErrMessage error, CloudAccountList []CloudAccount) {
	// Account为"" 返回 全量
	if CloudAccountID == "" {
		db.Find(&CloudAccountList)
	} else {
		db.Where("account_id= ?", CloudAccountID).Find(&CloudAccountList)
	}
	return Errmsg.SUCCESS, nil, CloudAccountList
}

// GetCloudAccountList 查询返回aliyun 账号清单
func GetCloudAccountList(query CloudAccountQueryInfo) (result []CloudAccount, resTotal int64, AccountListTotal int64) {
	query.GroupName = "%" + query.GroupName + "%"
	db.Where("group_name like ?", query.GroupName).Find(&result).Count(&AccountListTotal)
	if query.PageNum == 0 || query.PageSize == 0 {
		db.Where("group_name like ?", query.GroupName).Limit(-1).Order("rds_count desc").Find(&result)
		resTotal = int64(len(result))
	} else {
		db.Where("group_name like ?", query.GroupName).Order("rds_count desc").Limit(query.PageSize).Offset((query.PageNum - 1) * query.PageSize).Find(&result)
		resTotal = int64(len(result))
	}
	return result, resTotal, AccountListTotal
}

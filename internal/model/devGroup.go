package model

import (
	"Backend/internal/Errmsg"
	"github.com/google/uuid"
)

type InsertInfo struct {
	GroupName 	string
	GroupOwner 	string
	OwnerMail 	string
	SecPo 		string
	SecPoMail	string
}

type DevGroup struct {
	UUID				string 	`gorm:"primaryKey" json:"UUID"`
	GroupName 			string	`gorm:"type:varchar(50)" json:"GroupName"`  // 项目组名
	GroupOwner			string 	`gorm:"type:varchar(50)" json:"GroupOwner"` // 项目组负责人
	OwnerMail			string 	`gorm:"type:varchar(50)" json:"OwnerMail"`	// 负责人邮箱
	SecPo				string 	`gorm:"type:varchar(50)" json:"SecPo"`  // 安全负责人
	SecPoMail			string	`gorm:"type:varchar(50)" json:"SecPoMail"`  // 安全负责人邮箱
//	LastUpdateTime      string  `gorm:"type:varchar(50)" json:"LastFinishTime"` //最近一次更新时间
}

type GroupQueryInfo struct {
	PageNum  		int
	PageSize 		int
}

// AddDevGroup 新建开发项目组
func AddDevGroup(data InsertInfo) int {
	var input DevGroup
	// 初始化值
	input.UUID = uuid.New().String()
	input.GroupName = data.GroupName
	input.GroupOwner = data.GroupOwner
	input.OwnerMail = data.OwnerMail
	input.SecPo = data.SecPo
	input.SecPoMail = data.SecPoMail
	err := db.Create(&input).Error
	if err!=nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
}

// GetDevGroup 获取项目组列表
func GetDevGroup(query GroupQueryInfo) ([]DevGroup, int64, int64) {
	var result []DevGroup
	var resTotal, groupTotal int64
	// 获取总数
	db.Find(&result).Limit(-1).Count(&groupTotal)
	// 分页处理
	if query.PageNum == 0 || query.PageSize == 0 {
		db.Limit(-1).Find(&result)
		resTotal = int64(len(result))
	} else {
		db.Limit(query.PageSize).Offset((query.PageNum - 1) * query.PageSize).Find(&result)
		resTotal = int64(len(result))
	}
	return result, resTotal, groupTotal
}

// CheckGroup true表示不存在 false表示不存在
func CheckGroup(data string) bool {
	var result []DevGroup
	var resTotal int64
	db.Where("group_name = ?", data).Find(&result).Count(&resTotal)
	if resTotal > 0 {
		// false表示存在
		return false
	} else {
		// true表示不存在
		return true
	}
}

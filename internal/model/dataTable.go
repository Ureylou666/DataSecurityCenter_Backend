package model

import (
	"Backend/internal/utils/Errmsg"
)

type DataTable struct {
	UUID                string `gorm:"primaryKey" json:"UUID"`
	GroupName           string `gorm:"type:varchar(50)" json:"GroupName"`
	CreationTime        string `gorm:"type:varchar(50)" json:"CreationTime,omitempty"`         // 创建该数据资产表的时间
	Id                  int    `gorm:"type:int" json:"Id,omitempty"`                           // 数据资产表的唯一标识ID
	InstanceDescription string `gorm:"type:varchar(200)" json:"InstanceDescription,omitempty"` // 实例的备注信息
	InstanceId          int    `gorm:"type:int" json:"InstanceId,omitempty"`                   // 数据资产表所属的资产实例ID
	InstanceName        string `gorm:"type:varchar(200)" json:"InstanceName,omitempty"`        // 数据资产表的实例名称
	Name                string `gorm:"type:varchar(200)" json:"Name,omitempty"`                // 数据资产表的名称
	Owner               string `gorm:"type:varchar(200)" json:"Owner,omitempty"`               // 拥有该数据资产表的阿里云账号
	ProductCode         string `gorm:"type:varchar(200)" json:"ProductCode,omitempty"`         // 数据资产表所属产品名称，取值：MaxCompute、OSS、ADS、OTS、RDS等
	ProductId           string `gorm:"type:varchar(200)" json:"ProductId,omitempty"`           // 数据资产表所属的产品ID。
	RiskLevelId         int    `gorm:"type:int" json:"RiskLevelId,omitempty"`                  // 数据资产表的风险等级ID。 每个风险等级ID都有对应的风险等级名称
	RiskLevelName       string `gorm:"type:varchar(200)" json:"RiskLevelName,omitempty"`       // 数据资产表的风险等级名称
	Sensitive           bool   `gorm:"type:boolean" json:"Sensitive,omitempty"`                // 数据资产表中是否包含敏感字段
	SensitiveCount      int    `gorm:"type:int" json:"SensitiveCount,omitempty"`               // 数据资产表中包含的敏感字段总数
	SensitiveRatio      string `gorm:"type:varchar(200)" json:"SensitiveRatio,omitempty"`      // 数据资产表中敏感字段所占的百分比
	TenantName          string `gorm:"type:varchar(200)" json:"TenantName,omitempty"`          // 租户名称
	TotalCount          int    `gorm:"type:int" json:"TotalCount,omitempty"`                   // 数据资产表包含的字段总数
}

type TabelQueryInfo struct {
	InstanceName string
	GroupName    string
	PageNum      int
	PageSize     int
}

func AddTables(data *DataTable) int {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
}

func GetTables(query TabelQueryInfo) ([]DataTable, int64, int64) {
	var result []DataTable
	var resTotal, tablesTotal int64
	// 分页处理
	db.Where("instance_name = ?", query.InstanceName).Find(&result).Count(&tablesTotal)
	if query.PageNum == 0 || query.PageSize == 0 {
		db.Where("instance_name = ?", query.InstanceName).Limit(-1).Order("risk_level_name desc, name").Find(&result)
		resTotal = int64(len(result))
	} else {
		db.Where("instance_name = ?", query.InstanceName).Limit(query.PageSize).Offset((query.PageNum - 1) * query.PageSize).Order("risk_level_name desc, name").Find(&result)
	}
	return result, resTotal, tablesTotal
}

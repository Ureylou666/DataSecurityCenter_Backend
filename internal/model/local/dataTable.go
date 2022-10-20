package local

import (
	"Backend/internal/utils/Errmsg"
)

type DataTable struct {
	UUID           string `gorm:"primaryKey" json:"UUID"`
	GroupName      string `gorm:"type:varchar(50)" json:"GroupName"`
	InstanceId     string `gorm:"type:varchar(50)" json:"InstanceId,omitempty"`     // 数据资产表所属的资产实例ID
	DatabaseName   string `gorm:"type:varchar(200)" json:"DatabaseName"`            //数据实例Database的名称。
	TableName      string `gorm:"type:varchar(200)" json:"TableName,omitempty"`     // 数据资产表的名称
	RiskLevelId    int    `gorm:"type:int" json:"RiskLevelId,omitempty"`            // 数据资产表的风险等级ID。 每个风险等级ID都有对应的风险等级名称
	RiskLevelName  string `gorm:"type:varchar(200)" json:"RiskLevelName,omitempty"` // 数据资产表的风险等级名称
	Sensitive      bool   `gorm:"type:boolean" json:"Sensitive,omitempty"`          // 数据资产表中是否包含敏感字段
	SensitiveCount int    `gorm:"type:int" json:"SensitiveCount,omitempty"`         // 数据资产表中包含的敏感字段总数
	TotalCount     int    `gorm:"type:int" json:"TotalCount,omitempty"`             // 数据资产表包含的字段总数
}

type TabelQueryInfo struct {
	InstanceID   string
	DatabaseName string
	PageNum      int
	PageSize     int
}

func InsertTableData(data *DataTable) int {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
}

func GetTableList(query TabelQueryInfo) ([]DataTable, int64, int64) {
	var result []DataTable
	var resTotal, tablesTotal int64
	// 分页处理
	db.Where("database_name = ? AND instance_id = ?", query.DatabaseName, query.InstanceID).Find(&result).Count(&tablesTotal)
	if query.PageNum == 0 || query.PageSize == 0 {
		db.Where("database_name = ? AND instance_id = ?", query.DatabaseName, query.InstanceID).Limit(-1).Find(&result)
		resTotal = int64(len(result))
	} else {
		db.Where("database_name = ? AND instance_id = ?", query.DatabaseName, query.InstanceID).Limit(query.PageSize).Offset((query.PageNum - 1) * query.PageSize).Find(&result)
	}
	return result, resTotal, tablesTotal
}

// DeleteTableData 删除RDS下Database中对应Table数据
func DeleteTableData(InstanceID string, DatabaseName string) {
	db.Where("instance_id = ? and database_name = ? ", InstanceID, DatabaseName).Delete(&DataTable{})
}

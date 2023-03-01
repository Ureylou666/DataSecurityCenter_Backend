package local

import (
	"Backend/internal/utils/Errmsg"
)

type MappingRules struct {
	UUID       string `gorm:"type:varchar(50)" json:"uuid"`
	ColumnUUID string `gorm:"type:varchar(50)" json:"ColumnUUID"`
	ColumnName string `gorm:"type:varchar(100)" json:"ColumnName"`
	RuleUUID   string `gorm:"type:varchar(50)" json:"RuleUUID"`
	RuleName   string `gorm:"type:varchar(100)" json:"RuleName"`
	Sample     string `gorm:"type:text"`
}

// InsertMappingRules 增
func InsertMappingRules(data *MappingRules) int {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
}

// DeleteMappingRules 删
func DeleteMappingRules(RuleUUID string) (ErrCode int, ErrMessage error) {
	ErrMessage = db.Where("Rule_uuid = ?", RuleUUID).Delete(MappingRules{}).Error
	if ErrMessage != nil {
		return Errmsg.ERROR, ErrMessage
	}
	return Errmsg.SUCCESS, nil
}

// 改

// 查

package local

import (
	"Backend/internal/utils/Errmsg"
)

type Rules struct {
	UUID           string `gorm:"type:varchar(50)" json:"UUID"`
	RuleName       string `gorm:"type:varchar(50)" json:"RuleName"`
	PreferColumn   string `gorm:"type:varchar(50)" json:"PreferColumn"`
	RegExpress     string `gorm:"type:varchar(255)" json:"RegExpress"`
	CategoryUUID   string `gorm:"type:varchar(50)" json:"CategoryUUID"`
	CategoryName   string `gorm:"type:varchar(50)" json:"CategoryName"`
	SensitiveGrade string `gorm:"type:varchar(10)" json:"SensitiveGrade"`
	Status         bool   `gorm:"type:boolean"`
	Comments       string `gorm:"type:text" json:"Comments"`
}

type RulesListQuery struct {
	RuleName       string
	CategoryName   string
	SensitiveGrade string
	PageNum        int
	PageSize       int
}

// CreateRule 增 新增规则
func CreateRule(data *Rules) (ErrCode int, ErrMessage error) {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR, err
	}
	return Errmsg.SUCCESS, nil
}

// DeleteRule 删 删除规则
func DeleteRule(uuid string) (ErrCode int, ErrMessage error) {
	ErrMessage = db.Where("uuid = ?", uuid).Delete(&Rules{}).Error
	if ErrMessage != nil {
		return Errmsg.ERROR, ErrMessage
	}
	return Errmsg.SUCCESS, nil
}

// UpdateRule 改 更新规则
func UpdateRule(data *Rules) (ErrCode int, ErrMessage error) {
	ErrMessage = db.Where("rule_uuid = ?", data.UUID).Update("status", true).Error
	if ErrMessage != nil {
		return Errmsg.ERROR, ErrMessage
	}
	return Errmsg.SUCCESS, nil
}

// 查
func CheckRuleNameUnique(inputName string) bool {
	var temp []Rules
	db.Where("rule_name = ?", inputName).Find(&temp)
	if len(temp) > 0 {
		return false
	}
	return true
}

func CheckRuleExist(uuid string) bool {
	var temp []Rules
	db.Where("uuid = ?", uuid).Find(&temp)
	if len(temp) > 0 {
		return true
	}
	return false
}

// ListRules 列出所有规则
func ListRules(input RulesListQuery) ([]Rules, int64, int64) {
	var result []Rules
	var resTotal, rulesTotal int64
	input.SensitiveGrade = "%" + input.SensitiveGrade + "%"
	input.RuleName = "%" + input.RuleName + "%"
	input.CategoryName = "%" + input.CategoryName + "%"
	// 分页处理
	db.Table("(?) as Y", db.Table("(?) as X", db.Model(Rules{}).Where("rule_name like ? ", input.RuleName)).Where("sensitive_grade like ?", input.SensitiveGrade)).Where("category_name like ?", input.CategoryName).Find(&result).Count(&rulesTotal)
	if input.PageNum == 0 || input.PageSize == 0 {
		db.Table("(?) as Y", db.Table("(?) as X", db.Model(Rules{}).Where("rule_name like ? ", input.RuleName)).Where("sensitive_grade like ?", input.SensitiveGrade)).Where("category_name like ? ", input.CategoryName).Limit(-1).Order("sensitive_grade").Find(&result)
		resTotal = int64(len(result))
	} else {
		db.Table("(?) as Y", db.Table("(?) as X", db.Model(Rules{}).Where("rule_name like ?", input.RuleName)).Where("sensitive_grade like ?", input.SensitiveGrade)).Where("category_name like ? ", input.CategoryName).Order("sensitive_grade").Limit(input.PageSize).Offset((input.PageNum - 1) * input.PageSize).Find(&result)
		resTotal = int64(len(result))
	}
	return result, resTotal, rulesTotal
}

func GetRule(input string) (result Rules) {
	db.Where("uuid = ?", input).Limit(1).Find(&result)
	return result
}

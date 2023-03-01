package Basic

import (
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"strings"
)

// validationRules 验证输入
func validationRules(input local.Rules) (validation bool) {
	// 输入校验
	validation = true
	// 检查必填字段是否为空
	if (input.RuleName == "") || (input.PreferColumn == "") || (input.CategoryUUID == "") || (input.SensitiveGrade == "") {
		validation = false
	}
	// RulesName 仅大小写英文字符
	if validation {
		validation, _ = regexp.MatchString("^[A-Za-z]+$", input.RuleName)
	}
	// RuleNames 是否重复
	if validation {
		validation = local.CheckRuleNameUnique(input.RuleName)
	}
	// PreferColumn 仅大小写英文字符 拆分
	if validation {
		n := 1
		for i := 0; i < len(input.PreferColumn); i++ {
			if string(input.PreferColumn[i]) == "," {
				n++
			}
		}
		inputColumn := make([]string, n)
		// 拆分preferColumn¬
		n = 0
		for i := 0; i < len(input.PreferColumn); i++ {
			if string(input.PreferColumn[i]) == " " {
				i++
			}
			if string(input.PreferColumn[i]) != "," {
				inputColumn[n] += string(input.PreferColumn[i])
			} else {
				n++
			}
		}
		/*
			// 对每个Column进行匹配
			for i := 0; i < len(inputColumn); i++ {
				matched, _ := regexp.MatchString("^[A-Za-z]+$", inputColumn[i])
				if validation && (matched == false) {
					validation = false
				}
			}
		*/
	}
	// Category 是否符合规范
	if validation {
		validation, _ = regexp.MatchString("[a-f\\d]{8}(-[a-f\\d]{4}){3}-[a-f\\d]{12}", input.CategoryUUID)
	}
	// Category 是否存在
	if validation {
		validation = local.CheckCategoryExist(input.CategoryUUID)
	}
	// SensitiveGrade 只能为S1/S2/S3/S4
	if validation {
		if (input.SensitiveGrade != "S1") && (input.SensitiveGrade != "S2") && (input.SensitiveGrade != "S3") && (input.SensitiveGrade != "S4") {
			validation = false
		}
	}
	// Comments 长度<255
	if validation {
		if len(input.Comments) > 255 {
			validation = false
		}
	}
	return validation
}

// CreateRules 创建规则
func CreateRules(input local.Rules) (ErrCode int, ErrMessage error) {
	// 先进行输入校验
	if !validationRules(input) {
		return Errmsg.ErrorQueryInput, nil
	}
	if input.UUID == "" {
		input.UUID = uuid.New().String()
	}
	input.CategoryName = local.CategoryUUIDtoName(input.CategoryUUID)
	input.Status = false
	ErrCode, ErrMessage = local.CreateRule(&input)
	//ErrCode, ErrMessage = EnforceRules(input)
	return ErrCode, ErrMessage
}

// EnforceRules 实施规则
func EnforceRules(input string) (ErrCode int, ErrMessage error) {
	var CheckColumns []local.DataColumn
	var query local.ColumnDetailsQueryInfo
	//var sampleData [10]string
	var inputMapping local.MappingRules
	var Sample string
	var rules local.Rules
	rules = local.GetRule(input)
	// 清空mappingRules
	ErrCode, ErrMessage = local.DeleteMappingRules(rules.UUID)
	// 整理preferred columnName
	n := 1
	for i := 0; i < len(rules.PreferColumn); i++ {
		if string(rules.PreferColumn[i]) == "," {
			n++
		}
	}
	preferredColumn := make([]string, n)
	n = 0
	for i := 0; i < len(rules.PreferColumn); i++ {
		if string(rules.PreferColumn[i]) == " " {
			i++
		}
		if string(rules.PreferColumn[i]) != "," {
			preferredColumn[n] += string(rules.PreferColumn[i])
		} else {
			n++
		}
	}
	// 获取检测 column
	Total := local.CountColumns()
	// 每100个为一组 从数据库中获取
	query.PageSize = 100
	for i := 1; i < Total/100+1; i++ {
		query.PageNum = i
		CheckColumns, _, _ = local.GetColumnDetails(query)
		// 对每一个column进行判断
		for j := 0; j < len(CheckColumns); j++ {
			score := 0
			// 判断columnName是否包含 preferred Column 满足+50
			for k := 0; k < len(preferredColumn); k++ {
				if strings.Contains(CheckColumns[j].ColumnName, preferredColumn[k]) {
					score = score + 50
				}
			}
			matched, _ := regexp.MatchString(rules.RegExpress, CheckColumns[j].SampleData)
			if matched {
				score = score + 10
				Sample = CheckColumns[j].SampleData
			}
			/*
				// 获取对应 sampleData 通过正则匹配是否满足 满足+10
				sampleData = local.GetSingleSampleData(CheckColumns[j].UUID)
				for k := 0; k < 10; k++ {
					matched, _ := regexp.MatchString(rules.RegExpress, sampleData[k])
					if matched {
						score = score + 10
						Sample = sampleData[k]
					}
				}
			*/
			// 大于60判定为真
			if score >= 60 {
				inputMapping.UUID = uuid.New().String()
				inputMapping.RuleUUID = rules.UUID
				inputMapping.RuleName = rules.RuleName
				inputMapping.ColumnUUID = CheckColumns[j].UUID
				inputMapping.ColumnName = CheckColumns[j].ColumnName
				inputMapping.Sample = Sample
				local.InsertMappingRules(&inputMapping)
				local.MappingUpdate(CheckColumns[j].UUID, rules.RuleName, rules.CategoryName, rules.SensitiveGrade)
				fmt.Println(CheckColumns[j].UUID, " ", CheckColumns[j].ColumnName, ":", score)
			}
		}
	}
	rules.Status = true
	ErrCode, ErrMessage = local.UpdateRule(&rules)
	return Errmsg.SUCCESS, nil
}

// UpdateRules 更新规则 含创建规则
func UpdateRules(input local.Rules) (ErrCode int, ErrMessage error) {
	// 先进行输入校验
	if !validationRules(input) {
		return Errmsg.ErrorQueryInput, nil
	}
	// 检测uuid是否存在 删除后新增
	if local.CheckRuleExist(input.UUID) {
		// 删除 原有
		ErrCode, ErrMessage = DeleteRules(input.UUID)
		// 新增
		ErrCode, ErrMessage = CreateRules(input)
	} else {
		return Errmsg.ErrorQueryInput, nil
	}
	return Errmsg.SUCCESS, nil
}

// DeleteRules 删除规则及mapping关系
func DeleteRules(uuid string) (ErrCode int, ErrMessage error) {
	// 简单校验 uuid
	validation := true
	validation, _ = regexp.MatchString("[a-f\\d]{8}(-[a-f\\d]{4}){3}-[a-f\\d]{12}", uuid)
	// UUID是否存在
	if validation {
		validation = local.CheckRuleExist(uuid)
	}
	if validation {
		// 删除规则
		ErrCode, ErrMessage = local.DeleteRule(uuid)
		// 删除mapping
		ErrCode, ErrMessage = local.DeleteMappingRules(uuid)

	} else {
		return Errmsg.ErrorQueryInput, nil
	}
	return ErrCode, ErrMessage
}

func ListRules(input local.RulesListQuery) (ErrCode int, ErrMessage string, result []local.Rules, resTotal int64, rulesTotal int64) {
	// 校验前端输入
	validation := true
	// 校验RuleName 过滤非字母数字
	if input.RuleName != "" {
		validation, _ = regexp.MatchString("^[A-Za-z0-9]+$", input.RuleName)
	}
	// 校验Category 过滤非字母数字
	if validation {
		if input.CategoryName != "" {
			validation, _ = regexp.MatchString("^[A-Za-z0-9]+$", input.CategoryName)
		}
	}
	// 校验SensitiveGrade 仅能为S4 S3 S2 S1
	if validation {
		if (input.SensitiveGrade == "") || (input.SensitiveGrade == "S1") || (input.SensitiveGrade == "S2") || (input.SensitiveGrade == "S3") || (input.SensitiveGrade == "S4") {
			validation = true
		} else {
			validation = false
		}
	}
	if !validation {
		return Errmsg.ErrorQueryInput, Errmsg.GetErrMsg(Errmsg.ErrorQueryInput), nil, 0, 0
	}
	// 调用数据库
	result, resTotal, rulesTotal = local.ListRules(input)
	return Errmsg.SUCCESS, "", result, resTotal, rulesTotal
}

// GetRule 获取特定的rules
func GetRule(RuleUUID string) (ErrCode int, ErrMessage string, result local.Rules) {
	// 校验前端输入
	validation := true
	// 校验RuleName 过滤非字母数字 且不能为空
	validation, _ = regexp.MatchString("[a-f\\d]{8}(-[a-f\\d]{4}){3}-[a-f\\d]{12}", RuleUUID)
	if !validation {
		return Errmsg.ErrorQueryInput, Errmsg.GetErrMsg(Errmsg.ErrorQueryInput), result
	}
	result = local.GetRule(RuleUUID)
	if result.UUID == "" {
		return Errmsg.ErrorQueryInput, Errmsg.GetErrMsg(Errmsg.ErrorQueryInput), result
	}
	return Errmsg.SUCCESS, "", result
}

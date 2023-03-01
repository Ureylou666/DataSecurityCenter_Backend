package Basic

import (
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"github.com/google/uuid"
	"regexp"
)

func validationCategory(input local.Category) (validation bool) {
	// 输入校验
	validation = true
	// 空字符串
	if input.CategoryName == "" {
		validation = false
	}
	//  CategoryName 小于100字符
	if validation {
		if len(input.CategoryName) > 100 {
			validation = false
		}
	}
	// CategoryName 仅大小写英文字符
	if validation {
		validation, _ = regexp.MatchString("^[A-Za-z]+$", input.CategoryName)
	}
	// category Name 是否唯一
	if validation {
		validation = local.CheckCategoryNameUnique(input.CategoryName)
	}
	// category name 和 level 是否匹配
	if validation {
		validation = local.CheckCategoryToLevel(input.ParentCategory, input.Level)
	}
	// ListLevel 大于数据库中最大level+1
	if validation && (input.Level != 0) {
		if input.Level > local.ListLevel()+1 {
			validation = false
		}
	}
	// ParentCategory 是为uuid
	if validation && (input.ParentCategory != "") {
		validation, _ = regexp.MatchString("[a-f\\d]{8}(-[a-f\\d]{4}){3}-[a-f\\d]{12}", input.ParentCategory)
	}
	// ParentCategory 是否存在
	if validation && (input.ParentCategory != "") {
		validation = local.CheckCategoryExist(input.UUID)
	}
	// Comments 小于 500 字符
	if validation {
		if len(input.Comments) > 500 {
			validation = false
		}
	}
	return validation
}

// CreateCategory 新增
func CreateCategory(input local.Category) (ErrCode int, ErrMessage error) {
	// 验证输入
	if !validationCategory(input) {
		return Errmsg.ErrorQueryInput, nil
	}
	if input.UUID == "" {
		input.UUID = uuid.New().String()
	}
	// 创建
	ErrCode, ErrMessage = local.CreateCategory(&input)
	return ErrCode, ErrMessage
}

// UpdateCategory 更新
func UpdateCategory(input local.Category) (ErrCode int, ErrMessage error) {
	// 验证输入
	if !validationCategory(input) {
		return Errmsg.ErrorQueryInput, nil
	}
	// 新增uuid是否存在
	if !local.CheckCategoryExist(input.UUID) {
		return Errmsg.ErrorQueryInput, nil
	} else {
		// 删除 原有的
		ErrCode, ErrMessage = DeleteCategory(input)
		// 新增
		ErrCode, ErrMessage = CreateCategory(input)
	}
	return Errmsg.SUCCESS, nil
}

// DeleteCategory 删除分类
func DeleteCategory(input local.Category) (ErrCode int, ErrMessage error) {
	// 验证输入
	if !validationCategory(input) {
		return Errmsg.ErrorQueryInput, nil
	}
	ErrCode, ErrMessage = local.DeleteCategory(input.UUID)
	return Errmsg.SUCCESS, nil
}

// ListCategory 列出分类
func ListCategory(input local.CategoryQuery) (ErrCode int, ErrMessage string, result []local.Category, resTotal int64, categoryTotal int64) {
	validation := true
	if input.CategoryName != "" {
		validation, _ = regexp.MatchString("^[A-Za-z]+$", input.CategoryName)
	}
	if !validation {
		return Errmsg.ErrorQueryInput, Errmsg.GetErrMsg(Errmsg.ErrorQueryInput), nil, 0, 0
	}
	result, resTotal, categoryTotal = local.ListCategory(input)
	return Errmsg.SUCCESS, "", result, resTotal, categoryTotal
}

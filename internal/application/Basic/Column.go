package Basic

import (
	"Backend/internal/model/local"
	"regexp"
)

func GetColumnDetails(input local.ColumnDetailsQueryInfo) (resData []local.DataColumn, resTotal int64, columnTotal int64) {
	validation := true
	if input.PageSize > 100 {
		validation = false
	}
	if validation && input.GroupName != "" {
		validation, _ = regexp.MatchString("^[A-Za-z0-9]+$", input.GroupName)
	}
	if validation && input.CategoryName != "" {
		validation, _ = regexp.MatchString("^[A-Za-z0-9]+$", input.CategoryName)
	}
	if validation && input.RuleName != "" {
		validation, _ = regexp.MatchString("^[A-Za-z0-9]+$", input.RuleName)
	}
	if !validation {
		return nil, 0, 0
	}
	return local.GetColumnDetails(input)
}

func GetColumnsFromTable(input local.QueryFromTable) (resData []local.DataColumn, resTotal int64, columnTotal int64) {
	validation := true
	if (input.PageSize > 50) || (input.TableName == "") {
		validation = false
	}
	if validation {
		validation, _ = regexp.MatchString("^[A-Za-z0-9-]+$", input.InstanceName)
	}
	if !validation {
		return nil, 0, 0
	}
	return local.GetColumnFromTable(input)
}

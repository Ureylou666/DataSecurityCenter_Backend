package Basic

import (
	"Backend/internal/model/local"
	"regexp"
)

func GetTableList(input local.TabelQueryInfo) (resData []local.DataTable, resTotal int64, TableListTotal int64) {
	validation := true
	if input.PageSize > 50 || input.DatabaseName == "" {
		validation = false
	}
	if validation {
		validation, _ = regexp.MatchString("^[0-9a-zA-Z_-]{20}$", input.InstanceID)
	}
	if !validation {
		return nil, 0, 0
	}
	return local.GetTableList(input)
}

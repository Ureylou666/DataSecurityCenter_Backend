package Basic

import (
	"Backend/internal/model/local"
	"regexp"
)

func GetDatabase(input local.DatabaseQueryInfo) (result []local.DataDatabase, resTotal int64, DatabaseTotal int64) {
	// 输入校验
	validation, _ := regexp.MatchString("^[0-9a-zA-Z_-]{20}$", input.InstanceID)
	if validation {
		if (input.PageSize > 50) || (input.InstanceID == "") {
			validation = false
		}
	}
	if !validation {
		return nil, 0, 0
	}
	return local.GetDatabaseAPI(input)
}

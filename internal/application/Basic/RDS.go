package Basic

import (
	"Backend/internal/model/local"
	"regexp"
)

func GeTRDSInventoryList(input local.RDSQueryInfo) (resData []local.DataInventory, resTotal int64, rdsInventoryTotal int64) {
	validation := true
	if input.PageSize > 50 || input.AccountID == "" {
		validation = false
	}
	if validation && input.AccountID != "All" {
		validation, _ = regexp.MatchString("^[0-9]+$", input.AccountID)
	}
	if !validation {
		return nil, 0, 0
	}
	if input.AccountID == "All" {
		input.AccountID = ""
	}
	return local.GetRDSInventory(input)
}

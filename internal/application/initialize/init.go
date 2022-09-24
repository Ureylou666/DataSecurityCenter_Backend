package initialize

import (
	"Backend/internal/application/aliyunSDK"
	"Backend/internal/model"
	"github.com/google/uuid"
)

/*
初始化
- 初始化cloud账号列表
- 初始化rds列表
- 初始化
*/

// InitCloudAccountList 初始化cloud账号
func InitCloudAccountList() {
	var input model.CloudAccount
	rawlistaccounts := aliyunSDK.ListAliyunCloudAccounts()
	for i := 0; i < len(rawlistaccounts); i++ {
		input.UUID = uuid.New().String()
		if rawlistaccounts[i].Type != nil {
			input.Type = *rawlistaccounts[i].Type
		}
		if rawlistaccounts[i].DisplayName != nil {
			input.DisplayName = *rawlistaccounts[i].DisplayName
		}
		if rawlistaccounts[i].JoinTime != nil {
			input.JoinTime = *rawlistaccounts[i].JoinTime
		}
		if rawlistaccounts[i].AccountId != nil {
			input.AccountId = *rawlistaccounts[i].AccountId
		}
		if rawlistaccounts[i].ModifyTime != nil {
			input.ModifyTime = *rawlistaccounts[i].ModifyTime
		}
		model.AddCloudAccount(&input)
	}
}

func InitRDSInstanceList() {

}

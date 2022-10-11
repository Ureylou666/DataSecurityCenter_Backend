package Basic

import (
	"Backend/internal/application/aliyunSDK"
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"github.com/google/uuid"
)

// UpdateCloudAccountList 初始化cloud账号
func UpdateCloudAccountList() (ErrCode int, ErrMessage error) {
	var input local.CloudAccount
	previousCode, previousMsg, rawlistaccounts := aliyunSDK.ListAliyunCloudAccounts()
	// 调用aliyun sdk 进行错误控制
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	// 数据入库
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
		ErrCode, ErrMessage = local.AddCloudAccount(&input)
		// 对于插入数据 进行错误控制
		if ErrCode != Errmsg.SUCCESS {
			return Errmsg.ERROR, ErrMessage
		}
	}
	return Errmsg.SUCCESS, nil
}

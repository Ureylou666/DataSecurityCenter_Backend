package aliyunSDK

/*
aliyun ResourceManager 相关sdk接口
- ListAliyunCloudAccounts
*/

import (
	"Backend/internal/utils/Errmsg"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	resourcemanager20200331 "github.com/alibabacloud-go/resourcemanager-20200331/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

func createRMClient() (ErrCode int, ErrMessage error, client *resourcemanager20200331.Client) {
	// 获取临时STS token
	previousCode, previousMsg, AccessKeyId, AccessKeySecret, SecurityToken := AssumeRole("1174006592814680")
	// 判断是否获取成功
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg, nil
	}
	// 成功获取临时AK SK
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: tea.String(AccessKeyId),
		// 您的 AccessKey Secret
		AccessKeySecret: tea.String(AccessKeySecret),
		// 您的 AccessToken
		SecurityToken: tea.String(SecurityToken),
	}
	config.Endpoint = tea.String("resourcemanager.aliyuncs.com")
	client = &resourcemanager20200331.Client{}
	client, ErrMessage = resourcemanager20200331.NewClient(config)
	// 错误判断
	if ErrMessage != nil {
		return Errmsg.ErrorCreateRMClient, ErrMessage, nil
	} else {
		return Errmsg.SUCCESS, nil, client
	}
}

// ListAliyunCloudAccounts 查看整个资源目录下的所有成员的信息
func ListAliyunCloudAccounts() (ErrCode int, ErrMessage error, ListAccountsResponseBodyAccountsAccount []resourcemanager20200331.ListAccountsResponseBodyAccountsAccount) {
	previousCode, previousMsg, client := createRMClient()
	// 判断客户端创建是否成功
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg, nil
	}
	listAccountsRequest := &resourcemanager20200331.ListAccountsRequest{
		PageNumber: tea.Int32(1),
		PageSize:   tea.Int32(100),
	}
	runtime := &util.RuntimeOptions{}
	listAccountsResponse, ErrMessage := client.ListAccountsWithOptions(listAccountsRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorListAliCloudAccounts, ErrMessage, nil
	}
	accountlist := make([]resourcemanager20200331.ListAccountsResponseBodyAccountsAccount, int(*listAccountsResponse.Body.TotalCount))
	for i := 0; i < len(listAccountsResponse.Body.Accounts.Account); i++ {
		accountlist[i] = *listAccountsResponse.Body.Accounts.Account[i]
	}
	// 当前list > 100 这部分代码待完善 目前<200 是否做循环意义不是很大
	listAccountsRequest = &resourcemanager20200331.ListAccountsRequest{
		PageNumber: tea.Int32(2),
		PageSize:   tea.Int32(100),
	}
	n := 100
	listAccountsResponse, ErrMessage = client.ListAccountsWithOptions(listAccountsRequest, runtime)
	if ErrMessage != nil {
		return Errmsg.ErrorListAliCloudAccounts, ErrMessage, nil
	}
	for i := 0; i < len(listAccountsResponse.Body.Accounts.Account); i++ {
		accountlist[n] = *listAccountsResponse.Body.Accounts.Account[i]
		n++
	}
	return Errmsg.SUCCESS, nil, accountlist
}

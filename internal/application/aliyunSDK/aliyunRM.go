package aliyunSDK

/*
aliyun ResourceManager 相关sdk接口
- ListAliyunCloudAccounts
*/

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	resourcemanager20200331 "github.com/alibabacloud-go/resourcemanager-20200331/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

func createRMClient() (_result *resourcemanager20200331.Client, _err error) {
	AccessKeyId, AccessKeySecret, SecurityToken := AssumeRole("1174006592814680")
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: tea.String(AccessKeyId),
		// 您的 AccessKey Secret
		AccessKeySecret: tea.String(AccessKeySecret),
		// 您的 AccessToken
		SecurityToken: tea.String(SecurityToken),
	}
	// 访问的域名
	config.Endpoint = tea.String("resourcemanager.aliyuncs.com")
	_result = &resourcemanager20200331.Client{}
	_result, _err = resourcemanager20200331.NewClient(config)
	return _result, _err
}

// ListAliyunCloudAccounts 查看整个资源目录下的所有成员的信息
func ListAliyunCloudAccounts() []resourcemanager20200331.ListAccountsResponseBodyAccountsAccount {
	client, _ := createRMClient()
	listAccountsRequest := &resourcemanager20200331.ListAccountsRequest{
		PageNumber: tea.Int32(1),
		PageSize:   tea.Int32(100),
	}
	runtime := &util.RuntimeOptions{}
	listAccountsResponse, _err := client.ListAccountsWithOptions(listAccountsRequest, runtime)
	accountlist := make([]resourcemanager20200331.ListAccountsResponseBodyAccountsAccount, int(*listAccountsResponse.Body.TotalCount))
	for i := 0; i < len(listAccountsResponse.Body.Accounts.Account); i++ {
		accountlist[i] = *listAccountsResponse.Body.Accounts.Account[i]
	}
	// 当前list > 100 这部分代码待完善
	listAccountsRequest = &resourcemanager20200331.ListAccountsRequest{
		PageNumber: tea.Int32(2),
		PageSize:   tea.Int32(100),
	}
	n := 100
	listAccountsResponse, _err = client.ListAccountsWithOptions(listAccountsRequest, runtime)
	for i := 0; i < len(listAccountsResponse.Body.Accounts.Account); i++ {
		accountlist[n] = *listAccountsResponse.Body.Accounts.Account[i]
		n++
	}
	fmt.Println(_err)
	return accountlist
}

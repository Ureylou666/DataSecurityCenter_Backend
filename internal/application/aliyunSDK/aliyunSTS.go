package aliyunSDK

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/viper"
)

/*
aliyun STS 相关sdk接口

*/

func createSTSClient() (_result *sts20150401.Client, _err error) {
	accessKeyId := viper.GetString("CNISDP.AccessKey")
	accessKeySecret := viper.GetString("CNISDP.AccessSecret")
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: tea.String(accessKeyId),
		// 您的 AccessKey Secret
		AccessKeySecret: tea.String(accessKeySecret),
	}
	// 访问的域名
	config.Endpoint = tea.String("sts.cn-shanghai.aliyuncs.com")
	_result = &sts20150401.Client{}
	_result, _err = sts20150401.NewClient(config)
	return _result, _err
}

func AssumeRole(accountId string) (string, string, string) {
	client, _ := createSTSClient()
	assumeRoleRequest := &sts20150401.AssumeRoleRequest{
		DurationSeconds: tea.Int64(900),
		RoleArn:         tea.String("acs:ram::" + accountId + ":role/cmsadmin"),
		RoleSessionName: tea.String("CNISDP-Audit999"),
	}
	runtime := &util.RuntimeOptions{}
	assumeRoleResponse, _ := client.AssumeRoleWithOptions(assumeRoleRequest, runtime)
	return *assumeRoleResponse.Body.Credentials.AccessKeyId, *assumeRoleResponse.Body.Credentials.AccessKeySecret, *assumeRoleResponse.Body.Credentials.SecurityToken
}

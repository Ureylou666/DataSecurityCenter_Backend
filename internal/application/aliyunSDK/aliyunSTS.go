package aliyunSDK

import (
	"Backend/internal/utils/Errmsg"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/viper"
)

/*
aliyun STS 相关sdk接口

*/

func createSTSClient() (ErrCode int, ErrMessage error, client *sts20150401.Client) {
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
	client = &sts20150401.Client{}
	client, ErrMessage = sts20150401.NewClient(config)
	if ErrMessage != nil {
		return Errmsg.ErrorCreateSTSClient, ErrMessage, nil
	} else {
		return Errmsg.SUCCESS, ErrMessage, client
	}
}

func AssumeRole(accountId string) (ErrCode int, ErrMessage error, AccessKeyId string, AccessKeySecret string, SecurityToken string) {
	previousCode, previousErrMsg, client := createSTSClient()
	if previousCode != Errmsg.SUCCESS {
		// 创建STS client失败
		return previousCode, previousErrMsg, "", "", ""
	} else {
		// 创建STS client成功
		assumeRoleRequest := &sts20150401.AssumeRoleRequest{
			DurationSeconds: tea.Int64(900),
			RoleArn:         tea.String("acs:ram::" + accountId + ":role/cmsadmin"),
			RoleSessionName: tea.String("CNISDP-Audit"),
		}
		runtime := &util.RuntimeOptions{}
		assumeRoleResponse, ErrMessage := client.AssumeRoleWithOptions(assumeRoleRequest, runtime)
		if ErrMessage != nil {
			fmt.Println(accountId)
			return Errmsg.ErrorAssumeRole, ErrMessage, "", "", ""
		} else {
			return Errmsg.SUCCESS, nil, *assumeRoleResponse.Body.Credentials.AccessKeyId, *assumeRoleResponse.Body.Credentials.AccessKeySecret, *assumeRoleResponse.Body.Credentials.SecurityToken

		}
	}
}

// This file is auto-generated, don't edit it. Thanks.
package aliyunSDDP

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	sddp20190103 "github.com/alibabacloud-go/sddp-20190103/client"
	"github.com/alibabacloud-go/tea/tea"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */

// Group 包括 ISDP DAP FMC APP
func InitData(Group string) int {
	// 清空数据库中该项目组数据
	//model.RestoreData(Group)
	// 获取key值
	/*
			AccessKey := viper.GetString(Group + ".AccessKey")
			AccessSecret := viper.GetString(Group + ".AccessSecret")
			//打开client
			client, _err := CreateClient(tea.String(AccessKey), tea.String(AccessSecret))
			initDescribeInstances(Group, client)      // 获取SDDP连接授权的MaxCompute、RDS、OSS数据资产实例列表
			initDescribeUnauthInstance(Group, client) // 获取未授权资产列表
			initDescribeTables(Group, client)         //获取表信息
			initDescribeColumns(Group, client)        //获取列信息
			if Group == "ISDP" {
				initdescribeRules(Group, client)
			} // 获取规则信息


		if _err != nil {
			return 500
		} else {
			return 200
		}
	*/
	return 200
}

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *sddp20190103.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("sddp.cn-zhangjiakou.aliyuncs.com")
	_result = &sddp20190103.Client{}
	_result, _err = sddp20190103.NewClient(config)
	return _result, _err
}

//获取数据资产 实例 列表
/*
// 获取SDDP连接授权的MaxCompute、RDS、OSS数据资产实例列表
func initDescribeInstances(Group string, client *sddp20190103.Client) {
	var input model.DataInventory // input 表示录入到数据库中的值
	describeInstancesRequest := &sddp20190103.DescribeInstancesRequest{
		PageSize:    tea.Int32(1000),
		CurrentPage: tea.Int32(1),
	}
	Auth_res, _err := client.DescribeInstances(describeInstancesRequest)
	if _err != nil {
		return
	}
	// rawdata指aliyunSDK返回的原始数据，需要将其简单处理入库
	rawdata := Auth_res.Body.Items
	for i := 0; i < len(rawdata); i++ {
		input.UUID = uuid.New().String()
		input.GroupName = Group
		if rawdata[i].CreationTime != nil {
			input.CreationTime = time.Unix(*rawdata[i].CreationTime/1000, 0).Format("2006-01-02 15:04:05")
		} else {
			input.CreationTime = ""
		}
		if rawdata[i].DepartName != nil {
			input.DepartName = *rawdata[i].DepartName
		} else {
			input.DepartName = ""
		}
		if rawdata[i].Id != nil {
			input.InstanceId = int(*rawdata[i].Id)
		} else {
			input.InstanceId = 0
		}
		if rawdata[i].InstanceDescription != nil {
			input.InstanceDescription = *rawdata[i].InstanceDescription
		} else {
			input.InstanceDescription = ""
		}
		if rawdata[i].Labelsec != nil {
			input.Labelsec = *rawdata[i].Labelsec
		}
		if rawdata[i].LastFinishTime != nil {
			input.LastFinishTime = time.Unix(*rawdata[i].LastFinishTime/1000, 0).Format("2006-01-02 15:04:05")
		} else {
			input.LastFinishTime = ""
		}
		if rawdata[i].Name != nil {
			input.Name = *rawdata[i].Name
		} else {
			input.Name = ""
		}
		if rawdata[i].Owner != nil {
			input.Owner = *rawdata[i].Owner
		} else {
			input.Owner = ""
		}
		if rawdata[i].ProductCode != nil {
			input.ProductCode = *rawdata[i].ProductCode
		} else {
			input.ProductCode = ""
		}
		if rawdata[i].ProductId != nil {
			input.ProductId = *rawdata[i].ProductId
		} else {
			input.ProductId = ""
		}
		input.Protection = true
		if rawdata[i].RiskLevelId != nil {
			input.RiskLevelId = int(*rawdata[i].RiskLevelId)
		} else {
			input.RiskLevelId = 0
		}
		if rawdata[i].RiskLevelName != nil {
			input.RiskLevelName = *rawdata[i].RiskLevelName
		} else {
			input.RiskLevelName = ""
		}
		if rawdata[i].RuleName != nil {
			input.RuleName = *rawdata[i].RuleName
		} else {
			input.RuleName = ""
		}
		if rawdata[i].Sensitive != nil {
			input.Sensitive = *rawdata[i].Sensitive
		}
		if rawdata[i].SensitiveCount != nil {
			input.SensitiveCount = int(*rawdata[i].SensitiveCount)
		} else {
			input.SensitiveCount = 0
		}
		if rawdata[i].TenantName != nil {
			input.TenantName = *rawdata[i].TenantName
		} else {
			input.TenantName = ""
		}
		if rawdata[i].TotalCount != nil {
			input.TotalCount = int(*rawdata[i].TotalCount)
		} else {
			input.TotalCount = 0
		}
		model.AddInventory(&input)
	}
}

// 获取未授权资产列表
func initDescribeUnauthInstance(Group string, client *sddp20190103.Client) {
	var input model.DataInventory
	// 获取SDDP未授权的
	describeInstanceSourcesRequest := &sddp20190103.DescribeInstanceSourcesRequest{
		PageSize:    tea.Int32(1000),
		CurrentPage: tea.Int32(1),
	}
	runtime := &util.RuntimeOptions{}
	// 复制代码运行请自行打印 API 的返回值
	Unauth_res, _err := client.DescribeInstanceSourcesWithOptions(describeInstanceSourcesRequest, runtime)
	if _err != nil {
		return
	}
	// tempdata获取返回的原始数据，将其整理入库
	tempdata := Unauth_res.Body.Items
	for i := 0; i < len(tempdata); i++ {
		input.UUID = uuid.New().String()
		input.GroupName = Group
		input.CreationTime = time.Unix(*tempdata[i].GmtCreate/1000, 0).Format("2006-01-02 15:04:05")
		if tempdata[i].Id != nil {
			input.InstanceId = int(*tempdata[i].Id)
		}
		if tempdata[i].InstanceDescription != nil {
			input.InstanceDescription = *tempdata[i].InstanceDescription
		} else {
			input.InstanceDescription = ""
		}
		if tempdata[i].InstanceId != nil {
			input.Name = *tempdata[i].InstanceId
		}
		if tempdata[i].UserName != nil {
			input.Owner = *tempdata[i].UserName
		}
		input.ProductCode = *tempdata[i].EngineType
		input.Protection = false
		model.AddInventory(&input)
	}
}

// 查询数据安全中心连接授权的MaxCompute、RDS等数据资产表中列的数据
func initDescribeColumns(Group string, client *sddp20190103.Client) {
	var input model.DataColumn
	// 由于ali sdk限制 每次最多返回1000条数据，故需要通过分次进行数据录入
	describeColumnsRequest := &sddp20190103.DescribeColumnsRequest{
		PageSize:    tea.Int32(0),
		CurrentPage: tea.Int32(0),
	}
	runtime := &util.RuntimeOptions{}
	response, err := client.DescribeColumnsWithOptions(describeColumnsRequest, runtime)
	if err != nil {
		return
	}
	TotalColumn := int(*response.Body.TotalCount)
	fmt.Println(TotalColumn)
	// 每次请求1000条数据，ali QPS限制 每次sleep 0.5秒
	for n := 1; n <= (TotalColumn/1000 + 1); n++ {
		describeColumnsRequest := &sddp20190103.DescribeColumnsRequest{
			PageSize:    tea.Int32(1000),
			CurrentPage: tea.Int32(int32(n)),
		}
		runtime := &util.RuntimeOptions{}
		response, err = client.DescribeColumnsWithOptions(describeColumnsRequest, runtime)
		if err != nil {
			return
		}
		tempdata := response.Body.Items
		for i := 0; i < len(tempdata); i++ {
			input.UUID = uuid.New().String()
			input.GroupName = Group
			input.CreationTime = time.Unix(*tempdata[i].CreationTime/1000, 0).Format("2006-01-02 15:04:05")
			if tempdata[i].DataType != nil {
				input.DataType = *tempdata[i].DataType
			} else {
				input.DataType = ""
			}
			if tempdata[i].Id != nil {
				input.Id = *tempdata[i].Id
			} else {
				input.Id = ""
			}
			if tempdata[i].InstanceId != nil {
				input.InstanceId = int(*tempdata[i].InstanceId)
			} else {
				input.InstanceId = -1
			}
			if tempdata[i].InstanceName != nil {
				input.InstanceName = *tempdata[i].InstanceName
			} else {
				input.InstanceName = ""
			}
			if tempdata[i].Name != nil {
				input.Name = *tempdata[i].Name
			} else {
				input.Name = ""
			}
			if tempdata[i].ProductCode != nil {
				input.ProductCode = *tempdata[i].ProductCode
			} else {
				input.ProductCode = ""
			}
			if tempdata[i].RevisionId != nil {
				input.RevisionId = int(*tempdata[i].RevisionId)
			} else {
				input.RevisionId = -1
			}
			if tempdata[i].RevisionStatus != nil {
				input.RevisionStatus = int(*tempdata[i].RevisionStatus)
			} else {
				input.RevisionStatus = -1
			}
			if tempdata[i].RiskLevelId != nil {
				input.RiskLevelId = int(*tempdata[i].RiskLevelId)
			} else {
				input.RiskLevelId = -1
			}
			if tempdata[i].RiskLevelName != nil {
				input.RiskLevelName = *tempdata[i].RiskLevelName
			} else {
				input.RiskLevelName = ""
			}
			if tempdata[i].RuleId != nil {
				input.RuleId = int(*tempdata[i].RuleId)
			} else {
				input.RuleId = -1
			}
			if tempdata[i].RuleName != nil {
				input.RuleName = *tempdata[i].RuleName
			} else {
				input.RuleName = ""
			}
			if tempdata[i].SensLevelName != nil {
				input.SensLevelName = *tempdata[i].SensLevelName
			} else {
				input.SensLevelName = ""
			}
			if tempdata[i].Sensitive != nil {
				input.Sensitive = *tempdata[i].Sensitive
			}
			if tempdata[i].TableId != nil {
				input.TableId = int(*tempdata[i].TableId)
			}
			if tempdata[i].TableName != nil {
				input.TableName = *tempdata[i].TableName
			} else {
				input.TableName = ""
			}
			model.AddColumn(&input)
		}
		if n%5 == 0 {
			time.Sleep(1)
		}
	}
}

// 调用本接口查询SDDP中敏感数据识别规则的列表。
func initdescribeRules(Group string, client *sddp20190103.Client) {
	var input model.DataRules
	describeRulesRequest := &sddp20190103.DescribeRulesRequest{
		PageSize:    tea.Int32(1000),
		CurrentPage: tea.Int32(1),
	}
	runtime := &util.RuntimeOptions{}
	response, err := client.DescribeRulesWithOptions(describeRulesRequest, runtime)
	if err != nil {
		return
	}
	tempdata := response.Body.Items
	for i := 0; i < len(tempdata); i++ {
		input.UUID = uuid.New().String()
		if tempdata[i].Category != nil {
			input.Category = int(*tempdata[i].Category)
		}
		if tempdata[i].CategoryName != nil {
			input.CategoryName = *tempdata[i].CategoryName
		}
		if tempdata[i].Content != nil {
			input.Content = *tempdata[i].Content
		}
		if tempdata[i].ContentCategory != nil {
			input.ContentCategory = *tempdata[i].CategoryName
		}
		if tempdata[i].CustomType != nil {
			input.CustomType = int(*tempdata[i].CustomType)
		}
		if tempdata[i].Description != nil {
			input.Description = *tempdata[i].Description
		}
		if tempdata[i].DisplayName != nil {
			input.DisplayName = *tempdata[i].DisplayName
		}
		if tempdata[i].GmtCreate != nil {
			input.GmtCreate = time.Unix(*tempdata[i].GmtCreate/1000, 0).Format("2006-01-02 15:04:05")
		}
		if tempdata[i].GmtModified != nil {
			input.GmtModified = time.Unix(*tempdata[i].GmtModified/1000, 0).Format("2006-01-02 15:04:05")
		}
		if tempdata[i].GroupId != nil {
			input.GroupId = *tempdata[i].GroupId
		}
		if tempdata[i].HitTotalCount != nil {
			input.HitTotalCount = int(*tempdata[i].HitTotalCount)
		}
		if tempdata[i].Id != nil {
			input.Id = int(*tempdata[i].Id)
		}
		if tempdata[i].LoginName != nil {
			input.LoginName = *tempdata[i].LoginName
		}
		if tempdata[i].MajorKey != nil {
			input.MajorKey = *tempdata[i].MajorKey
		}
		if tempdata[i].Name != nil {
			input.Name = *tempdata[i].Name
		}
		if tempdata[i].ProductCode != nil {
			input.ProductCode = *tempdata[i].ProductCode
		}
		if tempdata[i].ProductId != nil {
			input.ProductId = int(*tempdata[i].ProductId)
		}
		if tempdata[i].RiskLevelId != nil {
			input.RiskLevelId = int(*tempdata[i].RiskLevelId)
		}
		if tempdata[i].RiskLevelName != nil {
			input.RiskLevelName = *tempdata[i].RiskLevelName
		}
		if tempdata[i].StatExpress != nil {
			input.StatExpress = *tempdata[i].StatExpress
		}
		if tempdata[i].Status != nil {
			input.Status = int(*tempdata[i].Status)
		}
		if tempdata[i].Target != nil {
			input.StatExpress = *tempdata[i].Target
		}
		if tempdata[i].UserId != nil {
			input.UserId = int(*tempdata[i].UserId)
		}
		if tempdata[i].WarnLevel != nil {
			input.WarnLevel = int(*tempdata[i].WarnLevel)
		}
		model.AddRules(&input)
	}
}

func initDescribeTables(Group string, client *sddp20190103.Client) {
	var input model.DataTable
	// 获取tables
	describeTablesRequest := &sddp20190103.DescribeTablesRequest{
		PageSize:    tea.Int32(1000),
		CurrentPage: tea.Int32(1),
	}
	runtime := &util.RuntimeOptions{}
	response, err := client.DescribeTablesWithOptions(describeTablesRequest, runtime)
	tempdata := response.Body.Items
	if err != nil {
		return
	}
	for i := 0; i < len(tempdata); i++ {
		input.UUID = uuid.New().String()
		input.GroupName = Group
		input.CreationTime = time.Unix(*tempdata[i].CreationTime/1000, 0).Format("2006-01-02 15:04:05")
		if tempdata[i].Id != nil {
			input.Id = int(*tempdata[i].Id)
		}
		if tempdata[i].InstanceDescription != nil {
			input.InstanceDescription = *tempdata[i].InstanceDescription
		} else {
			input.InstanceDescription = ""
		}
		if tempdata[i].InstanceId != nil {
			input.InstanceId = int(*tempdata[i].InstanceId)
		}
		if tempdata[i].InstanceName != nil {
			input.InstanceName = *tempdata[i].InstanceName
		} else {
			input.InstanceName = ""
		}
		if tempdata[i].Name != nil {
			input.Name = *tempdata[i].Name
		} else {
			input.Name = ""
		}
		if tempdata[i].Owner != nil {
			input.Owner = *tempdata[i].Owner
		} else {
			input.Owner = ""
		}
		if tempdata[i].ProductCode != nil {
			input.ProductCode = *tempdata[i].ProductCode
		} else {
			input.ProductCode = ""
		}
		if tempdata[i].ProductId != nil {
			input.ProductId = *tempdata[i].ProductId
		} else {
			input.ProductId = ""
		}
		if tempdata[i].RiskLevelId != nil {
			input.RiskLevelId = int(*tempdata[i].RiskLevelId)
		}
		if tempdata[i].RiskLevelName != nil {
			input.RiskLevelName = *tempdata[i].RiskLevelName
		} else {
			input.RiskLevelName = ""
		}
		if tempdata[i].Sensitive != nil {
			input.Sensitive = *tempdata[i].Sensitive
		}
		if tempdata[i].SensitiveCount != nil {
			input.SensitiveCount = int(*tempdata[i].SensitiveCount)
		}
		if tempdata[i].SensitiveRatio != nil {
			input.SensitiveRatio = *tempdata[i].SensitiveRatio
		} else {
			input.SensitiveRatio = ""
		}
		if tempdata[i].TenantName != nil {
			input.TenantName = *tempdata[i].TenantName
		} else {
			input.TenantName = ""
		}
		if tempdata[i].TotalCount != nil {
			input.TotalCount = int(*tempdata[i].TotalCount)
		}
		model.AddTables(&input)
	}
}

*/

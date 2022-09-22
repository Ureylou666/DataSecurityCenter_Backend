package model

import (
	"Backend/internal/Errmsg"
)

// DataColumn 定义数据资产表中列的数据
type DataColumn struct {
	UUID				string 	`gorm:"primaryKey" json:"UUID,omitempty"` //唯一id
	GroupName        	string  `gorm:"type:varchar(50)" json:"GroupName,omitempty"` //所属项目组
	CreationTime       	string  `gorm:"type:varchar(50)" json:"CreationTime,omitempty"` // 创建时间
	DataType           	string 	`gorm:"type:varchar(50)" json:"DataType"`  //数据资产表中列数据的数据类型
	Id                  string 	`gorm:"type:varchar(50)" json:"Id,omitempty"` //数据资产表中列的唯一标识ID
	InstanceId          int  	`gorm:"type:int" json:"InstanceId,omitempty"` //数据资产表中列数据所属的资产实例ID。
	InstanceName        string 	`gorm:"type:varchar(100)" json:"InstanceName,omitempty"` //数据资产表中列数据所属的资产实例名称。
	Name                string 	`gorm:"type:varchar(100)" json:"Name,omitempty"` //数据资产表中列的名称。
	ProductCode         string 	`gorm:"type:varchar(100)" json:"ProductCode,omitempty"` //数据资产表中列数据所属产品名称。取值：MaxCompute、OSS、ADS、OTS、RDS等
	RevisionId          int  	`gorm:"type:int" json:"RevisionId,omitempty"` //订正记录ID。
	RevisionStatus      int  	`gorm:"type:int" json:"RevisionStatus,omitempty"` //订正状态。取值：1：已订正。0：未订正。
	RiskLevelId         int  	`gorm:"type:int" json:"RiskLevelId,omitempty"` //数据资产表中列数据的风险等级ID。取值：1：N/A 2：S1 3：S2 4：S3 5：S4。
	RiskLevelName       string 	`gorm:"type:varchar(100)" json:"RiskLevelName,omitempty"` //敏感等级名。
	RuleId              int  	`gorm:"type:int" json:"RuleId,omitempty"` //数据资产表中列数据命中的敏感数据识别规则ID。
	RuleName            string 	`gorm:"type:varchar(100)" json:"RuleName,omitempty"` //数据资产表中列数据命中的敏感数据识别规则名称。
	SensLevelName       string 	`gorm:"type:varchar(100)" json:"SensLevelName,omitempty"` //敏感等级名。
	Sensitive           bool   	`gorm:"type:boolean" json:"Sensitive,omitempty"` //数据资产表中列数据是否包含敏感数据。
	TableId             int  	`gorm:"type:int" json:"TableId,omitempty"` //数据资产表中列数据所属的资产表ID。
	TableName           string 	`gorm:"type:varchar(255)" json:"TableName,omitempty"` //订正目标列所属的表名。
}

type ColumnDetailsQueryInfo struct {
	GroupName 		string
	RiskLevelName 	string
	RuleName 		string
	PageNum  		int
	PageSize 		int
}

type ColumnsQueryInfo struct {
	TableName		string
	PageNum  		int
	PageSize 		int
}

// AddColumn 新增aliyun数据列资产
func AddColumn(data *DataColumn) int {
	err := db.Create(&data).Error
	if err!=nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
}

// GetColumnDetails 查询Column详情
func GetColumnDetails(query ColumnDetailsQueryInfo) ([]DataColumn, int64, int64) {
	var result []DataColumn
	var resTotal, columnTotal int64
	// 分页处理
	db.Where(&DataColumn{GroupName: query.GroupName,RuleName: query.RuleName,RiskLevelName: query.RiskLevelName}).Where("risk_level_name <> ?", "N/A").Limit(-1).Find(&result).Count(&columnTotal)
	if query.PageNum == 0 || query.PageSize == 0 {
		db.Where(&DataColumn{GroupName: query.GroupName,RuleName: query.RuleName,RiskLevelName: query.RiskLevelName}).Where("risk_level_name <> ?", "N/A").Limit(-1).Find(&result)
		resTotal = int64(len(result))
	} else {
		db.Where(&DataColumn{GroupName: query.GroupName,RuleName: query.RuleName,RiskLevelName: query.RiskLevelName}).Where("risk_level_name <> ?", "N/A").Limit(query.PageSize).Offset((query.PageNum - 1) * query.PageSize).Find(&result)
		resTotal = int64(len(result))
	}
	return result, resTotal, columnTotal
}

// 查询指定表Column信息
func GetTableColumn(query ColumnsQueryInfo) ([]DataColumn, int64, int64) {
	var result []DataColumn
	var resTotal, columnTotal int64
	// 分页处理
	db.Where("table_name = ?", query.TableName).Find(&result).Count(&columnTotal)
	if query.PageNum == 0 || query.PageSize == 0 {
		db.Where("table_name = ?", query.TableName).Limit(-1).Order("risk_level_name desc, name").Find(&result)
		resTotal = int64(len(result))
	} else {
		db.Where("table_name = ?", query.TableName).Limit(query.PageSize).Offset((query.PageNum - 1) * query.PageSize).Order("risk_level_name desc, name").Find(&result)
		resTotal = int64(len(result))
	}
	return result, resTotal, columnTotal
}

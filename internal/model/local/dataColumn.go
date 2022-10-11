package local

import (
	"Backend/internal/utils/Errmsg"
)

// DataColumn 定义数据资产表中列的数据
type DataColumn struct {
	UUID          string `gorm:"primaryKey" json:"UUID,omitempty"`              //唯一id
	GroupName     string `gorm:"type:varchar(50)" json:"GroupName,omitempty"`   //所属项目组
	DataType      string `gorm:"type:varchar(50)" json:"DataType"`              //数据资产表中列数据的数据类型
	InstanceId    string `gorm:"type:varchar(200)" json:"InstanceId,omitempty"` //数据资产表中列数据所属的资产实例ID。
	DatabaseName  string `gorm:"type:varchar(200)" json:"DatabaseName"`         //数据实例Database的名称。
	TableName     string `gorm:"type:varchar(200)" json:"TableName,omitempty"`  // 数据资产表的名称
	ColumnName    string `gorm:"type:varchar(200)" json:"ColumnName"`
	RuleId        int    `gorm:"type:int" json:"RuleId,omitempty"`                 //数据资产表中列数据命中的敏感数据识别规则ID。
	RuleName      string `gorm:"type:varchar(100)" json:"RuleName,omitempty"`      //数据资产表中列数据命中的敏感数据识别规则名称。
	CategoryName  string `gorm:"type:varchar(100)" json:"CategoryName,omitempty"`  //分类名
	SensLevelName string `gorm:"type:varchar(100)" json:"SensLevelName,omitempty"` //等级名。
	SampleData    string `gorm:"type:string" json:"SampleData"`                    // 样例
}

type ColumnDetailsQueryInfo struct {
	GroupName     string
	RiskLevelName string
	RuleName      string
	PageNum       int
	PageSize      int
}

type ColumnsQueryInfo struct {
	TableName string
	PageNum   int
	PageSize  int
}

// InsertColumn 新增aliyun数据列资产
func InsertColumn(data *DataColumn) int {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
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

func DeleteColumnData(InstanceID string, DatabaseName string) {
	db.Where("instance_id = ? and database_name = ? ", InstanceID, DatabaseName).Delete(&DataColumn{})

}

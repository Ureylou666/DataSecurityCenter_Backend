package model

import (
	"Backend/internal/utils/Errmsg"
)

type DataDatabase struct {
	UUID                string `gorm:"primaryKey" json:"UUID"`
	CreationTime        string `gorm:"type:varchar(50)" json:"CreationTime"`         //创建该数据库实例的时间。使用时间戳表示，单位：毫秒。
	DepartName          string `gorm:"type:varchar(200)" json:"DepartName"`          //数据库实例所属部门的名称。
	SDDPInstanceId      int    `gorm:"type:int" json:"Id"`                           //数据安全中心服务中记录的数据资产实例的唯一标识ID。
	RDSInstanceID       string `gorm:"type:varchar(50)" json:"RDSInstanceID"`        //数据库所属实例id
	DatabaseName        string `gorm:"type:varchar(200)" json:"DatabaseName"`        //数据资产实例的名称。
	DatabaseDescription string `gorm:"type:varchar(200)" json:"DatabaseDescription"` //数据库的描述信息。
	DatabaseStatus      string `gorm:"type:varchar(200)" json:"DatabaseStatus"`      //数据库的状态信息。
	DatabaseEngine      string `gorm:"type:varchar(200)" json:"DatabaseEngine"`      //数据库引擎 mysql / pgsql
	Labelsec            bool   `gorm:"type:boolean" json:"Labelsec"`                 //数据资产实例的安全状态
	LastFinishTime      string `gorm:"type:varchar(50)" json:"LastFinishTime"`       //最近一次扫描数据资产实例的完成时间。使用时间戳表示，单位：毫秒。
	Owner               string `gorm:"type:varchar(200)" json:"Owner"`               //拥有该数据资产实例的阿里云账号。
	ProductCode         string `gorm:"type:varchar(200)" json:"ProductCode"`         //数据资产实例所属产品的名称，包括MaxCompute、OSS、RDS等。关于支持的具体产品名称
	ProductId           string `gorm:"type:varchar(200)" json:"ProductId"`           //数据资产实例所属产品的ID。
	Protection          bool   `gorm:"type:boolean" json:"Protection"`               //数据资产实例的防护状态
	RiskLevelId         int    `gorm:"type:int" json:"RiskLevelId"`                  //数据资产实例的风险等级ID。风险等级ID越高，表示识别出的数据越敏感
	RiskLevelName       string `gorm:"type:varchar(200)" json:"RiskLevelName"`       //数据资产实例的风险等级名称
	Sensitive           bool   `gorm:"type:boolean" json:"Sensitive"`                //数据资产实例中是否包含敏感数据。
	SensitiveCount      int    `gorm:"type:int" json:"SensitiveCount"`               //数据资产实例中包含的敏感数据总数。例如：当数据资产为RDS时，表示该实例中数据库的敏感总表数
	TotalCount          int    `gorm:"type:int" json:"TotalCount"`                   //数据资产实例中的数据总数。例如：当数据资产为RDS时，表示该实例中数据库的总表数。
}

// AddDatabase 新增aliyun数据实例资产
func AddDatabase(data *DataDatabase) int {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
}

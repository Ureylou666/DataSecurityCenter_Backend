package local

import (
	"Backend/internal/utils/Errmsg"
)

type DataDatabase struct {
	UUID           string `gorm:"primaryKey" json:"UUID"`
	DepartName     string `gorm:"type:varchar(200)" json:"DepartName"`          //数据库实例所属部门的名称。
	InstanceID     string `gorm:"type:varchar(50)" json:"RDSInstanceID"`        //数据库所属实例id
	DatabaseName   string `gorm:"type:varchar(200)" json:"DatabaseName"`        //数据资产实例的名称。
	Description    string `gorm:"type:varchar(200)" json:"DatabaseDescription"` //数据库的描述信息。
	Status         string `gorm:"type:varchar(200)" json:"DatabaseStatus"`      //数据库的状态信息。
	Engine         string `gorm:"type:varchar(200)" json:"DatabaseEngine"`      //数据库引擎 mysql / pgsql
	Owner          string `gorm:"type:varchar(200)" json:"Owner"`               //拥有该数据资产实例的阿里云账号。
	RiskLevelId    int    `gorm:"type:int" json:"RiskLevelId"`                  //数据资产实例的风险等级ID。风险等级ID越高，表示识别出的数据越敏感
	RiskLevelName  string `gorm:"type:varchar(200)" json:"RiskLevelName"`       //数据资产实例的风险等级名称
	Sensitive      bool   `gorm:"type:boolean" json:"Sensitive"`                //数据资产实例中是否包含敏感数据。
	SensitiveCount int    `gorm:"type:int" json:"SensitiveCount"`               //数据资产实例中包含的敏感数据总数。例如：当数据资产为RDS时，表示该实例中数据库的敏感总表数
	TotalCount     int    `gorm:"type:int" json:"TotalCount"`                   //数据资产实例中的数据总数。例如：当数据资产为RDS时，表示该实例中数据库的总表数。
}

// AddDatabase 新增RDS Instance下 Database信息
func AddDatabase(data *DataDatabase) int {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
}

// DeleteDatabase 删除RDS Instance下 Database信息
func DeleteDatabase(RDSInstance string) {
	db.Where("instance_id = ?", RDSInstance).Delete(&DataDatabase{})
}

// GetDatabase 获取aliyun数据资产列表 并进行分页展示
func GetDatabaseList(RDSInstance string) []DataDatabase {
	var result []DataDatabase
	// 获取数据库中Instance总数
	db.Where("instance_id = ?", RDSInstance).Find(&result)
	return result
}

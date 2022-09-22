package model

import "Backend/internal/Errmsg"

type RulesQueryInfo struct {
	RuleName		string
	RiskLevelName 	string
	PageNum			int
	PageSize		int
}

type DataRules struct {
	UUID			string 		`gorm:"primaryKey" json:"UUID,omitempty"` 					//唯一id
	Category        int  		`gorm:"type:int" json:"Category,omitempty"` 				//敏感数据识别规则内容的类型。取值： 0：关键字。 2：正则表达式。
	CategoryName    string 		`gorm:"type:varchar(100)" json:"CategoryName,omitempty"`	//敏感数据识别规则内容所属类型名称。
	Content         string 		`gorm:"type:varchar(100)" json:"Content,omitempty"`			//敏感数据识别规则的内容。
	ContentCategory string 		`gorm:"type:varchar(100)" json:"ContentCategory,omitempty"`	//内容类型，取值：1：SQL注入尝试利用，2：SQL注入尝试绕过， 3：存储过程滥用，4：缓冲区溢出，5：基于报错的SQL注入等。
	CustomType      int  		`gorm:"type:int" json:"CustomType,omitempty"`				//敏感数据识别规则类型。 0：表示系统内置。 1：表示用户自定义。
	Description     string 		`gorm:"type:varchar(100)" json:"Description,omitempty"`		//敏感数据识别规则的描述信息。
	DisplayName     string 		`gorm:"type:varchar(100)" json:"DisplayName,omitempty"` 	//敏感数据识别规则的创建人账号显示名。
	GmtCreate       string  	`gorm:"type:varchar(50)" json:"GmtCreate,omitempty"`		//敏感数据识别规则的创建时间。
	GmtModified     string  	`gorm:"type:varchar(50)" json:"GmtModified,omitempty"`		//敏感数据识别规则的修改时间。格式：时间戳。单位：毫秒。
	GroupId         string 		`gorm:"type:varchar(100)" json:"GroupId,omitempty"`			//规则父类分组。
	HitTotalCount   int  		`gorm:"type:int" json:"HitTotalCount,omitempty"`			//命中次数。
	Id              int  		`gorm:"type:int" json:"Id,omitempty"`						//保存记录的敏感数据识别规则的唯一标识ID。
	LoginName       string 		`gorm:"type:varchar(100)" json:"LoginName,omitempty"`		//敏感数据识别规则的创建人账号登录名
	MajorKey        string 		`gorm:"type:varchar(100)" json:"MajorKey,omitempty"`		//主维度key。
	Name            string 		`gorm:"type:varchar(100)" json:"Name,omitempty"`			//敏感数据识别规则的名称。
	ProductCode     string 		`gorm:"type:varchar(100)" json:"ProductCode,omitempty"`		//数据资产所属的产品名称。取值：MaxCompute、OSS、ADS、OTS、RDS等。
	ProductId       int  		`gorm:"type:int" json:"ProductId,omitempty"`				//数据资产所属的产品名称对应的ID。取值：1：MaxCompute，2：OSS，3：ADS，4：OTS，5：RDS等。
	RiskLevelId     int  		`gorm:"type:int" json:"RiskLevelId,omitempty"`				//敏感数据识别规则的敏感等级ID。取值：
	RiskLevelName   string 		`gorm:"type:varchar(100)" json:"RiskLevelName,omitempty"`	//敏感数据识别规则的风险等级名称。取值： N/A：未识别到敏感数据。 S1：1级敏感数据。 S2：2级敏感数据。 S3：3级敏感数据。 S4：4级敏感数据。
	StatExpress     string 		`gorm:"type:varchar(100)" json:"StatExpress,omitempty"`		//统计表达式。
	Status          int  		`gorm:"type:int" json:"Status,omitempty"`					//敏感数据识别规则的检测状态。取值： 0：关闭。 1：开启。
	Target          string 		`gorm:"type:varchar(100)" json:"Target,omitempty"`			//数据资产所属的产品名称。取值：MaxCompute、OSS、ADS、OTS、RDS等。
	UserId          int  		`gorm:"type:int" json:"UserId,omitempty"`					//敏感数据识别规则的创建人账号ID。
	WarnLevel       int  		`gorm:"type:int" json:"WarnLevel,omitempty"`				//危险等级。 1：低危。 2：中危。 3：高危。
}

func AddRules(data *DataRules) int {
	err := db.Create(&data).Error
	if err!=nil {
		return Errmsg.ERROR
	}
	return Errmsg.SUCCESS
}

func GetRules(query RulesQueryInfo) ([]DataRules, int64, int64) {
	var result []DataRules
	var resTotal, rulesTotal int64
	//分页处理
	db.Where(&DataRules{Name: query.RuleName, RiskLevelName: query.RiskLevelName}).Where("name like ?", "%"+query.RuleName+"%").Find(&result).Limit(-1).Count(&rulesTotal)
	if query.PageNum == 0 || query.PageSize == 0 {
		db.Where(&DataRules{RiskLevelName: query.RiskLevelName }).Where("name like ?", "%"+query.RuleName+"%").Order("risk_level_name desc, name").Limit(-1).Find(&result)
		resTotal = int64(len(result))
	} else {
		db.Where(&DataRules{RiskLevelName: query.RiskLevelName }).Where("name like ?", "%"+query.RuleName+"%").Order("risk_level_name desc, name").Limit(query.PageSize).Offset((query.PageNum - 1) * query.PageSize).Find(&result)
		resTotal = int64(len(result))
	}
	return  result, resTotal, rulesTotal
}

package local

type DataSample struct {
	ColumnUUID   string `gorm:"primaryKey" json:"UUID,omitempty"` //唯一id
	InstanceID   string `gorm:"type:string"`
	DatabaseName string `gorm:"type:string"`
	TableName    string `gorm:"type:string"`
	SampleData0  string `gorm:"type:string"`
	SampleData1  string `gorm:"type:string"`
	SampleData2  string `gorm:"type:string"`
	SampleData3  string `gorm:"type:string"`
	SampleData4  string `gorm:"type:string"`
	SampleData5  string `gorm:"type:string"`
	SampleData6  string `gorm:"type:string"`
	SampleData7  string `gorm:"type:string"`
	SampleData8  string `gorm:"type:string"`
	SampleData9  string `gorm:"type:string"`
}

// InsertSampleData 新增数据样例
func InsertSampleData(data *DataSample) {
	db.Create(&data)
}

func DeleteSampleData(InstanceID string, DatabaseName string) {
	db.Where("instance_id = ? and database_name = ?", InstanceID, DatabaseName).Delete(&DataSample{})
}

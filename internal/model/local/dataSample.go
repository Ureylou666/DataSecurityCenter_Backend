package local

type DataSample struct {
	ColumnUUID   string `gorm:"primaryKey" json:"UUID,omitempty"` //唯一id
	InstanceID   string `gorm:"type:string"`
	DatabaseName string `gorm:"type:string"`
	TableName    string `gorm:"type:string"`
	SampleData0  string `gorm:"type:text"`
	SampleData1  string `gorm:"type:text"`
	SampleData2  string `gorm:"type:text"`
	SampleData3  string `gorm:"type:text"`
	SampleData4  string `gorm:"type:text"`
	SampleData5  string `gorm:"type:text"`
	SampleData6  string `gorm:"type:text"`
	SampleData7  string `gorm:"type:text"`
	SampleData8  string `gorm:"type:text"`
	SampleData9  string `gorm:"type:text"`
}

// InsertSampleData 新增数据样例
func InsertSampleData(data *DataSample) {
	db.Create(&data)
}

func DeleteSampleData(InstanceID string, DatabaseName string) {
	db.Where("instance_id = ? and database_name = ?", InstanceID, DatabaseName).Delete(&DataSample{})
}

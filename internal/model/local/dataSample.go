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

func GetSingleSampleData(uuid string) (result [10]string) {
	var temp DataSample
	db.Where("column_uuid = ?", uuid).Find(&temp)
	result[0] = temp.SampleData0
	result[1] = temp.SampleData1
	result[2] = temp.SampleData2
	result[3] = temp.SampleData3
	result[4] = temp.SampleData4
	result[5] = temp.SampleData5
	result[6] = temp.SampleData6
	result[7] = temp.SampleData7
	result[8] = temp.SampleData8
	result[9] = temp.SampleData9
	return result
}

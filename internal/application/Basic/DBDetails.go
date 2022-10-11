package Basic

import (
	"Backend/internal/application/aliyunSDK"
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"Backend/internal/utils/setting"
	rds20140815 "github.com/alibabacloud-go/rds-20140815/v2/client"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ConnectionString struct {
	Db         string
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DBPassword string
	SSLEnabled bool
}

// UpdateDBDetails 获取、更新数据库表 / 字段 详情
func UpdateDBDetails(InstanceID string) (ErrCode int, ErrMessage string) {
	// 错误控制，判断传入的InstanceID是否存在
	if local.CheckInventoryExist(InstanceID) != true {
		return Errmsg.ErrorCheckInventoryExist, Errmsg.GetErrMsg(Errmsg.ErrorCheckInventoryExist)
	}
	// RDSInstanceID存在 则可以使用sdk 解锁审计账号cnisdp
	// 创建aliyunRDS客户端
	// 获取CloudAccountID
	CloudAccountId := local.GetAccountID(InstanceID)
	// 使用STS 创建客户端
	previousCode, previousMsg, RDSClient := aliyunSDK.CreateRDSClient(CloudAccountId)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg.Error()
	}
	// 注意 若审计账号不存在，则会进行创建新账号
	previousCode, previousMsg = UnlockRDSAccount(InstanceID, RDSClient)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg.Error()
	}
	setting.LoadAuditAccount()
	// 获取连接rds connection string
	previousCode, previousMsg, DbConnectString := InitDBConnection(InstanceID)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg.Error()
	}
	// 获取连接所需的数据库列表
	DbList := local.GetDatabaseList(InstanceID)
	// 遍历连接每一个数据库 获取table及column信息
	for i := 0; i < len(DbList); i++ {
		// 初始化连接字符串
		DbConnectString.DbName = DbList[i].DatabaseName
		DbConnectString.DbUser = setting.AccountName
		DbConnectString.DBPassword = setting.AccountPassword + InstanceID[len(InstanceID)-4:len(InstanceID)-1]
		// 对不同类型数据库 分别update
		switch DbList[i].Engine {
		case "PostgreSQL":
			{
				previousCode, previousMsg = UpdatePgsqlDetails(DbConnectString, InstanceID, DbList[i].DatabaseName)
			}
		case "MySQL":
			{

			}
		}
	}
	// 锁定审计账户
	previousCode, previousMsg = aliyunSDK.LockAccount(InstanceID, RDSClient)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg.Error()
	}
	return Errmsg.SUCCESS, ""
}

// UnlockRDSAccount 解锁数据库账号, 无账号则会进行创建
func UnlockRDSAccount(InstanceID string, RDSClient *rds20140815.Client) (ErrCode int, ErrMessage error) {
	// 判断审计用户是否存在，不存在进行创建
	if local.CheckAccountExist(InstanceID, "cnisdp") != true {
		previousCode, previousMsg := aliyunSDK.CreateRDSAccount(InstanceID, RDSClient)
		if previousCode != Errmsg.SUCCESS {
			return previousCode, previousMsg
		}
		var input local.RDSAccount
		input.RDSInstanceID = InstanceID
		input.AccountStatus = "Available"
		input.AccountDescription = setting.AccountDescription
		input.AccountType = setting.AccountType
		input.AccountName = setting.AccountName
		local.AddDatabaseAccount(&input)
	}
	// 解锁用户账号
	previousCode, previousMsg := aliyunSDK.UnlockAccount(InstanceID, RDSClient)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	return Errmsg.SUCCESS, nil
}

// InitDBConnection 判断数据库连接
func InitDBConnection(InstanceID string) (ErrCode int, ErrMessage error, DbConnectString ConnectionString) {
	var RDSDetails local.DataInventory
	RDSDetails = local.GetConnectString(InstanceID)
	// 赋值
	DbConnectString.Db = RDSDetails.RDSEngine
	DbConnectString.DbHost = RDSDetails.RDSConnectionString
	DbConnectString.DbPort = RDSDetails.RDSConnectionPort
	DbConnectString.SSLEnabled = RDSDetails.SSLEnabled
	return Errmsg.SUCCESS, nil, DbConnectString
}

// UpdatePgsqlDetails 更新pgsql 详情
func UpdatePgsqlDetails(connectionString ConnectionString, InstanceID string, DatabaseName string) (ErrCode int, ErrMessage error) {
	var TableName []interface{}
	var ColumnName []interface{}
	var columnDataType []interface{}
	var inputTable local.DataTable
	var inputColumn local.DataColumn
	var dsn string
	// 判断是否使用了ssl，生成特定对orm连接语句
	if connectionString.SSLEnabled {
		dsn = "host=pgm-uf61xs38ffvzk6z14o.pg.rds.aliyuncs.com" + " user=" + connectionString.DbUser + " password=" + connectionString.DBPassword + " dbname=" + connectionString.DbName + " port=" + connectionString.DbPort + " sslmode=verify-ca sslrootcert=" + setting.SSLCert
	} else {
		dsn = "host=pgm-uf61xs38ffvzk6z14o.pg.rds.aliyuncs.com" + " user=" + connectionString.DbUser + " password=" + connectionString.DBPassword + " dbname=" + connectionString.DbName + " port=" + connectionString.DbPort + " sslmode=disables"
	}
	// 创建数据库客户端
	targetDB, ErrMessage := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if ErrMessage != nil {
		return Errmsg.ErrorUpdatePgsqlDetailsConnection, ErrMessage
	}
	// 删除本地数据库中对应数据库中数据
	local.DeleteTableData(InstanceID, DatabaseName)
	local.DeleteColumnData(InstanceID, DatabaseName)
	local.DeleteSampleData(InstanceID, DatabaseName)
	// 获取数据库table信息
	targetDB.Raw("SELECT tablename FROM pg_tables WHERE tablename NOT LIKE 'pg%' AND tablename NOT LIKE 'sql_%'").Scan(&TableName)
	// 对每个table获取column信息
	for i := 0; i < len(TableName); i++ {
		inputTable.UUID = uuid.New().String()
		inputTable.InstanceId = InstanceID
		inputTable.DatabaseName = DatabaseName
		inputTable.TableName = TableName[i].(string)
		// 添加到数据库
		local.InsertTableData(&inputTable)
		// 初始化并遍历Column
		targetDB.Raw("SELECT column_name FROM information_schema.columns WHERE table_schema='public' AND table_name = ?", TableName[i].(string)).Scan(&ColumnName)
		for j := 0; j < len(ColumnName); j++ {
			// 初始化sampleData
			columnSampleData := make([]interface{}, 10)
			inputSampleData := new(local.DataSample)
			targetDB.Table(TableName[i].(string)).Distinct(ColumnName[j].(string)).Limit(10).Scan(&columnSampleData)
			targetDB.Raw("select data_type from information_schema.columns where table_name = ? and column_name =?", TableName[i].(string), ColumnName[j].(string)).Scan(&columnDataType)
			// 更新数据列
			inputColumn.UUID = uuid.New().String()
			inputColumn.DataType = columnDataType[0].(string)
			inputColumn.InstanceId = InstanceID
			inputColumn.DatabaseName = DatabaseName
			inputColumn.TableName = TableName[i].(string)
			inputColumn.ColumnName = ColumnName[j].(string)
			if columnSampleData != nil {
				inputColumn.SampleData = setting.InterfaceToString(columnSampleData[0])
			}
			local.InsertColumn(&inputColumn)
			// 更新数据样例
			inputSampleData.ColumnUUID = inputColumn.UUID
			inputSampleData.InstanceID = InstanceID
			inputSampleData.DatabaseName = DatabaseName
			inputSampleData.TableName = TableName[i].(string)
			if len(columnSampleData) > 0 {
				inputSampleData.SampleData0 = setting.InterfaceToString(columnSampleData[0])
			}
			if len(columnSampleData) > 1 {
				inputSampleData.SampleData1 = setting.InterfaceToString(columnSampleData[1])
			}
			if len(columnSampleData) > 2 {
				inputSampleData.SampleData2 = setting.InterfaceToString(columnSampleData[2])
			}
			if len(columnSampleData) > 3 {
				inputSampleData.SampleData3 = setting.InterfaceToString(columnSampleData[3])
			}
			if len(columnSampleData) > 4 {
				inputSampleData.SampleData4 = setting.InterfaceToString(columnSampleData[4])
			}
			if len(columnSampleData) > 5 {
				inputSampleData.SampleData5 = setting.InterfaceToString(columnSampleData[5])
			}
			if len(columnSampleData) > 6 {
				inputSampleData.SampleData6 = setting.InterfaceToString(columnSampleData[6])
			}
			if len(columnSampleData) > 7 {
				inputSampleData.SampleData7 = setting.InterfaceToString(columnSampleData[7])
			}
			if len(columnSampleData) > 8 {
				inputSampleData.SampleData8 = setting.InterfaceToString(columnSampleData[8])
			}
			if len(columnSampleData) > 9 {
				inputSampleData.SampleData9 = setting.InterfaceToString(columnSampleData[9])
			}
			local.InsertSampleData(inputSampleData)
		}
	}
	return Errmsg.SUCCESS, nil
}

// UpdateMysqlDetails 更新Mysql 详情
/*
func UpdateMysqlDetails(connectionString ConnectionString, InstanceID string, DatabaseName string) (ErrCode int, ErrMessage error) {
	var dsn string
	// 判断是否使用了ssl，生成特定对orm连接语句
	if connectionString.SSLEnabled {
		dsn = "host=pgm-uf61xs38ffvzk6z14o.pg.rds.aliyuncs.com" + " user=" + connectionString.DbUser + " password=" + connectionString.DBPassword + " dbname=" + connectionString.DbName + " port=" + connectionString.DbPort + " sslmode=verify-ca sslrootcert=" + setting.SSLCert
	} else {
		dsn = "host=pgm-uf61xs38ffvzk6z14o.pg.rds.aliyuncs.com" + " user=" + connectionString.DbUser + " password=" + connectionString.DBPassword + " dbname=" + connectionString.DbName + " port=" + connectionString.DbPort + " sslmode=disables"
	}
	// 创建数据库客户端
	// 删除本地数据库中对应数据库中数据
	// 获取数据库table信息
	// 对每个table获取column信息
	return Errmsg.SUCCESS, nil
}


*/

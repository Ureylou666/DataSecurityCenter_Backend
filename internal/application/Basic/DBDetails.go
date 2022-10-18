package Basic

import (
	"Backend/internal/application/aliyunSDK"
	"Backend/internal/model/local"
	"Backend/internal/utils"
	"Backend/internal/utils/Errmsg"
	"Backend/internal/utils/setting"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	rds20140815 "github.com/alibabacloud-go/rds-20140815/v2/client"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"log"
)

type ConnectionString struct {
	Db         string
	DbHost     string
	DbIPAddr   string
	DbPort     string
	DbName     string
	DbUser     string
	DBPassword string
	SSLEnabled bool
}

// UpdateDBDetails 获取、更新数据库表 / 字段 详情
func UpdateDBDetails(InstanceID string) (ErrCode int, ErrMessage error) {
	// 错误控制，判断传入的InstanceID是否存在
	if local.CheckInventoryExist(InstanceID) != true {
		return Errmsg.ErrorCheckInventoryExist, errors.New(Errmsg.GetErrMsg(Errmsg.ErrorCheckInventoryExist))
	}
	// RDSInstanceID存在 则可以使用sdk 解锁审计账号cnisdp
	// 创建aliyunRDS客户端
	// 获取CloudAccountID
	CloudAccountId := local.GetAccountID(InstanceID)
	// 获取连接所需的数据库列表, 因为要获取数据库类型，故前置处理
	DbList := local.GetDatabaseList(InstanceID)
	// 使用STS 创建客户端
	previousCode, previousMsg, RDSClient := aliyunSDK.CreateRDSClient(CloudAccountId)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	// 删除审计账户
	_, _ = aliyunSDK.DeleteAccount(InstanceID, RDSClient)
	// 密码为随机生成的25位数字、 大小写字母、特殊字符字符串
	TempPasswd := utils.GeneratePassword()
	previousCode, previousMsg = InitRDSAccount(InstanceID, DbList[0].Engine, TempPasswd, RDSClient)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	setting.LoadAuditAccount()
	// 获取连接rds connection string
	previousCode, previousMsg, DbConnectString, NewPublicAddr := InitDBConnection(InstanceID, RDSClient)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	// 遍历连接每一个数据库 获取table及column信息
	for i := 0; i < len(DbList); i++ {
		// 删除本地数据库中对应数据库中数据
		local.DeleteTableData(InstanceID, DbList[i].DatabaseName)
		local.DeleteColumnData(InstanceID, DbList[i].DatabaseName)
		local.DeleteSampleData(InstanceID, DbList[i].DatabaseName)
		// 初始化连接字符串
		DbConnectString.DbName = DbList[i].DatabaseName
		DbConnectString.DbUser = setting.AccountName
		DbConnectString.DBPassword = TempPasswd
		// 获取域名解析记录
		_, _, DbConnectString.DbHost, DbConnectString.DbIPAddr = aliyunSDK.DescribeDBInstanceNetInfo(InstanceID, RDSClient)
		// 对不同类型数据库 分别update
		switch DbList[i].Engine {
		case "PostgreSQL":
			{
				previousCode, previousMsg = UpdatePgsqlDetails(DbConnectString, InstanceID, DbList[i].DatabaseName)
			}
		case "MySQL":
			{
				previousCode, previousMsg = aliyunSDK.GrantAccountPrivilege(InstanceID, DbList[i].DatabaseName, RDSClient)
				previousCode, previousMsg = UpdateMysqlDetails(DbConnectString, InstanceID, DbList[i].DatabaseName)
			}
		case "SQLServer":
			{
				previousCode, previousMsg = aliyunSDK.GrantAccountPrivilege(InstanceID, DbList[i].DatabaseName, RDSClient)
				previousCode, previousMsg = UpdateSqlServerDetails(DbConnectString, InstanceID, DbList[i].DatabaseName)
			}
		}
	}
	// 删除审计账户
	previousCode, previousMsg = aliyunSDK.DeleteAccount(InstanceID, RDSClient)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	// 删除白名单
	previousCode, previousMsg = aliyunSDK.ModifySecurityIps(InstanceID, "Delete", RDSClient)
	// 释放外网地址
	if NewPublicAddr {
		previousCode, previousMsg = aliyunSDK.ReleaseInstancePublicConnection(InstanceID, DbConnectString.DbHost, RDSClient)
	}
	return Errmsg.SUCCESS, nil
}

// InitRDSAccount 解锁数据库账号, 无账号则会进行创建
func InitRDSAccount(InstanceID string, DbEngine string, TempPasswd string, RDSClient *rds20140815.Client) (ErrCode int, ErrMessage error) {
	// 创建rds account
	previousCode, previousMsg := aliyunSDK.CreateRDSAccount(InstanceID, DbEngine, TempPasswd, RDSClient)
	if previousCode != Errmsg.SUCCESS {
		return previousCode, previousMsg
	}
	return Errmsg.SUCCESS, nil
}

// InitDBConnection 判断数据库连接
func InitDBConnection(InstanceID string, RDSClient *rds20140815.Client) (ErrCode int, ErrMessage error, DbConnectString ConnectionString, NewPublicAddr bool) {
	var RDSDetails local.DataInventory
	// 赋值基本信息
	RDSDetails = local.GetConnectString(InstanceID)
	// 赋值
	DbConnectString.Db = RDSDetails.RDSEngine
	DbConnectString.DbPort = RDSDetails.RDSConnectionPort
	NewPublicAddr = false
	// 查询并判断是否存在外网地址
	ErrCode, ErrMessage, DbConnectString.DbHost, _ = aliyunSDK.DescribeDBInstanceNetInfo(InstanceID, RDSClient)
	if ErrMessage != nil {
		return ErrCode, ErrMessage, DbConnectString, NewPublicAddr
	}
	// 存在外网地址 判断是否开启了ssl； 不存在则创建外网连接地址
	if DbConnectString.DbHost != "" {
		DbConnectString.SSLEnabled = RDSDetails.SSLEnabled
	} else {
		// 创建外网连接地址
		ErrCode, ErrMessage, DbConnectString.DbHost = aliyunSDK.AllocateInstancePublicConnection(InstanceID, DbConnectString.DbPort, RDSClient)
		NewPublicAddr = true
	}
	// 开启外网访问白名单
	ErrCode, ErrMessage = aliyunSDK.ModifySecurityIps(InstanceID, "Append", RDSClient)
	return Errmsg.SUCCESS, nil, DbConnectString, NewPublicAddr
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
		dsn = "host=" + connectionString.DbIPAddr + " user=" + connectionString.DbUser + " password=" + connectionString.DBPassword + " dbname=" + connectionString.DbName + " port=" + connectionString.DbPort + " sslmode=verify-ca sslrootcert=" + setting.SSLCert
	} else {
		dsn = "host=" + connectionString.DbIPAddr + " user=" + connectionString.DbUser + " password=" + connectionString.DBPassword + " dbname=" + connectionString.DbName + " port=" + connectionString.DbPort + " sslmode=disables"
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
			// 当该列中有数据 才有意义进行更新数据列
			if columnSampleData[0] != nil {
				inputColumn.UUID = uuid.New().String()
				inputColumn.DataType = columnDataType[0].(string)
				inputColumn.InstanceId = InstanceID
				inputColumn.DatabaseName = DatabaseName
				inputColumn.TableName = TableName[i].(string)
				inputColumn.ColumnName = ColumnName[j].(string)
				inputColumn.SampleData = utils.InterfaceToString(columnSampleData[0])
				local.InsertColumn(&inputColumn)
				// 更新数据样例，只有大于10行数据才有意义
				if len(columnSampleData) == 10 {
					inputSampleData.ColumnUUID = inputColumn.UUID
					inputSampleData.InstanceID = InstanceID
					inputSampleData.DatabaseName = DatabaseName
					inputSampleData.TableName = TableName[i].(string)
					inputSampleData.SampleData0 = utils.InterfaceToString(columnSampleData[0])
					inputSampleData.SampleData1 = utils.InterfaceToString(columnSampleData[1])
					inputSampleData.SampleData2 = utils.InterfaceToString(columnSampleData[2])
					inputSampleData.SampleData3 = utils.InterfaceToString(columnSampleData[3])
					inputSampleData.SampleData4 = utils.InterfaceToString(columnSampleData[4])
					inputSampleData.SampleData5 = utils.InterfaceToString(columnSampleData[5])
					inputSampleData.SampleData6 = utils.InterfaceToString(columnSampleData[6])
					inputSampleData.SampleData7 = utils.InterfaceToString(columnSampleData[7])
					inputSampleData.SampleData8 = utils.InterfaceToString(columnSampleData[8])
					inputSampleData.SampleData9 = utils.InterfaceToString(columnSampleData[9])
					local.InsertSampleData(inputSampleData)
				}
			}
		}
	}
	return Errmsg.SUCCESS, nil
}

// UpdateMysqlDetails 更新Mysql 详情
func UpdateMysqlDetails(connectionString ConnectionString, InstanceID string, DatabaseName string) (ErrCode int, ErrMessage error) {
	// 判断是否使用了SSL进行连接
	if connectionString.SSLEnabled {
		mysqlSSLconfig()
	}
	// 初始化连接客户端
	targetDb, ErrMessage := sql.Open("mysql", connectionString.DbUser+":"+connectionString.DBPassword+"@tcp("+connectionString.DbIPAddr+":"+connectionString.DbPort+")/"+connectionString.DbName)
	if ErrMessage != nil {
		return Errmsg.ErrorUpdateMysqlDetailsConnection, ErrMessage
	}
	defer targetDb.Close()
	// 获取数据库表名
	sqlStrTables := "show tables"
	TableRows, ErrMessage := targetDb.Query(sqlStrTables)
	var TableName, ColumnName, DataType string
	for TableRows.Next() {
		TableRows.Scan(&TableName)
		inputTable := new(local.DataTable)
		inputTable.UUID = uuid.New().String()
		inputTable.InstanceId = InstanceID
		inputTable.DatabaseName = DatabaseName
		inputTable.TableName = TableName
		// 添加到数据库
		local.InsertTableData(inputTable)
		// 获取每个表下的列名
		sqlStrColumns := "select column_name from information_schema.columns where TABLE_NAME = ?"
		ColumnRows, _ := targetDb.Query(sqlStrColumns, TableName)
		for ColumnRows.Next() {
			ColumnRows.Scan(&ColumnName)
			// 查看列的数据类型
			sqlStrDataType := "select DATA_TYPE from information_schema.columns where TABLE_NAME = " + TableName + " and COLUMN_NAME = " + ColumnName
			_ = targetDb.QueryRow(sqlStrDataType).Scan(&DataType)
			// 更新数据列
			inputColumn := new(local.DataColumn)
			inputColumn.UUID = uuid.New().String()
			inputColumn.InstanceId = InstanceID
			inputColumn.DatabaseName = DatabaseName
			inputColumn.TableName = TableName
			inputColumn.ColumnName = ColumnName
			inputColumn.DataType = DataType
			// 获取每一列的sample data
			sqlStrSampleData := "select distinct " + ColumnName + " from " + TableName + " where " + ColumnName + " is not NULL limit 10"
			SampleDataRows, _ := targetDb.Query(sqlStrSampleData)
			// 初始化sampleData
			inputSampleData := new(local.DataSample)
			inputSampleData.ColumnUUID = inputColumn.UUID
			inputSampleData.InstanceID = InstanceID
			inputSampleData.DatabaseName = DatabaseName
			inputSampleData.TableName = TableName
			i := 0
			SampleData := make([]interface{}, 10)
			for SampleDataRows.Next() {
				SampleDataRows.Scan(&SampleData[i])
				switch i {
				case 0:
					inputSampleData.SampleData0 = utils.InterfaceToString(SampleData[i])
					if inputSampleData.SampleData0 != "" {
						inputColumn.SampleData = inputSampleData.SampleData0
						local.InsertColumn(inputColumn)
					}
				case 1:
					inputSampleData.SampleData1 = utils.InterfaceToString(SampleData[i])
				case 2:
					inputSampleData.SampleData2 = utils.InterfaceToString(SampleData[i])
				case 3:
					inputSampleData.SampleData3 = utils.InterfaceToString(SampleData[i])
				case 4:
					inputSampleData.SampleData4 = utils.InterfaceToString(SampleData[i])
				case 5:
					inputSampleData.SampleData5 = utils.InterfaceToString(SampleData[i])
				case 6:
					inputSampleData.SampleData6 = utils.InterfaceToString(SampleData[i])
				case 7:
					inputSampleData.SampleData7 = utils.InterfaceToString(SampleData[i])
				case 8:
					inputSampleData.SampleData8 = utils.InterfaceToString(SampleData[i])
				case 9:
					inputSampleData.SampleData9 = utils.InterfaceToString(SampleData[i])
				}
				i++
			}
			if inputSampleData.SampleData9 != "" {
				local.InsertSampleData(inputSampleData)
			}
			defer SampleDataRows.Close()
		}
		defer ColumnRows.Close()
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer TableRows.Close()
	return Errmsg.SUCCESS, nil
}

// 配置mysqlSSL证书
func mysqlSSLconfig() {
	setting.LoadDatabase()
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(setting.SSLCert)
	if err != nil {
		log.Fatal(err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append PEM.")
	}
	err = mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	})
}

func UpdateSqlServerDetails(connectionString ConnectionString, InstanceID string, DatabaseName string) (ErrCode int, ErrMessage error) {
	dsn := "server=" + connectionString.DbIPAddr + ";user id=" + connectionString.DbUser + ";password=" + connectionString.DBPassword + ";port=" + connectionString.DbPort + ";database=" + connectionString.DbName + ";encrypt=disable"
	targetDB, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	var TableName []interface{}
	var ColumnName []interface{}
	var columnDataType []interface{}
	var inputTable local.DataTable
	var inputColumn local.DataColumn
	// 获取表名
	targetDB.Raw("SELECT name FROM SysObjects Where XType='U' ORDER BY Name").Scan(&TableName)
	// 对每个table获取column信息
	for i := 0; i < len(TableName); i++ {
		if TableName[i].(string) != "OrionSnapshot" {
			inputTable.UUID = uuid.New().String()
			inputTable.InstanceId = InstanceID
			inputTable.DatabaseName = DatabaseName
			inputTable.TableName = TableName[i].(string)
			// 添加到数据库
			local.InsertTableData(&inputTable)
			// 初始化并遍历Column
			targetDB.Raw("SELECT name FROM SysColumns WHERE id = Object_Id(?)", TableName[i].(string)).Scan(&ColumnName)
			for j := 0; j < len(ColumnName); j++ {
				// 初始化sampleData
				columnSampleData := make([]interface{}, 10)
				inputSampleData := new(local.DataSample)
				querySql := "SELECT DISTINCT TOP 10 " + ColumnName[j].(string) + " FROM [" + inputTable.DatabaseName + "]" + ".dbo." + TableName[i].(string) + " WHERE " + ColumnName[j].(string) + " is not NULL"
				targetDB.Raw(querySql).Scan(&columnSampleData)
				targetDB.Raw("select DATA_TYPE from information_schema.columns where TABLE_NAME = ? and COLUMN_NAME =?", TableName[i].(string), ColumnName[j].(string)).Scan(&columnDataType)
				// 更新数据列, 只有存在数据的列才有意义
				if columnSampleData[0] != nil {
					inputColumn.UUID = uuid.New().String()
					inputColumn.DataType = columnDataType[0].(string)
					inputColumn.InstanceId = InstanceID
					inputColumn.DatabaseName = DatabaseName
					inputColumn.TableName = TableName[i].(string)
					inputColumn.ColumnName = ColumnName[j].(string)
					inputColumn.SampleData = utils.InterfaceToString(columnSampleData[0])
					local.InsertColumn(&inputColumn)
					// 更新数据样例 若数据样例数为10更新
					if len(columnSampleData) == 10 {
						inputSampleData.ColumnUUID = inputColumn.UUID
						inputSampleData.InstanceID = InstanceID
						inputSampleData.DatabaseName = DatabaseName
						inputSampleData.TableName = TableName[i].(string)
						inputSampleData.SampleData0 = utils.InterfaceToString(columnSampleData[0])
						inputSampleData.SampleData1 = utils.InterfaceToString(columnSampleData[1])
						inputSampleData.SampleData2 = utils.InterfaceToString(columnSampleData[2])
						inputSampleData.SampleData3 = utils.InterfaceToString(columnSampleData[3])
						inputSampleData.SampleData4 = utils.InterfaceToString(columnSampleData[4])
						inputSampleData.SampleData5 = utils.InterfaceToString(columnSampleData[5])
						inputSampleData.SampleData6 = utils.InterfaceToString(columnSampleData[6])
						inputSampleData.SampleData7 = utils.InterfaceToString(columnSampleData[7])
						inputSampleData.SampleData8 = utils.InterfaceToString(columnSampleData[8])
						inputSampleData.SampleData9 = utils.InterfaceToString(columnSampleData[9])
						local.InsertSampleData(inputSampleData)
					}
				}
			}
		}
	}
	return Errmsg.SUCCESS, nil
}

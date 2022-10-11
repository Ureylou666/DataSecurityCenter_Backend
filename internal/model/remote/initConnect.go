package remote

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var targetDB *gorm.DB
var err error

func InitConnection(Db string, DbHost string, DbPort string, DbUser string, DbPassword string, DbName string) {
	//var tablename []interface{}
	var column_name []interface{}
	switch Db {
	case "PostgreSQL":
		{
			dsn := "host=" + DbHost + " user=" + DbUser + " password=" + DbPassword + " dbname=" + DbName + " port=" + DbPort + " sslmode=disable"
			targetDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
			})
			//targetDB.Raw("SELECT tablename FROM pg_tables WHERE tablename NOT LIKE 'pg%' AND tablename NOT LIKE 'sql_%'").Scan(&tablename)
			targetDB.Raw("SELECT table_name FROM information_schema.columns").Scan(&column_name)
		}
	case "MySQL":
		{
			dsn := DbUser + ":" + DbPassword + "@tcp(" + DbHost + ":" + DbPort + ")/" + DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
			targetDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
			})
		}
	case "SQLServer":
		{
			dsn := "server=" + DbHost + ";user id=" + DbUser + ";password=" + DbPassword + ";port=" + DbPort + ";database=" + DbName + ";encrypt=disable"
			targetDB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
			})
		}
	}
	if err != nil {
		fmt.Println("连接数据库失败，请检查参数", err)
	}
}

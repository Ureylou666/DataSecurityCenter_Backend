package setting

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	//Server相关
	AppMode  string
	HttpPort string
	//Database相关
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	SSLCert    string
	// 阿里云AKSK
	AKSK Aliyun
	// cnisdp 审计账号
	AccountName        string
	AccountPassword    string
	AccountDescription string
	AccountType        string
)

type Aliyun struct {
	ISDP AliyunKeys
}

type AliyunKeys struct {
	AccessKey    string
	AccessSecret string
}

func init() {
	viper.SetConfigFile("/Users/ureylou/Downloads/golang/complianceCenter/backend/config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("配置文件读取错误：", err)
	}
	fmt.Println(viper.Get("Server.HttpPort"))
	LoadServer()
	LoadDatabase()
}

// 读取服务器配置
func LoadServer() {
	AppMode = viper.GetString("Server.Appmode")
	HttpPort = viper.GetString("Server.HttpPort")
}

// 读取数据库配置
func LoadDatabase() {
	Db = viper.GetString("Database_RDS_Dev.Db")
	DbHost = viper.GetString("Database_RDS_Dev.DbHost")
	DbPort = viper.GetString("Database_RDS_Dev.DbPort")
	DbUser = viper.GetString("Database_RDS_Dev.DbUser")
	DbPassword = viper.GetString("Database_RDS_Dev.DbPassword")
	DbName = viper.GetString("Database_RDS_Dev.DbName")
	SSLCert = viper.GetString("Database_RDS_Dev.SSLCert")
}

func LoadAuditAccount() {
	AccountName = viper.GetString("AuditAccount.AccountName")
	AccountPassword = viper.GetString("AuditAccount.AccountPassword")
	AccountDescription = viper.GetString("AuditAccount.AccountDescription")
	AccountType = viper.GetString("AuditAccount.AccountType")
}

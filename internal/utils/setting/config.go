package setting

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	//Server相关
	AppMode  string
	HttpPort string
	IPAddr   string
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
	// 密钥相关
	Salt string
	Key  string
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

// LoadServer 读取服务器配置
func LoadServer() {
	AppMode = viper.GetString("Server.Appmode")
	HttpPort = viper.GetString("Server.HttpPort")
	IPAddr = viper.GetString("Server.IPAddr")
}

// LoadDatabase 读取数据库配置
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

func LoadEncryption() {
	Salt = viper.GetString("Encryption.Salt")
	Key = viper.GetString("Encryption.Key")
}

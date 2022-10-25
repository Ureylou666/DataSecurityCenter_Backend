package main

import (
	"Backend/internal/model/local"
	"Backend/internal/router"
)

func main() {
	// 引用连接数据库
	local.InitDb()
	router.InitRouter()
	//errcode, err := Basic.UpdateRDS("1214179580276348")
	//errcode, err := Basic.UpdateDBDetails("pgm-uf64xc9s753fcgz6")
	//errcode, err, RDSClient := aliyunSDK.CreateRDSClient("1214179580276348")
	//aliyunSDK.DeleteAccount("pgm-uf64xc9s753fcgz6", RDSClient)
	//errcode, err = aliyunSDK.ModifySecurityIps("pgm-uf64xc9s753fcgz6", "Delete", RDSClient)
	//errcode, err = aliyunSDK.ReleaseInstancePublicConnection("pgm-uf64xc9s753fcgz6", "v4w44tnl7p.pg.rds.aliyuncs.com", RDSClient)
	//fmt.Println(errcode, " ", err)
}

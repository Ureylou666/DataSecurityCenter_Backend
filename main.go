package main

import (
	"Backend/internal/model/local"
	"Backend/internal/router"
)

func main() {
	// 引用连接数据库
	local.InitDb()
	router.InitRouter()
	//errcode, err := Basic.UpdateRDS("1709821521553093")
	//errcode, err, RDSClient := aliyunSDK.CreateRDSClient("")
	//aliyunSDK.DeleteAccount("pgm-uf67y7om9d12does", RDSClient)
	//errcode, err = aliyunSDK.ModifySecurityIps("pgm-uf6fowj28a15797x", "Delete", RDSClient)
	//errcode, err = aliyunSDK.ReleaseInstancePublicConnection("pgm-uf69q84sxp25nvj6", "p2q80scwl1.pg.rds.aliyuncs.com", RDSClient)
}

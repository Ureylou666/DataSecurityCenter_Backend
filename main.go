package main

import (
	"Backend/internal/application/initialize"
	"Backend/internal/model"
)

func main() {
	// 引用连接数据库
	model.InitDb()
	// 重置数据库中数据
	//model.RestoreData()
	// 获取初始化数据
	//aliyunSDDP.Entry("ISDP")
	initialize.InitCloudAccountList()
	//	router.InitRouter()
}

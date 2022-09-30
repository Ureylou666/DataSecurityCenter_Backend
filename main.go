package main

import (
	"Backend/internal/application/Basic"
	"Backend/internal/model"
	"fmt"
)

func main() {
	// 引用连接数据库
	model.InitDb()
	// 重置数据库中数据
	//model.RestoreData()
	// 获取初始化数据
	//aliyunSDDP.Entry("ISDP")
	//Basic.InitCloudAccountList()
	//model.GetCloudAccount("")
	//Basic.UpdateCloudAccountList()
	ErrCode, ErrMsg := Basic.UpdateRDS("")
	fmt.Println(ErrCode, ":  ", ErrMsg)
	//model.DeleteInventory("rm-uf6ys93126n35orln")
	//	router.InitRouter()
}

package main

import (
	"Backend/internal/application/Basic"
	"Backend/internal/model/local"
	"fmt"
)

func main() {
	// 引用连接数据库
	local.InitDb()
	errcode, err := Basic.UpdateRDS("1190073115559253")
	fmt.Println(errcode)
	fmt.Println(err)
}

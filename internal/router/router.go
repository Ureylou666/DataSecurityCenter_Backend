package router

import (
	v1 "Backend/api/v1"
	"Backend/internal/middleware"
	"Backend/internal/utils/setting"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(setting.AppMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	api := r.Group("api/v1")
	// Home Page > Data Security > Identification > Inventory
	api.POST("/Inventory/CloudAccount", v1.GetCloudAccount)
	api.POST("/Inventory/RDSList", v1.GetRDSInventoryList)
	api.POST("/Inventory/Database", v1.GetDatabaseList)
	api.POST("/Inventory/Table", v1.GetTableList)
	api.POST("/Inventory/Column", v1.GetColumns)
	// Home Page > Data Security > Identification > DataField
	api.POST("/DataFields", v1.GetColumnDetails)
	/*
		// 数据安全 > 存储阶段 > 数据资产 Inventory模块
		api.POST("/inventory", v1.GetInventory)
		// 数据安全 > 存储阶段 > 数据资产 Table模块
		api.POST("/tables", v1.GetTables)
		api.POST("/tables/column", v1.GetColumns)
		// 数据安全 > 存储阶段 > 数据清单 Column模块
		api.POST("/column", v1.GetColumnDetails)
		// 数据安全 > 存储阶段 > 分级规则 rules模块
		api.POST("/rules", v1.GetRules)
		// 系统管理 > 项目配置
		api.POST("/group", v1.GetGroup)
		api.POST("/group/create", v1.CreateGroup)
		api.POST("/group/initData", v1.InitGroupData)

	*/
	r.Run(setting.HttpPort)
}

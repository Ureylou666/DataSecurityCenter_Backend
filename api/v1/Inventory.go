package v1

import (
	"Backend/internal/Errmsg"
	"Backend/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取实例
func GetInventory(c *gin.Context) {
	var queryinfo model.InventoryQueryInfo
	err := c.ShouldBindJSON(&queryinfo)
	if (queryinfo.GroupName == "All") { queryinfo.GroupName = "" }
	if (queryinfo.InstanceType == "All") { queryinfo.InstanceType = "" }
	// 判断用户输入 限定InstanceType
	if (err != nil) || ((queryinfo.InstanceType != "") && (queryinfo.InstanceType != "OSS") && (queryinfo.InstanceType != "RDS")) || (queryinfo.PageSize > 50) {
		c.JSON(http.StatusOK, gin.H{
			"Code": Errmsg.ERROR,
			"Message": Errmsg.GetErrMsg(Errmsg.Error_QueryInfo),
		})
		return
	}
	res_data, res_total, inventory_total := model.GetInventory(queryinfo)
	// 未获取到对应数据
	if res_data == nil {
		c.JSON(http.StatusOK, gin.H{
			"Code": Errmsg.Error_NotFound,
			"Message": Errmsg.GetErrMsg(Errmsg.Error_NotFound),
		})
		return
	}
	// 正常恢复
	c.JSON(http.StatusOK, gin.H{
		"Code": Errmsg.SUCCESS,
		"Message": Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"Inventory_Total": inventory_total,
		"Res_Total": res_total,
		"Res_Data": res_data,
	})
}

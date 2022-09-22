package v1

import (
	"Backend/internal/Errmsg"
	"Backend/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTables(c *gin.Context) {
	var queryinfo model.TabelQueryInfo
	err := c.ShouldBindJSON(&queryinfo)
	// 判断用户输入
	if (err != nil) || (queryinfo.PageSize > 50) || (queryinfo.InstanceName == "") {
		c.JSON(http.StatusOK, gin.H{
			"Code": Errmsg.ERROR,
			"Message": Errmsg.GetErrMsg(Errmsg.Error_QueryInfo),
		})
		return
	}
	resData, resTotal, tableTotal := model.GetTables(queryinfo)
	// 未获取到对应数据
	if resData == nil {
		c.JSON(http.StatusOK, gin.H{
			"Code": Errmsg.Error_NotFound,
			"Message": Errmsg.GetErrMsg(Errmsg.Error_NotFound),
		})
		return
	}
	// 正常回复
	c.JSON(http.StatusOK, gin.H{
		"Code":         Errmsg.SUCCESS,
		"Message":      Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"Tables_Total": tableTotal,
		"Res_Total":    resTotal,
		"Res_Data":     resData,
	})
}

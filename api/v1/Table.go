package v1

import (
	"Backend/internal/application/Basic"
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTableList(c *gin.Context) {
	var queryinfo local.TabelQueryInfo
	err := c.ShouldBindJSON(&queryinfo)
	// 判断用户输入 DatabaseName
	if (err != nil) || (queryinfo.PageSize > 50) || queryinfo.DatabaseName == "" {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
		return
	}
	resData, resTotal, TableListTotal := Basic.GetTableList(queryinfo)
	// 未获取到对应数据
	if resData == nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.SUCCESS,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorNotFound),
		})
		return
	}
	// 正常恢复
	c.JSON(http.StatusOK, gin.H{
		"Code":           Errmsg.SUCCESS,
		"Message":        Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"TableListTotal": TableListTotal,
		"Res_Total":      resTotal,
		"Res_Data":       resData,
	})
}

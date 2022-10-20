package v1

import (
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDatabaseList(c *gin.Context) {
	var queryinfo local.DatabaseQueryInfo
	err := c.ShouldBindJSON(&queryinfo)
	// 判断用户输入 AccountID
	if (err != nil) || (queryinfo.PageSize > 50) || queryinfo.InstanceID == "" {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
		return
	}
	resData, resTotal, DatabaseTotal := local.GetDatabaseAPI(queryinfo)
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
		"Code":          Errmsg.SUCCESS,
		"Message":       Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"DatabaseTotal": DatabaseTotal,
		"Res_Total":     resTotal,
		"Res_Data":      resData,
	})
}

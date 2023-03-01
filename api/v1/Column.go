package v1

import (
	"Backend/internal/application/Basic"
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetColumnDetails(c *gin.Context) {
	var queryinfo local.ColumnDetailsQueryInfo
	err := c.ShouldBindJSON(&queryinfo)
	if queryinfo.GroupName == "All" {
		queryinfo.GroupName = ""
	}
	if queryinfo.SensLevelName == "All" {
		queryinfo.SensLevelName = ""
	}
	// 判断用户输入
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ERROR,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
		return
	}
	resData, resTotal, columnTotal := Basic.GetColumnDetails(queryinfo)
	// 未获取到对应数据
	if resData == nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorNotFound,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorNotFound),
		})
		return
	}
	// 正常恢复
	c.JSON(http.StatusOK, gin.H{
		"Code":         Errmsg.SUCCESS,
		"Message":      Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"Column_Total": columnTotal,
		"Res_Total":    resTotal,
		"Res_Data":     resData,
	})
}

// GetColumns 通过table获取column
func GetColumns(c *gin.Context) {
	var queryinfo local.QueryFromTable
	err := c.ShouldBindJSON(&queryinfo)
	// 判断用户输入
	if err != nil || queryinfo.TableName == "" || queryinfo.PageSize > 50 {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
		return
	}
	resData, resTotal, columnTotal := Basic.GetColumnsFromTable(queryinfo)
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
		"Code":         Errmsg.SUCCESS,
		"Message":      Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"Column_Total": columnTotal,
		"Res_Total":    resTotal,
		"Res_Data":     resData,
	})
}

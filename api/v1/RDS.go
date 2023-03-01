package v1

import (
	"Backend/internal/application/Basic"
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetRDSInventoryList 获取云账号列表
func GetRDSInventoryList(c *gin.Context) {
	var queryinfo local.RDSQueryInfo
	err := c.ShouldBindJSON(&queryinfo)
	// 判断用户输入 AccountID
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
		return
	}
	resData, resTotal, rdsInventoryTotal := Basic.GeTRDSInventoryList(queryinfo)
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
		"Code":              Errmsg.SUCCESS,
		"Message":           Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"rdsInventoryTotal": rdsInventoryTotal,
		"Res_Total":         resTotal,
		"Res_Data":          resData,
	})
}

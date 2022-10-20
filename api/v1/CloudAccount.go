package v1

import (
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetCloudAccount 获取云账号列表
func GetCloudAccount(c *gin.Context) {
	var queryinfo local.CloudAccountQueryInfo
	err := c.ShouldBindJSON(&queryinfo)
	if queryinfo.GroupName == "All" {
		queryinfo.GroupName = ""
	}
	// 判断用户输入 限定InstanceType
	if (err != nil) || (queryinfo.PageSize > 50) {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
		return
	}
	resData, resTotal, accountListTotal := local.GetCloudAccountList(queryinfo)
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
		"Code":             Errmsg.SUCCESS,
		"Message":          Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"AccountListTotal": accountListTotal,
		"Res_Total":        resTotal,
		"Res_Data":         resData,
	})
}

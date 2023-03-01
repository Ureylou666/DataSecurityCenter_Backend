package v1

import (
	"Backend/internal/application/Basic"
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// NewCategory - 增
func NewCategory(c *gin.Context) {
	var data local.Category
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
	}
	ErrCode, ErrMessage := Basic.CreateCategory(data)
	// 创建失败
	if ErrMessage != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    ErrCode,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorCreate),
		})
		return
	}
	// 正常响应
	c.JSON(http.StatusOK, gin.H{
		"Code":    Errmsg.SUCCESS,
		"Message": Errmsg.GetErrMsg(Errmsg.SUCCESS),
	})
}

// DeleteCategory - 删
func DeleteCategory(c *gin.Context) {
	var data local.Category
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
	}
	ErrCode, ErrMessage := Basic.DeleteCategory(data)
	// 创建失败
	if ErrMessage != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    ErrCode,
			"Message": Errmsg.GetErrMsg(ErrCode),
		})
		return
	}
	// 正常响应
	c.JSON(http.StatusOK, gin.H{
		"Code":    Errmsg.SUCCESS,
		"Message": Errmsg.GetErrMsg(Errmsg.SUCCESS),
	})
}

// UpdateCategory - 改
func UpdateCategory(c *gin.Context) {
	var data local.Category
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
	}
	ErrCode, ErrMessage := Basic.UpdateCategory(data)
	// 创建失败
	if ErrMessage != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    ErrCode,
			"Message": Errmsg.GetErrMsg(ErrCode),
		})
		return
	}
	// 正常响应
	c.JSON(http.StatusOK, gin.H{
		"Code":    Errmsg.SUCCESS,
		"Message": Errmsg.GetErrMsg(Errmsg.SUCCESS),
	})
}

// ListCategory - 查
func ListCategory(c *gin.Context) {
	var data local.CategoryQuery
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
	}
	ErrCode, ErrMessage, result, resTotal, categoryTotal := Basic.ListCategory(data)
	// 更新失败
	if ErrMessage != "" {
		c.JSON(http.StatusOK, gin.H{
			"Code":    ErrCode,
			"Message": ErrMessage,
		})
		return
	}
	// 正常响应
	c.JSON(http.StatusOK, gin.H{
		"Code":          Errmsg.SUCCESS,
		"Message":       Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"Data":          result,
		"resTotal":      resTotal,
		"categoryTotal": categoryTotal,
	})
}

func ListCategoryName(c *gin.Context) {
	result, resTotal := local.ListCategoryName()
	jsonResult, _ := json.Marshal(result)
	c.JSON(http.StatusOK, gin.H{
		"Code":     Errmsg.SUCCESS,
		"Message":  Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"Data":     json.RawMessage(jsonResult),
		"resTotal": resTotal,
	})
}

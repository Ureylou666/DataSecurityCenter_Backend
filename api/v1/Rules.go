package v1

import (
	"Backend/internal/application/Basic"
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

// NewRule - 增
func NewRule(c *gin.Context) {
	var data local.Rules
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
	}
	ErrCode, ErrMessage := Basic.CreateRules(data)
	// 创建失败
	if ErrMessage != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    ErrCode,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
		return
	}
	// 正常响应
	c.JSON(http.StatusOK, gin.H{
		"Code":    Errmsg.SUCCESS,
		"Message": Errmsg.GetErrMsg(Errmsg.SUCCESS),
	})
}

// DeleteRule - 删
func DeleteRule(c *gin.Context) {
	var data local.Rules
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
	}
	ErrCode, ErrMessage := Basic.DeleteRules(data.UUID)
	// 删除失败
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

// UpdateRule - 改
func UpdateRule(c *gin.Context) {
	var data local.Rules
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
	}
	ErrCode, ErrMessage := Basic.UpdateRules(data)
	// 更新失败
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

// ListRule - 查
func ListRules(c *gin.Context) {
	var data local.RulesListQuery
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
	}
	ErrCode, ErrMessage, result, resTotal, rulesTotal := Basic.ListRules(data)
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
		"Code":       Errmsg.SUCCESS,
		"Message":    Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"Data":       result,
		"resTotal":   resTotal,
		"rulesTotal": rulesTotal,
	})
}

// GetRule 获取特定单一rule
func GetRule(c *gin.Context) {
	input := c.Param("RuleUUID")
	ErrCode, ErrMessage, result := Basic.GetRule(input)
	if ErrMessage != "" {
		c.JSON(http.StatusOK, gin.H{
			"Code":    ErrCode,
			"Message": ErrMessage,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"Code":    ErrCode,
			"Message": ErrMessage,
			"Data":    result,
		})
	}
}

func EnforceRule(c *gin.Context) {
	var data local.Rules
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
	}
	// 输入校验
	validation, _ := regexp.MatchString("[a-f\\d]{8}(-[a-f\\d]{4}){3}-[a-f\\d]{12}", data.UUID)
	if validation {
		validation = local.CheckRuleExist(data.UUID)
	}
	if !validation {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ErrorQueryInput,
			"Message": Errmsg.GetErrMsg(Errmsg.ErrorQueryInput),
		})
		return
	}
	// 部署
	ErrCode, ErrMessage := Basic.EnforceRules(data.UUID)
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

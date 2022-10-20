package v1

/*
import (
	"Backend/internal/application/Basic/aliyunSDDP"
	"Backend/internal/model/local"
	"Backend/internal/utils/Errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InitInfo struct {
	Groupname string
}

// GetGroup 查询Group信息
func GetGroup(c *gin.Context) {
	var queryinfo local.GroupQueryInfo
	err := c.ShouldBindJSON(&queryinfo)
	// 判断用户输入
	if (err != nil) || (queryinfo.PageSize > 50) {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ERROR,
			"Message": Errmsg.GetErrMsg(Errmsg.Error_QueryInfo),
		})
		return
	}
	resData, resTotal, groupTotal := local.GetDevGroup(queryinfo)
	// 未获取到对应数据
	if resData == nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.Error_NotFound,
			"Message": Errmsg.GetErrMsg(Errmsg.Error_NotFound),
		})
		return
	}
	// 正常回复
	c.JSON(http.StatusOK, gin.H{
		"Code":        Errmsg.SUCCESS,
		"Message":     Errmsg.GetErrMsg(Errmsg.SUCCESS),
		"Group_Total": groupTotal,
		"Res_Total":   resTotal,
		"Res_Data":    resData,
	})
}

// CreateGroup 新增项目组 暂不在前端显示！！！
func CreateGroup(c *gin.Context) {
	var insertData local.InsertInfo
	err := c.ShouldBindJSON(&insertData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ERROR,
			"Message": Errmsg.GetErrMsg(Errmsg.ERROR),
		})
		return
	}
	if local.AddDevGroup(insertData) != Errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.Error_InsertInfo,
			"Message": Errmsg.GetErrMsg(Errmsg.Error_InsertInfo),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.SUCCESS,
			"Message": Errmsg.GetErrMsg(Errmsg.SUCCESS),
		})
		return
	}
}

// InitGroupData 初始化Group数据
func InitGroupData(c *gin.Context) {
	var data InitInfo
	err := c.ShouldBindJSON(&data)
	if (err != nil) || (local.CheckGroup(data.Groupname)) {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ERROR,
			"Message": Errmsg.GetErrMsg(Errmsg.ERROR),
		})
		return
	}
	if aliyunSDDP.InitData(data.Groupname) != 200 {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.ERROR_InitError,
			"Message": Errmsg.GetErrMsg(Errmsg.ERROR_InitError),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"Code":    Errmsg.SUCCESS,
			"Message": Errmsg.GetErrMsg(Errmsg.SUCCESS),
		})
		return
	}
}


*/

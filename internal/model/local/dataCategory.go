package local

import (
	"Backend/internal/utils/Errmsg"
)

type Category struct {
	UUID           string `gorm:"type:varchar(50)" json:"UUID"`
	CategoryName   string `gorm:"type:varchar(50)" json:"CategoryName"`
	Level          int    `gorm:"type:int" json:"Level"`
	ParentCategory string `gorm:"type:varchar(50)" json:"ParentCategory"`
	Comments       string `gorm:"type:text" json:"Comments"`
}

type CategoryQuery struct {
	CategoryName string
	PageNum      int
	PageSize     int
}

type CategoryName struct {
	Value string `json:"value"`
	Label string `json:"label"`
	UUID  string `json:"UUID"`
}

// CreateCategory 增 新增分类
func CreateCategory(data *Category) (ErrCode int, ErrMessage error) {
	err := db.Create(&data).Error
	if err != nil {
		return Errmsg.ERROR, err
	}
	return Errmsg.SUCCESS, nil
}

// DeleteCategory 删 删除分类
func DeleteCategory(uuid string) (ErrCode int, ErrMessage error) {
	ErrMessage = db.Where("uuid = ?", uuid).Delete(&Category{}).Error
	if ErrMessage != nil {
		return Errmsg.ERROR, ErrMessage
	}
	return Errmsg.SUCCESS, nil
}

// UpdateCategory 改 更新分类
func UpdateCategory(data *Category) (ErrCode int, ErrMessage error) {
	ErrMessage = db.Save(&data).Error
	if ErrMessage != nil {
		return Errmsg.ERROR, ErrMessage
	}
	return Errmsg.SUCCESS, nil
}

// ListLevel 查 返回当前系统中最大level
func ListLevel() (MaxLevel int) {
	result := db.Find(&Category{})
	MaxLevel = 0
	if result.RowsAffected == 0 {
		return -1
	}
	if result.RowsAffected > 0 {
		db.Raw("select max(level) from Category").Scan(&MaxLevel)
	}
	return MaxLevel
}

// ListLevelCategory 查 通过level 进行查找
func ListLevelCategory(CategoryLevel string) (ErrCode int, ErrMessage error, CategoryList []Category) {
	ErrMessage = db.Where("level = ?", CategoryLevel).Scan(CategoryList).Error
	if ErrMessage != nil {
		return Errmsg.ERROR, ErrMessage, nil
	}
	return Errmsg.SUCCESS, nil, CategoryList
}

// CheckCategoryToLevel 校验 新建/更新时 若存在parentCate 进行校验 确保level和Category匹配
func CheckCategoryToLevel(ParentCategory string, level int) bool {
	var temp []Category
	if (level == 0) && (ParentCategory == "") {
		return true
	}
	if (level == 0) && (ParentCategory != "") {
		return false
	}
	if (level != 0) && (ParentCategory == "") {
		return false
	}
	if (level != 0) && (ParentCategory != "") {
		db.Where("UUID = ?", ParentCategory).Find(&temp)
		// UUID 不存在 或 不是唯一
		if len(temp) != 1 {
			return false
		}
		// 对比等级
		if (temp[0].UUID != "") && (temp[0].Level == level-1) {
			return true
		}
	}
	// 其他情况
	return false
}

func CheckCategoryNameUnique(CategoryName string) bool {
	var temp []Category
	db.Where("category_name = ?", CategoryName).Find(&temp)
	if len(temp) > 0 {
		return false
	}
	return true
}

func CheckCategoryExist(input string) bool {
	var temp []Category
	db.Where("uuid = ?", input).Find(&temp)
	if len(temp) == 0 {
		return false
	}
	return true
}

func CategoryUUIDtoName(categoryUUID string) (CategoryName string) {
	var temp []Category
	db.Where("uuid = ?", categoryUUID).Find(&temp)
	return temp[0].CategoryName
}

// ListCategory 查 查询
func ListCategory(input CategoryQuery) (result []Category, resTotal int64, categoryTotal int64) {
	input.CategoryName = "%" + input.CategoryName + "%"
	// 分页处理
	db.Where("category_name like ?", input.CategoryName).Find(&result).Count(&categoryTotal)
	if input.PageNum == 0 || input.PageSize == 0 {
		db.Where("category_name like ?", input.CategoryName).Limit(-1).Find(&result)
		resTotal = int64(len(result))
	} else {
		db.Where("category_name like ?", input.CategoryName).Limit(input.PageSize).Offset((input.PageNum - 1) * input.PageSize).Find(&result)
		resTotal = int64(len(result))
	}
	return result, resTotal, categoryTotal
}

func ListCategoryName() (result []map[string]interface{}, resTotal int64) {
	var temp []Category
	db.Find(&temp)
	resTotal = int64(len(temp))
	var m1 map[string]interface{}
	for i := 0; i < int(resTotal); i++ {
		m1 = make(map[string]interface{})
		m1["Value"] = temp[i].CategoryName
		m1["Label"] = temp[i].CategoryName
		m1["UUID"] = temp[i].UUID
		result = append(result, m1)
	}
	return result, resTotal
}

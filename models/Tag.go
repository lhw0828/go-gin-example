package models

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
	DeletedOn  int    `json:"deleted_on"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func GatTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return

}

// ExistTagByName 是否存在同名的标签
func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name =?", name).First(&tag)

	if tag.ID > 0 {
		return true
	}

	return false

}

// AddTag 添加标签
func AddTag(name string, state int, createdBy string) bool {
	tag := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}

	if err := db.Create(&tag).Error; err != nil {
		return false
	}

	return true
}

func ExistTagById(id int) bool {
	var tag Tag
	db.Select("id").Where("id =?", id).First(&tag)

	if tag.ID > 0 {
		return true
	}

	return false
}

func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})

	return true
}

func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id =?", id).Updates(data)

	return true
}

// CleanAllTag 清除所有
func CleanAllTag() bool {
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{})

	return true
}

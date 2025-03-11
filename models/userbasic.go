package models

import (
	"im-master/utils"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string
	Email         string
	ClientIP      string
	ClientPost    string
	LoginTime     uint64
	HeartBeatTime uint64
	LogOutTime    uint64
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

// GetUserList 返回所有用户列表
func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	return data
}

// CreateUser 创建新用户
func CreateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Create(user)
}

// DeleteUser 删除用户
func DeleteUser(user *UserBasic) *gorm.DB {
	return utils.DB.Delete(user)
}

// UpdateUser 更新用户信息
func UpdateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Model(user).Updates(UserBasic{Name: user.Name, Password: user.Password})
}

// FindUserByName 通过用户名查找用户
func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

// FindUserByPhone 通过手机号查找用户
func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

// FindUserByEmail 通过邮箱查找用户
func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

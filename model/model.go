package model

import (
	"errors"
	"fmt"
	"software_bottle/dao"
)

const (
	DefaultImg = "default.jpg"
	IMAG_PATH  = "图片存放的路径"
)

type User struct {
	UserId    int    `json:"user_id" gorm:"primary_key;AUTO_INCREMENT"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
	Sex       string `json:"sex"`
	ImagePath string `json:"image_path"`
}

type Friend struct {
	UserId   int `json:"user_id"`
	FriendId int `json:"friend_id"`
}

func CreateAUser(u *User) error {
	return dao.DB.Create(u).Error
}

func GetAUser(u *User, username string) error {
	return dao.DB.Where("user_name=?", username).First(u).Error
}

func GetAUserById(u *User, userId int) error {
	return dao.DB.Where("user_id=?", userId).First(u).Error
}

func GetAllUser(u []*User) error {
	return dao.DB.Find(u).Error
}

func UpdateUser(u *User) error {
	return dao.DB.Save(u).Error
}

func DeleteAUser(userId int) error {
	return dao.DB.Where("user_name=?", userId).Delete(User{}).Error
}

func UserIsExists(username string) error {
	err := dao.DB.Where("user_name=?", username).Take(&User{}).Error
	if err == nil {
		return errors.New(fmt.Sprintf("The user already exists with UserName:%s", username))
	} else {
		return nil
	}
}

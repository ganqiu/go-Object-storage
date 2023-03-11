// Package User
// CREATED 2023-3-9
// LAST-MODIFIED 2023-3-10
// CREATOR Junixor
// VERSION 2.0
// 包含各种用户操作
// 创建数据库表：
// CREATE TABLE users (
//   id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
//   created_at TIMESTAMP,
//   updated_at TIMESTAMP,
//   deleted_at TIMESTAMP,
//   uuid VARCHAR(40) UNIQUE,
//   name VARCHAR(20) UNIQUE,

//   pass_word VARCHAR(30) NOT NULL
// );

package User

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// UUIDError : 错误类型，生成 uuid 失败错误
var UUIDError = errors.New("failed to generate uuid")

// WrongPassWord : 错误类型，密码与数据库中密码不符，无法获得权限
var WrongPassWord = errors.New("wrong password")

// User 用户信息
type User struct {
	gorm.Model
	Uuid     string `gorm:"UNIQUE"`   // 唯一标识
	Name     string `gorm:"UNIQUE"`   // 用户名
	PassWord string `gorm:"NOT NULL"` // 密码
}

// generateUUID 生成唯一标识
// 若错误，返回空 uuid 和错误信息，否则返回生成的 uuid 和空
func generateUUID() (string, error) {
	ID, err := uuid.NewV7()
	if err != nil {
		return "", UUIDError
	}
	key := ID.String()
	return key, nil
}

// CreateUser 创建用户
// db 数据库指针，userName 用户名，passWord 密码
// 返回一个创建好的用户结构体和错误信息
// INSERT INTO users(uuid, created_at, updated_at, deleted_at, name, pass_word) VALUES(...)
func CreateUser(db *gorm.DB, userName string, passWord string) (User, error) {
	id, _ := generateUUID()
	user := User{
		Uuid:     id,
		Name:     userName,
		PassWord: passWord,
	}
	log := db.Create(&user)
	return user, log.Error
}

// UpdateUserName 更新用户名
// db 数据库指针，name 原用户名，newName 新用户名
// 返回新用户信息结构体和错误信息
// UPDATE users SET name="newName" WHERE name="name" and passWord="passWord"
func UpdateUserName(db *gorm.DB, userName string, passWord string, newName string) (User, error) {
	_, err1 := ReadUser(db, userName, passWord)
	if err1 == WrongPassWord {
		fmt.Println("No access: WrongPassWord")
		return User{}, WrongPassWord
	}
	log := db.Model(&User{}).Where("name = ?", userName).Update("name", newName)
	if log.Error != nil {
		fmt.Println(log.Error)
		return User{}, log.Error
	}
	fmt.Printf("Successfully updated userName")
	user, err2 := ReadUser(db, newName, passWord)
	return user, err2
}

// UpdatePassWord 更新密码
// db 数据库指针，name 用户名，newPassWord 新密码
// 返回新用户信息结构体和错误信息
// UPDATE users SET pass_word="newPassWord" WHERE name="name" passWord="passWord"
func UpdatePassWord(db *gorm.DB, userName string, oldPassWord, newPassWord string) (User, error) {
	_, err1 := ReadUser(db, userName, oldPassWord)
	if err1 == WrongPassWord {
		fmt.Println("No access: WrongPassWord")
		return User{}, WrongPassWord
	}
	log := db.Model(&User{}).Where("name = ?", userName).Update("pass_word", newPassWord)
	if log.Error != nil {
		fmt.Println(log.Error)
		return User{}, log.Error
	}
	fmt.Printf("Successfully updated passWord")
	user, err2 := ReadUser(db, userName, newPassWord)
	return user, err2
}

// ReadUser 获取用户信息
// db 数据库指针，name 用户名, passWord 密码
// 返回用户信息结构体和错误信息
// SELECT * FROM users WHERE name="name" and pass_word="passWord"
func ReadUser(db *gorm.DB, userName string, passWord string) (User, error) {
	var result User
	log := db.Where("name = ?", userName).First(&result)
	if passWord != result.PassWord {
		return User{}, WrongPassWord
	}
	return result, log.Error
}

// DeleteUser 删除用户
// db 数据库指针，name 用户名，passWord 密码
// 返回错误信息
// DELETE FROM users WHERE name="name" and pass_word="passWord"
func DeleteUser(db *gorm.DB, userName string, passWord string) error {
	_, err := ReadUser(db, userName, passWord)
	if err == WrongPassWord {
		fmt.Println("No access: WrongPassWord")
		return WrongPassWord
	}
	log := db.Where("name = ?", userName).Delete(&User{})
	return log.Error
}

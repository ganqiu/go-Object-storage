// Package Object
// CREATED 2023-3-9
// LAST-MODIFIED 2023-3-11
// CREATOR Junixor
// VERSION 1.0
// 包含各种对象操作
// 创建对象表：
// CREATE TABLE objects (
//   id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
//   created_at TIMESTAMP,
//   updated_at TIMESTAMP,
//   deleted_at TIMESTAMP,
//   uuid VARCHAR(40) UNIQUE,
//   name VARCHAR(20) NOT NULL,
//   user_id INT UNSIGNED NOT NULL,
//   bucket_id INT UNSIGNED NOT NULL,
//   CONSTRAINT object_user FOREIGN KEY(user_id) REFERENCES users(id),
//   CONSTRAINT object_bucket FOREIGN KEY(bucket_id) REFERENCES buckets(id)
// );

package Object

import (
	"GoOSS/Bucket"
	"GoOSS/User"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Object struct {
	gorm.DB
	Uuid     string  `gorm:"unique"`
	Name     string  `gorm:"NOT NULL"`                                      // 对象名称
	UserId   uint    `gorm:"constraint:OnUpdate:CASCADE, OnDelete:DELETE;"` // 所属创建者
	BucketId uint    `gorm:"constraint:OnUpdate:CASCADE, OnDelete:DELETE;"` // 所属桶
	Size     string  `gorm:"NOT NULL"`                                      // 对象大小
	Version  float32 `gorm:"NOT NULL"`                                      // 对象版本号
}

var InvalidChoice = errors.New("invalid choice")
var UserCancelled = errors.New("cancelled by user")

// generateUUID 创建 uuid 唯一标识
func generateUUID() (string, error) {
	ID, err := uuid.NewV7()
	if err != nil {
		return "", User.UUIDError
	}
	key := ID.String()
	return key, nil
}

func CreateObject(db *gorm.DB, userName string, passWord string, bucketName string, objectName string, size string, version float32) (Object, error) {
	var choice byte
	user, err1 := User.ReadUser(db, userName, passWord)
	if err1 != nil {
		fmt.Printf("Failed to create object: %v", err1)
		return Object{}, err1
	}
	bucket, err2 := Bucket.ReadBucket(db, userName, passWord, bucketName)
	if err2 != nil {
		fmt.Printf("Failed to create object: %v", err2)
	}
	if bucket.StorageType == Bucket.ColdArchive {
		_, _ = fmt.Scanf("Mode ColdArchive: Need Admin Interference?(y/N)%v", choice)
		switch choice {
		case 'y':
		case 'Y':
			fmt.Printf("Redirecting to manual processing.")
			return Object{}, Bucket.ManualHandle
		case 'n':
		case 'N':
			fmt.Printf("Operation Canceled.")
			return Object{}, UserCancelled
		default:
			fmt.Printf("Error: %v", InvalidChoice)
			return Object{}, InvalidChoice
		}
	}
	id, err3 := generateUUID()
	if err3 != nil {
		fmt.Printf("Failed to create object: %v", err3)
	}
	object := Object{
		Uuid:     id,
		Name:     objectName,
		UserId:   user.ID,
		BucketId: bucket.ID,
		Size:     size,
		Version:  version,
	}
	log := db.Create(&object)
	return object, log.Error
}

func ReadObject(db *gorm.DB, userName string, passWord string, bucketName string, objectName string) (Object, error) {
	var result Object
	user, err1 := User.ReadUser(db, userName, passWord)
	if err1 != nil {
		fmt.Printf("Failed to read object: %v", err1)
		return Object{}, err1
	}
	bucket, err2 := Bucket.ReadBucket(db, userName, passWord, bucketName)
	if err2 != nil {
		fmt.Printf("Failed to read object: %v", err2)
		return Object{}, err2
	}
	log := db.Where("name = ? AND user_id = ? AND bucket_id = ?", objectName, user.ID, bucket.ID).Find(&result)
	return result, log.Error
}

func UpdateObject(db *gorm.DB, userName string, passWord string, bucketName string, objectName string)

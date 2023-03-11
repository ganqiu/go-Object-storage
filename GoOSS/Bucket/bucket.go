// Package Bucket
// CREATED 2023-3-9
// LAST_MODIFIED 2023-3-11
// CREATOR Junixor
// VERSION 1.0
// 包含桶的各种操作
// 创建桶表：
// CREATE TABLE buckets (
//	 id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
//   created_at TIMESTAMP,
//   updated_at TIMESTAMP,
//   deleted_at TIMESTAMP,
//   uuid VARCHAR(40) UNIQUE,
//   name VARCHAR(20) NOT NULL,
//   version FLOAT NOT NULL,
//   user_id INT UNSIGNED NOT NULL,
//   storage_type TINYINT NOT NULL,
//   CONSTRAINT bucket_user FOREIGN KEY(user_id) references users(id)
// );

package Bucket

import (
	"GoOSS/User"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

// DuplicatedName 名称重复错误
// 由于所有用户的桶都在同一张表中，为确保用户桶名称不重复引入
var DuplicatedName = errors.New("duplicated name")

// ManualHandle 需要管理员手动介入
var ManualHandle = errors.New("need manual interference")

// 存储模式
const (
	_           uint8 = iota
	Standard          // 标准存储
	Infrequent        // 低频存储
	Archive           // 归档存储
	ColdArchive       // 冷归档存储
)

// Bucket 对象桶
type Bucket struct {
	gorm.Model
	Uuid        string  `gorm:"unique"`
	Name        string  `gorm:"NOT NULL"`                                     // 桶名称
	Version     float32 `gorm:"NOT NULL"`                                     // 版本号
	UserId      uint    `gorm:"constraint:OnUpdate:CASCADE, OnDelete:DELETE"` // 创建者
	StorageType uint8   `gorm:"NOT NULL"`                                     // 存储模式
}

// generateUUID 创建 uuid 唯一标识
// 返回一个 v7 版本的 uuid 和错误信息
func generateUUID() (string, error) {
	ID, err := uuid.NewV7()
	if err != nil {
		return "", User.UUIDError
	}
	key := ID.String()
	return key, nil
}

// CreateBucket 创建桶
// db 数据库指针，userName 用户名，passWord 密码，bucketName 对象桶名称，version 版本号，storageType 存储类型
// 若桶名已存在，则返回名称重复错误
// 返回一个包含对象桶信息的结构体和错误信息
func CreateBucket(db *gorm.DB, userName string, passWord string, bucketName string, version float32, storageType uint8) (Bucket, error) {
	user, err1 := User.ReadUser(db, userName, passWord)
	if err1 != nil {
		fmt.Printf("Failed to create bucket: %v", err1)
		return Bucket{}, err1
	}
	temp, _ := ReadBucket(db, userName, passWord, bucketName)
	if temp.Name != "" {
		fmt.Printf("Failed to create bucket: %v", DuplicatedName)
		return Bucket{}, DuplicatedName
	}
	id, err3 := generateUUID()
	if err3 != nil {
		fmt.Printf("Failed to create bucket: %v", err3)
		return Bucket{}, err3
	}
	bucket := Bucket{
		Uuid:        id,
		Name:        bucketName,
		Version:     version,
		UserId:      user.ID,
		StorageType: storageType,
	}
	log := db.Create(&bucket)
	if log.Error != nil {
		fmt.Printf("Failed to create bucket: %v", log.Error)
		return Bucket{}, log.Error
	}
	fmt.Printf("Successfully created bucket {name: %v, version: %v, storage_type: %v} at %v", bucket.Name, bucket.Version, bucket.StorageType, bucket.CreatedAt)
	return bucket, nil
}

// FindAllBuckets 获取用户创建的所有桶
// db 数据库指针，userName 用户名，passWord 密码
// 返回用户创建的所有桶信息和错误信息
func FindAllBuckets(db *gorm.DB, userName string, passWord string) (Bucket, error) {
	user, err := User.ReadUser(db, userName, passWord)
	if err != nil {
		fmt.Printf("No access: %v", err)
		return Bucket{}, err
	}
	var result Bucket
	log := db.Where("user_id = ?", user.ID).Find(&result)
	if log.Error != nil {
		fmt.Printf("Failed to find results related to %v: %v", userName, log.Error)
		return Bucket{}, log.Error
	}
	fmt.Printf("Successfully Found %v results related to %v at %v", log.RowsAffected, userName, time.Now())
	return result, nil
}

// ReadBucket 读取对象桶信息
// db 数据库指针，userName 用户名，passWord 密码，bucketName 桶名称
// 返回符合桶名称为 bucketName 且用户名为 userName 的桶的信息
// 返回所有符合条件的桶和错误信息
func ReadBucket(db *gorm.DB, userName string, passWord string, bucketName string) (Bucket, error) {
	var result Bucket
	user, err := User.ReadUser(db, userName, passWord)
	if err != nil {
		fmt.Printf("Failed to read bucket: %v", err)
		return Bucket{}, err
	}
	log := db.Where("name = ? AND user_id = ?", bucketName, user.ID).First(&result)
	if log.Error != nil {
		fmt.Printf("Failed to find results related to %v : %v", bucketName, log.Error)
		return Bucket{}, log.Error
	}
	fmt.Printf("Successfully Found %v results related to %v at %v", log.RowsAffected, bucketName, time.Now())
	return result, nil
}

// UpdateBucketName 修改桶名称
// db 数据库指针，userName 用户名，passWord 密码，bucketName 桶名称，newName 桶的新名称，newVersion，桶的新版本信息
// 更新桶的名称和版本信息，返回更新后的桶信息结构体和错误信息
func UpdateBucketName(db *gorm.DB, userName string, passWord string, bucketName string, newName string, newVersion float32) (Bucket, error) {
	user, err1 := User.ReadUser(db, userName, passWord)
	if err1 != nil {
		fmt.Printf("Failed to read bucket: %v", err1)
		return Bucket{}, err1
	}
	temp, _ := ReadBucket(db, userName, passWord, bucketName)
	if temp.Name != "" {
		fmt.Printf("Failed to create bucket: %v", DuplicatedName)
		return Bucket{}, DuplicatedName
	}
	log := db.Model(&Bucket{}).Where("name = ? AND user_id = ?", bucketName, user.ID).Updates(Bucket{Name: newName, Version: newVersion})
	if log.Error != nil {
		fmt.Println(log.Error)
		return Bucket{}, log.Error
	}
	bucket, _ := ReadBucket(db, userName, passWord, newName)
	fmt.Printf("Successfully updated bucketName to %v at %v", newName, bucket.UpdatedAt)
	return bucket, nil
}

// UpdateBucketStorageType 更新桶存储模式
// db 数据库指针，userName 用户名，passWord 密码，bucketName 桶名称，newStorageType 新存储模式，newVersion 新版本号
// 更新桶的存储模式和八本信息，返回更新后的桶信息结构体和错误信息
func UpdateBucketStorageType(db *gorm.DB, userName string, passWord string, bucketName string, newStorageType uint8, newVersion float32) (Bucket, error) {
	user, err1 := User.ReadUser(db, userName, passWord)
	if err1 != nil {
		fmt.Printf("Failed to read bucket: %v", err1)
		return Bucket{}, err1
	}
	temp, err2 := ReadBucket(db, userName, passWord, bucketName)
	if err2 != nil {
		fmt.Printf("Failed to read bucket: %v", err2)
		return Bucket{}, err2
	}
	if temp.StorageType == ColdArchive {
		fmt.Printf("Mode ColdArchive: Redirecting to manual processing.")
		return Bucket{}, ManualHandle
	}
	log := db.Model(&Bucket{}).Where("name = ? AND user_id = ?", bucketName, user.ID).Updates(Bucket{StorageType: newStorageType, Version: newVersion})
	if log.Error != nil {
		fmt.Println(log.Error)
		return Bucket{}, log.Error
	}
	bucket, _ := ReadBucket(db, userName, passWord, bucketName)
	fmt.Printf("Successfully updated bucket storageType to %v at %v", newStorageType, bucket.UpdatedAt)
	return bucket, nil
}

// DeleteBucket 删除桶
// db 数据库指针，userName 用户名，passWord 密码，bucketName 桶名称
// 删除名称为 bucketName 且 所属用户名称为 userName 的桶，返回错误信息。
func DeleteBucket(db *gorm.DB, userName string, passWord string, bucketName string) error {
	_, err := User.ReadUser(db, userName, passWord)
	if err != nil {
		fmt.Printf("Failed to delete bucket: %v", err)
		return err
	}
	log := db.Where("name = ?", bucketName).Delete(&Bucket{})
	if log.Error != nil {
		fmt.Printf("Failed to delete bucket: %v", log.Error)
		return log.Error
	}
	fmt.Printf("Successfully deleted bucket %v at %v", bucketName, time.Now())
	return log.Error
}


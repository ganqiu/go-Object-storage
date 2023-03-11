package main

import (
	"GoOSS/Bucket"
	"GoOSS/User"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func testUser() {
	dbName := "test"
	db := Connect(dbName)
	userName := "junixor"
	passWord := "123456"
	user, err1 := User.CreateUser(db, userName, passWord)
	fmt.Println(err1)
	result1, err2 := User.ReadUser(db, user.Name, user.PassWord)
	fmt.Println(result1, err2)
	updatedUserName, err3 := User.UpdateUserName(db, userName, user.PassWord, "DuskyJuniXor")
	fmt.Println(updatedUserName, err3)
	updatedPassword, err4 := User.UpdatePassWord(db, updatedUserName.Name, updatedUserName.PassWord, "654321")
	fmt.Println(updatedPassword, err4)
	result2, err5 := User.ReadUser(db, updatedUserName.Name, updatedPassword.PassWord)
	fmt.Println(result2, err5)
	err6 := User.DeleteUser(db, updatedUserName.Name, updatedPassword.PassWord)
	fmt.Println(err6)
	return
}

func testBucket() {
	dbName := "test"
	db := Connect(dbName)
	userName := "junixor"
	passWord := "123456"
	user, err1 := User.CreateUser(db, userName, passWord)
	fmt.Println(err1)
	result1, err2 := User.ReadUser(db, user.Name, user.PassWord)
	fmt.Println(result1, err2)
	bucketName := "bucket1"
	version := float32(0.0)
	storageType := Bucket.Standard
	bucket, err1 := Bucket.CreateBucket(db, userName, passWord, bucketName, version, storageType)
	fmt.Println(bucket, err1)
	bucket2, err2 := Bucket.UpdateBucketName(db, userName, passWord, "bucket1", "bucket2", 0.1)
	fmt.Println(bucket2, err2)
	bucket3, err3 := Bucket.UpdateBucketStorageType(db, userName, passWord, "bucket2", Bucket.ColdArchive, 0.2)
	fmt.Println(bucket3, err3)
	bucket4, err4 := Bucket.UpdateBucketStorageType(db, userName, passWord, "bucket2", Bucket.Standard, 0.3)
	fmt.Println(bucket4, err4)
	bucket5, err5 := Bucket.CreateBucket(db, userName, passWord, "bucket3", 0.0, storageType)
	fmt.Println(bucket5, err5)
	result, err6 := Bucket.FindAllBuckets(db, userName, passWord)
	fmt.Println(result, err6)
	return
}

// Connect 连接到本地 mysql 数据库
// 返回连接的数据库指针
func Connect(dbName string) *gorm.DB {
	dsn := "root:DY-20050206-dkj@tcp(127.0.0.1:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database, error: " + err.Error())
		return nil
	}
	return gormDB
}

func main() {
	testUser()
	testBucket()
	return
}

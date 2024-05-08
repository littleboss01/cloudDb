package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

var (
	Db *gorm.DB
)

// 创建数据库
func CreateDb() {
	fmt.Println("创建数据库")
	//创建数据库
	Db.Exec("CREATE DATABASE IF NOT EXISTS cloudDB")

}

// 连接数据库
func Connect() {
	fmt.Println("连接数据库")
	dsn := "root:root@tcp(localhost:3306)/?charset=utf8&parseTime=True&loc=Local"
	gorm.Open("mysql", dsn)
	//设置数据库连接池
	//gorm.DB().SetMaxIdleConns(10)
}

// 创建表
func DbCreateTable(tableName string) error {
	fmt.Println("创建表")
	//创建表 默认的有id,created_at,updated_at,deleted_at,projectThis
	err := Db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (id int(11) NOT NULL AUTO_INCREMENT,created_at datetime DEFAULT NULL,updated_at datetime DEFAULT NULL,deleted_at datetime DEFAULT NULL,projectThis varchar(255) DEFAULT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	if err != nil {
		fmt.Println(err)
		return err.Error
	}
	return nil
}

func DbAddUserColumn(tableName, columnName string) error {
	fmt.Println("添加列")
	//添加列
	err := Db.Exec("ALTER TABLE " + tableName + " ADD " + columnName + " varchar(255) DEFAULT NULL;")
	if err != nil {
		fmt.Println(err)
		return err.Error
	}
	return nil
}

func DbAddData(tableName string, dataList []string) error {

	err := Db.Exec("INSERT INTO " + tableName + " (projectThis) VALUES ('" + strings.Join(dataList, "','") + "')")
	if err != nil {
		fmt.Println(err)
		return err.Error
	}
	return nil

}

func DbGetData(tableName, project string) ([]string, error) {
	db := Db.Exec("SELECT * FROM "+tableName+" WHERE projectThis = ?", project)
	if db != nil {
		fmt.Println(db)
		return nil, db.Error
	}
	//db取出数据组成数组返回

	//标记project为used
	db = db.Exec("UPDATE "+tableName+" SET projectThis = 'used' WHERE projectThis = ?", project)
	return nil, db.Error
}

func DbUpdate(tableName, id string, dataList []string) error {
	err := Db.Exec("UPDATE "+tableName+" SET "+strings.Join(dataList, ",")+" WHERE id = ?", id)
	if err != nil {
		fmt.Println(err)
		return err.Error
	}
	return nil
}

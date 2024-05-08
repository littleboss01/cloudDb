package main

import (
	//gin
	"github.com/gin-gonic/gin"
	"strings"
)

var (
	//表名缓存
	tableNameCache = make(map[string]bool)
)

// 创建表
func CreateTable(c *gin.Context) {
	tabalName := c.Query("tableName")
	if tabalName == "" {
		c.JSON(200, gin.H{
			"message": "表名不能为空",
		})
		return
	}
	//创建表
	if DbCreateTable(tabalName) == nil {
		c.JSON(200, gin.H{
			"message": "创建表成功",
		})

	} else {
		c.JSON(200, gin.H{
			"message": "创建表失败",
		})
	}
}

// 添加使用者列
func AddUserColumn(c *gin.Context) {
	tabalName := c.Query("tableName")
	//查询表名是否在tableNameCache中
	if _, ok := tableNameCache[tabalName]; !ok {

	}
	//添加列
	columnName := c.Query("columnName")
	if DbAddUserColumn(tabalName, columnName) == nil {
		c.JSON(200, gin.H{
			"message": "添加列成功",
		})

	} else {
		c.JSON(200, gin.H{
			"message": "添加列失败",
		})
	}
}

// 添加数据
func AddData(c *gin.Context) {
	tabalName := c.Query("tableName")
	//查询表名是否在tableNameCache中
	if _, ok := tableNameCache[tabalName]; !ok {

	}
	//添加数据
	data := c.Query("data")
	dataList := strings.Split(data, "----")
	if DbAddData(tabalName, dataList) == nil {
		c.JSON(200, gin.H{
			"message": "添加数据成功",
		})

	} else {
		c.JSON(200, gin.H{
			"message": "添加数据失败",
		})
	}
}

// 取出数据
func GetData(c *gin.Context) {
	tabalName := c.Query("tableName")
	Project := c.Query("project")
	//查询表名是否在tableNameCache中
	if _, ok := tableNameCache[tabalName]; !ok {

	}
	//取出数据
	dataList, err := DbGetData(tabalName, Project)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "取出数据失败",
		})
		return

	}
	c.JSON(200, gin.H{
		"message": "取出数据成功",
		"data":    dataList,
	})
}

// 修改数据
func UpdateData(c *gin.Context) {
	tabalName := c.Query("tableName")
	id := c.Query("id")
	data := c.Query("data")
	//查询表名是否在tableNameCache中
	if _, ok := tableNameCache[tabalName]; !ok {

	}
	dataList := strings.Split(data, "----")
	//修改数据
	err := DbUpdate(tabalName, id, dataList)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "修改数据失败",
		})
		return

	}
	c.JSON(200, gin.H{
		"message": "修改数据成功",
		"data":    dataList,
	})
}

// 备份数据库发送到邮箱
func BackupToEmail(c *gin.Context) {

}

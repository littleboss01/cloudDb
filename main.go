package main

import (
	"github.com/gin-gonic/gin"
	"main/cloudDb"
)

var (
	gToken string
)

func main() {
	dsn, DataName, ginPort, _token := cloudDb.SetDbConfig()
	gToken = _token
	cloudDb.Connect(dsn, DataName)

	//注册api
	e := gin.Default()
	e.Use(cloudDb.CustomMiddleware(gToken))

	e.GET("/CreateTable", cloudDb.CreateTable)
	e.GET("/AddUserColumn", cloudDb.AddUserColumn)
	e.GET("/AddData", cloudDb.AddData)
	e.POST("/AddDataJson", cloudDb.AddDataJson)
	e.GET("/GetData", cloudDb.GetData)
	e.GET("/UpdateData", cloudDb.UpdateData)
	e.GET("/BackupToEmail", cloudDb.BackupToEmail)

	//运行
	e.Run("127.0.0.1:" + ginPort)
}

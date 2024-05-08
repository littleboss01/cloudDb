package cloudDb

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
)

// 生成默认配置文件
func GenerateConfig() {
	//检测同目录下是否有config.ini文件
	//如果没有则生成一个
	_, err := ini.Load("config.ini")
	if err == nil {
		return
	}
	cfg := ini.Empty()
	//生成默认配置
	cfg.Section("database").Key("host").SetValue("localhost")
	cfg.Section("database").Key("port").SetValue("3306")
	cfg.Section("database").Key("user").SetValue("root")
	cfg.Section("database").Key("password").SetValue("root")
	cfg.Section("database").Key("DataName").SetValue("cloudDb")

	cfg.Section("gin").Key("port").SetValue("8080")
	cfg.Section("gin").Key("release").SetValue("false")
	cfg.Section("gin").Key("token").SetValue("")
	cfg.SaveTo("config.ini")

}

// 设置数据库配置
func SetDbConfig() (dns, DataName, ginPort, token string) {
	GenerateConfig()
	//从config.ini读取数据库地址端口账号密码
	//然后设置到Db
	config, err := ini.Load("config.ini")
	if err != nil {
		panic(err)
	}
	dbHost := config.Section("database").Key("host").String()
	dbPort := config.Section("database").Key("port").String()
	dbUser := config.Section("database").Key("user").String()
	dbPassword := config.Section("database").Key("password").String()
	DataName = config.Section("database").Key("DataName").String()
	dns = dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + DataName + "?charset=utf8mb4&parseTime=True&loc=Local"

	// 读取 Gin 配置
	ginPort = config.Section("gin").Key("port").String()
	isRelease, _ := config.Section("gin").Key("release").Bool()
	token = config.Section("gin").Key("token").String()
	if isRelease {
		gin.SetMode(gin.ReleaseMode)
	}
	return dns, DataName, ginPort, token
}

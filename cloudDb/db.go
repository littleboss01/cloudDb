package cloudDb

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gorm.io/driver/mysql"
	"strings"
)

var (
	Db *gorm.DB
)
var (
	//表名缓存
	tableNameCache = make(map[string]bool)

	//操作map的锁
	lock = new(sync.Mutex)
)

// 查询表名是否在tableNameCache中
func IsTableExist(tableName string) bool {
	if _, ok := tableNameCache[tableName]; !ok {
		//在锁里添加表,并更新到map
		lock.Lock()
		if DbCreateTable(tableName) == nil {
			tableNameCache[tableName] = true
		}
		lock.Unlock()
	}
	return tableNameCache[tableName]
}

// 创建数据库
func CreateDb(DataName string) error {
	fmt.Println("创建数据库")
	//创建数据库,如果存在则不创建
	err := Db.Exec("CREATE DATABASE IF NOT EXISTS " + DataName)
	if err != nil {
		fmt.Println(err)
		return err.Error
	}

	return nil
}

// 选择库
func SelectDb(DataName string) error {
	fmt.Println("选择数据库")
	//选择数据库
	err := Db.Exec("USE " + DataName).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Println("选择数据库成功")
	return nil
}

// 连接数据库
func Connect(dsn, dbName string) {
	fmt.Println("连接数据库")
	var err error
	//dsn := "root:root@tcp(localhost:3306)/?charset=utf8&parseTime=True&loc=Local"
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	CreateDb(dbName)
	SelectDb(dbName)
	//gorm.DB().SetMaxIdleConns(10)
}

// 创建表 todo 增加字段,过期时间,状态(比如锁定,密码错误)
func DbCreateTable(tableName string) error {
	fmt.Println("创建表")
	//创建表 默认的有id,created_at,updated_at,deleted_at,projectThis
	tx := Db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (id int(11) NOT NULL AUTO_INCREMENT,projectThis varchar(255) DEFAULT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return tx.Error
	}

	return nil
}

func DbAddUserColumn(tableName, columnNames string) error {
	columns := strings.Split(columnNames, ",") // 将列名字符串分割成数组

	// 获取表中所有现有的列名
	var columnsInTable []string
	Db.Select("COLUMN_NAME").Table("information_schema.COLUMNS").Where("TABLE_NAME = ?", tableName).Scan(&columnsInTable)
	for _, v := range columns {
		// 判断是否存在
		for _, v2 := range columnsInTable {
			if v == v2 {
				fmt.Println("列名已存在")
				return errors.New("列名已存在")
			}
		}
	}

	sql := "ALTER TABLE " + tableName + " ADD " + strings.Join(columns, " varchar(255) DEFAULT NULL, ADD ") + " varchar(255) DEFAULT NULL;"
	err := Db.Exec(sql).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func DbAddData(tableName string, data string) error {
	//sql := "INSERT INTO " + tableName + " (projectThis) VALUES ('" + data + "')"
	//err := Db.Exec(sql)
	//if err != nil {
	//	fmt.Println(err)
	//	return err.Error
	//}
	tx := Db.Table(tableName).Create(map[string]interface{}{"projectThis": data})
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return tx.Error
	}
	return nil
}

func DbAddDatas(tableName, datas, split string) (int, error) {
	IsTableExist(tableName)
	list := strings.Split(datas, split)
	//开启事务插入数据
	var countOk = 0
	Db.Begin()
	for _, v := range list {
		tx := Db.Table(tableName).Create(map[string]interface{}{"projectThis": v})
		if tx.Error != nil {
			fmt.Println(tx.Error)
			Db.Rollback()
			log.Println("插入数据失败", tx.Error)
		}
		countOk++
	}
	Db.Commit()
	return countOk, nil
}

func DbGetData(tableName, project string) (string, error) {
	var data []map[string]interface{}
	tx := Db.Table(tableName).Where(project + " is null OR " + project + " = ''").Limit(1).Scan(&data)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return "", tx.Error
	}
	//db取出id和projectThis
	var id int32
	if len(data) == 0 {
		return "", errors.New("没有数据")
	}
	id = data[0]["id"].(int32)

	//标记project为used
	tx = Db.Table(tableName).Where("id = ?", id).Update(project, "used")
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return "", tx.Error
	}
	return data[0]["projectThis"].(string), nil
}

func DbUpdate(tableName, id string, dataList []string) error {
	err := Db.Exec("UPDATE "+tableName+" SET "+strings.Join(dataList, ",")+" WHERE id = ?", id)
	if err != nil {
		fmt.Println(err)
		return err.Error
	}
	return nil
}

func DbBackupToEmail(dataname string) error {
	//备份dataname

	DbBackupTable("outlook")

	log.Println("备份数据库成功")
	// 发送邮件
	return nil
}

// 备份一个表
func DbBackupTable(tableName string) error {
	//备份tableName
	// 打开文件用于写入备份数据
	sqlFile, err := os.Create("backup" + tableName + ".sql")
	if err != nil {
		panic("failed to create backup file")
	}
	defer sqlFile.Close()
	// 用查询语句从数据库中导出数据到sql文件,然后打包压缩
	writer := bufio.NewWriter(sqlFile)

	//查询表所有数据
	rows, err := Db.Table(tableName).Rows()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()
	//遍历数据
	for rows.Next() {
		var id int
		var projectThis string
		err = rows.Scan(&id, &projectThis)
		if err != nil {
			fmt.Println(err)
			return err
		}
		//写入数据
		writer.WriteString(fmt.Sprintf("INSERT INTO %s (id, projectThis) VALUES (%d, %s);\n", tableName, id, projectThis))
	}
	writer.Flush()

	//打包压缩
	tarFile, err := os.Create("backup.tar.gz")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer tarFile.Close()

	// 创建ZIP文件
	zipFile, err := os.Create("export.zip")
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	// 创建ZIP写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 将SQL文件添加到ZIP中
	sqlFilePath := filepath.Join(".", "backup"+tableName+".sql")
	zipEntry, err := zipWriter.Create(sqlFilePath)
	if err != nil {
		panic(err)
	}

	// 读取SQL文件并写入ZIP
	sqlFileContent, err := os.ReadFile(sqlFilePath)
	if err != nil {
		panic(err)
	}
	zipEntry.Write(sqlFileContent)

	// 删除SQL文件
	err = os.Remove(sqlFilePath)
	if err != nil {
		panic(err)
	}
	log.Println("备份表成功")
	// 发送邮件
	return err
}

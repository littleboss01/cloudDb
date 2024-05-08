package cloudDb

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"

	//gin
	"github.com/gin-gonic/gin"
	"strings"
)

// done 中间件,如果强求你的数据有Encrypted 字段,则解密,否则不解密
// 解密函数
func decrypt(key, ciphertext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 填充模式
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 初始化向量
	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}

// 加密函数
func encrypt(key, plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 初始化向量
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	// 填充模式
	plaintextBytes := []byte(plaintext)
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(plaintextBytes, plaintextBytes)

	// 返回加密后的字符串
	return base64.StdEncoding.EncodeToString(append(iv, plaintextBytes...)), nil
}

// 自定义中间件
type customResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (c *customResponseWriter) Write(b []byte) (int, error) {
	return c.body.Write(b)
}

func CustomMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		if c.Request.Header.Get("token") != token {
			c.AbortWithStatus(401)
			return
		}
		//如果header里面有 EncrypData  字段,则解密,否则不解密
		if c.Request.Header.Get("EncrypData") == "" {
			c.Next()
			return
		}
		if c.Request.Method == "GET" {
			body = []byte(c.Query("EncrypData"))
		} else if c.Request.Method == "POST" {
			c.Request.Body.Read(body)
			ContentType := c.Request.Header.Get("Content-Type")
			if strings.Contains(ContentType, "application/x-www-form-urlencoded") {
				body = []byte(c.PostForm("EncrypData"))
			} else if strings.Contains(ContentType, "application/json") {
				//从json中取出EncrypData
				body = []byte(c.PostForm("EncrypData"))
			} else {

			}
		}

		decryptedBody, err := decrypt("124149449", string(body))
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBufferString(decryptedBody))

		// Create a custom ResponseWriter
		writer := &customResponseWriter{
			ResponseWriter: c.Writer,
			body:           new(bytes.Buffer),
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Encrypt response body
		encryptedBody, err := encrypt("124149449", writer.body.String())
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		c.Writer.Write([]byte(encryptedBody))
	}
}

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
	err := DbCreateTable(tabalName)
	if err == nil {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "创建表成功",
		})

	} else {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
	}
}

// 添加使用者列
func AddUserColumn(c *gin.Context) {
	tabalName := c.Query("tableName")

	columnNames := c.Query("columnNames")
	err := DbAddUserColumn(tabalName, columnNames)
	if err == nil {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "添加列成功",
		})

	} else {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
	}
}

// 添加数据
func AddData(c *gin.Context) {
	tabalName := c.Query("tableName")

	//添加数据
	data := c.Query("data")
	err := DbAddData(tabalName, data)
	if err == nil {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "添加数据成功",
		})
	} else {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
	}
}

// 取出数据
func GetData(c *gin.Context) {
	tabalName := c.Query("tableName")
	Project := c.Query("project")

	//取出数据
	data, err := DbGetData(tabalName, Project)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return

	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "ok",
		"data":    data,
	})
}

// 修改数据
func UpdateData(c *gin.Context) {
	tabalName := c.Query("tableName")
	id := c.Query("id")
	data := c.Query("data")

	dataList := strings.Split(data, "----")
	//修改数据
	err := DbUpdate(tabalName, id, dataList)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return

	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "修改数据成功",
		"data":    dataList,
	})
}

// Post json上传文件实现批量添加数据
func AddDataJson(c *gin.Context) {
	tabalName := c.Query("tableName")
	data := c.PostForm("data")
	//data是文件base64编码
	//解码
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}

	//通过事务批量添加数据
	DbAddDatas(tabalName, string(dataByte), "\n")
}

// 备份数据库发送到邮箱
func BackupToEmail(c *gin.Context) {
	//备份数据库
	err := DbBackupToEmail("cloudDb")
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "备份成功",
	})

}

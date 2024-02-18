package respond

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jsasg/gopkg/encrypt"
	"github.com/jsasg/gopkg/utils"
)

const (
	defaultSuccessMessage = "success"
	defaultFailMessage    = "fail"
)

type responser struct {
	Code          int         `json:"code"`
	Data          interface{} `json:"data"`
	Message       string      `json:"message"`
	HttpStatus    int         `json:"-"`
	Error         error       `json:"-"`
	ignoreEncrypt bool        `json:"-"`
	publicSSL     string      `json:"-"`
	privateSSL    string      `json:"-"`
}

func crypto(data interface{}, publicSSL, privateSSL string) interface{} {
	if publicSSL != "" && privateSSL != "" {
		// 随机生成 AES 密钥
		key := utils.RandString(16)

		// 使用 AES 加密返回数据
		aes := encrypt.NewAES([]byte(key))
		jsonDataBytes, _ := json.Marshal(data)
		ciphertext, _ := aes.EncryptString(string(jsonDataBytes))

		// 使用 RSA 加密保护 AES 密钥
		rsa := encrypt.NewRSA()
		rsa.WritePublicString(publicSSL)
		rsa.WritePrivateString(privateSSL)
		cipherkey, _ := rsa.EncryptString(utils.Strrev(string(key)))

		data = fmt.Sprintf("%s.%s", ciphertext, cipherkey)
	}

	return data
}

// SuccessJSON 成功返回
func SuccessJSON(c *gin.Context, options ...Option) {
	var p = &responser{
		Code:       0,
		Data:       struct{}{},
		Message:    defaultSuccessMessage,
		HttpStatus: http.StatusOK,
	}

	for _, option := range options {
		option(p)
	}
	if p.Error != nil {
		c.Error(p.Error)
	}
	if !p.ignoreEncrypt {
		p.Data = crypto(p.Data, p.publicSSL, p.privateSSL)
	}
	c.Abort()
	c.JSON(p.HttpStatus, p)
}

// FailJSON 失败返回
func FailJSON(c *gin.Context, options ...Option) {
	var p = &responser{
		Code:       1,
		Data:       struct{}{},
		Message:    defaultFailMessage,
		HttpStatus: http.StatusOK,
	}

	for _, option := range options {
		option(p)
	}
	if p.Error != nil {
		c.Error(p.Error)
	}
	if !p.ignoreEncrypt {
		p.Data = crypto(p.Data, p.publicSSL, p.privateSSL)
	}
	c.Abort()
	c.JSON(p.HttpStatus, p)
}

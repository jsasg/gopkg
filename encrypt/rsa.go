package encrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
)

type encryptRSA struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

// NewRSA 默认空密钥对
func NewRSA() *encryptRSA {
	return &encryptRSA{}
}

// BuildRSA 生成rsa密钥对
func BuildRSA() (r *encryptRSA, err error) {
	r = NewRSA()

	// 生成密钥对
	r.Private, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		err = fmt.Errorf("failed to build rsa, error: %w", err)
		return
	}
	r.Public = &r.Private.PublicKey

	return
}

// Validate 验证rsa合法
func (r *encryptRSA) Validate() (ok bool) {
	ok = reflect.DeepEqual(r, &encryptRSA{})
	return
}

// PrivateString 取私钥字符串
func (r *encryptRSA) PrivateString() (p string) {
	var buf bytes.Buffer

	// 处理私钥
	x509PrivateKey, _ := x509.MarshalPKCS8PrivateKey(r.Private)
	block := &pem.Block{Type: "PRIVATE KEY", Bytes: x509PrivateKey}
	private := pem.EncodeToMemory(block)

	buf.Write(private)

	return buf.String()
}

// PublicString 取公钥字符串
func (r *encryptRSA) PublicString() (p string) {
	var buf bytes.Buffer

	// 处理公钥
	x509PublicKey, _ := x509.MarshalPKIXPublicKey(r.Public)
	publicBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: x509PublicKey}
	public := pem.EncodeToMemory(publicBlock)

	buf.Write(public)

	return buf.String()
}

// WritePublicBytes 把[]byte类型的公钥写回rsa中
func (r *encryptRSA) WritePublicBytes(public []byte) (err error) {
	var (
		block        *pem.Block
		pubInterface any
		ok           bool
	)
	block, _ = pem.Decode(public)
	if block == nil {
		err = errors.New("please confirm that the public key is correct")
		return
	}
	pubInterface, _ = x509.ParsePKIXPublicKey(block.Bytes)
	if r.Public, ok = pubInterface.(*rsa.PublicKey); !ok {
		err = fmt.Errorf("failed to write the bytes, error: %w", err)
	}
	return
}

// WritePublicString 把string类型的公钥写回rsa中
func (r *encryptRSA) WritePublicString(public string) (err error) {
	var buf bytes.Buffer
	buf.WriteString(public)
	err = r.WritePublicBytes(buf.Bytes())
	return
}

// WritePrivateBytes 把[]byte类型的私钥写回rsa
func (r *encryptRSA) WritePrivateBytes(private []byte) (err error) {
	var (
		block        *pem.Block
		priInterface any
		ok           bool
	)
	block, _ = pem.Decode(private)
	if block == nil {
		err = errors.New("please confirm that the public key is correct")
		return
	}
	priInterface, _ = x509.ParsePKCS8PrivateKey(block.Bytes)
	if r.Private, ok = priInterface.(*rsa.PrivateKey); !ok {
		err = fmt.Errorf("failed to write the bytes, error: %w", err)
	}
	return
}

// WritePrivateString 把string类型的私钥写回rsa
func (r *encryptRSA) WritePrivateString(private string) (err error) {
	var buf bytes.Buffer
	buf.WriteString(private)
	err = r.WritePrivateBytes(buf.Bytes())
	return
}

// Encrypt 使用rsa加密
func (r *encryptRSA) Encrypt(plainbytes []byte) (cipherbytes []byte, err error) {
	var b []byte
	if r.Public == nil {
		err = errors.New("please confirm that the public key is correct")
		return
	}
	if b, err = rsa.EncryptPKCS1v15(rand.Reader, r.Public, plainbytes); err != nil {
		err = fmt.Errorf("failed to encrypt the plainbytes: %v, error: %w", plainbytes, err)
	}
	cipherbytes = make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(cipherbytes, b)
	return
}

// EncryptString 使用rsa加密
func (r *encryptRSA) EncryptString(plaintext string) (ciphertext string, err error) {
	var (
		cipherbytes []byte
		buf         bytes.Buffer
	)
	buf.WriteString(plaintext)
	if cipherbytes, err = r.Encrypt(buf.Bytes()); err != nil {
		return
	}
	buf.Reset()
	buf.Write(cipherbytes)
	ciphertext = buf.String()
	return
}

// Decrypt 使用rsa解密
func (r *encryptRSA) Decrypt(cipherbytes []byte) (plainbytes []byte, err error) {
	var buf bytes.Buffer
	if r.Private == nil {
		err = errors.New("please confirm that the private key is correct")
		return
	}
	buf.Write(cipherbytes)
	dst, _ := base64.StdEncoding.DecodeString(buf.String())
	if plainbytes, err = rsa.DecryptPKCS1v15(rand.Reader, r.Private, dst); err != nil {
		err = fmt.Errorf("failed to decrypt the cipherbytes: %v, error: %w", cipherbytes, err)
	}
	return
}

// DecryptString 使用rsa解密
func (r *encryptRSA) DecryptString(ciphertext string) (plaintext string, err error) {
	var (
		plainbytes []byte
		buf        bytes.Buffer
	)
	buf.WriteString(ciphertext)
	if plainbytes, err = r.Decrypt(buf.Bytes()); err != nil {
		return
	}
	buf.Reset()
	buf.Write(plainbytes)
	plaintext = buf.String()
	return
}

// ShardingEncryptString 分片加密
func (r *encryptRSA) ShardingEncryptString(plaintext string, options ...any) (ciphertext string, err error) {
	var (
		cipher           string
		plaintexts       []string
		limit, separ, ok = 81, ",", true
	)
	if len(options) >= 1 {
		if limit, ok = options[0].(int); !ok {
			err = fmt.Errorf("please confirm that the type of sharding size parameter: %v is int", options[0])
			return
		}
	}
	if len(options) >= 2 {
		if separ, ok = options[1].(string); !ok {
			err = fmt.Errorf("please confirm that the type of sharding size parameter: %v is string", options[1])
			return
		}
	}
	start := 0
	runes := []rune(plaintext)
	// 总字数
	wordage := len(runes)
	// 切片数
	copies := int(math.Ceil(float64(wordage) / float64(limit)))
	for i := 0; i < copies; i++ {
		// 字数不够时，直接取剩余字数
		end := limit
		if surplus := wordage - start; limit > surplus {
			end = surplus
		}
		// 分片
		text := string(runes[start : start+end])
		// 给分片加密comma
		if cipher, err = r.EncryptString(text); err != nil {
			return
		}
		plaintexts = append([]string{cipher}, plaintexts...)
		start = i*end + end
	}
	ciphertext = strings.Join(plaintexts, separ)
	return
}

// ShardingDecryptString 分片解密
func (r *encryptRSA) ShardingDecryptString(ciphertext string, options ...any) (plaintext string, err error) {
	var (
		plaintexts, ciphertexts []string
		text, cipher            string
		separ, ok               = ",", true
	)
	if len(options) >= 1 {
		if separ, ok = options[0].(string); !ok {
			err = fmt.Errorf("please confirm that the type of sharding size parameter: %v is string", options[0])
			return
		}
	}
	ciphertexts = strings.Split(ciphertext, separ)
	for _, cipher = range ciphertexts {
		if text, err = r.DecryptString(cipher); err != nil {
			return
		}
		plaintexts = append([]string{text}, plaintexts...)
	}
	plaintext = strings.Join(plaintexts, "")
	return
}

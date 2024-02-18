package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"math"
	"math/rand"
	"time"
)

var (
	ErrPlainTextEmpty    = errors.New("明文为空")
	ErrCipherTextEmpty   = errors.New("密文为空")
	ErrCipherLength      = errors.New("密文长度不正确")
	ErrCipherSaltInvalid = errors.New("密文混淆字符串无效")
	ErrCipherBlockLength = errors.New("密文不是块大小的倍数")
)

type Aes struct {
	Key       []byte
	Salt      []byte
	BlockSize int
	iv        []byte
}

type AesOption func(a *Aes)

// WithAesSalt 设置混淆字符串
func WithAesSalt(salt string) AesOption {
	return func(a *Aes) {
		a.Salt = []byte(salt)
	}
}

// NewAES 实例一个aes
func NewAES(aesKey []byte, options ...AesOption) *Aes {
	iv := aesKey[:aes.BlockSize]
	aes := &Aes{
		Key:       aesKey,
		BlockSize: aes.BlockSize,
		iv:        iv,
	}
	for _, opt := range options {
		opt(aes)
	}
	return aes
}

// EncryptString 加密字符串
func (a *Aes) EncryptString(plaintext string) (ciphertext string, err error) {
	if plaintext == "" {
		return "", ErrPlainTextEmpty
	}
	if cipherbytes, err := a.Encrypt([]byte(plaintext)); err == nil {
		ciphertext = base64.StdEncoding.EncodeToString(cipherbytes)
	}
	return
}

// Encrypt 加密
func (a *Aes) Encrypt(plaintext []byte) (ciphertext []byte, err error) {
	if plaintext == nil {
		return nil, ErrPlainTextEmpty
	}
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	var plaintextPaded []byte
	if a.Salt != nil {
		ptSaltOffset := len(plaintext) / 2
		ptPadedSalt := append(plaintext[:ptSaltOffset], append(a.Salt, plaintext[ptSaltOffset:]...)...)
		ptPadedSaltLenOffset := len(ptPadedSalt) / 2
		ptPadedSaltAndSaltLen := append(ptPadedSalt[:ptPadedSaltLenOffset], append(packN(uint32(len(a.Salt))), ptPadedSalt[ptPadedSaltLenOffset:]...)...)
		plaintextPaded = append(packN(1), ptPadedSaltAndSaltLen...)
	} else {
		plaintextPaded = append(packN(0), plaintext...)
	}
	plaintextPadedRandomStr := append(randomString(a.BlockSize), plaintextPaded...)
	text := pkCS7Padding(plaintextPadedRandomStr, a.BlockSize)
	if len(text)%a.BlockSize != 0 {
		err = ErrCipherBlockLength
		return
	}
	ciphertext = make([]byte, len(text))
	blockMode := cipher.NewCBCEncrypter(block, a.Key[:a.BlockSize])
	blockMode.CryptBlocks(ciphertext, text)
	return
}

// DecryptString 解密字符串
func (a *Aes) DecryptString(ciphertext string) (plaintext string, err error) {
	if ciphertext == "" {
		return "", ErrCipherTextEmpty
	}
	text, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	plainbytes, err := a.Decrypt(text)
	if err == nil {
		plaintext = string(plainbytes)
	}
	return
}

// Decrypt 解密
func (a *Aes) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	if ciphertext == nil {
		return nil, ErrCipherTextEmpty
	}
	if len(ciphertext[:a.BlockSize])%a.BlockSize != 0 {
		err = ErrCipherBlockLength
		return
	}
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	pt := make([]byte, len(ciphertext))
	blockMode := cipher.NewCBCDecrypter(block, a.Key[:a.BlockSize])
	blockMode.CryptBlocks(pt, ciphertext)
	unpadedPt := pkCS7UnPadding(pt)
	// 删除随机字符
	ptUnRanom := unpadedPt[a.BlockSize:]
	// 是否具有salt
	saltExists := unpackN(ptUnRanom[:4])
	ptUnSaltExists := ptUnRanom[4:]
	if saltExists == 0 {
		plaintext, err = ptUnSaltExists, nil
		return
	}
	if saltExists == 1 {
		// 获取salt长度
		saltLenOffset := int(math.Floor(float64(len(ptUnSaltExists))/2.0)) - 2
		saltLen := int(unpackN(ptUnSaltExists[saltLenOffset : saltLenOffset+4]))
		// 去除salt长度后的内容
		ptUnSaltLen := append(ptUnSaltExists[:saltLenOffset], ptUnSaltExists[saltLenOffset+4:]...)
		// 获取salt
		saltOffset := int(math.Floor(float64(len(ptUnSaltLen))/2.0 - float64(saltLen)/2.0))
		if salt := append([]byte{}, ptUnSaltLen[saltOffset:saltOffset+saltLen]...); !bytes.Equal(a.Salt, salt) {
			err = ErrCipherSaltInvalid
			return
		}
		// 明文内容
		plaintext = append(ptUnSaltLen[:saltOffset], ptUnSaltLen[saltOffset+saltLen:]...)
	}
	return
}

// 解密补位
func pkCS7UnPadding(ciphertext []byte) []byte {
	length := len(ciphertext)
	unpadding := int(ciphertext[length-1])
	return ciphertext[:(length - unpadding)]
}

// 加密补位
func pkCS7Padding(plantText []byte, blockSize int) []byte {
	padding := blockSize - len(plantText)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plantText, padtext...)
}

// randomString 生成随机字符
func randomString(length int) []byte {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = byte(letters[r.Intn(len(letters))])
	}
	return b
}

// packN 32位无符号长整型pack
func packN(arg uint32) []byte {
	var b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, arg)
	return b
}

// unpackN 32位无符号长整unpack
func unpackN(arg []byte) uint32 {
	return binary.BigEndian.Uint32(arg)
}

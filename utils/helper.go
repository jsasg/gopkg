package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// RandString 生成随机字符串
func RandString(n int) (s string) {
	r := make([]byte, n/2)
	rdom := rand.New(rand.NewSource(time.Now().UnixNano()))
	rdom.Read(r)
	s = hex.EncodeToString(r)
	return
}

// Strrev 反转字符串
func Strrev(str string) (rstr string) {
	r := []rune(str)
	l := len(r) - 1
	for from, to := 0, l; from < to; from, to = from+1, to-1 {
		r[from], r[to] = r[to], r[from]
	}
	rstr = string(r)
	return
}

// Unicode2utf8 unicode字符串转utf8字符串
func Unicode2utf8(from string) (to string) {
	res := []string{""}
	str := strings.Split(from, "\\u")
	context := ""
	for _, v := range str {
		additional := ""
		if len(v) < 1 {
			continue
		}
		if len(v) > 4 {
			rs := []rune(v)
			v = string(rs[:4])
			additional = string(rs[4:])
		}
		tmp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			context += v
		}
		context += fmt.Sprintf("%c", tmp)
		context += additional
	}
	res = append(res, context)
	return strings.Join(res, "")
}

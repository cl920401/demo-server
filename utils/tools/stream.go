package tools

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
	"unsafe"

	tsgutils "github.com/typa01/go-utils"
)

var (
	letters     = []byte("0123456789abcdefghijklmnopqrstuvwxyz_+-")
	lettersSize = len(letters)
)

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Bytes2str : Bytes2str高效转换
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//func CurrentTimeMillis() float64 {
//	return float64(time.Now().UnixNano()) / 1000000000 // ns to s
//}

// Exist : file is exist
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// FormatFloat : float2string
func FormatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}

func Base64DecodeToString(str string) string {
	deStr, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return Bytes2str(deStr)
}

//Base64EncodeToString base64编码
func Base64EncodeToString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64DecodeToByte(str string) []byte {
	deStr, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return []byte{}
	}
	return deStr
}

//func UrlBase64EncodeToString(data []byte) string {
//	return base64.URLEncoding.EncodeToString(data)
//}
//func SignBase64EncodeToString(data []byte) string {
//	base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
//	return base64.StdEncoding.EncodeToString(data)
//}
//
//func UrlBase64DecodeToByte(str string) ([]byte, error) {
//	deStr, err := base64.URLEncoding.DecodeString(str)
//	if err != nil {
//		return nil, err
//	}
//	return deStr, nil
//}

// SHA1 : SHA1 encode
// TODO : h.Write error listener
func SHA1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// MD5 : MD5 encode
func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func Hmac(key, data string) string {
	hmacObj := hmac.New(md5.New, []byte(key))
	hmacObj.Write([]byte(data))
	return hex.EncodeToString(hmacObj.Sum([]byte("")))
}

// HmacSha256 : Sha256 encode
// TODO : hmacObj.Write error listener
func HmacSha256(message string, secret string) string {
	hmacObj := hmac.New(sha256.New, []byte(secret))
	hmacObj.Write([]byte(message))
	return hex.EncodeToString(hmacObj.Sum(nil))
}

//func Str2Int(s string) int {
//	i, _ := strconv.Atoi(s)
//	return i
//}

// GenerateGUID : 生成UUID
func GenerateGUID() string {
	return tsgutils.GUID()
}

// GetRandomString : 当前时间戳 + 长度为n的随机字符串
func GetRandomString(n int) string {
	st := strconv.FormatInt(time.Now().Unix(), 10)
	res := []byte(st)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		res = append(res, letters[r.Intn(lettersSize)])
	}
	return string(res)
}

//Encrypt 对称加密
func Encrypt(data []byte, key string) string {
	md5Key := Str2bytes(MD5(key))

	dataLen := len(data)
	x := 0

	var char = make([]byte, 0, dataLen)
	for i := 0; i < dataLen; i++ {
		if x == 32 {
			x = 0
		}
		char = append(char, md5Key[x])
		x++
	}
	var data1 = make([]byte, 0, dataLen)
	for i := 0; i < dataLen; i++ {
		data1 = append(data1, uint8(int(data[i]+char[i])%256))
	}
	return Base64EncodeToString(data1)
}

//Decrypt 对称解密
func Decrypt(data string, key string) string {
	md5Key := Str2bytes(MD5(key))
	dataByte := Base64DecodeToByte(data)
	dataLen := len(dataByte)
	x := 0

	var char = make([]byte, 0, dataLen)

	for i := 0; i < dataLen; i++ {
		if x == 32 {
			x = 0
		}
		char = append(char, md5Key[x])
		x++
	}
	var data1 = make([]byte, 0, dataLen)

	for i := 0; i < dataLen; i++ {
		if dataByte[i] < char[i] {
			data1 = append(data1, uint8(int(dataByte[i])+256-int(char[i])))
		} else {
			data1 = append(data1, uint8(int(dataByte[i])-int(char[i])))
		}
	}
	return Bytes2str(data1)
}

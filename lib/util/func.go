package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"unsafe"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	// 如果是指针，则获取其所指向的元素
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		k := v.Type().Field(i).Tag.Get("map")
		if k == "" {
			continue
		}
		data[k] = v.Field(i).Interface()
	}
	return data
}

//Str2bytes 字符串转byte 高效
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//Bytes2str byte转字符串 高效
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//Exist 判断文件|目录是否存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//FormatFloat 格式化float
func FormatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}

//Base64DecodeToString base64 解密
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

//Base64DecodeToByte base64解密
func Base64DecodeToByte(str string) []byte {
	deStr, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return []byte{}
	}
	return deStr
}

//SHA1 SHA1加密
func SHA1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

//MD5 MD5加密
func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

//RandStringBytesMaskImprSrcBase 随机生成字符串
func RandStringBytesMaskImprSrcBase(n int, seed rand.Source, bs string) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, seed.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = seed.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(bs) {
			b[i] = bs[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return Bytes2str(b)
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

//计算文件md5值
func MD5sum(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	md5sum := md5.New()

	if _, err := io.Copy(md5sum, f); err != nil {
		return ""
	}
	return hex.EncodeToString(md5sum.Sum(nil))
}

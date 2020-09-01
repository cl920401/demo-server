package util

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterNumberBytes = "1234567890"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

//RandString 上称随机字母字符串
func RandString(n int) string {
	return randStringBytesMaskImprSrc(n, letterBytes)
}

//RandNumber 生成随机数字字符串
func RandNumber(n int) string {
	return randStringBytesMaskImprSrc(n, letterNumberBytes)
}

func randStringBytesMaskImprSrc(n int, letter string) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letter) {
			b[i] = letter[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

package crud

import (
	cryptoRand "crypto/rand"
	"math/big"
	"math/rand"
	"strconv"
	"time"
	"unsafe"

	"github.com/google/uuid"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func GetRandomStr(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

type IdGenerator[T any] interface {
	GetNewId() T
}

func CreateNewIdGenerator[T any]() IdGenerator[T] {
	return &DefaultIdGenerator[T]{}
}

type DefaultIdGenerator[T any] struct {
}

func (service *DefaultIdGenerator[T]) GetNewId() T {
	var val any = *new(T)
	if _, ok := val.(string); ok {
		val = uuid.NewString()
		return val.(T)
	} else if _, ok := val.(int32); ok {
		nBig, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(100000000000))
		val = int32(nBig.Uint64() & 0xFFFF0000)
		return val.(T)
	} else if _, ok := val.(int64); ok {
		nBig, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(100000000000))
		val = nBig.Uint64()
		return val.(T)
	} else {
		return *new(T)
	}
}

func Parse[T any](str string) T {
	var val any = *new(T)
	if _, ok := val.(string); ok {
		val = str
		return val.(T)
	} else if _, ok := val.(int32); ok {
		parsedInt64, _ := strconv.ParseInt(str, 10, 32)
		val = int32(parsedInt64)
		return val.(T)
	} else if _, ok := val.(int64); ok {
		parsedInt64, _ := strconv.ParseInt(str, 10, 64)
		val = parsedInt64
		return val.(T)
	} else if _, ok := val.(int16); ok {
		parsedInt64, _ := strconv.ParseInt(str, 10, 16)
		val = int16(parsedInt64)
		return val.(T)
	} else if _, ok := val.(int8); ok {
		parsedInt64, _ := strconv.ParseInt(str, 10, 8)
		val = int8(parsedInt64)
		return val.(T)
	} else if _, ok := val.(uint32); ok {
		parsedInt64, _ := strconv.ParseUint(str, 10, 32)
		val = uint32(parsedInt64)
		return val.(T)
	} else if _, ok := val.(uint64); ok {
		parsedInt64, _ := strconv.ParseUint(str, 10, 64)
		val = parsedInt64
		return val.(T)
	} else if _, ok := val.(uint16); ok {
		parsedInt64, _ := strconv.ParseUint(str, 10, 16)
		val = uint16(parsedInt64)
		return val.(T)
	} else if _, ok := val.(uint8); ok {
		parsedInt64, _ := strconv.ParseUint(str, 10, 8)
		val = uint8(parsedInt64)
		return val.(T)
	} else if _, ok := val.(bool); ok {
		parsedBool, _ := strconv.ParseBool(str)
		val = parsedBool
		return val.(T)
	} else {
		return *new(T)
	}
}

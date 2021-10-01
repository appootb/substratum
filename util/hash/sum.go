package hash

import (
	"fmt"
	"reflect"
	"strconv"
)

// Sum returns the checksum of the key.
func Sum(key interface{}) (val int64) {
	var err error

	switch reflect.TypeOf(key).Kind() {
	case reflect.Int:
		val = int64(key.(int))
	case reflect.Int8:
		val = int64(key.(int8))
	case reflect.Int16:
		val = int64(key.(int16))
	case reflect.Int32:
		val = int64(key.(int32))
	case reflect.Int64:
		val = key.(int64)
	case reflect.Uint:
		val = int64(key.(uint))
	case reflect.Uint8:
		val = int64(key.(int8))
	case reflect.Uint16:
		val = int64(key.(uint16))
	case reflect.Uint32:
		val = int64(key.(uint32))
	case reflect.Uint64:
		val = int64(key.(uint64))
	case reflect.String:
		val, err = strconv.ParseInt(key.(string), 10, 64)
		if err != nil {
			return checkSum(key.(string))
		}
	case reflect.Ptr:
		return Sum(reflect.ValueOf(key).Elem().Interface())
	default:
		val = checkSum(fmt.Sprintf("%v", key))
	}
	return
}

func checkSum(key string) int64 {
	var hash uint32 = 0
	for i := 0; i < len(key); i++ {
		hash += uint32(key[i])
		hash += hash << 10
		hash ^= hash >> 6
	}
	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15

	return int64(hash)
}

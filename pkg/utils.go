package util

import (
	"container/list"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"sync"
)

func Repr(v interface{}) string {
	if v == nil {
		return ""
	}

	// if func (v *Type) String() string, we can't use Elem()
	switch vt := v.(type) {
	case fmt.Stringer:
		return vt.String()
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	switch vt := val.Interface().(type) {
	case bool:
		return strconv.FormatBool(vt)
	case error:
		return vt.Error()
	case float32:
		return strconv.FormatFloat(float64(vt), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(vt, 'f', -1, 64)
	case fmt.Stringer:
		return vt.String()
	case int:
		return strconv.Itoa(vt)
	case int8:
		return strconv.Itoa(int(vt))
	case int16:
		return strconv.Itoa(int(vt))
	case int32:
		return strconv.Itoa(int(vt))
	case int64:
		return strconv.FormatInt(vt, 10)
	case string:
		return vt
	case uint:
		return strconv.FormatUint(uint64(vt), 10)
	case uint8:
		return strconv.FormatUint(uint64(vt), 10)
	case uint16:
		return strconv.FormatUint(uint64(vt), 10)
	case uint32:
		return strconv.FormatUint(uint64(vt), 10)
	case uint64:
		return strconv.FormatUint(vt, 10)
	case []byte:
		return string(vt)
	default:
		return fmt.Sprint(val.Interface())
	}
}

func InterfaceToSlice(vals []interface{}) []string {
	ret := make([]string, len(vals))
	for i, val := range vals {
		if val == nil {
			ret[i] = ""
		} else {
			switch val := val.(type) {
			case string:
				ret[i] = val
			default:
				ret[i] = Repr(val)
			}
		}
	}
	return ret
}

func SyncMapLen(sm *sync.Map) int {
	length := 0
	f := func(key, value interface{}) bool {
		length++
		return true
	}
	one := length
	length = 0
	sm.Range(f)
	if one != length {
		one = length
		length = 0
		sm.Range(f)
		if one < length {
			return length
		}

	}
	return one
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}

	return y
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}

	return y
}

func SyncMapWalk(key, value interface{}) bool {
	fmt.Println("Key=", key, "Value=", value)
	return true
}

func ListPrint(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func SliceRemove(list []string, target string) (result []string) {
	result = list
	for i, elem := range list {
		if elem == target {
			result = append(list[:i], list[i+1:]...)
			break
		}
	}
	return
}

func SliceSplit(srcSlice []string, batchSize int) (s [][]string) {
	s = make([][]string, 0, (len(srcSlice)+batchSize-1)/batchSize)
	for i := 0; i < len(srcSlice); i += batchSize {
		var (
			tmp []string
		)
		if len(srcSlice) < batchSize {
			tmp = srcSlice[:len(srcSlice)]
		} else if i+batchSize > len(srcSlice) {
			tmp = srcSlice[i:]
		} else {
			tmp = srcSlice[i : i+batchSize]
		}
		s = append(s, tmp)
	}
	return
}

func PrintMemStats(symbol string) {
	runtime.GC()
	var m runtime.MemStats

	runtime.ReadMemStats(&m)
	var unit uint64 = 1024
	log.Printf("NumGC = %v,symbol = %s,Alloc = %v KB,"+
		"TotalAlloc = %v KB, HeapInuse = %v KB,Sys = %v KB,HeapObjects = %v\n",
		m.NumGC, symbol, m.Alloc/unit, m.TotalAlloc/unit,
		m.HeapInuse/unit, m.Sys/unit, m.HeapObjects)
}

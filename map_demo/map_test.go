package main

import (
	util "github.com/tti6o/go_demo/pkg"
	"log"
	"testing"
)

//map里的value赋值给变量再修改，并不会改变map,需要重新赋值回去
func BenchmarkMapValueModify(b *testing.B) {
	baseContTag := make(map[string][]interface{})
	baseContTag["111"] = []interface{}{"xyz"}
	tagIDList := baseContTag["111"]
	tagIDList = tagIDList[:len(tagIDList)-1]
	//baseContTag["111"] = tagIDList
	log.Println(tagIDList)
	log.Println(baseContTag["111"])
	log.Println(baseContTag)
}

//测试全局map需要主动赋值nil才会被gc释放内存
var intMap = make(map[int]interface{})

func TestMapRelease(t *testing.T) {
	log.Printf("initMap before")
	util.PrintMemStats("TestMapRelease-1")
	var cnt = 100000
	for i := 0; i < cnt; i++ {
		intMap[i] = i
	}
	log.Println(len(intMap))
	log.Printf("initMap after")
	util.PrintMemStats("TestMapRelease-2")
	for i := 0; i < cnt-1; i++ {
		delete(intMap, i)
	}
	log.Println(len(intMap))
	log.Printf("delete after")
	util.PrintMemStats("TestMapRelease-3")
	intMap = nil
	log.Printf("set nil after")
	util.PrintMemStats("TestMapRelease-4")
}

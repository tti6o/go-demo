package main

import (
	"encoding/json"
	"fmt"
)

type ObjWithSliceMember struct {
	Name      string
	Children  *[]string
	Children2 interface{}
}

func main() {
	var member ObjWithSliceMember
	var arr []string
	var arr2 []string
	arr = append(arr, "a")
	arr2 = append(arr2, "a")
	member.Name = "a"
	member.Children = &arr
	member.Children2 = arr2
	//fmt.Printf("&arr:%p\n",&arr)
	//fmt.Printf("member.Children address:%p\n",member.Children)
	//fmt.Printf("member.Children2 Type:%T\n",member.Children2)
	mjson, _ := json.Marshal(member)
	ajson, _ := json.Marshal(arr)
	ajson2, _ := json.Marshal(arr2)
	fmt.Println(string(mjson), string(ajson), string(ajson2))
	jsonStr := `{"Name":"a","Children":["bbbb"],"Children2":["cccc"]}`
	err := json.Unmarshal([]byte(jsonStr), &member)
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	mjson, _ = json.Marshal(member)
	ajson, _ = json.Marshal(arr)
	ajson2, _ = json.Marshal(arr2)
	fmt.Printf("&arr:%p\n", &arr)
	fmt.Printf("member.Children address :%p\n", member.Children)
	fmt.Printf("member.Children2 Type:%T\n", member.Children2)
	fmt.Println(string(mjson), string(ajson), string(ajson2))
}

//func getType(a interface{}) {
//	switch a.(type) {
//	case int:
//		fmt.Println("the type of a is int")
//	case string:
//		fmt.Println("the type of a is string")
//	case float64:
//		fmt.Println("the type of a is float")
//	default:
//		fmt.Println("unknown type")
//	}
//}

package main

import (
	"fmt"
	"testing"
)

//二维数组赋值时要特别注意
func TestBFS(t *testing.T) {
	queue := make([]int, 0)
	//queue[1] = 1 panic
	queue = append(queue, 1)
	fmt.Println(queue)
	result := make([][]int, 0) //初始化
	//result[0] = make([]int, 0) panic
	tmpList := make([]int, 0)
	result = append(result, tmpList)
	fmt.Println(result)
}

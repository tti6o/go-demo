package main

import "fmt"

func main() {
	var arr []string
	arr = make([]string, 0, 10)
	printLenCap(arr)
	for i := 0; i < 10; i++ {
		arr = append(arr, "aaa")
		printLenCap(arr)
	}
	arr2 := arr[:5]
	printLenCap(arr2)
	arr2[0] = "bbb"
	fmt.Println(arr)
	fmt.Println(arr2)
}

func printLenCap(arr []string) {
	fmt.Printf("len: %d, cap: %d %v\n", len(arr), cap(arr), arr)
}

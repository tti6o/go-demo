package main

import (
	"fmt"
	"strings"
)

func main() {
	//判断同名是先拿父节点的路径名再加tagName跟已有的节点路径名对比的，这里父节点不存在直接就返回false了
	names := make(map[string]string, 0)
	namePath := names["111"]
	lv := strings.Count(namePath, "/")
	fmt.Println(lv)
}

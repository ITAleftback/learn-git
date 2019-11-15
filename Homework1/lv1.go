package main

import (
	"fmt"
	"time"
)
func main(){
	var a int64
	fmt.Println("input")
	fmt.Scanf("%d",&a)
	fmt.Println("output")
	fmt.Println(time.Unix(a,0))

}

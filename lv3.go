package main

import (
	"fmt"
	"time"
)
func prime(){
	for i:=2;i<=10000;i++{
		for n:=2;n<=i;n++{
			if n==i{
				fmt.Println(i)
			}
			if i%n==0&&n<i{
				break
			}
		}
	}
}
func main(){
	go prime()
	time.Sleep(time.Second)
}





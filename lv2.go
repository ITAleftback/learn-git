package main

import (
	"fmt"
	"time"
)

func factorial(n int,c chan int) {
	var res = 1
	for i := 1; i <=n; i++ {
		res *= i
		c<-res
	}
 close(c)
}
func main(){
c:=make(chan int,20)
go factorial(cap(c),c)
for i:=range c{
	fmt.Println(i)
}
time.Sleep(time.Second)
}
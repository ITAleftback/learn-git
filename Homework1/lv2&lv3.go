package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	//"strconv"
	"sync"
	//""
)


var j,i int64
var WEB string
var lock sync.Mutex
func main()  {
	file, err1 := os.Create("list.txt");
	if err1 != nil {
		fmt.Println(err1);
	}
	client := &http.Client{}

	for i = 0;;i++ {
		gogo()
		resp, err := client.Get(WEB)
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		a := regexp.MustCompile(`201921[\d]{4}[\p{Han}]+`)
		c := regexp.MustCompile(`[\p{Han}]+`)
		name := strings.Join(c.FindStringSubmatch(strings.Join(a.FindStringSubmatch(string(body)),``)),``)

		fmt.Println(name)
		if err != nil {
			fmt.Println(err)
		}
		file.Write([]byte(name));
		file.Write([]byte(" "))
		fmt.Println(i+1)
	}
	file.Close()
}
func gogo()  {
	WEB = "http://jwzx.cqupt.edu.cn/kebiao/kb_stu.php?xh="
	lock.Lock()
	j = 2019210000 + i
	s :=strconv.FormatInt(j, 10)
	WEB += s
	lock.Unlock()
	return
}

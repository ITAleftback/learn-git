package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func MatchString(pattern string, s string) (matched bool, err error){

}
func main() {
	data, err := ioutil.ReadFile("test.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	findNumOrLowerLetter(data)
	fmt.Println(regexp.Match("H.* ", "Hello World!"))

}
func Validate(str string) int {
	reg1 := regexp.MustCompile(`[0-9]+`)  //正则匹配数字，匹配到数字函数返回1
	reg2 := regexp.MustCompile(`[\p{Han}]+`)//正则匹配汉字，匹配到汉字函数返回2
	reg3 := regexp.MustCompile(`[a-z]+`)//正则匹配拼音，匹配到拼音函数返回3
	if reg1.MatchString(str) {
		return 1
	}
	if reg2.MatchString(str) {
		return 2
	}
	if reg3.MatchString(str) {
		return 3
	}
	return -1
}
func findNumOrLowerLetter(data string) {
	str := data
	reg := regexp.MustCompile("[\\d|a-z]+")
	fmt.Println(reg.FindAllString(str, -1))
	//[00az abc09 ab 99]
}


type student struct {
	name string
}
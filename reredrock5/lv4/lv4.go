package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Vedio struct {
	name string
	vedios string
}
///   写的爬b站的视频排行榜
func main() {
	Parse()
}

func Parse()  {
	//网页数据获取

	vedio:=[]Vedio{}

	userAgent := `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36`
	c := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	url := "https://www.bilibili.com/ranking?spm_id_from=333.851.b_7072696d61727950616765546162.3"

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", userAgent)
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Failed to get the website information")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	//处理数据
	reBody := strings.ReplaceAll(string(body), "\n", "")



	diReg :=regexp.MustCompile(`<ul class="rank-list">(.*?)</ul>`)
	diList :=diReg.FindAllStringSubmatch(reBody,-1)

	liReg :=regexp.MustCompile(`<li class="rank-item">(.*?)</li>`)
	liList :=liReg.FindAllString(diList[0][0],-1)

	for _,v :=range liList{
		//爬名字
		naReg :=regexp.MustCompile(`<img alt="(.*?)"`)
		naInfo :=naReg.FindStringSubmatch(v)
		// 爬vedio
		veReg :=regexp.MustCompile(`<a href="(.*?)" target="_blank">`)
		veInfo :=veReg.FindStringSubmatch(v)
		vedio=append(vedio,Vedio{
			name: naInfo[1],
			vedios: veInfo[1],
		})
	}
	for k,v:=range vedio{
		fmt.Printf("No.%d\n",k+1);
		fmt.Println("名称:",v.name,"\n","视频:",v.vedios)
	}

}

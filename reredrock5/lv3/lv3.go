package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)
type Search struct {
	result string
}
func main()  {
	router:=gin.Default()
	router.POST("/query",Query)
	router.Run(":8080")
}
func Query(c *gin.Context)  {

	keyword:=c.PostForm("keyword")
	value:=parse(keyword)
	//打印结果
	Result(value)
}
func parse(keyword string)(Searchs []Search){
	userAgent := `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36`
	c := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	url:="https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=1&rsv_idx=1&tn=baidu&wd="+keyword+"&oq=python%25E8%25AE%25BF%25" +
		"E9%2597%25AE%25E7%2599%25BE%25E5%25BA%25A6%25E6%2590%259C%25E7%25B4%25A2%25E6%258E%25A5%25E5%258F%25A3&rsv_pq=ecb5cd4b0002b9" +
		"a4&rsv_t=0200YI4Jcg6vohyCpg%2FUOvB8rI59hFP6yd9iordlpfk58Z5UrEWH7Kv1sUQ&rqlang=cn&rsv_enter=0&rsv_d" +
		"l=tb&inputT=5015&rsv_sug3=170&rsv_sug1=70&rsv_sug7=101&rsv_sug2=0&rsv_sug4=5832"

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
	rereBody:=strings.ReplaceAll(reBody," ","")
	divReg:=regexp.MustCompile(`<divid="content_left">(.*?)</div>`)
	divList:=divReg.FindAllStringSubmatch(rereBody,-1)
	//爬结果
	liReg := regexp.MustCompile(`<divclass="resultc-container"id="\d"srcid="1599"tpl="se_com_default"data-click="{"rsv_bdr":"0","p5":\d}">(.*?)</div>`)
	liList := liReg.FindAllString(divList[0][0], -1)
	search:=[]Search{}
	for _,v:=range liList{
		imgReg:=regexp.MustCompile(`<atarget="_blank"href="(.*?)"class="c-showurl>"style="text-decoration:none;">`)
		imgInfo:=imgReg.FindStringSubmatch(v)

		search=append(search,Search{
			result:imgInfo[1],
		})
	}
	Searchs=search
	return
}
//打印结果
func Result(value []Search){
	for k,v:=range value{
		fmt.Printf("No.%d\n",k+1);
		fmt.Println("搜索结果为：",v.result)
	}
}

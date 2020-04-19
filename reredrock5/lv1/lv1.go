package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Movie struct {
	Img string
	Name string
	Commentnum string
	Comment string
	Director string
}
//  因为最后一页有三个电影没有评价  故我爬了TOP225 的电影
func parse()(Movies []Movie) {
	//网页数据获取
	userAgent := `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36`
	c := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	movies :=[]Movie{}
	for i := 0; i < 9; i++ {
		url := "https://movie.douban.com/top250?start=" + strconv.Itoa(i*25) + "&filter="

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
		olReg := regexp.MustCompile(`<ol class="grid_view">(.*?)</ol>`)
		olList := olReg.FindAllStringSubmatch(reBody, -1)
		//fmt.Println(olList)
		liReg := regexp.MustCompile(`<li>(.*?)</li>`)

		liList := liReg.FindAllString(olList[0][0], -1)
		for _, v := range liList {
			//爬 名字 图片
			imgReg := regexp.MustCompile(`<img width="\d+" alt="(.*?)" src="(.*?)" class="">`)
			imgInfo := imgReg.FindStringSubmatch(v)
			//爬评分
			imgReg2 := regexp.MustCompile(`<span class="rating_num" property="v:average">(.*?)</span>`)
			imgInfo2 := imgReg2.FindStringSubmatch(v)
			//爬评价

			imgReg3:=regexp.MustCompile(`<span class="inq">(.*?)</span>`)
			imgInfo3:=imgReg3.FindStringSubmatch(v)

			//爬导演
			imgReg4 := regexp.MustCompile(`导演:(.*?)&`)
			imgInfo4:= imgReg4.FindStringSubmatch(v)

			//M[i].Name=imgInfo[1]
			//M[i].Img=imgInfo[2]
			//M[i].Commentnum=imgInfo2[1]
			//M[i].Director=imgInfo4[1]
			m:=Movie{
				Name:       imgInfo[1],
				Img:        imgInfo[2],
				Commentnum: imgInfo2[1],
				Comment: imgInfo3[1],
				Director:   imgInfo4[1],
			}
			movies = append(movies, m)
		}
	}
	Movies=movies
	return
}
func main()  {
	router:=gin.Default()
	router.POST("/query",Query)
	router.Run(":8080")
}
func Query(c *gin.Context)  {
	value:=parse()
	for k,v:=range value{
		fmt.Printf("No.%d\n",k+1);
		fmt.Println("名称为：", v.Name, "\n", "Img：", v.Img, "\n", "评分：", v.Commentnum, "\n", "评价：",v.Comment,"\n", "导演：", v.Director)
		c.JSON(200,gin.H{"status":http.StatusOK,"message":v})
	}
}
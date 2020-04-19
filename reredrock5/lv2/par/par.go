package par

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	. "reredrock5/lv2/init"
	"reredrock5/lv2/str"
	"strconv"
	"strings"
)

func Parse(num int)(Messages  []str.Student)  {
	//网页数据获取
	message:=[]str.Lesson{}
	messages:=[]str.Student{}
	userAgent := `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36`
	c := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	url := "http://jwc.cqupt.edu.cn/kebiao/kb_stu.php?xh=" + strconv.Itoa(2019210000+num)

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
	reBody := strings.ReplaceAll(string(body), "\r\n", "")
	rerebody := strings.ReplaceAll(reBody, " ", "")
	rererebody := strings.ReplaceAll(rerebody, "\t", "")


	//爬 名字学号的
	divReg := regexp.MustCompile(`<divstyle="float:left;">(.*?)</div>`)
	divList := divReg.FindAllStringSubmatch(rererebody, -1)
	olReg := regexp.MustCompile(`<ul>(.*?)</ul>`)
	olList := olReg.FindAllStringSubmatch(divList[0][0], -1)
	// 名字 学号都在这
	liReg := regexp.MustCompile(`<li>(.*?)>>(\d+)(.*?)</li>`)
	liInfo := liReg.FindStringSubmatch(olList[0][0])
	//至此  名字学号爬完


	//限定条件
	tdReg := regexp.MustCompile(`<tdrowspan='\d'>(.*?)<tdrowspan='\d'></td></tr>`)
	tdList := tdReg.FindAllString(rererebody, -1)
	for _, v := range tdList {

		//爬课程
		lessonReg := regexp.MustCompile(`<tdrowspan='\d'>(.*?)</td>`)
		lessonInfo := lessonReg.FindStringSubmatch(v)
		////爬教学班
		classReg := regexp.MustCompile(`<tdrowspan='\d'>(\w+)</td>`)
		classInfo := classReg.FindStringSubmatch(v)
		////爬类型 老师
		ttReg := regexp.MustCompile(`<tdrowspan='\d'>\w+</td><tdrowspan='\d'>(.*?)</td><tdrowspan='\d'align='center'>(.*?)</td><td>(.*?)</td>`)
		ttInfo := ttReg.FindStringSubmatch(v)
		//
		////爬 上课时间 地点
		tpReg := regexp.MustCompile(`<td>(.*?)</td><td>(.*?)</td><td>(.*?)</td>`)
		tpInfo := tpReg.FindStringSubmatch(v)

		p:=str.Lesson{
			Lesson:  lessonInfo[1],
			Class:   classInfo[1],
			Typee:   ttInfo[1],
			Teacher: ttInfo[3],
			Time:    tpInfo[2],
			Place:   tpInfo[3],
		}
		message = append(message, p)
		DB.Create(&p)

	}
	s:=str.Student{
		Stunum:   liInfo[2],
		Name:    liInfo[3],
		Lesson: message,
	}

	messages = append(messages,s)
	Messages=messages
	return
}
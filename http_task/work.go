package http_task

import (
	"bytes"
	"cron_demo/const"
	"cron_demo/http_request"
	"cron_demo/task"
	"cron_demo/work"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/**
Work 代表一个协程内具体执行任务工作者
*/
type Work struct {
	work.HttpWork
	manager  *Manager
	QPS      int
	closeSig bool
	num      int
}


func (this *Work) Init(t task.Task) {
	this.QPS = 10
	this.num = 0
	this.closeSig = false
}

/**

 */
func (this *Work) RunWorker(t task.Task) {
	for !this.closeSig {
		var throttle <-chan time.Time
		if this.QPS > 0 {
			throttle = time.Tick(time.Duration(1e6/(this.QPS)) * time.Microsecond)
		}

		if this.QPS > 0 {
			<-throttle
		}
		this.worker(t)
	}
}

//不同的请求

func (this *Work) worker(t task.Task) {
	this.num++
	//println("=======site - 启动接口 是否升级，启动图=======")
	//http_request.Launch()
	//println("=======site - 版头图及首页精品周更=======")
	//http_request.GetTop()
	//println("=======site - 首页小图标=======")
	//http_request.GetIcons()
	//println("=======site - 分类根信息=======")
	//http_request.GetCatalogRoot()
	//println("=======site - 有声漫画=======")
	//http_request.GetCatalogHomepage()
	if (this.num == _const.Num) {
		panic("运行结束")
	}
	println("=======message - 发送评论=======")
	http_request.UnlockTest()
}
func (this *Work) Close(t task.Task) {
	this.closeSig = true
	fmt.Println(fmt.Sprintf("num : %d", this.num))
}

//登录
func WxReadLogin() string {

	params := make(map[string]interface{})
	params["deviceToken"] = "f457f844d5a8d276ff68a5e51325a7c2bf8b46d1f9c140f2df4633d589c8f06c"
	params["refreshToken"] = "onb3MjrYiTYHVJty-8XVrmttJLOE@SA4lqeNuR_Co8LPauWOAewAA"
	params["deviceId"] = "ad15368d332e644562c5e887d0f977a2"
	params["refCgi"] = "/mobileSync"
	params["wxToken"] = 1
	params["inBackground"] = 0
	bytesData, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err.Error() )
	}
	reader := bytes.NewReader(bytesData)

	client := &http.Client{}
	request, _ := http.NewRequest("POST","https://i.weread.qq.com/login", reader)

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("User-Agent", "PostmanRuntime/7.4.0")
	request.Header = header
	result, err := client.Do(request)
	defer result.Body.Close()

	if err != nil {
		println("error")
	}

	if result.StatusCode != 200 {
		println("error")
	}

	body, _ := ioutil.ReadAll(result.Body)
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(body), &dat); err == nil {
		return dat["skey"].(string)
	} else {
		fmt.Println(err)
	}

	return ""
}




//获取详细
func getChapters(bookId string, skey string) string {
	client := &http.Client{}
	total := getChaptersSize(bookId,skey)
	//total = 2
	//count := 0
	for i:= int64(1);i < total ; i++ {
		time.Sleep(time.Duration(55000)*time.Microsecond)
		url := "https://i.weread.qq.com/book/chapterdownload?bookId="+bookId+"&bookType=txt&bookVersion=0&chapters="+ strconv.FormatInt(i, 10) +"&pf=weread_wx-2001-iap-2001-iphone&pfkey=pfkey&release=1&screenSize=16x9&synckey=0&zoneId=0"
		request, _ := http.NewRequest("GET", url, nil)
		header := make(http.Header)
		header.Set("User-Agent", "PostmanRuntime/7.4.0")
		header.Set("vid", "238361576")
		header.Set("skey", skey)
		request.Header = header
		result, err := client.Do(request)
		defer result.Body.Close()
		if err != nil {
			println("getChapters.error")
		}

		if result.StatusCode != 200 {
			println("getChapters.StatusCode.error"+skey)
			println(url)
		}

		body, _ := ioutil.ReadAll(result.Body)
		//errmsg := gjson.Get(string(body), "errmsg")
		//println("错误"+string(body))

		res := string(body)
		tmp := strings.Split(res, "ustar")
		//println(len(tmp))
		//println(res)
		if (len(tmp) < 3) {
			db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/charles")
			if err != nil {
				panic(err)
			}
			stmt, err := db.Prepare("UPDATE book SET status = 2 WHERE bookId="+bookId)
			if err != nil {
				log.Fatal(err)
			}
			res, err := stmt.Exec()
			lastId, err := res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
			return ""
		}
		tmp1 := strings.Split(tmp[1], "}]}")
		titleS := tmp1[0] + "}]}"

		title := gjson.Get(titleS, "chapters.0.title")
		chapterUid := gjson.Get(titleS, "chapters.0.chapterUid")
		wordCount := gjson.Get(titleS, "chapters.0.wordCount")
		updateTime := gjson.Get(titleS, "chapters.0.updateTime")
		content := tmp[2]

		sqlStr := "INSERT INTO `chapters` (`bookId`, `content`, `chapterUid`, `title`, `wordCount`, `updateTime`) VALUES ('"+bookId+"', '"+content+"', '"+chapterUid.String()+"', '"+title.String()+"', '"+wordCount.String()+"', "+updateTime.String()+")"

		//fmt.Print(sqlStr)

		c, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("Connect to redis error", err)
		}
		c.Do("lpush", bookId, sqlStr)
		println("============领取成功================" + bookId)
		defer c.Close()
	}


	return "123"
}

//获取章节数
func getChaptersSize(bookId string, skey string) int64 {
	client := &http.Client{}
	url := "https://i.weread.qq.com/book/info?bookId="+bookId+"&myzy=1"
	request, _ := http.NewRequest("GET", url, nil)
	header := make(http.Header)
	header.Set("User-Agent", "PostmanRuntime/7.4.0")
	header.Set("vid", "238361576")
	header.Set("skey", skey)

	request.Header = header
	result, err := client.Do(request)

	defer result.Body.Close()

	if err != nil {
		println("getChaptersSize.error")
	}

	if result.StatusCode != 200 {
		println("getChaptersSize.error"+skey)
	}

	body, _ := ioutil.ReadAll(result.Body)
	//println(string(body))
	//var dat map[string]interface{}
	//if err := json.Unmarshal([]byte(body), &dat); err == nil {
	//
	//	return int(dat["chapterSize"].(float64))
	//
	//} else {
	//	fmt.Println(err)
	//}
	total := gjson.Get(string(body), "chapterSize").Int()
	return total
}



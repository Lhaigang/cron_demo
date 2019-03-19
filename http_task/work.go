package http_task

import (
	"cron_demo/task"
	"cron_demo/work"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
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
	GetData()
}
func (this *Work) Close(t task.Task) {
	this.closeSig = true
	fmt.Println(fmt.Sprintf("num : %d", this.num))
}

func GetData() {
	client := &http.Client{}
	resp, err := client.Get("https://www.baidu.com")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

//登录
func WxReadLogin() string {
	timeStr:=time.Now().Format("2006-01-02 15:04:05")

	println(timeStr + "=======message - 发送评论=======")

	//params := make(map[string]interface{})
	//params["deviceToken"] = "f457f844d5a8d276ff68a5e51325a7c2bf8b46d1f9c140f2df4633d589c8f06c"
	//params["refreshToken"] = "onb3MjrYiTYHVJty-8XVrmttJLOE@SA4lqeNuR_Co8LPauWOAewAA"
	//params["deviceId"] = "ad15368d332e644562c5e887d0f977a2"
	//params["refCgi"] = "/mobileSync"
	//params["wxToken"] = 1
	//params["inBackground"] = 0
	//bytesData, err := json.Marshal(params)
	//if err != nil {
	//	fmt.Println(err.Error() )
	//}
	//reader := bytes.NewReader(bytesData)
	//
	//client := &http.Client{}
	//request, _ := http.NewRequest("POST","https://i.weread.qq.com/login", reader)
	//
	//header := make(http.Header)
	//header.Set("Content-Type", "application/json")
	//header.Set("User-Agent", "PostmanRuntime/7.4.0")
	//request.Header = header
	//result, err := client.Do(request)
	//defer result.Body.Close()
	//
	//if err != nil {
	//	println("error")
	//}
	//
	//if result.StatusCode != 200 {
	//	println("error")
	//}
	//
	//body, _ := ioutil.ReadAll(result.Body)
	//var dat map[string]interface{}
	//if err := json.Unmarshal([]byte(body), &dat); err == nil {
	//	return dat["skey"].(string)
	//} else {
	//	fmt.Println(err)
	//}
	//
	return  ""
}







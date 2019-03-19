package http_request

import (
	"cron_demo/const"
	"fmt"
	"github.com/tidwall/gjson"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 添加评论
func MessageAdd(id int64) {
	var ok, content string
	fmt.Printf("是否参与评论[y/n]: ")
	fmt.Scanln(&ok)
	if (strings.EqualFold(ok,"n")) {
		println("======= 跳过评论 =======")
	} else {
		//fmt.Printf("请输入评论内容: ")
		//fmt.Scanln(&content)
		content = "测试评论23333"  //默认值
		url := host + "/message/add"
		params := make(map[string]string)
		params["c_type"] = "1"
		params["element_id"] = strconv.FormatInt(id, 10)
		params["comment_content"] = content
		body := PostFormRequest(url, params)
		success := gjson.Get(string(body), "success").Bool()
		//panic(body)
		if (success == false) {
			panic("=======失败=======  " + gjson.Get(string(body), "info").String())
		}
	}
}

// 添加弹幕
func MessageAddDm() {
	var ok, content string
	fmt.Printf("是否添加弹幕[y/n]: ")
	fmt.Scanln(&ok)
	if (strings.EqualFold(ok,"n")) {
		println("======= 跳过添加弹幕 =======")
	} else {
		//fmt.Printf("请输入评论内容: ")
		//fmt.Scanln(&content)
		content = "测试弹幕23333"  //默认值
		url := host + "/message/add-dm"
		params := make(map[string]interface{})
		params["token"] = ""
		params["sound_id"] = ""
		params["stime"] = ""
		params["text"] = content
		body := PostRequest(url, params)
		code := gjson.Get(string(body), "code").Int()
		if (code != int64(0)) {
			panic("=======失败=======" + "添加弹幕失败")
		}
	}
}

func LiveMessage()  {
	url := "https://fm.uat.missevan.com/api/chatroom/message/send"
	params := make(map[string]interface{})
	params["room_id"] = _const.RoomId
	params["message"] = "加油" + GetRandomString(10) + time.Now().Format("2006年01月02日 15：04：05")
	params["msg_id"] = "34243979-c079-4af4-9efa-d88e116870f6"
	//params["msg_id"] =  uuid.Must(uuid.NewV4()).String()
	//panic(params["msg_id"])
	body := PostRequest(url, params)
	code := gjson.Get(string(body), "code").Int()
	panic(body)
	if (code != int64(0)) {
		panic("=======失败=======" + "添加失败")
	}
}

func GetRandomString(l int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}


func LiveGift() {
	url := "https://fm.uat.missevan.com/api/chatroom/gift/send"
	params := make(map[string]interface{})
	params["room_id"] = _const.RoomId
	params["gift_id"] = 1
	params["gift_num"] = 1

	body := PostRequest(url, params)
	code := gjson.Get(string(body), "code").Int()
	panic(body)
	if (code != int64(0)) {
		panic("=======失败=======" + "礼物发送失败")
	}
}

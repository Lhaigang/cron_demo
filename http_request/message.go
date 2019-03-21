package http_request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 发送模板消息

var (
	send_template_url    = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
)


type TemplateMsg struct {
	Touser      string        `json:"touser"`      //接收者的OpenID
	TemplateID  string        `json:"template_id"` //模板消息ID
	URL         string        `json:"url"`         //点击后跳转链接
	Data        *TemplateData `json:"data"`
}
type Miniprogram struct {
	AppID    string `json:"appid"`
	Pagepath string `json:"pagepath"`
}

type TemplateData struct {
	First    KeyWordData `json:"first,omitempty"`
	Keyword1 KeyWordData `json:"keyword1,omitempty"`
	Keyword2 KeyWordData `json:"keyword2,omitempty"`
	Keyword3 KeyWordData `json:"keyword3,omitempty"`
	Keyword4 KeyWordData `json:"keyword4,omitempty"`
	Keyword5 KeyWordData `json:"keyword5,omitempty"`
	Remark   KeyWordData `json:"remark,omitempty"`
}

type KeyWordData struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

type SendTemplateResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	MsgID   int `json:"msgid"`
}

type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func SendTemplate() {
	templateMsg := TemplateMsg{}
	tempData := TemplateData{}
	tempData.First.Value = "测试模板消息"
	tempData.Keyword1.Value = "大家记得买票啊"
	tempData.Keyword2.Value = "马上就要放假了，大家记得买回家的票啊"
	tempData.Remark.Value = "remark"
	templateMsg.Data = &tempData

	templateMsg.Touser = "oban11Jt_4OTpxEjC7TBf6FVoMrc"
	templateMsg.TemplateID = "QLOSTXWVpRwqqOHdfE6VbHp4z9Gl_Sl0ccNimjTju3I"
	response, err := _SendTemplate(&templateMsg)
	if err != nil {
		fmt.Println("发送模板消息失败", err.Error())
	}

	fmt.Println(response)
}



// SendTemplate 发送模板消息
func _SendTemplate(msg *TemplateMsg) (*SendTemplateResponse, error) {
	accessToken := "19_kBHc3krhXjX9m0EI7E-m7UdUwnyTtbRT56c1vVCM6IbSDl4p8P2pRQGBGOhxiKjkDYpbKz9uVDZzID8lzRZiVOg0CD1tK9TuFMxvdE6lKDXnAbexBYIGaOAZI60DsJomyFDUhvqPzFe1LNK-IFPdABAAXH"
	url := fmt.Sprintf(send_template_url, accessToken)
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var templateResponse SendTemplateResponse
	err = json.Unmarshal(body, &templateResponse)
	if err != nil {
		return nil, err
	}
	return &templateResponse, nil
}


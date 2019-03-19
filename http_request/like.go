package http_request

import "github.com/tidwall/gjson"

//猜你喜欢 获取画像 ID
func GetPersona() int64 {
	url := host + "/you-might-like/get-persona"
	body := GetRequest(url)
	info := gjson.Get(string(body), "info").Int()

	return info
}
// you_might_like - 猜你喜欢推荐音
func GetSounds(id int64) {
	url := host + "/you-might-like/get-persona"
	body := GetRequest(url)
	code := gjson.Get(string(body), "code").Int()
	info := gjson.Get(string(body), "info").Array()
	if (code != int64(0)) {
		panic("=======失败=======" + "you_might_like - 猜你喜欢推荐音")
	}
	for _, v := range info {
		front_cover := gjson.Get(v.String(), "front_cover").String()
		cover_image := gjson.Get(v.String(), "cover_image").String()
		soundstr := gjson.Get(v.String(), "soundstr").String()
		url := gjson.Get(v.String(), "url").String()
		if (front_cover == "" || cover_image == "" || CheckResources(front_cover) == false || CheckResources(cover_image)) {
			panic("=======you_might_like - 猜你喜欢推荐音 资源错误=======")
		}
		if ( soundstr  == "" || url == "") {
			panic("=======you_might_like - 猜你喜欢推荐音 标题或跳转链接错误=======")
		}
	}
}

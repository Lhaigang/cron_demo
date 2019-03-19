package http_request

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

const host = "http://app.uat.missevan.com"


//启动接口 是否升级，启动图
func Launch() {
	url := host + "/site/launch"
	body := GetRequest(url)
	code := gjson.Get(string(body), "code").Int()
	pic_url := gjson.Get(string(body), "info.pic_url").String()
	sound_url := gjson.Get(string(body), "info.sound_url").String()
	v := gjson.Get(string(body), "info.new_version.v").String()

	if (code != int64(0)) {
		panic("=======失败=======" + "site - 启动接口 是否升级，启动图");
	}
	if (pic_url == "" || CheckResources(pic_url) == false) {
		println("=======启动图资源缺失=======");
	}
	if (sound_url == "" ||CheckResources(sound_url) == false) {
		panic("=======启动音频资源缺失=======");
	}
	if (v == "") {
		panic("=======版本信息错误=======");
	}
}

//site - 版头图及首页精品周更
func GetTop()  {
	url := host + "/site/get-top"
	body := GetRequest(url)
	code := gjson.Get(string(body), "code").Int()
	banners := gjson.Get(string(body), "info.banners").Array() //banner
	weeklydrama := gjson.Get(string(body), "info.weeklydrama").Array() //精品周更
	reward_rank := gjson.Get(string(body), "info.reward_rank").Array() //广播剧打赏周榜

	if (code != int64(0) ||  len(banners) == 0) {
		panic("=======失败=======" + "site - 版头图及首页精品周更")
	}
	for _, banner := range banners {
		url := gjson.Get(banner.String(), "url").String()
		pic := gjson.Get(banner.String(), "pic").String()
		if (url == "" || pic == "" || CheckResources(url) == false || CheckResources(pic) == false) {
			println("=======版头图及首页精品周更 资源错误=======")
		}
	}

	for _, v := range weeklydrama {
		cover := gjson.Get(v.String(), "cover").String()
		author := gjson.Get(v.String(), "author").String()
		name := gjson.Get(v.String(), "name").String()
		if (cover == "" || CheckResources(cover) == false) {
			println("=======精品周更 资源错误=======")
		}
		if (author == "" || name  == "") {
			panic("=======精品周更 标题或作者错误=======")
		}
	}

	for _, rank := range reward_rank {
		cover := gjson.Get(rank.String(), "cover").String()
		name := gjson.Get(rank.String(), "name").String()
		if (cover == "" || CheckResources(cover) == false) {
			println("=======广播剧打赏周榜 资源错误=======")
		}
		if ( name  == "") {
			panic("=======广播剧打赏周榜 标题或作者错误=======")
		}
	}
}

//site - 首页小图标
func GetIcons()  {
	url := host + "/site/icons"
	body := GetRequest(url)
	code := gjson.Get(string(body), "code").Int()
	info := gjson.Get(string(body), "info").Array()
	if (code != int64(0)) {
		panic("=======失败=======" + "site - 首页小图标")
	}
	for _, v := range info {
		icon := gjson.Get(v.String(), "icon").String()
		title := gjson.Get(v.String(), "title").String()
		url := gjson.Get(v.String(), "url").String()
		if (icon == "" || CheckResources(icon) == false) {
			println("=======site - 首页小图标 资源错误=======")
		}
		if ( title  == "" || url == "") {
			panic("=======site - 首页小图标 标题或跳转链接错误=======")
		}
	}

}

//site - 分类根信息
func GetCatalogRoot()  {
	url := host + "/site/catalog-root?dark=0&extra=1"
	body := GetRequest(url)
	code := gjson.Get(string(body), "code").Int()
	info := gjson.Get(string(body), "info").Array()
	if (code != int64(0)) {
		panic("=======失败=======" + "site - 分类根信息")
	}

	for _, v := range info {
		pic := gjson.Get(v.String(), "pic").String()
		title := gjson.Get(v.String(), "title").String()
		if (pic == "" || CheckResources(pic) == false) {
			println("=======site - 分类根信息 资源错误=======")
		}
		if ( title  == "" ) {
			panic("=======site - 分类根信息 标题错误=======")
		}
	}

}

//site - 根据分类 ID 获取分类数据
func GetCatalogHomepage()  {
	url := host + "/site/catalog-homepage?cid=46&dark=0"
	body := GetRequest(url)
	code := gjson.Get(string(body), "code").Int()
	banners := gjson.Get(string(body), "info.banner").Array() //banner
	sounds := gjson.Get(string(body), "info.sounds").Array() //banner
	if (code != int64(0) ||  len(banners) == 0) {
		panic("=======失败=======" + "site - 根据分类 ID 获取分类数据")
	}
	for _, banner := range banners {
		url := gjson.Get(banner.String(), "url").String()
		title := gjson.Get(banner.String(), "title").String()
		pic_app := gjson.Get(banner.String(), "pic_app").String()
		if (url == "" || pic_app == "" || CheckResources(url) == false || CheckResources(pic_app) == false) {
			println("=======根据分类 ID 获取分类数据 资源错误=======")
		}
		if ( title  == "" ) {
			panic("=======site - 根据分类 ID 获取分类数据 标题错误=======")
		}
	}

	for _, sound := range sounds {
		name := gjson.Get(sound.String(), "name").String()
		items := gjson.Get(sound.String(), "item").Array()
		if ( name  == "" ) {
			panic("=======site - 根据分类 ID 获取分类数据 标题错误=======")
		}
		id := int64(0)
		n := 0
		for _, item := range items {
			if n > 0 {
				break
			}
			id = gjson.Get(item.String(), "id").Int()
			soundstr := gjson.Get(item.String(), "soundstr").String()
			view_count := gjson.Get(item.String(), "view_count").String()
			all_comments := gjson.Get(item.String(), "all_comments").String()
			soundurl := gjson.Get(item.String(), "soundurl").String()
			front_cover := gjson.Get(item.String(), "front_cover").String()
			if (!strings.Contains(soundurl,"https")) {
				soundurl = "https://static.missevan.com/sound/" + soundurl
			}
			if (soundurl == "" || CheckResources(soundurl) == false) {
				println("=======根据分类 ID 获取分类数据 资源错误======="+soundurl)
			}
			if (front_cover == "" || CheckResources(front_cover) == false) {
				println("=======根据分类 ID 获取分类数据 资源错误=======")
			}
			if ( soundstr  == "" || view_count == "" || all_comments == "") {
				panic("=======site - 根据分类 ID 获取分类数据 标题错误=======")
			}
			n++
		}

		detailUrl := host + "/sound/sound?sound_id=" + strconv.FormatInt(id,10)
		detail := GetRequest(detailUrl)
		detailCode := gjson.Get(string(detail), "code").Int()
		if (detailCode != int64(0)) {
			panic("=======失败=======" + "播放详情")
		}
		username := gjson.Get(string(detail), "info.username").String()
		soundurl1 := gjson.Get(string(detail), "info.soundurl").String()
		soundstr := gjson.Get(string(detail), "info.soundstr").String()
		if (soundurl1 == "" || CheckResources(soundurl1) == false) {
			println("=======根据分类 ID 获取分类数据 资源错误=======")
		}
		if ( username  == "" || soundstr == "" ) {
			panic("=======site - 根据分类 ID 获取分类数据 标题错误=======")
		}
		GetDrame(id)
	}
}

//剧集详情
func GetDrame(id int64) {

	url := host + "/drama/get-drama-by-sound?sound_id=" + strconv.FormatInt(id,10)
	body := GetRequest(url)
	code := gjson.Get(string(body), "code").Int()
	episodes := gjson.Get(string(body), "info.episodes.episode").Array()
	cvs := gjson.Get(string(body), "info.cvs").Array()
	if (code != int64(0)) {
		panic("=======失败=======" + "drama - 根据单音 ID 获取剧集 ID")
	}
	if (len(episodes) > 0) { // 有剧集，就有推荐剧集
		GetSoundRecommend(id)
	}
	if (len(cvs) == 0) {
		panic("=======失败=======" + "剧集没有cvs（配音）")
	}

	for _, v := range cvs {
		character := gjson.Get(v.String(), "character").String()
		name := gjson.Get(v.String(), "cvinfo.name").String()
		icon := gjson.Get(v.String(), "cvinfo.icon").String()
		if ( character == "") {
			panic("=======site - drama - 剧集 cvs 角色信息缺失 =======")
		}
		if ( name  == "" ) {
			panic("=======site - drama - 剧集 cvs 名称信息缺失 =======")
		}
		if (icon == "" || CheckResources(icon) == false) {
			println("=======剧集 cvs 作者头像缺失 =======")
		}
	}
	MessageAdd(id)
}

func GetSoundRecommend(id int64)  {
	url := host + "/sound/recommend?sound_id=" + strconv.FormatInt(id,10)
	body := GetRequest(url)
	code := gjson.Get(string(body), "code").Int()
	if (code != int64(0)) {
		panic("=======失败=======" + "推荐剧集获取")
	}
	drama := gjson.Get(string(body), "info.drama").Array() //剧集
	sound := gjson.Get(string(body), "info.sound").Array() //相似音频
	album := gjson.Get(string(body), "info.album").Array() //音单
	if (len(drama) == 0) {
		panic("=======失败=======" + "推荐剧集获取为空")
	}
	if (len(sound) == 0) {
		panic("=======失败=======" + "相似音频获取为空")
	}
	if (len(album) == 0) {
		panic("=======失败=======" + "剧集包含音单获取为空")
	}
	for _, v := range drama {
		name := gjson.Get(v.String(), "name").String()
		front_cover := gjson.Get(v.String(), "front_cover").String()
		if ( name  == "") {
			panic("======= 推荐剧集标题缺失 =======")
		}
		if (front_cover == "" || CheckResources(front_cover) == false) {
			println("======= 推荐剧集图片缺失 =======")
		}
	}
	for _, sv := range sound {
		soundstr := gjson.Get(sv.String(), "soundstr").String()
		front_cover := gjson.Get(sv.String(), "front_cover").String()
		if ( soundstr  == "") {
			panic("======= 相似音频标题缺失 =======")
		}
		if (front_cover == "" || CheckResources(front_cover) == false) {
			println("======= 相似音频图片缺失 =======")
		}
	}
	for _, av := range album {
		title := gjson.Get(av.String(), "title").String()
		front_cover := gjson.Get(av.String(), "front_cover").String()
		if ( title  == "") {
			panic("======= 剧集包含音单标题缺失 =======")
		}
		if (front_cover == "" || CheckResources(front_cover) == false) {
			println("======= 剧集包含音单图片缺失 =======")
		}
	}
}

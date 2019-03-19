package http_request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// get请求
func GetRequest(url string) string {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	header := make(http.Header)
	header.Set("User-Agent", "MissEvanApp/4.2.4 (iOS;12.1;iPhone9,1)")
	header.Set("equip-code", "93a34a10-76bc-458a-9b66-4b93feba24bf")
	header.Set("Accept", "application/json")
	header.Set("Cookie", "equip_code=93a34a10-76bc-458a-9b66-4b93feba24bf; token=5bea88c7d602601d70c6a7b8|e2226092542acb9a|1542097095|22598632d81c01e4; Hm_lvt_91a4e950402ecbaeb38bd149234eb7cc=1543301210")

	request.Header = header
	result, err := client.Do(request)

	defer result.Body.Close()

	if err != nil {
		println(url + ".error")
	}

	if result.StatusCode != 200 {
		println(url + ".error")
	}

	body, _ := ioutil.ReadAll(result.Body)

	return string(body)
}

// post请求
func PostRequest(url string, params map[string]interface{}) string {
	bytesData, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err.Error() )
	}
	reader := bytes.NewReader(bytesData)

	client := &http.Client{}
	request, _ := http.NewRequest("POST", url, reader)

	header := make(http.Header)
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
	header.Set("equip-code", "93a34a10-76bc-458a-9b66-4b93feba24bf")
	header.Set("Accept", "application/json")
	header.Set("Content-Type", "application/json")
	header.Set("Cookie", "Hm_lpvt_91a4e950402ecbaeb38bd149234eb7cc=1544538498; Hm_lvt_91a4e950402ecbaeb38bd149234eb7cc=1544535734; FM_SESS=20181211|7wm3rqxlhdtygtjg3xppg13yn; FM_SESS.sig=RLsOZ57YLOBBa2-KUWrT51RLLjY; token=5c0fc96df49b4745e16d20a4%7Ccd451b66b6662292%7C1544538477%7C9e1752a33a52a9d4; MUATSESSID=8eec5c16873dba647a736ea8e7d9e79c")
	//header.Set("token", "5c0dd4b7f49b4745e16d206d|6b391eda92b8f289|1544410295|87de4f1f4deacccd")

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
	return string(body)
}

// post请求
func PostFormRequest(url1 string, params map[string]string) string {

	postValue := url.Values{}
	for key, value := range params{
		postValue.Set(key, value)
	}
	client := &http.Client{}
	request, _ := http.NewRequest("POST", url1, strings.NewReader(postValue.Encode()))

	header := make(http.Header)
	header.Set("User-Agent", "MissEvanApp/4.2.4 (iOS;12.1;iPhone9,1)")
	header.Set("equip-code", "93a34a10-76bc-458a-9b66-4b93feba24bf")
	header.Set("Accept", "application/json")
	header.Set("Cookie", "Hm_lpvt_91a4e950402ecbaeb38bd149234eb7cc=1544538498; Hm_lvt_91a4e950402ecbaeb38bd149234eb7cc=1544535734; FM_SESS=20181211|7wm3rqxlhdtygtjg3xppg13yn; FM_SESS.sig=RLsOZ57YLOBBa2-KUWrT51RLLjY; token=5c0fc96df49b4745e16d20a4%7Ccd451b66b6662292%7C1544538477%7C9e1752a33a52a9d4; MUATSESSID=8eec5c16873dba647a736ea8e7d9e79c")
	header.Set("token", "5c0dd4b7f49b4745e16d206d|6b391eda92b8f289|1544410295|87de4f1f4deacccd")
	header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

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
	return string(body)
}

//判读资源是否存在
func CheckResources(url string) bool {

	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	header := make(http.Header)
	header.Set("User-Agent", "MissEvanApp/4.2.4 (iOS;12.1;iPhone9,1)")
	header.Set("equip-code", "93a34a10-76bc-458a-9b66-4b93feba24bf")
	header.Set("Accept", "application/json")
	header.Set("Cookie", "equip_code=93a34a10-76bc-458a-9b66-4b93feba24bf; token=5bea88c7d602601d70c6a7b8|e2226092542acb9a|1542097095|22598632d81c01e4; Hm_lvt_91a4e950402ecbaeb38bd149234eb7cc=1543301210")

	request.Header = header
	result, err := client.Do(request)

	defer result.Body.Close()

	if err != nil {
		println( "error")
	}

	if result.StatusCode != 200 {
		println("=========资源获取失败=======" + url)
		return false
	}
	return true
}

package http_request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
	header.Set("Accept", "application/json")
	header.Set("Content-Type", "application/x-www-form-urlencoded")

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




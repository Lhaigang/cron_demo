package utils

import (
	"encoding/base64"
	"encoding/json"
)

// 生成头部信息加密
func EncodeString(content map[string]interface{}) string {

	b, _ := json.Marshal(content)
	input := []byte(b)
	encodeString := base64.StdEncoding.EncodeToString(input)

	return encodeString
}

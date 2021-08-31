package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// https://www.sojson.com/
// https://mholt.github.io/json-to-go/
type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}
type Items struct {
	ID2   int    `json:"Id2,omitempty"`
	Title string `json:"Title"`
	ID    int    `json:"Id,omitempty"`
	Test  int    `json:"Test,omitempty"`
}

type Data struct {
	Items []*Items `json:"items"`
}

func main() {

	// json字符串数组,转换成切片
	response := `{
		"code": 0,
		"message": "操作成功",
		"data": {
			"items": [{
				"Id": 100,
				"Title": "木华黎"
			}, {
				"Id": 200,
				"Title": "耶律楚才"
			}, {
				"Id": 300,
				"Title": "纳呀啊",
				"Test": 100
			}]
		}
	}`

	res := Response{}
	data := Data{}

	err := json.Unmarshal([]byte(response), &res)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	fmt.Println(res)

	err = json.Unmarshal(res.Data, &data)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(reflect.TypeOf(data.Items), data.Items, *data.Items[0])
}

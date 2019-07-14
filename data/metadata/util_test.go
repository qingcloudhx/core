package metadata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
* @Author: hexing
* @Date: 19-7-14 下午9:44
 */
func TestStructToMap(t *testing.T) {
	in := struct {
		Id   string                   `md:"id"`
		Data []map[string]interface{} `md:"data"`
	}{
		Id: "hexing",
		Data: []map[string]interface{}{
			{
				"id":   "1212",
				"name": "hexing",
			},
			{

				"id":   "1212",
				"name": "hexing",
			},
		},
	}
	out := StructToMap(in)
	t.Log(out)
}
func TestMapToStruct(t *testing.T) {
	in := map[string]interface{}{
		"id": "xxxx",
		"data": []map[string]interface{}{
			{
				"id":   "1212",
				"name": "hexing",
			},
			{
				"id":   "1212",
				"name": "hexing",
			},
		},
	}
	out := struct {
		Id   string                   `md:"id"`
		Data []map[string]interface{} `md:"data"`
	}{}
	err := MapToStruct(in, &out, true)
	assert.Nil(t, err)
	t.Log(out)
}

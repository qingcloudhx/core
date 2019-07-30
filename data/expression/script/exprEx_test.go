package script

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
* @Author: hexing
* @Date: 19-7-30 下午5:48
 */

func TestGetJSRuntime(t *testing.T) {
	v := GetJSRuntime()
	r, err := v.RunString("1 + 2")
	assert.Nil(t, err)
	t.Log(r.Export())
}

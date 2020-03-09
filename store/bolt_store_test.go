/**
 * @Author: hexing
 * @Description:
 * @File:  bolt_store_test
 * @Version: 1.0.0
 * @Date: 20-3-9 上午10:32
 */

package store

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNewBoltStore(t *testing.T) {
	s, err := NewBoltStore("./hexing")
	assert.Nil(t, err)
	for i := 0; i < 100; i++ {
		err = s.Set([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)+"hexing"))
	}
	for i := 0; i < 100; i++ {
		value, err := s.Get([]byte(strconv.Itoa(i)))
		assert.Nil(t, err)
		t.Log(string(value))
	}
	err = s.ForEach(func(k, v []byte) error {
		t.Log(string(k), string(v))
		return nil
	})
	assert.Nil(t, err)
	for i := 0; i < 100; i++ {
		err := s.Del([]byte(strconv.Itoa(i)))
		assert.Nil(t, err)
	}
}

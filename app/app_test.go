package app

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	//conf := &Config{}
	//if _, err := New(conf, nil, nil); err != nil {
	//	t.Log(err)
	//}
	start := time.Now()
	fmt.Println(time.Since(start))
}

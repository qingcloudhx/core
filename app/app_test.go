package app

import "testing"

func TestNew(t *testing.T) {
	conf := &Config{}
	if _, err := New(conf, nil, nil); err != nil {
		t.Log(err)
	}

}

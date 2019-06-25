package sample

import (
	"testing"

	"github.com/qingcloudhx/core/support"
	"github.com/qingcloudhx/core/trigger"
	"github.com/stretchr/testify/assert"
)

func TestTrigger_Register(t *testing.T) {

	ref := support.GetRef(&Trigger{})
	f := trigger.GetFactory(ref)
	assert.NotNil(t, f)
}

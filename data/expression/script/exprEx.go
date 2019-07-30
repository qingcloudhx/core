package script

/**
* @Author: hexing
* @Date: 19-7-30 下午5:41
 */
import (
	"github.com/dop251/goja"
)

var vm *goja.Runtime

func init() {
	vm = goja.New()
}
func GetJSRuntime() *goja.Runtime {
	return vm
}

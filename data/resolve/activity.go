package resolve

import (
	"fmt"
	"github.com/qingcloudhx/core/data"
)

/**
* @Author: hexing
* @Date: 19-8-2 下午3:08
 */
var dynamicItemResolver = NewResolverInfo(false, false)

type ActivityResolver struct {
}

func (r *ActivityResolver) GetResolverInfo() *ResolverInfo {
	return dynamicItemResolver
}

func (r *ActivityResolver) Resolve(scope data.Scope, itemName, valueName string) (interface{}, error) {

	value, exists := scope.GetValue("_A." + itemName + "." + valueName)
	if !exists {
		return nil, fmt.Errorf("failed to resolve activity attr: '%s', not found in activity '%s'", valueName, itemName)
	}

	return value, nil
}

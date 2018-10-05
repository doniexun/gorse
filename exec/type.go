package exec

import (
	"github.com/ZhangZhenghao/gorse/core"
	"github.com/yuin/gopher-lua"
)

// Convert LTable to parameters
func LTAsParameters(table *lua.LTable) core.Parameters {
	params := core.Parameters{}
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		name := lua.LVAsString(key)
		params[name] = value
	})
	return params
}

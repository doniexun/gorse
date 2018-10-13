package exec

import (
	"github.com/yuin/gopher-lua"
	"github.com/zhenghaoz/gorse/core"
)

type _Model func(core.Parameters) core.Model

var models = map[string]_Model{
	"random": func(params core.Parameters) core.Model { return core.NewRandom(params) },
}

func NewExecutor() *lua.LState {
	L := lua.NewState()
	// Add models
	for name, creator := range models {
		L.SetGlobal(name, L.NewFunction(func(L *lua.LState) int {
			userData := L.NewUserData()
			userData.Value = creator(nil)
			L.Push(userData)
			return 1
		}))
	}
	// Add functions
	for name, function := range functions {
		L.SetGlobal(name, L.NewFunction(function))
	}
	return L
}

package exec

import (
	"github.com/yuin/gopher-lua"
	"github.com/zhenghaoz/gorse/core"
	"gonum.org/v1/gonum/stat"
	"runtime"
)

var functions = map[string]lua.LGFunction{
	// load(name)
	//   name	- The name of data set.
	"load": func(L *lua.LState) int {
		name := L.ToString(1)
		dataSet := L.NewUserData()
		dataSet.Value = core.LoadDataFromBuiltIn(name)
		L.Push(dataSet)
		return 1
	},

	// len(data):
	//	 data	- The data set.
	"len": func(L *lua.LState) int {
		dataset := L.ToUserData(1).Value.(core.DataSet)
		L.Push(lua.LNumber(dataset.Length()))
		return 1
	},

	// fit(model, data)
	//   model	- The model.
	//   data	- The train data set.
	"fit": func(L *lua.LState) int {
		estimator := L.ToUserData(1).Value.(core.Model)
		trainSet := L.ToUserData(2).Value.(core.TrainSet)
		estimator.Fit(trainSet)
		return 0
	},

	// predict(model, u, i)
	//   model	- The model.
	//   u      - The user's ID.
	//   i      - The item's ID.
	"predict": func(L *lua.LState) int {
		estimator := L.ToUserData(1).Value.(core.Model)
		userId := L.ToInt(2)
		itemId := L.ToInt(3)
		L.Push(lua.LNumber(estimator.Predict(userId, itemId)))
		return 1
	},

	// cv(model)
	//   model	- The model.
	"cv": func(L *lua.LState) int {
		estimator := L.ToUserData(1).Value.(core.Model)
		dataSet := L.ToUserData(2).Value.(core.DataSet)
		cv := core.CrossValidate(estimator, dataSet, []core.Evaluator{core.RMSE},
			core.NewKFoldSplitter(5), 0, nil, runtime.NumCPU())
		L.Push(lua.LNumber(stat.Mean(cv[0].Tests, nil)))
		return 1
	},
}

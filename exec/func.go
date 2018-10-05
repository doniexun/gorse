package exec

import (
	"fmt"
	"github.com/ZhangZhenghao/gorse/core"
	"github.com/yuin/gopher-lua"
	"gonum.org/v1/gonum/stat"
	"reflect"
	"runtime"
)

func SetGlobalUserData(L *lua.LState, name string, value interface{}) {
	lValue := L.NewUserData()
	lValue.Value = value
	L.SetGlobal(name, lValue)
}

func NewEngine() *lua.LState {
	L := lua.NewState()
	SetGlobalUserData(L, "rmse", core.RMSE)
	SetGlobalUserData(L, "MAE", core.MAE)
	L.SetGlobal("load", L.NewFunction(Load))
	L.SetGlobal("len", L.NewFunction(Len))
	L.SetGlobal("cv", L.NewFunction(CV))
	L.SetGlobal("estimator", L.NewFunction(Estimator))
	return L
}

/* Data Set */

// load(name)
// 	 name	- THe name of data set.
func Load(L *lua.LState) int {
	// Retrieve arguments
	dataSetName := L.ToString(1)
	dataSet := L.NewUserData()
	dataSet.Value = core.LoadDataFromBuiltIn(dataSetName)
	L.Push(dataSet)
	return 1
}

// len(dataSet)
//	 dataSet	- The data set.
func Len(L *lua.LState) int {
	lDataSet := L.ToUserData(1)
	switch dataSet := lDataSet.Value.(type) {
	case core.DataSet:
		L.Push(lua.LNumber(dataSet.Length()))
		return 1
	default:
		L.RaiseError("Expect core.DataSet, receive %v", reflect.TypeOf(dataSet))
		return 0
	}
}

/* Estimator */

// estimator(name, parameters)
//   name	    - The name of the estimator.
//	 parameters	- Parameters for the estimator.
func Estimator(L *lua.LState) int {
	// Retrieve parameters
	name := L.ToString(1)
	params := core.Parameters{}
	if L.GetTop() > 1 {
		params = LTAsParameters(L.ToTable(2))
	}
	estimator := L.NewUserData()
	switch name {
	case "random":
		estimator.Value = core.NewRandom(params)
	case "baseline":
		estimator.Value = core.NewBaseLine(params)
	case "svd":
		estimator.Value = core.NewSVD(params)
	case "svd++":
		estimator.Value = core.NewSVDpp(params)
	case "nmf":
		estimator.Value = core.NewNMF(params)
	case "knn":
		estimator.Value = core.NewKNN(params)
	case "slopOne":
		estimator.Value = core.NewSlopOne(params)
	case "coClustering":
		estimator.Value = core.NewCoClustering(params)
	}
	L.Push(estimator)
	return 1
}

// fit(estimator, trainSet)
//   estimator	- The estimator.
//   trainSet	- The train data set.
func Fit(L *lua.LState) int {
	// Retrieve arguments
	estimator := L.ToUserData(1).Value.(core.Estimator)
	trainSet := L.ToUserData(2).Value.(core.TrainSet)
	estimator.Fit(trainSet)
	return 0
}

// cv(estimator, dataSet)
func CV(L *lua.LState) int {
	estimator := L.ToUserData(1).Value.(core.Estimator)
	dataSet := L.ToUserData(2).Value.(core.DataSet)
	//evaluators := make([]core.Evaluator, 0)
	L.ToTable(3).ForEach(func(key lua.LValue, value lua.LValue) {
		//evaluators = append(evaluators, lua.LUserData(value))
		fmt.Println(value)
	})
	cv := core.CrossValidate(estimator, dataSet, []core.Evaluator{core.RMSE},
		core.NewKFoldSplitter(5), 0, nil, runtime.NumCPU())
	L.Push(lua.LNumber(stat.Mean(cv[0].Tests, nil)))
	return 1
}

// predict(estimator, userId, itemId)
func Predict(L *lua.LState) int {
	lEstimator := L.ToUserData(1)
	userId := L.ToInt(2)
	itemId := L.ToInt(3)
	estimator := lEstimator.Value.(core.Estimator)
	L.Push(lua.LNumber(estimator.Predict(userId, itemId)))
	return 1
}

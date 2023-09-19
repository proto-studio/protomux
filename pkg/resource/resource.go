package resource

import (
	"errors"
	"fmt"

	"proto.zip/studio/mux/pkg/tokenizer"
)

type Resource[RequestHandlerType any] struct {
	methods  map[string]RequestHandlerType
	paramMap map[string][]tokenizer.Token
}

func New[H any]() *Resource[H] {
	return &Resource[H]{
		methods:  make(map[string]H),
		paramMap: make(map[string][]tokenizer.Token),
	}
}

func (rh *Resource[H]) Method(methodName string) (H, bool) {
	handler, existing := rh.methods[string(methodName)]
	return handler, existing
}

func (rh *Resource[H]) Methods() []string {
	keys := make([]string, 0, len(rh.methods))
	for k := range rh.methods {
		keys = append(keys, k)
	}
	return keys
}

func (rh *Resource[H]) HandleMethod(methodName string, handler H) {
	nameStr := methodName
	_, existing := rh.methods[nameStr]

	if existing {
		panic(errors.New("can only be called once per method"))
	}

	rh.methods[nameStr] = handler
}

func (rh *Resource[H]) SetParamNames(methodName string, paramNames []tokenizer.Token) {
	nameStr := string(methodName)
	_, existing := rh.paramMap[nameStr]

	if existing {
		panic(errors.New("can only be called once per method"))
	}

	rh.paramMap[nameStr] = paramNames
}

func (rh *Resource[H]) ParamMap(methodName string, paramValues []tokenizer.Token) map[string]string {
	paramNames, ok := rh.paramMap[string(methodName)]

	if !ok && len(paramValues) == 0 {
		return nil
	}

	if len(paramNames) != len(paramValues) {
		panic(fmt.Errorf("mismatched parameter length: configured with %d name(s) got %d value(s)", len(paramNames), len(paramValues)))
	}

	paramMap := make(map[string]string, len(paramNames))

	for paramIdx, paramName := range paramNames {
		paramMap[string(paramName)] = string(paramValues[paramIdx])
	}

	return paramMap
}

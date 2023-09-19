package host

import (
	"fmt"
	"strings"

	"proto.zip/studio/mux/internal/routetree"
	"proto.zip/studio/mux/internal/tokenizers"
	"proto.zip/studio/mux/pkg/resource"
	"proto.zip/studio/mux/pkg/tokenizer"
)

type Host[RequestHandlerType any, ErrorHandlerType any] struct {
	routes       routetree.Node[resource.Resource[RequestHandlerType]]
	params       []tokenizer.Token
	ErrorHandler ErrorHandlerType
}

func New[RH any, EH any]() *Host[RH, EH] {
	return &Host[RH, EH]{
		routes: routetree.NewWildcardNode[resource.Resource[RH]](),
	}
}

func NewWithParams[RH any, EH any](params []tokenizer.Token) *Host[RH, EH] {
	return &Host[RH, EH]{
		params: params,
		routes: routetree.NewWildcardNode[resource.Resource[RH]](),
	}
}

func (h *Host[RH, EH]) Resource(path []byte) (*resource.Resource[RH], []tokenizer.Token) {
	tok := tokenizers.NewPathTokenizer(path)

	var paramValues []tokenizer.Token

	node := h.routes
	token, _, _ := tok.Next()
	for node != nil && token != nil {
		node = node.Child(token)

		if node != nil && node.Dynamic() {
			if paramValues == nil {
				paramValues = make([]tokenizer.Token, 0, 1)
			}
			paramValues = append(paramValues, token)
		}

		token, _, _ = tok.Next()
	}

	if node == nil {
		return nil, nil
	}

	return node.Value(), paramValues
}

func (h *Host[RH, EH]) NewResource(pathPattern []byte) (*resource.Resource[RH], []tokenizer.Token, error) {
	tok := tokenizers.NewPathPatternTokenizer(pathPattern)

	node := h.routes
	token, tokenType, err := tok.Next()
	if err != nil {
		return nil, nil, err
	}

	var paramNames []tokenizer.Token

	for token != nil {
		parent := node

		if tokenType == tokenizer.TokenTypeLabel {
			if paramNames == nil {
				paramNames = make([]tokenizer.Token, 0, 1)
			}
			paramNames = append(paramNames, token)
		}

		node = node.Child(token)
		if node == nil {
			if tokenType == tokenizer.TokenTypeLabel {
				node = routetree.NewWildcardNode[resource.Resource[RH]]()
			} else {
				node = routetree.NewLiteralNode[resource.Resource[RH]](token)
			}
			parent.AddChild(node)
		}

		token, tokenType, err = tok.Next()
		if err != nil {
			return nil, nil, err
		}
	}

	r := node.Value()
	if r == nil {
		r = resource.New[RH]()
		node.SetValue(r)
	}
	return r, paramNames, nil
}

func (h *Host[RH, EH]) Handle(method, path string, handler RH) {
	resource, paramNames, err := h.NewResource([]byte(path))

	if err != nil {
		panic(err)
	}

	// User supplied input so we convert to upper case for ease of use.
	methodUpper := strings.ToUpper(method)

	if len(paramNames) > 0 {
		resource.SetParamNames(methodUpper, paramNames)
	}

	resource.HandleMethod(methodUpper, handler)
}

func (h *Host[RH, EH]) ParamMap(paramValues []tokenizer.Token) map[string]string {
	if len(h.params) == 0 {
		return nil
	}

	if len(h.params) != len(paramValues) {
		panic(fmt.Errorf("mismatched parameter length: configured with %d name(s) got %d value(s)", len(h.params), len(paramValues)))
	}

	paramMap := make(map[string]string, len(h.params))

	for paramIdx, paramName := range h.params {
		paramMap[string(paramName)] = string(paramValues[paramIdx])
	}

	return paramMap
}

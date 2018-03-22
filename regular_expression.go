package cucumberexpressions

import (
	"fmt"
	"regexp"
)

type RegularExpression struct {
	expressionRegexp      *regexp.Regexp
	parameterTypeRegistry *ParameterTypeRegistry
	treeRegexp            *TreeRegexp
}

func NewRegularExpression(expressionRegexp *regexp.Regexp, parameterTypeRegistry *ParameterTypeRegistry) *RegularExpression {
	return &RegularExpression{
		expressionRegexp:      expressionRegexp,
		parameterTypeRegistry: parameterTypeRegistry,
		treeRegexp:            NewTreeRegexp(expressionRegexp),
	}
}

func (r *RegularExpression) Match(text string) ([]*Argument, error) {
	parameterTypes := []*ParameterType{}
	for _, groupBuilder := range r.treeRegexp.GroupBuilder().Children() {
		parameterTypeRegexp := groupBuilder.Source()
		parameterType, err := r.parameterTypeRegistry.LookupByRegexp(parameterTypeRegexp)
		if err != nil {
			return nil, err
		}
		if parameterType == nil {
			parameterType, err = NewParameterType(
				parameterTypeRegexp,
				[]*regexp.Regexp{regexp.MustCompile(parameterTypeRegexp)},
				"string",
				func(arg3 ...string) interface{} {
					return arg3[0]
				},
				false,
				false,
			)
			if err != nil {
				return nil, err
			}
		}
		parameterTypes = append(parameterTypes, parameterType)
	}
	args := BuildArguments(r.treeRegexp, text, parameterTypes)
	fmt.Printf("%#v\n", args)
	return args, nil
}

func (r *RegularExpression) Regexp() *regexp.Regexp {
	return r.expressionRegexp
}

func (r *RegularExpression) Source() string {
	return r.expressionRegexp.String()
}

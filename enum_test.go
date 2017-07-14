package gee

import "testing"

type ComplexLiteral struct {
	*literalBase
}

func (o *ComplexLiteral) Order() int {
	enum := o.Enum.(*complexEnums)
	return enum.orders[o.i]
}

func (o *ComplexLiteral) IsLitName1() bool {
	return o == _complexEnums.LitName1()
}

func (o *ComplexLiteral) IsLitName2() bool {
	return o == _complexEnums.LitName2()
}

type complexEnums struct {
	*enumBase
	orders []int
}

var _complexEnums = &complexEnums{orders: []int{7, 15}}
var _enum = NewEnum([]string{"LitName1", "LitName2"},
	func(enum *enumBase) Enum {
		_complexEnums.enumBase = enum
		return _complexEnums
	},
	func(literal *literalBase) Literal {
		return &ComplexLiteral{literalBase: literal}
	})

func ComplexEnums() *complexEnums {
	return _complexEnums
}

func (o *complexEnums) LitName1() *ComplexLiteral {
	return o.Literals()[0].(*ComplexLiteral)
}

func (o *complexEnums) LitName2() *ComplexLiteral {
	return o.Literals()[1].(*ComplexLiteral)
}

func TestComplexEnums(t *testing.T) {
	lit := ComplexEnums().LitName1()
	println(lit.Name())
	println(lit.Order())
}

package enum

import (
	"testing"
	"fmt"
)

type ComplexLiteral struct {
	*literalBase
	vendor string
}

func (o *ComplexLiteral) Vendor() string {
	return o.vendor
}

func (o *ComplexLiteral) IsLitName1() bool {
	return o == _complexEnums.LitName1()
}

func (o *ComplexLiteral) IsLitName2() bool {
	return o == _complexEnums.LitName2()
}

type complexEnums struct {
	*enumBase
	vendors []string
}

var vendors = []string{"Error", "NoError"}
var _complexEnums = &complexEnums{NewEnum([]string{"LitName1Vendor", "VendorLitName2"},
	func(literal *literalBase) Literal {
		return &ComplexLiteral{literalBase: literal, vendor: vendors[literal.ordinal]}
	}), vendors}

func ComplexEnums() *complexEnums {
	return _complexEnums
}

func (o *complexEnums) LitName1() *ComplexLiteral {
	return o.Get(0)
}

func (o *complexEnums) LitName2() *ComplexLiteral {
	return o.Get(1)
}

func (o *complexEnums) Get(ordinal int) *ComplexLiteral {
	return o.Literals()[ordinal].(*ComplexLiteral)
}

func TestComplexEnums(t *testing.T) {
	lit := ComplexEnums().LitName1()
	println(lit.Name())
	println(lit.Vendor())
	println(fmt.Sprintf("%v", ComplexEnums().Get(1)))
}

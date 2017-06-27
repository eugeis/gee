package gee

type ComplexEnum struct {
	name  string
	order int
}

type ComplexEnums struct {
	values []*ComplexEnum
}

var _complexEnums = &ComplexEnums{values: []*ComplexEnum{
	{name: "", order: 0},
}}

func ComplexEnumss() *ComplexEnums {
	return _complexEnums
}

func (o *ComplexEnums) LitName1() *ComplexEnum {
	return o.values[0]
}

func (o *ComplexEnums) Values() []*ComplexEnum {
	return o.values
}

func (o *ComplexEnum) IsLitName1() bool {
	return o == _complexEnums.LitName1()
}
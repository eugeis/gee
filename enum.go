package gee

type ComplexEnum struct {
	name  string
	ordinal int
	order int
}

func (o *ComplexEnum) Name() string {
	return o.name
}

func (o *ComplexEnum) Ordinal() int {
	return o.ordinal
}

func (o *ComplexEnum) Order() int {
	return o.order
}

func (o *ComplexEnum) IsLitName1() bool {
	return o == _complexEnums.LitName1()
}

func (o *ComplexEnum) IsLitName2() bool {
	return o == _complexEnums.LitName2()
}

type complexEnums struct {
	values []*ComplexEnum
}

var _complexEnums = &complexEnums{values: []*ComplexEnum{
	{name: "LitName1", ordinal: 0, order: 0},
	{name: "LitName2", ordinal: 1, order: 1}},
}

func ComplexEnums() *complexEnums {
	return _complexEnums
}

func (o *complexEnums) Values() []*ComplexEnum {
	return o.values
}

func (o *complexEnums) LitName1() *ComplexEnum {
	return _complexEnums.values[0]
}

func (o *complexEnums) LitName2() *ComplexEnum {
	return _complexEnums.values[1]
}

func (o *complexEnums) ParseComplexEnum(name string) (ret *ComplexEnum, ok bool) {
	switch name {
	case "LitName1":
		ret = o.LitName1()
	case "LitName2":
		ret = o.LitName2()
	}
	return
}




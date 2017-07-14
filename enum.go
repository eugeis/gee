package gee

import "strings"

type Literal interface {
	Name() string
	Ordinal() int
}

type Enum interface {
	Literals() []Literal
	Parse(name string) Literal
}
func Parse(name string, enum Enum) (ret Literal, ok bool) {
	for _, lit := range enum.Literals() {
		if strings.EqualFold(lit.Name(), name) {
			return lit, true
		}
	}
	return nil, true
}



//simple
type LiteralBase struct {
	name    string
	ordinal int
}

func NewLiteralBase(name string) *LiteralBase {
	return &LiteralBase{name: name}
}

func (o *LiteralBase) Name() string {
	return o.name
}

func (o *LiteralBase) Ordinal() int {
	return o.ordinal
}

type EnumBase struct {
	literals []Literal
}

func NewEnumBase(literals []Literal) *EnumBase {
	return &EnumBase{literals: literals}
}

func (o *EnumBase) First() Literal {
	return o.literals[0]
}

func (o *EnumBase) Last() Literal {
	return o.literals[len(o.literals)-1]
}

func (o *EnumBase) Literals() []Literal {
	return o.literals
}

func (o *EnumBase) Parse(name string) (ret Literal, ok bool) {
	for _, lit := range o.literals {
		if strings.EqualFold(lit.Name(), name) {
			return lit, true
		}
	}
	return nil, true
}


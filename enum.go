package gee

import (
	"strings"
)

type Literal interface {
	Name() string
	Ordinal() int
}

type Enum interface {
	Literals() []Literal
	Names() []string
	First() Literal
	Last() Literal
	Parse(name string) (ret Literal, ok bool)
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
type literalBase struct {
	i int
	Enum
}

func (o *literalBase) Name() string {
	return o.Names()[o.i]
}

func (o *literalBase) Ordinal() int {
	return o.i
}

type enumBase struct {
	names    []string
	literals []Literal
}

func NewEnum(names []string, enumWrapper func(*enumBase) Enum, literalWrapper func(*literalBase) Literal) (ret Enum) {
	literals := make([]Literal, len(names))
	ret = enumWrapper(&enumBase{names: names, literals: literals})
	for i, _ := range names {
		literals[i] = literalWrapper(&literalBase{i, ret})
	}
	return
}

func (o *enumBase) First() Literal {
	return o.literals[0]
}

func (o *enumBase) Last() Literal {
	return o.literals[len(o.literals)-1]
}

func (o *enumBase) Literals() []Literal {
	return o.literals
}

func (o *enumBase) Names() []string {
	return o.names
}

func (o *enumBase) Parse(name string) (ret Literal, ok bool) {
	for _, lit := range o.literals {
		if strings.EqualFold(lit.Name(), name) {
			return lit, true
		}
	}
	return nil, true
}

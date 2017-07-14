package enum

import (
	"strings"
)

type Literal interface {
	Name() string
	Ordinal() int
	String() string
}

type Enum interface {
	Literals() []Literal
	Names() []string
	Size() int
	Get(ordinal int) Literal
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
	name    string
	ordinal int
}

func (o *literalBase) Name() string {
	return o.name
}

func (o *literalBase) Ordinal() int {
	return o.ordinal
}

func (o *literalBase) String() string {
	return o.name
}

type enumBase struct {
	names    []string
	literals []Literal
}

func NewEnum(names []string, literalWrapper func(*literalBase) Literal) (ret *enumBase) {
	literals := make([]Literal, len(names))
	ret = &enumBase{names: names, literals: literals}
	for i, name := range names {
		literals[i] = literalWrapper(&literalBase{name: name, ordinal: i})
	}
	return
}

func (o *enumBase) Get(ordinal int) Literal {
	return o.literals[ordinal]
}

func (o *enumBase) First() Literal {
	return o.Get(0)
}

func (o *enumBase) Last() Literal {
	return o.Get(len(o.literals) - 1)
}

func (o *enumBase) Literals() []Literal {
	return o.literals
}

func (o *enumBase) Names() []string {
	return o.names
}

func (o *enumBase) Size() int {
	return len(o.names)
}

func (o *enumBase) Parse(name string) (ret Literal, ok bool) {
	for _, lit := range o.literals {
		if strings.EqualFold(lit.Name(), name) {
			return lit, true
		}
	}
	return nil, true
}

package gee

type Enum interface {
	Name() string
	Ordinal() int
}

type EnumLiterals struct {
	Values *[]Enum
}

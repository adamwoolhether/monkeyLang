package compiler

type SymbolScope string

const GlobalScope SymbolScope = "GLOBAL"

// Symbol holds necessary info about a symbol encountered in Monkey.
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable helps associate identifiers in the global scope with a
// given unique number. It also allows access the previously associated
// number for a given identifier.
type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)

	return &SymbolTable{store: s}
}
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Scope: GlobalScope,
		Index: s.numDefinitions,
	}

	s.store[name] = symbol
	s.numDefinitions++

	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]

	return obj, ok
}

package utils

type Momento int

const (
	MomentoDefault Momento = iota
	Desayuno
	Almuerzo
	Merienda
	Cena
)

// Método para convertir los enums en cadenas
func (momento Momento) String() string {
	return [...]string{"Indefinido", "Desayuno", "Almuerzo", "Merienda", "Cena"}[momento]
}

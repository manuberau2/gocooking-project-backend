package utils

type TipoComida int

const (
	TipoDefault TipoComida = iota
	Verdura
	Lacteo
	Queso
	Legumbre
	Carne
	Fruta
)

// Método para convertir los enums en cadenas
func (tipoComida TipoComida) String() string {
	return [...]string{"Indefinido", "Verdura", "Lacteo", "Queso", "Legumbre", "Carne", "Fruta"}[tipoComida]
}

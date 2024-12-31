package dto

import (
	"errors"
)

type ParametrosProductosCantidad struct {
	Tipo   int    `form:"tipo"`
	Nombre string `form:"nombre"`
}

func (parametros ParametrosProductosCantidad) Validate() error {
	if parametros.Tipo < 0 || parametros.Tipo > 6 {
		return errors.New("tipo de comida invalido")
	}
	return nil
}

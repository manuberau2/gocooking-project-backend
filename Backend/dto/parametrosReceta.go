package dto

import (
	"errors"
)

type ParametrosReceta struct {
	Momento int    `form:"momento"`
	Tipo    int    `form:"tipo"`
	Nombre  string `form:"nombre"`
}

// hay que corregir pq si no se les asigna valor arrancan en 0
func (parametros ParametrosReceta) Validate() error {
	// Contar cuántos campos son válidos
	count := 0

	// Verificar el campo Momento
	if parametros.Momento >= 1 && parametros.Momento <= 4 {
		count++
	}

	// Verificar el campo Tipo
	if parametros.Tipo >= 1 && parametros.Tipo <= 6 {
		count++
	}

	// Verificar el campo Nombre
	if parametros.Nombre != "" {
		count++
	}

	// Validar que al menos uno de los campos esté presente
	if count == 0 {
		return errors.New("debe proporcionar al menos uno de los parámetros (Momento, Tipo, Nombre)")
	}

	return nil
}

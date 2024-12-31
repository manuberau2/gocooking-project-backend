package dto

import (
	"GoCooking/Backend/model"
	"GoCooking/Backend/utils"
	"errors"
)

type Alimento struct {
	Id                string           `json:"id"`
	Nombre            string           `json:"nombre"`
	Tipo              utils.TipoComida `json:"tipo"`
	MomentosDeConsumo []utils.Momento  `json:"momentos_de_consumo"`
	PrecioUnitario    float64          `json:"precio_unitario"`
	CantidadActual    float64          `json:"cantidad_actual"`
	CantidadMinima    float64          `json:"cantidad_minima"`
	UsuarioID         string           `json:"usuario_id"`
}

func NewAlimento(alimento model.Alimento) *Alimento {
	return &Alimento{
		Id:                utils.GetStringIDFromObjectID(alimento.Id),
		Nombre:            alimento.Nombre,
		Tipo:              alimento.Tipo,
		MomentosDeConsumo: alimento.MomentosDeConsumo,
		PrecioUnitario:    alimento.PrecioUnitario,
		CantidadActual:    alimento.CantidadActual,
		CantidadMinima:    alimento.CantidadMinima,
		UsuarioID:         alimento.UsuarioID,
	}
}

func (alimento Alimento) GetModel() model.Alimento {
	return model.Alimento{
		Id:                utils.GetObjectIDFromStringID(alimento.Id),
		Nombre:            alimento.Nombre,
		Tipo:              alimento.Tipo,
		MomentosDeConsumo: alimento.MomentosDeConsumo,
		PrecioUnitario:    alimento.PrecioUnitario,
		CantidadActual:    alimento.CantidadActual,
		CantidadMinima:    alimento.CantidadMinima,
		UsuarioID:         alimento.UsuarioID,
	}
}

func (alimento Alimento) Validate() error {
	if alimento.Nombre == "" {
		return errors.New("el nombre del alimento no puede estar vacío")
	}
	if alimento.CantidadActual <= 0 {
		return errors.New("la cantidad actual del alimento debe ser mayor a cero")
	}
	if alimento.CantidadMinima <= 0 {
		return errors.New("la cantidad mínima del alimento debe ser mayor a cero")
	}
	if alimento.PrecioUnitario <= 0 {
		return errors.New("el precio unitario del alimento debe ser mayor a cero")
	}
	if alimento.Tipo < 1 || alimento.Tipo > utils.Fruta { // que el rango sea correcto
		return errors.New("tipo de comida inválido")
	}
	if len(alimento.MomentosDeConsumo) == 0 {
		return errors.New("debe haber al menos un momento de consumo definido")
	}
	return nil
}

package dto

import (
	"GoCooking/Backend/model"
	"GoCooking/Backend/utils"
	"errors"
)

type Receta struct {
	Id               string        `json:"id"`
	Nombre           string        `json:"nombre"`
	MomentoDeConsumo utils.Momento `json:"momento_consumo"`
	Ingredientes     []Ingrediente `json:"ingredientes"`
	UsuarioID        string        `json:"usuario_id"`
}

type Ingrediente struct {
	AlimentoId string  `json:"alimento_id"`
	Nombre     string  `json:"nombre"`
	Cantidad   float64 `json:"cantidad"`
}

func NewReceta(receta model.Receta) *Receta {
	// Mapear cada ingrediente del model a dto
	ingredientesDTO := make([]Ingrediente, len(receta.Ingredientes))
	for i, ing := range receta.Ingredientes {
		ingredientesDTO[i] = Ingrediente{
			AlimentoId: utils.GetStringIDFromObjectID(ing.AlimentoId),
			Cantidad:   ing.Cantidad,
			Nombre:     ing.Nombre,
		}
	}

	return &Receta{
		Id:               utils.GetStringIDFromObjectID(receta.Id),
		Nombre:           receta.Nombre,
		MomentoDeConsumo: receta.MomentoDeConsumo,
		Ingredientes:     ingredientesDTO,
		UsuarioID:        receta.UsuarioID,
	}
}
func (receta Receta) GetModel() model.Receta {
	// Mapear cada ingrediente del dto a model
	ingredientesModel := make([]model.Ingrediente, len(receta.Ingredientes))
	for i, ing := range receta.Ingredientes {
		ingredientesModel[i] = model.Ingrediente{
			AlimentoId: utils.GetObjectIDFromStringID(ing.AlimentoId),
			Cantidad:   ing.Cantidad,
			Nombre:     ing.Nombre,
		}
	}

	return model.Receta{
		Id:               utils.GetObjectIDFromStringID(receta.Id),
		Nombre:           receta.Nombre,
		MomentoDeConsumo: receta.MomentoDeConsumo,
		Ingredientes:     ingredientesModel,
		UsuarioID:        receta.UsuarioID,
	}

}

func (receta Receta) Validate() error {
	// Verifica que el nombre no esté vacío
	if receta.Nombre == "" {
		return errors.New("el nombre de la receta es obligatorio")
	}

	// Verifica que el momento de consumo sea válido
	if receta.MomentoDeConsumo < 1 || receta.MomentoDeConsumo > 4 {
		return errors.New("el momento de consumo no es válido")
	}

	// Verifica que haya al menos un ingrediente
	if len(receta.Ingredientes) == 0 {
		return errors.New("debe haber al menos un ingrediente en la receta")
	}

	// Verifica que cada ingrediente tenga una cantidad válida
	for _, ingrediente := range receta.Ingredientes {
		if ingrediente.AlimentoId == "" {
			return errors.New("el ID del alimento es obligatorio en todos los ingredientes")
		}
		if ingrediente.Cantidad <= 0 {
			return errors.New("la cantidad de cada ingrediente debe ser mayor que cero")
		}
		if ingrediente.Nombre == "" {
			return errors.New("el nombre del ingrediente es obligatorio")
		}
	}

	return nil
}

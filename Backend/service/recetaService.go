package service

import (
	"GoCooking/Backend/dto"
	"GoCooking/Backend/repositories"
	"GoCooking/Backend/utils"
)

type RecetaInterface interface {
	GetRecetas(usuarioID string) ([]*dto.Receta, *utils.AppError)
	GetRecetaById(id string) (*dto.Receta, *utils.AppError)
	InsertReceta(receta *dto.Receta) (bool, *utils.AppError)
	UpdateReceta(receta *dto.Receta) (bool, *utils.AppError)
	DeleteReceta(id string) (bool, *utils.AppError)
	GetRecetasByParameters(parametros dto.ParametrosReceta, usuarioID string) ([]*dto.Receta, *utils.AppError)
	GetCantidadRecetasPorMomento(usuarioID string) (map[string]int, *utils.AppError)
	GetCantidadRecetasPorTipoAlimento(usuarioID string) (map[string]int, *utils.AppError)
}

type RecetaService struct {
	recetaRepository repositories.RecetaRepositoryInterface
}

func NewRecetaService(recetaRepository repositories.RecetaRepositoryInterface) *RecetaService {
	return &RecetaService{
		recetaRepository: recetaRepository,
	}
}
func (service *RecetaService) GetRecetas(usuarioID string) ([]*dto.Receta, *utils.AppError) {
	recetasDB, err := service.recetaRepository.GetRecetas(usuarioID)
	if err != nil {
		return nil, utils.NewAppError("ERR_500", "Error al obtener las recetas")
	}
	if len(*recetasDB) == 0 {
		return nil, utils.NewAppError("ERR_404", "No se encontraron recetas")
	}
	var recetas []*dto.Receta
	for _, recetaDB := range *recetasDB {
		receta := dto.NewReceta(recetaDB)
		recetas = append(recetas, receta)
	}
	return recetas, nil
}

func (service *RecetaService) GetRecetaById(id string) (*dto.Receta, *utils.AppError) {
	recetaDB, err := service.recetaRepository.GetRecetaById(utils.GetObjectIDFromStringID(id))
	if err != nil {
		if err.Error() == "404" {
			return nil, utils.NewAppError("ERR_404", "La receta no fue encontrada")
		}
		return nil, utils.NewAppError("ERR_500", "Error al obtener la receta")
	}
	receta := dto.NewReceta(*recetaDB)
	return receta, nil
}

func (service *RecetaService) InsertReceta(receta *dto.Receta) (bool, *utils.AppError) {
	err := receta.Validate()
	if err != nil {
		return false, utils.NewAppError("ERR_400", err.Error())
	}
	resultado, err := service.recetaRepository.InsertReceta(receta.GetModel())
	if err != nil || resultado == nil {
		return false, utils.NewAppError("ERR_500", "Error al insertar la receta: "+err.Error())
	}
	return true, nil
}

func (service *RecetaService) UpdateReceta(receta *dto.Receta) (bool, *utils.AppError) {
	err := receta.Validate()
	if err != nil {
		return false, utils.NewAppError("ERR_400", err.Error())
	}
	resultado, err := service.recetaRepository.UpdateReceta(receta.GetModel())
	if err != nil || resultado == nil {
		if err.Error() == "404" {
			return false, utils.NewAppError("ERR_404", "La receta no fue encontrada")
		}
		return false, utils.NewAppError("ERR_500", "Error al actualizar la receta: "+err.Error())
	}
	return true, nil
}

func (service *RecetaService) DeleteReceta(id string) (bool, *utils.AppError) {
	resultado, err := service.recetaRepository.DeleteReceta(utils.GetObjectIDFromStringID(id))
	if err != nil || resultado == nil {
		if err.Error() == "404" {
			return false, utils.NewAppError("ERR_404", "La receta no fue encontrada")
		}
		return false, utils.NewAppError("ERR_500", "Error al eliminar la receta: "+err.Error())
	}
	return true, nil
}

func (service *RecetaService) GetRecetasByParameters(parametros dto.ParametrosReceta, usuarioID string) ([]*dto.Receta, *utils.AppError) {
	err := parametros.Validate()
	if err != nil {
		return nil, utils.NewAppError("ERR_400", err.Error())
	}
	recetasDB, err := service.recetaRepository.GetRecetasByParameters(parametros, usuarioID)
	if err != nil {
		return nil, utils.NewAppError("ERR_500", "Error al obtener las recetas: "+err.Error())
	}
	if len(recetasDB) == 0 {
		return nil, utils.NewAppError("ERR_404", "No se encontraron recetas")
	}
	var recetas []*dto.Receta
	for _, recetaDB := range recetasDB {
		receta := dto.NewReceta(recetaDB)
		recetas = append(recetas, receta)
	}
	return recetas, nil
}

func (service *RecetaService) GetCantidadRecetasPorMomento(usuarioID string) (map[string]int, *utils.AppError) {
	cantidadRecetasPorMomento, err := service.recetaRepository.GetCantidadRecetasPorMomento(usuarioID)
	if err != nil {
		return nil, utils.NewAppError("ERR_500", "Error al obtener la cantidad de recetas por momento")
	}
	if len(cantidadRecetasPorMomento) == 0 {
		return nil, utils.NewAppError("ERR_404", "No se encontraron recetas")
	}
	return cantidadRecetasPorMomento, nil
}

func (service *RecetaService) GetCantidadRecetasPorTipoAlimento(usuarioID string) (map[string]int, *utils.AppError) {
	cantidadRecetasPorTipoAlimento, err := service.recetaRepository.GetCantidadRecetasPorTipoAlimento(usuarioID)
	if err != nil {
		return nil, utils.NewAppError("ERR_500", "Error al obtener la cantidad de recetas por tipo de alimento")
	}
	if len(cantidadRecetasPorTipoAlimento) == 0 {
		return nil, utils.NewAppError("ERR_404", "No se encontraron recetas")
	}
	return cantidadRecetasPorTipoAlimento, nil
}

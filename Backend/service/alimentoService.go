package service

import (
	"GoCooking/Backend/dto"
	"GoCooking/Backend/repositories"
	"GoCooking/Backend/utils"
)

type AlimentoInterface interface {
	GetAlimentos(usuarioID string) ([]*dto.Alimento, *utils.AppError)
	GetAlimentoByID(id string) (*dto.Alimento, *utils.AppError)
	InsertAlimento(alimento *dto.Alimento) (bool, *utils.AppError)
	UpdateAlimento(alimento *dto.Alimento) (bool, *utils.AppError)
	DeleteAlimento(id string) (bool, *utils.AppError)
}
type AlimentoService struct {
	alimentoRepository repositories.AlimentoRepositoryInterface
}

func NewAlimentoService(alimentoRepository repositories.AlimentoRepositoryInterface) *AlimentoService {
	return &AlimentoService{
		alimentoRepository: alimentoRepository,
	}
}

func (service *AlimentoService) GetAlimentos(usuarioID string) ([]*dto.Alimento, *utils.AppError) {
	alimentosDB, err := service.alimentoRepository.GetAlimentos(usuarioID)
	if err != nil {
		return nil, utils.NewAppError("ERR_500", "Error al obtener los alimentos: "+err.Error())
	}

	// Verificar si la lista de alimentos está vacía
	if len(*alimentosDB) == 0 {
		return nil, utils.NewAppError("ERR_404", "No se encontraron alimentos disponibles")
	}

	var alimentos []*dto.Alimento
	for _, alimentoDB := range *alimentosDB {
		alimento := dto.NewAlimento(alimentoDB)
		alimentos = append(alimentos, alimento)
	}

	return alimentos, nil
}

func (service *AlimentoService) GetAlimentoByID(id string) (*dto.Alimento, *utils.AppError) {
	alimentoDB, err := service.alimentoRepository.GetAlimentoByID(utils.GetObjectIDFromStringID(id))
	if err != nil {
		if err.Error() == "404" {
			return nil, utils.NewAppError("ERR_404", "El alimento no fue encontrado")
		}
		return nil, utils.NewAppError("ERR_500", "Error al obtener el alimento: "+err.Error())
	}
	alimento := dto.NewAlimento(*alimentoDB)
	return alimento, nil
}

func (service *AlimentoService) InsertAlimento(alimento *dto.Alimento) (bool, *utils.AppError) {
	err := alimento.Validate()
	if err != nil {
		return false, utils.NewAppError("ERR_400", err.Error())
	}
	resultado, err := service.alimentoRepository.InsertAlimento(alimento.GetModel())
	if err != nil || resultado == nil {
		return false, utils.NewAppError("ERR_500", "Error al insertar el alimento: "+err.Error())
	}
	return true, nil
}

func (service *AlimentoService) UpdateAlimento(alimento *dto.Alimento) (bool, *utils.AppError) {
	err := alimento.Validate()
	if err != nil {
		return false, utils.NewAppError("ERR_400", err.Error())
	}
	resultado, err := service.alimentoRepository.UpdateAlimento(alimento.GetModel())
	if err != nil || resultado == nil {
		if err.Error() == "404" {
			return false, utils.NewAppError("ERR_404", "El alimento no fue encontrado")
		}
		return false, utils.NewAppError("ERR_500", "Error al actualizar el alimento: "+err.Error())
	}
	return true, nil
}

func (service *AlimentoService) DeleteAlimento(id string) (bool, *utils.AppError) {
	resultado, err := service.alimentoRepository.DeleteAlimento(utils.GetObjectIDFromStringID(id))
	if err != nil || resultado == nil {
		if err.Error() == "404" {
			return false, utils.NewAppError("ERR_404", "El alimento no fue encontrado")
		}
		return false, utils.NewAppError("ERR_500", "Error al eliminar el alimento "+err.Error())
	}
	return true, nil
}

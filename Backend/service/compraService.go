package service

import (
	"GoCooking/Backend/dto"
	"GoCooking/Backend/repositories"
	"GoCooking/Backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompraInterface interface {
	GetProductosPorCantidadMinima(parametros dto.ParametrosProductosCantidad, usuarioID string) ([]*dto.ProductoCompra, *utils.AppError) //dto para los productos?
	PostNuevaCompra(usuarioID string, idsComprasSeleccionadas []string) (*dto.Compra, *utils.AppError)
	GetCompras(usuarioID string) ([]*dto.Compra, *utils.AppError)
	GetCostoPromedioPorMesUltimoAnio(usuarioID string) (map[string]float64, *utils.AppError)
}

type CompraService struct {
	compraRepository repositories.CompraRepositoryInterface
}

func NewCompraService(compraRepository repositories.CompraRepositoryInterface) *CompraService {
	return &CompraService{
		compraRepository: compraRepository,
	}
}

func (service *CompraService) GetProductosPorCantidadMinima(parametros dto.ParametrosProductosCantidad, usuarioID string) ([]*dto.ProductoCompra, *utils.AppError) {
	err := parametros.Validate()
	if err != nil {
		return nil, utils.NewAppError("ERR_400", err.Error())
	}
	productosDB, err := service.compraRepository.GetProductosPorCantidadMinima(parametros, usuarioID)
	if err != nil {
		return nil, utils.NewAppError("ERR_500", "Error al obtener los productos: "+err.Error())
	}
	if len(*productosDB) == 0 {
		return nil, utils.NewAppError("ERR_404", "No se encontraron productos")
	}
	var productos []*dto.ProductoCompra
	for _, productoDB := range *productosDB {
		producto := dto.NewProductoCompra(productoDB)
		productos = append(productos, producto)
	}
	return productos, nil
}
func (service *CompraService) PostNuevaCompra(usuarioID string, idsComprasSeleccionadas []string) (*dto.Compra, *utils.AppError) {
	// Convertir los IDs de string a ObjectID
	var objectIDs []primitive.ObjectID
	for _, id := range idsComprasSeleccionadas {
		objectID := utils.GetObjectIDFromStringID(id)
		objectIDs = append(objectIDs, objectID)
	}

	compraModel, err := service.compraRepository.PostNuevaCompra(usuarioID, objectIDs)
	if err != nil {
		if len(compraModel.Productos) == 0 {
			return nil, utils.NewAppError("ERR_400", "No se puede realizar la compra, no hay productos seleccionados")
		}
		return nil, utils.NewAppError("ERR_500", "Error al crear la compra: "+err.Error())
	}
	compra := dto.NewCompra(*compraModel)
	return compra, nil
}

func (service *CompraService) GetCompras(usuarioID string) ([]*dto.Compra, *utils.AppError) {
	comprasDB, err := service.compraRepository.GetCompras(usuarioID)
	if err != nil {
		return nil, utils.NewAppError("ERR_500", "Error al obtener las compras: "+err.Error())
	}
	if len(*comprasDB) == 0 {
		return nil, utils.NewAppError("ERR_404", "No se encontraron compras")
	}
	var compras []*dto.Compra
	for _, compraDB := range *comprasDB {
		compra := dto.NewCompra(compraDB)
		compras = append(compras, compra)
	}
	return compras, nil
}

func (service *CompraService) GetCostoPromedioPorMesUltimoAnio(usuarioID string) (map[string]float64, *utils.AppError) {
	costoPromedioPorMes, err := service.compraRepository.GetCostoPromedioPorMesUltimoAnio(usuarioID)
	if err != nil {
		return nil, utils.NewAppError("ERR_500", "Error al obtener el costo promedio por mes: "+err.Error())
	}
	return costoPromedioPorMes, nil
}

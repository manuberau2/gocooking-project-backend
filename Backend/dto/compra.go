package dto

import (
	"GoCooking/Backend/model"
	"GoCooking/Backend/utils"
	"time"
)

type Compra struct {
	ID          string           `json:"id"`
	Productos   []ProductoCompra `json:"productos"`
	CostoTotal  float64          `json:"costo_total"`
	FechaCompra time.Time        `json:"fecha_compra"`
	UsuarioID   string           `json:"usuario_id"`
}

type ProductoCompra struct {
	AlimentoID string  `json:"alimento_id"`
	Nombre     string  `json:"nombre"`
	Cantidad   float64 `json:"cantidad"`
}

func NewCompra(compra model.Compra) *Compra {
	// Mapear cada producto del model a dto
	productosDTO := make([]ProductoCompra, len(compra.Productos))
	for i, prod := range compra.Productos {
		productosDTO[i] = ProductoCompra{
			AlimentoID: utils.GetStringIDFromObjectID(prod.AlimentoId),
			Cantidad:   prod.Cantidad,
			Nombre:     prod.Nombre,
		}
	}

	return &Compra{
		ID:          utils.GetStringIDFromObjectID(compra.Id),
		Productos:   productosDTO,
		CostoTotal:  compra.CostoTotal,
		FechaCompra: compra.FechaCreacion,
		UsuarioID:   compra.UsuarioID,
	}
}

func (compra Compra) GetModel() model.Compra {
	// Mapear cada producto del dto a model
	productosModel := make([]model.ProductoCompra, len(compra.Productos))
	for i, prod := range compra.Productos {
		productosModel[i] = model.ProductoCompra{
			AlimentoId: utils.GetObjectIDFromStringID(prod.AlimentoID),
			Cantidad:   prod.Cantidad,
			Nombre:     prod.Nombre,
		}
	}

	return model.Compra{
		Id:            utils.GetObjectIDFromStringID(compra.ID),
		Productos:     productosModel,
		FechaCreacion: compra.FechaCompra,
		CostoTotal:    compra.CostoTotal,
		UsuarioID:     compra.UsuarioID,
	}
}
func NewProductoCompra(prod model.ProductoCompra) *ProductoCompra {
	return &ProductoCompra{
		AlimentoID: utils.GetStringIDFromObjectID(prod.AlimentoId),
		Cantidad:   prod.Cantidad,
		Nombre:     prod.Nombre,
	}
}

func (prod ProductoCompra) GetModel() model.ProductoCompra {
	return model.ProductoCompra{
		AlimentoId: utils.GetObjectIDFromStringID(prod.AlimentoID),
		Cantidad:   prod.Cantidad,
		Nombre:     prod.Nombre,
	}
}

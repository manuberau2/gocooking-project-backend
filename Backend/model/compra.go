package model

import (
	"GoCooking/Backend/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Compra struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	Productos          []ProductoCompra   `bson:"lista_productos"`
	CostoTotal         float64            `bson:"costo_total"`
	FechaCreacion      time.Time          `bson:"fecha_creacion"`
	FechaActualizacion time.Time          `bson:"fecha_actualizacion"`
	UsuarioID          string             `bson:"id_usuario"`
}

type ProductoCompra struct {
	AlimentoId primitive.ObjectID `bson:"id_alimento"`
	Cantidad   float64            `bson:"cantidad_comprada"`
	Nombre     string             `bson:"nombre"`
	Tipo       utils.TipoComida   `bson:"tipo"`
}

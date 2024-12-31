package model

import (
	"GoCooking/Backend/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alimento struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	Nombre             string             `bson:"nombre"`
	Tipo               utils.TipoComida   `bson:"tipo"`
	MomentosDeConsumo  []utils.Momento    `bson:"momento"`
	PrecioUnitario     float64            `bson:"precio_unitario"`
	CantidadActual     float64            `bson:"cantidad_actual"`
	CantidadMinima     float64            `bson:"cantidad_minima"`
	UsuarioID          string             `bson:"id_usuario"`
	FechaCreacion      time.Time          `bson:"fecha_creacion"`
	FechaActualizacion time.Time          `bson:"fecha_actualizacion"`
}

package model

import (
	"GoCooking/Backend/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Receta struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	Nombre             string             `bson:"nombre"`
	MomentoDeConsumo   utils.Momento      `bson:"momento_consumo"`
	Ingredientes       []Ingrediente      `bson:"ingredientes"`
	FechaCreacion      time.Time          `bson:"fecha_creacion"`
	FechaActualizacion time.Time          `bson:"fecha_actualizacion"`
	UsuarioID          string             `bson:"id_usuario"`
}

type Ingrediente struct {
	AlimentoId primitive.ObjectID `bson:"id_alimento"`
	Nombre     string             `bson:"nombre"`
	Cantidad   float64            `bson:"cantidad"`
}

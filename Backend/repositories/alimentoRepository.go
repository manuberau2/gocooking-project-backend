package repositories

import (
	"GoCooking/Backend/model"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlimentoRepositoryInterface interface {
	GetAlimentos(usuarioID string) (*[]model.Alimento, error)
	GetAlimentoByID(id primitive.ObjectID) (*model.Alimento, error)
	InsertAlimento(alimento model.Alimento) (*mongo.InsertOneResult, error)
	UpdateAlimento(alimento model.Alimento) (*mongo.UpdateResult, error)
	DeleteAlimento(id primitive.ObjectID) (*mongo.DeleteResult, error)
}

type AlimentoRepository struct {
	db DB
}

func NewAlimentoRepository(db DB) *AlimentoRepository {
	return &AlimentoRepository{
		db: db,
	}
}

func (repository AlimentoRepository) GetAlimentos(usuarioID string) (*[]model.Alimento, error) {
	collection := repository.db.GetClient().Database("gocooking").Collection("alimentos")
	filtro := bson.M{
		"id_usuario": usuarioID,
	}
	cursor, err := collection.Find(context.TODO(), filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var alimentos []model.Alimento
	for cursor.Next(context.Background()) {
		var alimento model.Alimento
		err = cursor.Decode(&alimento)
		if err != nil {
			return nil, err
		}
		alimentos = append(alimentos, alimento)
	}
	return &alimentos, err
}

func (repository AlimentoRepository) GetAlimentoByID(id primitive.ObjectID) (*model.Alimento, error) {
	collection := repository.db.GetClient().Database("gocooking").Collection("alimentos")
	filtro := bson.M{"_id": id}
	cursor, err := collection.Find(context.TODO(), filtro)
	if err != nil {
		return &model.Alimento{}, err
	}
	defer cursor.Close(context.TODO())

	var alimento model.Alimento
	if cursor.Next(context.Background()) {
		err = cursor.Decode(&alimento)
		if err != nil {
			return &model.Alimento{}, err
		}
	} else {
		// Si no se encontró el alimento, devolver un error 404
		return &model.Alimento{}, errors.New("404")
	}
	return &alimento, nil
}

func (repository AlimentoRepository) InsertAlimento(alimento model.Alimento) (*mongo.InsertOneResult, error) {
	alimento.FechaCreacion = time.Now()
	collection := repository.db.GetClient().Database("gocooking").Collection("alimentos")
	resultado, err := collection.InsertOne(context.TODO(), alimento)
	return resultado, err
}

func (repository AlimentoRepository) UpdateAlimento(alimento model.Alimento) (*mongo.UpdateResult, error) {
	alimento.FechaActualizacion = time.Now()
	collection := repository.db.GetClient().Database("gocooking").Collection("alimentos")

	filtro := bson.M{"_id": alimento.Id}
	entidad := bson.M{
		"$set": bson.M{
			"nombre":              alimento.Nombre,
			"fecha_actualizacion": alimento.FechaActualizacion,
			"tipo":                alimento.Tipo,
			"momento":             alimento.MomentosDeConsumo,
			"precio_unitario":     alimento.PrecioUnitario,
			"cantidad_actual":     alimento.CantidadActual,
			"cantidad_minima":     alimento.CantidadMinima,
		},
	}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, entidad)
	if err != nil {
		return nil, err
	}

	// Si no se modificó nada, devolver un error indicando que no se encontró el alimento
	if resultado.ModifiedCount == 0 {
		return nil, errors.New("no se encontró el alimento")
	}

	return resultado, nil
}

func (repository AlimentoRepository) DeleteAlimento(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("gocooking").Collection("alimentos")

	filtro := bson.M{"_id": id}

	resultado, err := collection.DeleteOne(context.TODO(), filtro)
	if err != nil {
		return nil, err
	}

	// Verificar si se eliminó algún documento
	if resultado.DeletedCount == 0 {
		return nil, errors.New("404")
	}

	return resultado, nil
}

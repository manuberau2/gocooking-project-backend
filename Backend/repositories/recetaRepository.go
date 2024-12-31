package repositories

import (
	"GoCooking/Backend/dto"
	"GoCooking/Backend/model"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecetaRepositoryInterface interface {
	GetRecetas(usuarioID string) (*[]model.Receta, error)
	GetRecetaById(id primitive.ObjectID) (*model.Receta, error)
	InsertReceta(receta model.Receta) (*mongo.InsertOneResult, error)
	UpdateReceta(receta model.Receta) (*mongo.UpdateResult, error)
	DeleteReceta(id primitive.ObjectID) (*mongo.DeleteResult, error)
	GetRecetasByParameters(parametros dto.ParametrosReceta, usuarioID string) ([]model.Receta, error)
	GetCantidadRecetasPorMomento(usuarioID string) (map[string]int, error)
	GetCantidadRecetasPorTipoAlimento(usuarioID string) (map[string]int, error)
}

type RecetaRepository struct {
	db DB
}

func NewRecetaRepository(db DB) *RecetaRepository {
	return &RecetaRepository{
		db: db,
	}
}

func (repository RecetaRepository) GetRecetas(usuarioID string) (*[]model.Receta, error) {
	filtro := bson.M{
		"id_usuario": usuarioID,
	}
	cursor, err := repository.db.GetClient().Database("gocooking").Collection("recetas").Find(context.TODO(), filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var recetas []model.Receta
	for cursor.Next(context.Background()) {
		var receta model.Receta
		if err := cursor.Decode(&receta); err != nil {
			return nil, err
		}
		// por cada receta, buscamos que los ingredientes estén disponibles en la coleccion de alimentos
		disponible := true
		for _, ingrediente := range receta.Ingredientes {
			var alimento model.Alimento
			err := repository.db.GetClient().Database("gocooking").Collection("alimentos").FindOne(context.TODO(), bson.M{"_id": ingrediente.AlimentoId}).Decode(&alimento)
			if err != nil {
				return nil, err
			}
			if alimento.CantidadActual < ingrediente.Cantidad {
				disponible = false
				break
			}
		}
		if disponible {
			recetas = append(recetas, receta)
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return &recetas, nil
}

func (repository RecetaRepository) GetRecetaById(id primitive.ObjectID) (*model.Receta, error) {
	var receta model.Receta
	err := repository.db.GetClient().Database("gocooking").Collection("recetas").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&receta)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Retorna un error específico si no se encuentra la receta
			return nil, errors.New("404")
		}
		// Retorna cualquier otro error que pueda ocurrir
		return nil, err
	}

	return &receta, nil
}

func (repository RecetaRepository) InsertReceta(receta model.Receta) (*mongo.InsertOneResult, error) {
	receta.FechaCreacion = time.Now()
	for _, ingrediente := range receta.Ingredientes {
		var alimento model.Alimento

		// Buscar el alimento correspondiente al ingrediente en la colección "alimentos"
		err := repository.db.GetClient().Database("gocooking").Collection("alimentos").FindOne(context.TODO(), bson.M{"_id": ingrediente.AlimentoId}).Decode(&alimento)
		if err != nil {
			return nil, err
		}

		// Verificar si hay suficiente cantidad de alimento
		if alimento.CantidadActual < ingrediente.Cantidad {
			return nil, errors.New("no hay suficiente cantidad del alimento " + alimento.Nombre)
		}

		// Verificar que el alimento es adecuado para el momento de consumo de la receta
		alimentoAdecuado := false
		for _, momento := range alimento.MomentosDeConsumo {
			if momento == receta.MomentoDeConsumo {
				alimentoAdecuado = true
			}
		}
		if !alimentoAdecuado {
			return nil, errors.New("el alimento " + alimento.Nombre + " no es adecuado para el momento de consumo de la receta")
		}
	}

	// Realizar la inserción de la receta en la colección "recetas"
	resultado, err := repository.db.GetClient().Database("gocooking").Collection("recetas").InsertOne(context.TODO(), receta)
	if err != nil {
		return nil, errors.New("error al insertar la receta: " + err.Error())
	}

	// Restar las cantidades utilizadas a los alimentos en el almacén
	for _, ingrediente := range receta.Ingredientes {
		_, err := repository.db.GetClient().Database("gocooking").Collection("alimentos").UpdateOne(context.TODO(), bson.M{"_id": ingrediente.AlimentoId}, bson.M{"$inc": bson.M{"cantidad_actual": -ingrediente.Cantidad}})
		if err != nil {
			return nil, errors.New("error al actualizar la cantidad de alimento: " + err.Error())
		}
	}

	return resultado, nil
}

func (repository RecetaRepository) UpdateReceta(receta model.Receta) (*mongo.UpdateResult, error) {
	receta.FechaActualizacion = time.Now()
	for _, ingrediente := range receta.Ingredientes {
		var alimento model.Alimento
		err := repository.db.GetClient().Database("gocooking").Collection("alimentos").FindOne(context.TODO(), bson.M{"_id": ingrediente.AlimentoId}).Decode(&alimento)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("404")
			}
			return nil, err
		}

		// Verificar que el alimento tenga suficiente cantidad
		if alimento.CantidadActual < ingrediente.Cantidad {
			return nil, errors.New("no hay suficiente cantidad del alimento " + alimento.Nombre)
		}

		// Verificar que el alimento sea adecuado para el momento de la receta
		alimentoAdecuado := false
		for _, momento := range alimento.MomentosDeConsumo {
			if momento == receta.MomentoDeConsumo {
				alimentoAdecuado = true
			}
		}
		if !alimentoAdecuado {
			return nil, errors.New("el alimento " + alimento.Nombre + " no es adecuado para el momento de consumo de la receta")
		}
	}

	// Actualizar receta en la base de datos
	filter := bson.M{"_id": receta.Id}
	update := bson.M{
		"$set": receta,
	}
	result, err := repository.db.GetClient().Database("gocooking").Collection("recetas").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repository RecetaRepository) DeleteReceta(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	// Obtener la receta a eliminar
	receta, err := repository.GetRecetaById(id)
	if err != nil {
		if err.Error() == "404" {
			return nil, errors.New("404")
		}
		return nil, err
	}

	// Devolver las cantidades de los ingredientes al stock
	for _, ingrediente := range receta.Ingredientes {
		filter := bson.M{"_id": ingrediente.AlimentoId}
		update := bson.M{
			"$inc": bson.M{
				"cantidad_actual": ingrediente.Cantidad, // Sumar la cantidad utilizada
			},
		}
		_, err := repository.db.GetClient().Database("gocooking").Collection("alimentos").UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return nil, err
		}
	}

	// Eliminar la receta
	result, err := repository.db.GetClient().Database("gocooking").Collection("recetas").DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repository RecetaRepository) GetRecetasByParameters(parametros dto.ParametrosReceta, usuarioID string) ([]model.Receta, error) {
	var recetas []model.Receta
	filter := bson.M{
		"id_usuario": usuarioID,
	}
	log.Printf("momento parámetro: %v , nombre parametro: %v", parametros.Momento, parametros.Nombre)
	// Filtros opcionales
	if parametros.Momento >= 1 && parametros.Momento <= 4 {
		filter["momento_consumo"] = parametros.Momento // Usar el valor entero directamente
	}

	cursor, err := repository.db.GetClient().Database("gocooking").Collection("recetas").Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Validar que las recetas tengan suficientes cantidades de alimentos en stock
	for cursor.Next(context.TODO()) {
		var receta model.Receta
		err := cursor.Decode(&receta)
		if err != nil {
			return nil, err
		}

		// Verificar que haya stock suficiente para cada ingrediente de la receta
		disponible := true
		tipoCoincide := false   // Variable para comprobar si al menos un ingrediente coincide
		nombreCoincide := false // Inicialmente asumimos que coincide con el nombre

		// Recoger los IDs de los alimentos a buscar
		alimentoIDs := make([]primitive.ObjectID, len(receta.Ingredientes))
		for i, ingrediente := range receta.Ingredientes {
			alimentoIDs[i] = ingrediente.AlimentoId
		}

		// Buscar todos los alimentos de una vez
		var alimentos []model.Alimento
		cursorAlimentos, err := repository.db.GetClient().Database("gocooking").Collection("alimentos").Find(context.TODO(), bson.M{"_id": bson.M{"$in": alimentoIDs}})
		if err != nil {
			return nil, err
		}
		defer cursorAlimentos.Close(context.TODO())

		for cursorAlimentos.Next(context.TODO()) {
			var alimento model.Alimento
			err := cursorAlimentos.Decode(&alimento)
			if err != nil {
				return nil, err
			}
			alimentos = append(alimentos, alimento)
		}

		// Crear un mapa para acceder fácilmente a los alimentos por ID
		alimentoMap := make(map[primitive.ObjectID]model.Alimento)
		for _, alimento := range alimentos {
			alimentoMap[alimento.Id] = alimento
		}

		for _, ingrediente := range receta.Ingredientes {
			alimento, exists := alimentoMap[ingrediente.AlimentoId]
			if !exists {
				disponible = false
				break
			}

			// Verificar stock
			if alimento.CantidadActual < ingrediente.Cantidad {
				disponible = false
				break
			}
			// Comprobar tipo de alimento
			if parametros.Tipo >= 1 && parametros.Tipo <= 6 {
				if int(alimento.Tipo) == parametros.Tipo {
					tipoCoincide = true
				}
			}

			// Comprobar nombre de ingrediente
			if parametros.Nombre != "" {
				if strings.Contains(strings.ToLower(alimento.Nombre), strings.ToLower(parametros.Nombre)) {
					nombreCoincide = true
				}
			}
		}
		log.Printf("Receta: %s - disponible: %v, tipoCoincide: %v, nombreCoincide: %v", receta.Nombre, disponible, tipoCoincide, nombreCoincide)
		if disponible && (parametros.Tipo == 0 || tipoCoincide) && (parametros.Nombre == "" || nombreCoincide) {
			recetas = append(recetas, receta)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return recetas, nil
}

func (repository RecetaRepository) GetCantidadRecetasPorMomento(usuarioID string) (map[string]int, error) {
	recetas, err := repository.GetRecetas(usuarioID)
	if err != nil {
		return nil, err
	}

	cantidadRecetasPorMomento := make(map[string]int)
	for _, receta := range *recetas {
		cantidadRecetasPorMomento[receta.MomentoDeConsumo.String()]++
	}
	return cantidadRecetasPorMomento, nil
}

func (repository RecetaRepository) GetCantidadRecetasPorTipoAlimento(usuarioID string) (map[string]int, error) {
	// Obtener las recetas del usuario
	recetas, err := repository.GetRecetas(usuarioID)
	if err != nil {
		return nil, err
	}

	// Inicializar el mapa para almacenar los conteos
	cantidadRecetasPorTipoAlimento := make(map[string]int)

	for _, receta := range *recetas {
		// Usar un mapa local para evitar contar tipos duplicados en una receta
		tiposContados := make(map[string]bool)

		for _, ingrediente := range receta.Ingredientes {
			var alimento model.Alimento
			// Obtener el alimento de la base de datos
			err := repository.db.GetClient().Database("gocooking").Collection("alimentos").FindOne(
				context.TODO(),
				bson.M{"_id": ingrediente.AlimentoId},
			).Decode(&alimento)
			if err != nil {
				return nil, err
			}

			// Contar el tipo de alimento si no fue contado en esta receta
			if !tiposContados[alimento.Tipo.String()] {
				cantidadRecetasPorTipoAlimento[alimento.Tipo.String()]++
				tiposContados[alimento.Tipo.String()] = true
			}
		}
	}

	return cantidadRecetasPorTipoAlimento, nil
}

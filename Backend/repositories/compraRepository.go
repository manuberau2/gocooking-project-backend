package repositories

import (
	"GoCooking/Backend/dto"
	"GoCooking/Backend/model"
	"GoCooking/Backend/utils"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompraRepositoryInterface interface {
	GetProductosPorCantidadMinima(parametros dto.ParametrosProductosCantidad, usuarioID string) (*[]model.ProductoCompra, error)
	PostNuevaCompra(usuarioID string, idsComprasSeleccionadas []primitive.ObjectID) (*model.Compra, error)
	GetCompras(usuarioID string) (*[]model.Compra, error)
	getAlimentoByID(id primitive.ObjectID) (*model.Alimento, error)
	GetCostoPromedioPorMesUltimoAnio(usuarioID string) (map[string]float64, error)
}

type CompraRepository struct {
	db DB
}

func NewCompraRepository(db DB) *CompraRepository {
	return &CompraRepository{
		db: db,
	}
}

func (repository CompraRepository) GetProductosPorCantidadMinima(parametros dto.ParametrosProductosCantidad, usuarioID string) (*[]model.ProductoCompra, error) {
	log.Printf("parametros: %v", parametros)
	// Conectar a la colección 'alimentos'
	collection := repository.db.GetClient().Database("gocooking").Collection("alimentos")

	// Filtro inicial que se aplica en la consulta a MongoDB
	filtro := bson.M{
		"$expr": bson.M{
			"$lt": []interface{}{"$cantidad_actual", "$cantidad_minima"},
		},
		"id_usuario": usuarioID,
	}

	// Ejecutar la consulta con el filtro inicial
	cursor, err := collection.Find(context.TODO(), filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var productos []model.ProductoCompra

	// Procesar cada producto encontrado
	for cursor.Next(context.Background()) {
		var alimento model.Alimento
		if err := cursor.Decode(&alimento); err != nil {
			return nil, err
		}

		// Crear un producto con la cantidad que falta para llegar al mínimo
		producto := model.ProductoCompra{
			AlimentoId: alimento.Id,
			Cantidad:   alimento.CantidadMinima - alimento.CantidadActual,
			Nombre:     alimento.Nombre,
			Tipo:       alimento.Tipo,
		}

		// Agregar el producto a la lista
		productos = append(productos, producto)
	}
	log.Printf("productos: %v", productos)
	// Aplicar filtros opcionales en Go, si se proporcionaron
	var productosFiltrados []model.ProductoCompra

	for _, producto := range productos {
		// Filtrar por tipo, solo si el parámetro `Tipo` es distinto de 0
		log.Printf("producto tipo: %v, parametro tipo: %v", producto.Tipo, parametros.Tipo)
		if parametros.Tipo != 0 && producto.Tipo != utils.TipoComida(parametros.Tipo) {
			continue
		}

		// Filtrar por nombre usando aproximación, solo si el parámetro `Nombre` no está vacío
		log.Printf("producto nombre: %v, parametro nombre: %v", producto.Nombre, parametros.Nombre)
		if parametros.Nombre != "" && !strings.Contains(strings.ToLower(producto.Nombre), strings.ToLower(parametros.Nombre)) {
			continue
		}

		productosFiltrados = append(productosFiltrados, producto)
	}

	// Si no hay filtros adicionales (tipo y nombre), devolvemos todos los productos que cumplen la condición inicial
	if parametros.Tipo == 0 && parametros.Nombre == "" {
		return &productos, nil
	}

	return &productosFiltrados, nil
}
func (repository CompraRepository) PostNuevaCompra(usuarioID string, idsComprasSeleccionadas []primitive.ObjectID) (*model.Compra, error) {
	// Llamar al método para obtener productos cuya cantidad mínima sea menor a la cantidad actual
	parametros := dto.ParametrosProductosCantidad{}
	productos, err := repository.GetProductosPorCantidadMinima(parametros, usuarioID)
	if err != nil {
		return nil, err
	}

	// Verificar si hay productos disponibles
	if len(*productos) == 0 {
		return nil, errors.New("no hay productos con cantidad menor a la mínima")
	}

	// Filtrar productos basados en los IDs seleccionados
	var productosFiltrados []model.ProductoCompra
	for _, producto := range *productos {
		for _, idSeleccionado := range idsComprasSeleccionadas {
			if producto.AlimentoId == idSeleccionado {
				productosFiltrados = append(productosFiltrados, producto)
				break
			}
		}
	}

	// Verificar si hay productos seleccionados después del filtrado
	if len(productosFiltrados) == 0 {
		return nil, errors.New("no se seleccionaron productos válidos para la compra")
	}

	var costoTotal float64

	// Calcular el costo total de la compra
	for _, producto := range productosFiltrados {
		// Obtener el alimento correspondiente para acceder al precio unitario
		alimento, err := repository.getAlimentoByID(producto.AlimentoId)
		if err != nil {
			return nil, err
		}

		cantidadComprada := producto.Cantidad
		costoTotal += float64(cantidadComprada) * alimento.PrecioUnitario
	}

	// Crear la estructura de la compra
	compra := model.Compra{
		FechaCreacion: time.Now(),
		CostoTotal:    costoTotal,
		Productos:     productosFiltrados,
		UsuarioID:     usuarioID,
	}

	// Insertar la compra en la colección 'compras'
	collectionCompras := repository.db.GetClient().Database("gocooking").Collection("compras")
	resultado, err := collectionCompras.InsertOne(context.TODO(), compra)
	if err != nil {
		return nil, err
	}

	// Actualizar el ID de la compra con el generado por MongoDB
	compra.Id = resultado.InsertedID.(primitive.ObjectID)

	// Actualizar la cantidad de productos en la colección 'alimentos' para que vuelva a la cantidad mínima
	for _, producto := range productosFiltrados {
		// Obtener el alimento correspondiente
		alimento, err := repository.getAlimentoByID(producto.AlimentoId)
		if err != nil {
			return nil, err
		}

		// Calcular la nueva cantidad
		_, err = repository.db.GetClient().Database("gocooking").Collection("alimentos").UpdateOne(
			context.TODO(),
			bson.M{"_id": producto.AlimentoId},
			bson.M{"$set": bson.M{"cantidad_actual": alimento.CantidadMinima * 2}},
		)
		if err != nil {
			return nil, err
		}
	}

	return &compra, nil
}

func (repository CompraRepository) GetCompras(usuarioID string) (*[]model.Compra, error) {
	collection := repository.db.GetClient().Database("gocooking").Collection("compras")
	filtro := bson.M{
		"id_usuario": usuarioID,
	}
	cursor, err := collection.Find(context.TODO(), filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var compras []model.Compra

	for cursor.Next(context.Background()) {
		var compra model.Compra
		if err := cursor.Decode(&compra); err != nil {
			return nil, err
		}

		// Agregar la compra a la lista
		compras = append(compras, compra)
	}

	return &compras, nil
}
func (repository CompraRepository) getAlimentoByID(id primitive.ObjectID) (*model.Alimento, error) {
	// Conectar a la colección 'alimentos'
	collection := repository.db.GetClient().Database("gocooking").Collection("alimentos")

	// Buscar el alimento por ID
	var alimento model.Alimento
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&alimento)
	if err != nil {
		return nil, err
	}

	return &alimento, nil
}

func (repository CompraRepository) GetCostoPromedioPorMesUltimoAnio(usuarioID string) (map[string]float64, error) {
	collection := repository.db.GetClient().Database("gocooking").Collection("compras")

	// Obtener el primer día del año actual
	fechaInicio := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Local)

	filtro := bson.M{
		"id_usuario": usuarioID,
		"fecha_creacion": bson.M{
			"$gte": fechaInicio,
		},
	}

	cursor, err := collection.Find(context.TODO(), filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	costoPromedioPorMes := make(map[string]float64)
	cantidadComprasPorMes := make(map[string]int)

	// Inicializar todos los meses del año actual con 0
	for mes := time.January; mes <= time.December; mes++ {
		mesKey := time.Date(time.Now().Year(), mes, 1, 0, 0, 0, 0, time.Local).Format("01-2006")
		costoPromedioPorMes[mesKey] = 0
		cantidadComprasPorMes[mesKey] = 0
	}

	for cursor.Next(context.Background()) {
		var compra model.Compra
		if err := cursor.Decode(&compra); err != nil {
			return nil, err
		}

		// Obtener el mes y año de la compra en formato "MM-YYYY"
		mes := compra.FechaCreacion.Format("01-2006")

		costoPromedioPorMes[mes] += compra.CostoTotal
		cantidadComprasPorMes[mes]++
	}

	// Calcular el costo promedio por mes
	for mes, costoTotal := range costoPromedioPorMes {
		if cantidadComprasPorMes[mes] > 0 {
			costoPromedioPorMes[mes] = costoTotal / float64(cantidadComprasPorMes[mes])
		}
	}

	return costoPromedioPorMes, nil
}

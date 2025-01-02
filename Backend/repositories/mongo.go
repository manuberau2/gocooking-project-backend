package repositories

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB struct para mantener la conexión con MongoDB
type MongoDB struct {
	Client *mongo.Client
}

// NewMongoDB crea una nueva instancia de la conexión de MongoDB
func NewMongoDB() (*MongoDB, error) {
	instancia := &MongoDB{}
	err := instancia.Connect()
	if err != nil {
		log.Printf("Error al conectar a MongoDB: %v", err) // Log de error
		return nil, err
	}
	log.Println("Conexión a MongoDB establecida exitosamente") // Log de éxito
	return instancia, nil
}

// GetClient devuelve el cliente de MongoDB
func (mongoDB *MongoDB) GetClient() *mongo.Client {
	return mongoDB.Client
}

// Connect establece la conexión con MongoDB
func (mongoDB *MongoDB) Connect() error {
	if mongoDB.Client != nil {
		// Si ya existe una conexión, la desconectamos antes de volver a intentar
		log.Println("Desconectando la conexión previa a MongoDB...") // Log de desconexión previa
		if err := mongoDB.Disconnect(); err != nil {
			log.Printf("Error al desconectar la conexión anterior: %v", err) // Log de error
			return err
		}
	}

	// Configurar las opciones de conexión
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	log.Printf("Conectando a MongoDB con URI: %s", clientOptions.GetURI()) // Log de URI

	// Intentar conectar a MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Printf("Error al intentar conectar a MongoDB: %v", err) // Log de error
		return err
	}

	// Comprobar la conexión
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // timeout de 10 segundos
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Error al hacer ping a MongoDB: %v", err) // Log de error de ping
		return err
	}

	// Asignar el cliente al campo de MongoDB
	mongoDB.Client = client
	log.Println("Conexión con MongoDB exitosa") // Log de éxito
	return nil
}

// Disconnect cierra la conexión con MongoDB
func (mongoDB *MongoDB) Disconnect() error {
	if mongoDB.Client != nil {
		log.Println("Desconectando de MongoDB...") // Log de desconexión
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := mongoDB.Client.Disconnect(ctx)
		if err != nil {
			log.Printf("Error al desconectar de MongoDB: %v", err) // Log de error al desconectar
			return err
		}
		log.Println("Desconexión de MongoDB exitosa") // Log de éxito al desconectar
	}
	return nil
}

package main

import (
	"GoCooking/Backend/clients"
	"GoCooking/Backend/handlers"
	"GoCooking/Backend/middlewares"
	"GoCooking/Backend/repositories"
	"GoCooking/Backend/service"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	router           *gin.Engine
	alimentosHandler *handlers.AlimentoHandler
	recetasHandler   *handlers.RecetaHandler
	compraHandler    *handlers.CompraHandler
)

func main() {
	router = gin.Default()
	dependencies()
	mappingRoutes()
	log.Println("Iniciando servidor en el puerto 8080...")
	router.Run(":8080")

}

func dependencies() {
	var database repositories.DB

	var alimentosRepository repositories.AlimentoRepositoryInterface
	var recetasRepository repositories.RecetaRepositoryInterface
	var comprasRepository repositories.CompraRepositoryInterface

	var alimentosService service.AlimentoInterface
	var recetasService service.RecetaInterface
	var comprasService service.CompraInterface
	//Inyectar repositorios
	database = repositories.NewMongoDB()
	alimentosRepository = repositories.NewAlimentoRepository(database)
	recetasRepository = repositories.NewRecetaRepository(database)
	comprasRepository = repositories.NewCompraRepository(database)
	//Inyectar servicios
	alimentosService = service.NewAlimentoService(alimentosRepository)
	recetasService = service.NewRecetaService(recetasRepository)
	comprasService = service.NewCompraService(comprasRepository)
	//Inyectar handlers
	alimentosHandler = handlers.NewAlimentoHandler(alimentosService)
	recetasHandler = handlers.NewRecetaHandler(recetasService)
	compraHandler = handlers.NewCompraHandler(comprasService)

}

func mappingRoutes() {
	authClients := clients.NewAuthClient()
	authMiddleware := middlewares.NewAuthMiddleware(authClients)
	router.Use(middlewares.CORSMiddleware())
	router.Use(authMiddleware.ValidateToken)
	//Ruta alimentos
	groupAlimentos := router.Group("/alimentos")

	groupAlimentos.GET("/", alimentosHandler.GetAlimentos)
	groupAlimentos.GET("/:id", alimentosHandler.GetAlimentoByID)
	groupAlimentos.POST("/", alimentosHandler.InsertAlimento)
	groupAlimentos.PUT("/:id", alimentosHandler.UpdateAlimento)
	groupAlimentos.DELETE("/:id", alimentosHandler.DeleteAlimento)

	//Ruta recetas
	groupRecetas := router.Group("/recetas")

	groupRecetas.GET("/", recetasHandler.GetRecetas)
	groupRecetas.GET("/:id", recetasHandler.GetRecetaByID)
	groupRecetas.GET("/buscar", recetasHandler.GetRecetasByParameters)
	groupRecetas.POST("/", recetasHandler.InsertReceta)
	groupRecetas.PUT("/:id", recetasHandler.UpdateReceta)
	groupRecetas.DELETE("/:id", recetasHandler.DeleteReceta)

	//Ruta compras
	groupCompras := router.Group("/compras")

	groupCompras.GET("/", compraHandler.GetCompras)
	groupCompras.GET("/productos-cantidad", compraHandler.GetProductosPorCantidadMinima)
	groupCompras.POST("/", compraHandler.PostNuevaCompra)

	groupReportes := router.Group("/reportes")

	groupReportes.GET("/recetas-momento", recetasHandler.GetCantidadRecetasPorMomento)
	groupReportes.GET("/recetas-tipo-alimento", recetasHandler.GetCantidadRecetasPorTipoAlimento)
	groupReportes.GET("/costo-promedio-mes", compraHandler.GetCostoPromedioPorMesUltimoAnio)

}

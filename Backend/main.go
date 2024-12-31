package handler

import (
	"gocooking-backend/clients"
	"gocooking-backend/handlers"
	"gocooking-backend/middlewares"
	"gocooking-backend/repositories"
	"gocooking-backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Global variables
var (
	router           *gin.Engine
	alimentosHandler *handlers.AlimentoHandler
	recetasHandler   *handlers.RecetaHandler
	compraHandler    *handlers.CompraHandler
)

func init() {
	router = gin.Default()
	dependencies()
	mappingRoutes()
}

// dependencies initializes the dependencies like database and services
func dependencies() {
	var database repositories.DB

	var alimentosRepository repositories.AlimentoRepositoryInterface
	var recetasRepository repositories.RecetaRepositoryInterface
	var comprasRepository repositories.CompraRepositoryInterface

	var alimentosService service.AlimentoInterface
	var recetasService service.RecetaInterface
	var comprasService service.CompraInterface

	// Initialize repositories
	database = repositories.NewMongoDB()
	alimentosRepository = repositories.NewAlimentoRepository(database)
	recetasRepository = repositories.NewRecetaRepository(database)
	comprasRepository = repositories.NewCompraRepository(database)

	// Initialize services
	alimentosService = service.NewAlimentoService(alimentosRepository)
	recetasService = service.NewRecetaService(recetasRepository)
	comprasService = service.NewCompraService(comprasRepository)

	// Initialize handlers
	alimentosHandler = handlers.NewAlimentoHandler(alimentosService)
	recetasHandler = handlers.NewRecetaHandler(recetasService)
	compraHandler = handlers.NewCompraHandler(comprasService)
}

// mappingRoutes defines all the routes for the application
func mappingRoutes() {
	authClients := clients.NewAuthClient()
	authMiddleware := middlewares.NewAuthMiddleware(authClients)
	router.Use(middlewares.CORSMiddleware())
	router.Use(authMiddleware.ValidateToken)

	// Alimentos routes
	groupAlimentos := router.Group("/alimentos")
	groupAlimentos.GET("/", alimentosHandler.GetAlimentos)
	groupAlimentos.GET("/:id", alimentosHandler.GetAlimentoByID)
	groupAlimentos.POST("/", alimentosHandler.InsertAlimento)
	groupAlimentos.PUT("/:id", alimentosHandler.UpdateAlimento)
	groupAlimentos.DELETE("/:id", alimentosHandler.DeleteAlimento)

	// Recetas routes
	groupRecetas := router.Group("/recetas")
	groupRecetas.GET("/", recetasHandler.GetRecetas)
	groupRecetas.GET("/:id", recetasHandler.GetRecetaByID)
	groupRecetas.GET("/buscar", recetasHandler.GetRecetasByParameters)
	groupRecetas.POST("/", recetasHandler.InsertReceta)
	groupRecetas.PUT("/:id", recetasHandler.UpdateReceta)
	groupRecetas.DELETE("/:id", recetasHandler.DeleteReceta)

	// Compras routes
	groupCompras := router.Group("/compras")
	groupCompras.GET("/", compraHandler.GetCompras)
	groupCompras.GET("/productos-cantidad", compraHandler.GetProductosPorCantidadMinima)
	groupCompras.POST("/", compraHandler.PostNuevaCompra)

	// Reportes routes
	groupReportes := router.Group("/reportes")
	groupReportes.GET("/recetas-momento", recetasHandler.GetCantidadRecetasPorMomento)
	groupReportes.GET("/recetas-tipo-alimento", recetasHandler.GetCantidadRecetasPorTipoAlimento)
	groupReportes.GET("/costo-promedio-mes", compraHandler.GetCostoPromedioPorMesUltimoAnio)
}

// Vercel will call this function
func Handler(w http.ResponseWriter, r *http.Request) {
	// Use Gin to serve the request
	router.ServeHTTP(w, r)
}

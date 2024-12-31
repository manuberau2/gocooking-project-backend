package handlers

import (
	"GoCooking/Backend/dto"
	"GoCooking/Backend/service"
	"GoCooking/Backend/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CompraHandler struct {
	compraService service.CompraInterface
}

func NewCompraHandler(compraService service.CompraInterface) *CompraHandler {
	return &CompraHandler{
		compraService: compraService,
	}
}
func (handler *CompraHandler) GetProductosPorCantidadMinima(c *gin.Context) {
	usuario := utils.GetUserInfoFromContext(c)
	log.Printf("[handler:CompraHandler][method:GetProductosPorCantidadMinima][status:before_service_call][user:%s]", usuario)
	var parametros dto.ParametrosProductosCantidad
	err := c.ShouldBindQuery(&parametros)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parametros invalidos"})
		return
	}
	log.Printf("parametros: %+v", parametros)
	productos, appErr := handler.compraService.GetProductosPorCantidadMinima(parametros, usuario.Codigo)
	log.Printf("[handler:CompraHandler][method:GetProductosPorCantidadMinima][status:after_service_call][cantidad:%d][user:%s]", len(productos), usuario)
	if appErr != nil {
		if appErr.Codigo == "ERR_404" {
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}
	c.JSON(http.StatusOK, productos)
}
func (handler *CompraHandler) PostNuevaCompra(c *gin.Context) {
	usuario := utils.GetUserInfoFromContext(c)

	log.Printf("[handler:CompraHandler][method:PostNuevaCompra][status:before_service_call][user: %s]", usuario)

	var requestBody struct {
		IdsComprasSeleccionadas []string `json:"ids_compras_seleccionadas"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("[handler:CompraHandler][method:PostNuevaCompra][status:error_parsing_request][error: %s]", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida, verifica el formato del JSON."})
		return
	}

	compra, appErr := handler.compraService.PostNuevaCompra(usuario.Codigo, requestBody.IdsComprasSeleccionadas)
	if appErr != nil {
		if appErr.Codigo == "ERR_400" {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}

	log.Printf("[handler:CompraHandler][method:PostNuevaCompra][status:success][compra_id: %s][user: %s]", compra.ID, usuario)

	c.JSON(http.StatusOK, compra)
}

func (handler *CompraHandler) GetCompras(c *gin.Context) {
	usuario := utils.GetUserInfoFromContext(c)
	log.Printf("[handler:CompraHandler][method:GetCompras][status:before_service_call][user:%s]", usuario)
	compras, appErr := handler.compraService.GetCompras(usuario.Codigo)
	log.Printf("[handler:CompraHandler][method:GetCompras][status:after_service_call][cantidad:%d][user:%s]", len(compras), usuario)
	if appErr != nil {
		if appErr.Codigo == "ERR_404" {
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}
	c.JSON(http.StatusOK, compras)
}
func (handler *CompraHandler) GetCostoPromedioPorMesUltimoAnio(c *gin.Context) {
	usuario := utils.GetUserInfoFromContext(c)
	log.Printf("[handler:CompraHandler][method:GetCostoPromedioPorMesUltimoAnio][status:before_service_call][user:%s]", usuario)
	costoPromedioPorMes, appErr := handler.compraService.GetCostoPromedioPorMesUltimoAnio(usuario.Codigo)
	log.Printf("[handler:CompraHandler][method:GetCostoPromedioPorMesUltimoAnio][status:after_service_call][user:%s]", usuario)
	if appErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}
	c.JSON(http.StatusOK, costoPromedioPorMes)
}

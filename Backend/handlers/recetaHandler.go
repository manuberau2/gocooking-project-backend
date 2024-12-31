package handlers

import (
	"GoCooking/Backend/dto"
	"GoCooking/Backend/service"
	"GoCooking/Backend/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecetaHandler struct {
	recetaService service.RecetaInterface
}

func NewRecetaHandler(recetaService service.RecetaInterface) *RecetaHandler {
	return &RecetaHandler{
		recetaService: recetaService,
	}
}

func (handler *RecetaHandler) GetRecetas(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:RecetaHandler][method:GetRecetas][status:before_service_call][user:%s]", usuario.Codigo)
	recetas, err := handler.recetaService.GetRecetas(usuario.Codigo)
	log.Printf("[handler:RecetaHandler][method:GetRecetas][status:after_service_call][cantidad:%d][user:%s]", len(recetas), usuario.Codigo)
	if err != nil {
		if err.Codigo == "ERR_404" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Mensaje})
		return
	}
	c.JSON(http.StatusOK, recetas)
}
func (handler *RecetaHandler) GetRecetaByID(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:RecetaHandler][method:GetRecetaByID][status:before_service_call][user:%s]", usuario.Codigo)
	id := c.Param("id")
	receta, err := handler.recetaService.GetRecetaById(id)
	log.Printf("[handler:RecetaHandler][method:GetRecetaByID][status:after_service_call][receta:%s][user:%s]", receta.Id, usuario.Codigo)
	if err != nil {
		if err.Codigo == "ERR_404" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Mensaje})
		return
	}
	c.JSON(http.StatusOK, receta)
}
func (handler *RecetaHandler) InsertReceta(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:RecetaHandler][method:InsertReceta][status:before_service_call][user:%s]", usuario.Codigo)
	var receta dto.Receta
	err := c.BindJSON(&receta)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	receta.UsuarioID = usuario.Codigo
	_, appErr := handler.recetaService.InsertReceta(&receta)
	log.Printf("[handler:RecetaHandler][method:InsertReceta][status:after_service_call][user:%s]", usuario.Codigo)
	if appErr != nil {
		if appErr.Codigo == "ERR_400" {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}
	c.JSON(http.StatusCreated, receta)
}
func (handler *RecetaHandler) UpdateReceta(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:RecetaHandler][method:UpdateReceta][status:before_service_call][user:%s]", usuario.Codigo)
	var receta dto.Receta
	err := c.BindJSON(&receta)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	id := c.Param("id")
	receta.Id = id
	receta.UsuarioID = usuario.Codigo
	_, appErr := handler.recetaService.UpdateReceta(&receta)
	log.Printf("[handler:RecetaHandler][method:UpdateReceta][status:after_service_call][user:%s]", usuario.Codigo)
	if appErr != nil {
		if appErr.Codigo == "ERR_400" {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return

	}
	c.JSON(http.StatusCreated, receta)
}
func (handler *RecetaHandler) DeleteReceta(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:RecetaHandler][method:DeleteReceta][status:before_service_call][user:%s]", usuario.Codigo)
	id := c.Param("id")
	_, appErr := handler.recetaService.DeleteReceta(id)
	log.Printf("[handler:RecetaHandler][method:DeleteReceta][status:after_service_call][user:%s]", usuario.Codigo)
	if appErr != nil {
		if appErr.Codigo == "ERR_404" {
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Receta eliminada"})
}
func (handler *RecetaHandler) GetRecetasByParameters(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:RecetaHandler][method:GetRecetasByParameters][status:before_service_call][user:%s]", usuario.Codigo)

	var parametros dto.ParametrosReceta
	err := c.ShouldBindQuery(&parametros)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	recetas, appErr := handler.recetaService.GetRecetasByParameters(parametros, usuario.Codigo)
	log.Printf("[handler:RecetaHandler][method:GetRecetasByParameters][status:after_service_call][cantidad:%d][user:%s]", len(recetas), usuario.Codigo)

	if appErr != nil {
		if appErr.Codigo == "ERR_400" {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}

	c.JSON(http.StatusOK, recetas)
}

func (handler *RecetaHandler) GetCantidadRecetasPorMomento(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:RecetaHandler][method:GetCantidadRecetasPorMomento][status:before_service_call][user:%s]", usuario.Codigo)

	cantidadRecetasPorMomento, appErr := handler.recetaService.GetCantidadRecetasPorMomento(usuario.Codigo)
	log.Printf("[handler:RecetaHandler][method:GetCantidadRecetasPorMomento][status:after_service_call][user:%s]", usuario.Codigo)

	if appErr != nil {
		if appErr.Codigo == "ERR_404" {
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}

	c.JSON(http.StatusOK, cantidadRecetasPorMomento)
}

func (handler *RecetaHandler) GetCantidadRecetasPorTipoAlimento(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:RecetaHandler][method:GetCantidadRecetasPorTipoAlimento][status:before_service_call][user:%s]", usuario.Codigo)

	cantidadRecetasPorTipoAlimento, appErr := handler.recetaService.GetCantidadRecetasPorTipoAlimento(usuario.Codigo)
	log.Printf("[handler:RecetaHandler][method:GetCantidadRecetasPorTipoAlimento][status:after_service_call][user:%s]", usuario.Codigo)

	if appErr != nil {
		if appErr.Codigo == "ERR_404" {
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}

	c.JSON(http.StatusOK, cantidadRecetasPorTipoAlimento)
}

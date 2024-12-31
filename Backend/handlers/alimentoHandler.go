package handlers

import (
	"GoCooking/Backend/dto"
	"GoCooking/Backend/service"
	"GoCooking/Backend/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AlimentoHandler struct {
	alimentoService service.AlimentoInterface
}

func NewAlimentoHandler(alimentoService service.AlimentoInterface) *AlimentoHandler {
	return &AlimentoHandler{
		alimentoService: alimentoService,
	}
}
func (handler *AlimentoHandler) GetAlimentos(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:AlimentoHandler][method:GetAlimentos][status:before_service_call][user:%s]", usuario.Codigo)
	alimentos, err := handler.alimentoService.GetAlimentos(usuario.Codigo)
	log.Printf("[handler:AlimentoHandler][method:GetAlimentos][status:after_service_call][cantidad:%d][user:%s]", len(alimentos), usuario.Codigo)
	if err != nil {
		if err.Codigo == "ERR_404" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Mensaje})
		return
	}

	c.JSON(http.StatusOK, alimentos)
}
func (handler *AlimentoHandler) GetAlimentoByID(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:AlimentoHandler][method:GetAlimentoByID][status:before_service_call][user:%s]", usuario.Codigo)
	id := c.Param("id")
	alimento, err := handler.alimentoService.GetAlimentoByID(id)
	log.Printf("[handler:AlimentoHandler][method:GetAlimentoByID][status:after_service_call][alimento:%s][user:%s]", alimento.Id, usuario.Codigo)
	if err != nil {
		if err.Codigo == "ERR_404" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Mensaje})
		return
	}
	c.JSON(http.StatusOK, alimento)
}
func (handler *AlimentoHandler) InsertAlimento(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:AlimentoHandler][method:InsertAlimento][status:before_service_call][user:%s]", usuario.Codigo)
	var alimento dto.Alimento
	err := c.BindJSON(&alimento)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el body"})
		return
	}
	alimento.UsuarioID = usuario.Codigo
	success, appErr := handler.alimentoService.InsertAlimento(&alimento)
	log.Printf("[handler:AlimentoHandler][method:InsertAlimento][status:after_service_call][success:%t][user:%s]", success, usuario.Codigo)
	if appErr != nil {
		if appErr.Codigo == "ERR_400" {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": success})
}
func (handler *AlimentoHandler) UpdateAlimento(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:AlimentoHandler][method:UpdateAlimento][status:before_service_call][user:%s]", usuario.Codigo)
	var alimento dto.Alimento
	err := c.BindJSON(&alimento)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el body"})
		return
	}
	id := c.Param("id")
	alimento.Id = id
	alimento.UsuarioID = usuario.Codigo
	success, appErr := handler.alimentoService.UpdateAlimento(&alimento)
	log.Printf("[handler:AlimentoHandler][method:UpdateAlimento][status:after_service_call][success:%t][user:%s]", success, usuario.Codigo)
	if appErr != nil {
		if appErr.Codigo == "ERR_400" {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Mensaje})
			return
		}
		if appErr.Codigo == "ERR_404" {
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": success})
}
func (handler *AlimentoHandler) DeleteAlimento(c *gin.Context) {
	usuario := dto.NewUsuario(utils.GetUserInfoFromContext(c))
	log.Printf("[handler:AlimentoHandler][method:DeleteAlimento][status:before_service_call][user:%s]", usuario.Codigo)
	id := c.Param("id")
	success, appErr := handler.alimentoService.DeleteAlimento(id)
	log.Printf("[handler:AlimentoHandler][method:DeleteAlimento][status:after_service_call][success:%t][user:%s]", success, usuario.Codigo)
	if appErr != nil {
		if appErr.Codigo == "ERR_404" {
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Mensaje})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Mensaje})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": success})
}

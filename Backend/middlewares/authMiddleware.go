package middlewares

import (
	"GoCooking/Backend/clients"
	"GoCooking/Backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authClient clients.AuthClientInterface
}

func NewAuthMiddleware(authClient clients.AuthClientInterface) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

// Este middleware se ejecuta en el grupo de rutas privadas.
func (auth *AuthMiddleware) ValidateToken(c *gin.Context) {
	//Se obtiene el header necesario con nombre "Authorization"
	authToken := c.GetHeader("Authorization")

	if authToken == "" {
		//log.Printf("[service:AulaService][method:ObtenerAulaPorId][reason:NOT_FOUND][id:%s]", id)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token no encontrado"})
		return
	}

	//Obtener la informacion del usuario a partir del token desde el servicio externo
	user, err := auth.authClient.GetUserInfo(authToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	//Validar que el usuario tenga alguno de todos los roles que yo quiero en mi aplicacion.

	//Seteamos los datos del usuario logueado en el contexto de GIN.
	utils.SetUserInContext(c, user)

	c.Next()
}

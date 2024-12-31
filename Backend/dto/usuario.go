package dto

import (
	"GoCooking/Backend/clients/responses"
)

type Usuario struct {
	Codigo   string `json:"codigo"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Rol      string `json:"rol"`
}

func NewUsuario(usuario *responses.UsuarioInfo) *Usuario {
	usuarioNew := Usuario{}
	if usuario != nil {
		usuarioNew.Codigo = usuario.Codigo
		usuarioNew.Email = usuario.Email
		usuarioNew.Username = usuario.Username
		usuarioNew.Rol = usuario.Rol
	}
	return &usuarioNew
}

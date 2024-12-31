package responses

type UsuarioInfo struct {
	Codigo   string `json:"codigo"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Rol      string `json:"rol"`
}

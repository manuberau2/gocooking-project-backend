package utils

type AppError struct {
	Codigo  string
	Mensaje string
}

func NewAppError(codigo, mensaje string) *AppError {
	return &AppError{
		Codigo:  codigo,
		Mensaje: mensaje,
	}
}

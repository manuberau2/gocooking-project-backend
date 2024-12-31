package clients

import (
	"GoCooking/Backend/clients/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AuthClientInterface interface {
	GetUserInfo(token string) (*responses.UsuarioInfo, error)
}

type AuthClient struct {
}

func NewAuthClient() *AuthClient {
	return &AuthClient{}
}

func (auth *AuthClient) GetUserInfo(token string) (*responses.UsuarioInfo, error) {
	//Ruta donde apunta esta invocacion
	apiUrl := "http://w230847.ferozo.com/tp_prog2/api/account/UserInfo"

	// Crear un cliente HTTP
	client := &http.Client{}

	// Crear una solicitud GET
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		fmt.Println("Error al crear la solicitud GET:", err)
		return nil, err
	}

	// Agregar encabezado personalizado
	req.Header.Add("Authorization", token)

	// Realizar la solicitud GET
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al realizar la solicitud GET:", err)
		return nil, err
	}

	defer response.Body.Close()
	// Lee el cuerpo de la respuesta como una cadena
	responseBody, err := ioutil.ReadAll(response.Body)
	//Si el codigo es distinto de 200, es porque dio un error.
	if response.StatusCode != 200 {
		fmt.Println("Error al realizar la solicitud GET:", responseBody)
		return nil, errors.New("la peticion respondio con error")
	}

	if err != nil {
		fmt.Println("Error al leer el cuerpo de la respuesta:", err)
		return nil, err
	}

	// Convierte el cuerpo de la respuesta a una cadena
	bodyString := string(responseBody)

	var userInfo responses.UsuarioInfo

	if err := json.Unmarshal([]byte(bodyString), &userInfo); err != nil {
		fmt.Println("Error al deserializar el JSON:", err)
		return nil, err
	}

	fmt.Println("Código de estado:", response.Status)

	return &userInfo, nil
}

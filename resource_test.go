package resty

import (
	"log"
	"net/http"
	"testing"
)

type UsersController struct {
}

func (controller UsersController) Index(response http.ResponseWriter, request *http.Request, params map[string]string) {
	log.Printf("id = %v", params["id"])
}

func (controller UsersController) Show(response http.ResponseWriter, request *http.Request, params map[string]string) {
	log.Printf("id = %v", params["id"])
}

func TestRegisteringResource(t *testing.T) {
	userHandler := ResourceHandler{Name: "users", Controller: UsersController{}}
	userHandler.RegisterRoutes()

	http.ListenAndServe(":8080", nil)
}

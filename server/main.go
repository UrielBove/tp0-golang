package main

import (
	"net/http"
	"server/utils"
)

func main() {
	mux := http.NewServeMux() //Crea un router

	mux.HandleFunc("/paquetes", utils.RecibirPaquetes) //Define una ruta para recibir paquetes
	mux.HandleFunc("/mensaje", utils.RecibirMensaje)

	//panic("no implementado!")
	err := http.ListenAndServe(":8080", mux) //Inicia el servidor escuchando en el puerto 8080
	if err != nil {
		panic(err)
	}
}

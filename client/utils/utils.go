package utils

import (
	"bufio"
	"bytes"
	"client/globals"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

type Paquete struct {
	Valores []string `json:"valores"`
}

func IniciarConfiguracion(filePath string) *globals.Config {
	var config *globals.Config
	configFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}

func LeerConsola() {
	// Leer de la consola
	reader := bufio.NewReader(os.Stdin)
	//bufio.NewReader(os.Stdin): crea un lector de entrada estándar (la consola).

	log.Println("Ingrese los mensajes(ENTER vacío para salir):")

	for {
		text, err := reader.ReadString('\n') //ReadString('\n'): lee hasta que se presione ENTER.
		if err != nil {
			log.Printf("Error leyendo de consola: %s", err.Error())
			break
		}

		//Si solo se presiono ENTER, salir
		if text == "\n" {
			log.Println("Entrada vacía detectada. Saliendo del programa.")
			break
		}

		log.Printf("Mensaje ingresado: %s", text)
		//log.Print(text):muestra (y guarda en el log) lo que se leyó.
	}
}

func GenerarYEnviarPaquete() {
	paquete := Paquete{}
	// Leemos y cargamos el paquete

	reader := bufio.NewReader(os.Stdin)

	log.Println("Ingrese las lineas del paquete (ENTER vacio para salir)")

	for {
		text, _ := reader.ReadString('\n')

		if text == "\n" {
			break
		}
		paquete.Valores = append(paquete.Valores, text[:len(text)-1]) //Saco el /n del final
	}

	log.Printf("paquete a enviar: %+v", paquete)
	// Enviamos el paqute

	if globals.ClientConfig == nil {
		log.Println("No se pudo enviar el paquete: Configuracion no cargada.")
		return
	}
	EnviarPaquete(globals.ClientConfig.Ip, globals.ClientConfig.Puerto, paquete)
}

func EnviarMensaje(ip string, puerto int, mensajeTxt string) {
	mensaje := Mensaje{Mensaje: mensajeTxt}
	body, err := json.Marshal(mensaje)
	if err != nil {
		log.Printf("error codificando mensaje: %s", err.Error())
	}

	url := fmt.Sprintf("http://%s:%d/mensaje", ip, puerto)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("error enviando mensaje a ip:%s puerto:%d", ip, puerto)
	}

	log.Printf("respuesta del servidor: %s", resp.Status)
}

func EnviarPaquete(ip string, puerto int, paquete Paquete) {
	body, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("error codificando mensajes: %s", err.Error())
	}

	url := fmt.Sprintf("http://%s:%d/paquetes", ip, puerto)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("error enviando mensajes a ip:%s puerto:%d", ip, puerto)
	}

	log.Printf("respuesta del servidor: %s", resp.Status)
}

func ConfigurarLogger() {
	logFile, err := os.OpenFile("tp0.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	//WRONLY solo escritura, para logs
	if err != nil {
		log.Fatalf("Error al abrir archivo de log: %v", err)
		//Fatal me deja ponerle un texto en vez de panic
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//Flags primero para fecha, segundo hora, tercero muestra el nombre del archivo y
	//línea del log (útil para debug).
}

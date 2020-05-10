package main

import (
        "log"
		"fmt"	
		"github.com/gorilla/mux"
		"net/http"	
		"encoding/json"
		//"math/rand"
		//"strconv"
		"github.com/streadway/amqp"
)

type Caso struct {
	Nombre string `json:"nombre"`
	Departamento string `json:"departamento"`
	Edad string `json:"edad"`
	FormadeContagio string `json:"forma"`
	Estado string `json:"estado"`
}
func failOnError(err error, msg string) {
	if err != nil {
			log.Fatalf("%s: %s", msg, err)
	}
}
func envio_datos(caso Caso) {
	conn, err := amqp.Dial("amqp://guest:guest@35.237.178.105:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
			"task_queue", // name
			true,         // durable
			false,        // delete when unused
			false,        // exclusive
			false,        // no-wait
			nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body , _:= json.Marshal(caso)
	err = ch.Publish(
			"",           // exchange
			q.Name,       // routing key
			false,        // mandatory
			false,
			amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/json",
					Body:         []byte(body),
			})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body) 
}

func ingreso(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "ingreso.html")
}
func createEntrada(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var caso Caso
	_ = json.NewDecoder(r.Body).Decode(&caso)
	envio_datos(caso)
	json.NewEncoder(w).Encode(&caso)
}
var router = mux.NewRouter()

func main() {
    router.HandleFunc("/ingreso", ingreso).Methods("GET")
	router.HandleFunc("/ingreso", createEntrada).Methods("POST")
    http.Handle("/", router)
	fmt.Println("Servidor corriendo en http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}



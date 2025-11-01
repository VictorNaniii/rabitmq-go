package main

import (
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleRiderWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrader error:", err)
		return
	}
	defer conn.Close()
	userID := r.URL.Query().Get("userId")
	if userID != "" {
		log.Println("No user ID providers")
		return
	}
	for {
		_, mesage, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}
		log.Println("Recived message: ", mesage)
	}
}

func handleDriverWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrader error: ", err)
	}
	defer conn.Close()
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("No user ID provider")
		return
	}
	packageSlug := r.URL.Query().Get("packageSLug")
	if packageSlug == "" {
		log.Println("No slug providers")
		return
	}

	type Driver struct {
		Id             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profilePicture"`
		CarPlate       string `json:"carPlate"`
		PackageSlug    string `json:"packageSlug"`
	}
	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			Id:             userID,
			Name:           "VCitor",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate:       "ABFS31",
			PackageSlug:    packageSlug,
		},
	}
	if err := conn.WriteJSON(msg); err != nil {
		log.Println("Error sending message:", err)
		return
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Print("Error reading message:", err)
			break
		}
		log.Println("Recived message: ", message)
	}
}

package main

import (
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_client"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/proto/driver"

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
	ctx := r.Context()
	driverService, err := grpc_client.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		driverService.Client.UnregisterDriver(ctx, &driver.RegisterDriverRequest{
			DriverID:    userID,
			PackageSlug: packageSlug,
		})
		driverService.Close()
		log.Println("Driver unregister: ", userID)
	}()
	driverData, err := driverService.Client.RegisterDriver(ctx, &driver.RegisterDriverRequest{
		DriverID:    userID,
		PackageSlug: packageSlug,
	})
	if err != nil {
		log.Println("Error register driver: %v", err)
		return
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: driverData.Driver,
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

package main

import (
	"log"
	"media-net/appremove"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	//funcframework.RegisterHTTPFunction("/", appremove.GetAppRemoveEvents)
	//funcframework.RegisterHTTPFunction("/", appremove.ReadFireStore1)
	//funcframework.RegisterHTTPFunction("/", appremove.ReadFireStore)
	//funcframework.RegisterHTTPFunction("/", appremove.SaveUninstallData)
	funcframework.RegisterHTTPFunction("/", appremove.ReadDataStore)

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}

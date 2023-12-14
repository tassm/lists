package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/tassm/lists/internal/api"
	"github.com/tassm/lists/internal/data"
)

const (
	dynamoListItemsTable = "list_items"
	dynamoListNamesTable = "list_names"

	// base URL path for list entities
	basePath = "/api/v1/list/"
)

func main() {
	// setup and connect to aws
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Panicf("failed to load AWS configration!")
	}
	client := dynamodb.NewFromConfig(config) //dynamodb.(dynamodb.Options{}, nil)
	svc := data.NewDynamoListService(client, dynamoListItemsTable, dynamoListNamesTable)

	mux := api.HandlerMux(svc, basePath)

	slog.Info("serving on :8080")
	http.ListenAndServe(":8080", mux)
}

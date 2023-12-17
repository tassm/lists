package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/tassm/lists/internal/data"
)

const (
	dynamoTableName = "list_items"

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
	itemListService := data.NewDynamoListService(client, dynamoTableName)

	http.HandleFunc(basePath, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		listId := path[len(basePath):]
		switch r.Method {
		case http.MethodGet:
			// GET /list
			// Respond with a list of all items in the list
			w.Header().Set("Content-Type", "application/json")
			if res, err := itemListService.GetListItems(r.Context(), listId); err == nil {
				if json, err := json.Marshal(res); err == nil {
					w.Write(json)
					return
				}
			} else {
				log.Printf("failed to retrieve items: %s", err.Error())
				http.Error(w, "failed to save item", http.StatusInternalServerError)
				return
			}
		case http.MethodPost:
			// POST /list
			// Add a new item to the list
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var item data.ListItem
			err = json.Unmarshal(body, &item)
			if err != nil || item.Item == "" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			uuid, _ := uuid.NewRandom()
			item.ID = uuid.String()
			item.Done = false
			err = itemListService.CreateListItem(r.Context(), &item)
			if err != nil {
				log.Printf("failed to save item: %s", err.Error())
				http.Error(w, "failed to save item", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var items []data.ListItem
			err = json.Unmarshal(body, &items)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// TODO: make this transactional
			for _, item := range items {
				err := itemListService.UpdateListItem(r.Context(), &item)
				if err != nil {
					log.Printf("failed to update item: %s", err.Error())
					http.Error(w, "failed to update item", http.StatusInternalServerError)
					return
				}
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Serve a folder of web content at the root path /
	http.Handle("/", http.FileServer(http.Dir("./web")))

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

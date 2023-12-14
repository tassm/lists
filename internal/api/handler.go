package api

import (
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/tassm/lists/internal/data"
)

func HandlerMux(itemSvc *data.DynamoListItemService, basePath string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(basePath, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		listId := path[len(basePath):]
		switch r.Method {
		case http.MethodGet:
			// validate the list exists
			err := itemSvc.IsValidList(r.Context(), listId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			// Respond with a list of all items in the list
			w.Header().Set("Content-Type", "application/json")
			if res, err := itemSvc.GetListItems(r.Context(), listId); err == nil {
				if res == nil {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("[]"))
					return
				} else {
					log.Println(res)
					if json, err := json.Marshal(res); err == nil {
						w.WriteHeader(http.StatusOK)
						w.Write(json)
						return
					}
				}
			} else {
				slog.Error("failed to retrieve items", "error", err.Error())
				http.Error(w, "failed to save item", http.StatusInternalServerError)
				return
			}
		case http.MethodPost:
			// POST /list
			// Add a new item to the list
			body, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("bad request", "error", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var item data.ListItem
			err = json.Unmarshal(body, &item)
			if err != nil || item.Item == "" {
				slog.Error("bad request", "error", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// validate the supplied list name
			err = itemSvc.IsValidList(r.Context(), item.ListID)
			if err != nil {
				slog.Error("bad request", "error", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			uuid, _ := uuid.NewRandom()
			item.ID = uuid.String()
			item.Done = false
			err = itemSvc.CreateListItem(r.Context(), &item)
			if err != nil {
				slog.Error("something went wrong", "error", err.Error())
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
			var listId string
			// TODO: make this transactional
			for i, item := range items {
				// validate the supplied list name
				if i == 0 {
					err = itemSvc.IsValidList(r.Context(), item.ListID)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					listId = item.ListID
				} else {
					if listId != item.ListID {
						slog.Error("bad request", "error", err.Error())
						http.Error(w, "trying to update items for multiple lists! very naughty...", http.StatusBadRequest)
						return
					}
				}
				// update the item
				err := itemSvc.UpdateListItem(r.Context(), &item)
				if err != nil {
					slog.Error("something went wrong", "error", err.Error())
					http.Error(w, "failed to update item", http.StatusInternalServerError)
					return
				}
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			slog.Warn("request not handled - not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	serverRoot, err := fs.Sub(data.WebFs, "web")
	if err != nil {
		log.Fatal(err)
	}

	// Serve a folder of web content at the root path from the embedded fs /
	mux.Handle("/", http.FileServer(http.FS(serverRoot)))
	return mux
}

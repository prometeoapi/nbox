package api

import (
	"encoding/json"
	"fmt"
	"nbox/internal/domain"
	"nbox/internal/domain/models"
	"net/http"
)

type Entrypoint interface {
	SetUp()
}

//type BoxHandler interface {
//	//GetBox(w http.ResponseWriter, r *http.Request)
//	//CreateBox(w http.ResponseWriter, r *http.Request)
//}

type BoxHandler struct {
	store domain.StoreOperations
}

func NewBoxHandler(store domain.StoreOperations) *BoxHandler {
	return &BoxHandler{store: store}
}

func (b *BoxHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	command := &models.Command[models.Box]{}
	//box :=  &models.Box{}
	if err := json.NewDecoder(r.Body).Decode(command); err != nil {
		http.Error(w, fmt.Sprintf(`{"message": "%s"}`, err), http.StatusBadRequest)
		return
	}

	box, err := b.store.CreateBox(&command.Payload)
	fmt.Printf("box result: %+v \n", box)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message": "%s"}`, err), http.StatusBadRequest)
		return
	}

	w.Write([]byte(`{"message": "ok"}`))
}

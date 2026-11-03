package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
)

const (
	//todo!! => must move to dotenv
	dbName = "hotelReservation_test"
)

type userStore struct {
	//embed interface into struct <<mock on demand>> pattern to
	//to satisfy db.UserStore without implementing every method

	db.UserStore
}

func newMongoUserStore(t *testing.T) *userStore {
	client := db.NewMongoDbClient()
	return &userStore{
		UserStore: db.NewMongoUserStore(client, dbName),
	}
}

func (tdb *userStore) tearDown(t *testing.T) {
	fmt.Println("--- dropping collection ---")
	if err := tdb.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("--- dropped collection successfully ---")
}

func TestPostUser(t *testing.T) {
	userStore := newMongoUserStore(t)
	defer userStore.tearDown(t)

	app := fiber.New()

	userHandler := NewUserHandler(userStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "Heliya",
		LastName:  "Zarei",
		Email:     "helisunz@gmail.com",
		Password:  "123456789",
	}
	marshalBody, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(marshalBody))
	req.Header.Add("Content-Type", "application/json")

	callResponse, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if callResponse.StatusCode != 200 {
		t.Fatalf("expected status 200 but got %s", callResponse.Status)
	}

	var user *types.User
	json.NewDecoder(callResponse.Body).Decode(&user)

	if len(user.Id) == 0 {
		t.Errorf("expected an userId but nothing found")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected the EncryptedPassword not to be included in json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName to be %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastName to be %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email to be %s but got %s", params.Email, user.Email)
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"context"
	"os"
	"time"
	"log"
	"math/rand"
	"strconv"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// user struct (Model)
type User struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Password string `json:"password"`
}

// post struct (Model)
type Post struct {
	ID     string  `json:"id"`
	Caption   string  `json:"caption"`
	ImageURL  string  `json:"img"`
	Timestamp string `json:"timestamp"`
}

// Init users var as a slice User struct
var users []User
// Init posts var as a slice Post struct
var posts []Post

// Get all user
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var people []User

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.TODO(), clientOptions)
	col := client.Database("First_Database").Collection("First Collection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, _ := col.Find(ctx, bson.M{})
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var person User
		cursor.Decode(&person)
		people = append(people, person)
	}
	json.NewEncoder(w).Encode(people)
}


// Add new user
func createUser(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.TODO(), clientOptions)

	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	bs, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost) //ecnrypting the user password then storing it
	user.Password = string(bs)
	// fmt.Println(user.Password)
	col := client.Database("First_Database").Collection("First Collection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, insertErr := col.InsertOne(ctx, user)
	if insertErr != nil {
		fmt.Println("InsertONE Error:", insertErr)
		os.Exit(1)
	}

	json.NewEncoder(w).Encode(result)
}


// Get single user
func getSingleUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person User
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.TODO(), clientOptions)
	col := client.Database("First_Database").Collection("First Collection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// id := string(params["id"])
	err := col.FindOne(ctx, User{ID: id}).Decode(&person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	} else {
		json.NewEncoder(w).Encode(person)
	}
}

// Update user
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			var user User
			_ = json.NewDecoder(r.Body).Decode(&user)
			user.ID = params["id"]
			users = append(users, user)
			json.NewEncoder(w).Encode(user)
			return
		}
	}
}

// Delete user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(users)
}


// Add new post
func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post Post
	_ = json.NewDecoder(r.Body).Decode(&post)
	post.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe
	posts = append(posts, post)
	json.NewEncoder(w).Encode(post)
}

// Get all posts
func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through posts and find one with the id from the params
	for _, item := range posts {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
		}
	}
	json.NewEncoder(w).Encode(&Post{})
}

// Pagination middleware is used to extract the next page id from the url query
func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PageID := r.URL.Query().Get(string(PageIDKey))
		intPageID := 0
		var err error
		if PageID != "" {
			intPageID, err = strconv.Atoi(PageID)
			if err != nil {
				_ = render.Render(w, r, types.ErrInvalidRequest(fmt.Errorf("couldn't read %s: %w", PageIDKey, err)))
				return
			}
		}
		ctx := context.WithValue(r.Context(), PageIDKey, intPageID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func TestGetEntryByID(t *testing.T) {

	req, err := http.NewRequest("GET", "/entry", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("id", "1")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetEntryByID)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"id":1,"first_name":"Lokesh","last_name":"Mishra","email_address":"lokesh@gmail.com","phone_number":"999999999"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateEntry(t *testing.T) {

	var jsonStr = []byte(`{"id":4,"first_name":"xyz","last_name":"pqr","email_address":"xyz@pqr.com","phone_number":"1234567890"}`)

	req, err := http.NewRequest("POST", "/entry", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateEntry)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":4,"first_name":"xyz","last_name":"pqr","email_address":"xyz@pqr.com","phone_number":"1234567890"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// Main function
func main() {

	// Init router
	r := mux.NewRouter()

	// Hardcoded data - @todo: add database
	users = append(users, User{ID: "1", Name: "Lokesh", Email: "user_One@gmail.com", Password: "yoyo"})
	users = append(users, User{ID: "2", Name: "Aman", Email: "user_Two@gmail.com", Password: "oyoy"})
	users = append(users, User{ID: "3", Name: "Samar", Email: "user_Three@gmail.com", Password: "oyyo"})

	posts = append(posts, Post{ID: "3", Caption: "Samar here", ImageURL: "samar_image@gmail.com", Timestamp: "20:20"})
	posts = append(posts, Post{ID: "3", Caption: "Samar here 2", ImageURL: "samar_image2@gmail.com", Timestamp: "08:45"})
	posts = append(posts, Post{ID: "1", Caption: "Lokesh here", ImageURL: "lokesh_image@gmail.com", Timestamp: "08:07"})

	// Route handles & endpoints
	r.HandleFunc("/all_users", getAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getSingleUser).Methods("POST")  //using post method here to pass password also in the request body for verification
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	r.HandleFunc("/posts", createPost).Methods("POST")
	r.HandleFunc("/posts/{id}", getPosts).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}




package main

import (

	"net/http"
	"encoding/json"
	"sync"
	"io/ioutil"
	"fmt"
	"time"
	"strings"
	"context"
	// "log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

const (
	DBName = "appointy"
	usersCollection = "users"

	URI = "mongodb://localhost:27017/"
)

type User struct {
	Id string `json:"Id"`
	Name string `json:"Name"`
	Email string `json:"Email"`
	Password string `json:"Password"`
}

type userHandlers struct {
	sync.Mutex
	store map[string]User
}

func (h *userHandlers) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w,r)
		return
	case "POST":
		h.post(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}


func (h *userHandlers) get(w http.ResponseWriter, r *http.Request){
	users := make([]User, len(h.store))

	h.Lock()
	i := 0
	for _, user := range h.store {
		users[i] = user
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(users)
	if err != nil {
		 w.WriteHeader(http.StatusInternalServerError)
		 w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *userHandlers) getUser(w http.ResponseWriter, r *http.Request){
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	user, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *userHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	
	// ct := r.Header.Get("content-type")
	// if ct != "application/json"{
	// 	w.WriteHeader(http.StatusUnsupportedMediaType)
	// 	w.Write([]byte(fmt.Sprintf("need content type application/json, but got '%s' ", ct)))
	// 	return
	// }

	var user User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return 
	}

	user.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[user.Id] = user
	defer h.Unlock()
}

func newUserHandlers() *userHandlers {
	return &userHandlers{
		store: map[string]User{},
	}
}

// FOR POSTS

type Posts struct {
	Id string `json:"Id"`
	Caption string `json:"Caption"`
	ImageURL string `json:"ImageURL"`
	TimeStamp string `json:"TimeStamp"`
}

type postsHandlers struct {
	sync.Mutex
	store map[string]Posts
}

func (h *postsHandlers) postss(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w,r)
		return
	case "POST":
		h.post(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *postsHandlers) get(w http.ResponseWriter, r *http.Request){
	postss := make([]Posts, len(h.store))

	h.Lock()
	i := 0
	for _, posts := range h.store {
		postss[i] = posts
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(postss)
	if err != nil {
		 w.WriteHeader(http.StatusInternalServerError)
		 w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *postsHandlers) getPosts(w http.ResponseWriter, r *http.Request){
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	posts, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(posts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *postsHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	
	// ct := r.Header.Get("content-type")
	// if ct != "application/json"{
	// 	w.WriteHeader(http.StatusUnsupportedMediaType)
	// 	w.Write([]byte(fmt.Sprintf("need content type application/json, but got '%s' ", ct)))
	// 	return
	// }

	var posts Posts
	err = json.Unmarshal(bodyBytes, &posts)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return 
	}

	posts.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[posts.Id] = posts
	defer h.Unlock()
}

func newPostsHandlers() *postsHandlers {
	return &postsHandlers{
		store: map[string]Posts{},
	}
}

func main() {
	userHandlers := newUserHandlers()
	http.HandleFunc("/users", userHandlers.users)
	http.HandleFunc("/users/",userHandlers.getUser)

	postsHandlers := newPostsHandlers()
	http.HandleFunc("/posts", postsHandlers.postss)
	http.HandleFunc("/posts/",postsHandlers.getPosts)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

	// DATABASE CONNECTION
	clientOpts := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(context.TODO(), clientOpts)

	if err != nil {
		fmt.Println(err)
		return
	}

	// ctx, _ := context.WithTimeout(context.Background(), 15*time.second)

	db := client.Database(DBName)
	coll :=db.Collection(usersCollection)

	fmt.Println(db.Name())
	fmt.Println(coll.Name())

	}
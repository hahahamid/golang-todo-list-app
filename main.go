package main

import(
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"context"
	"os"
	"os/signal"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevaddams/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


var rnd *renderer.Render
var db *mgo.Database

const (
	hostName 			string = "localhost:27017"
	dbName				string = "demo_todo"
	collectionName		string = "todo"
	port				string = ":8081"
)

type(
	todoModel struct {
		ID 			bson.ObjectId 	`bson:"_id,omitempty"`
		Title 		string 			`bson:"title"`
		Completed	bool 			`bson:"completed"`
		CreatedAt 	time.Time 		`bson:"createdAt"`
	}

	todo struct{
		ID			string 		`json:"id"`
		Title		string 		`json:"title"`
		Completed	bool 		`json:"completed"`
		CreatedAt 	time.Time 	`json:"created_at"`
	}
)

func init(){
	rnd = renderer.New()
	sess, err:= mgo.Dial(hostName)
	checkErr(err)
	sess.SetMode(mgo.Monotonic, true)
	db = sess.DB(dbName)
}


func homeHandler(w http.ResponseWriter, r *http.Request){
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	checkErr(err)
}


func main(){
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandlers())
	
	
	srv := &http.Server{
		Addr: port,
		Handler: r,
		ReadTimeout: 60*time.Second,
		WriteTimeout: 60*time.Second,
		IdleTimeout: 60*time.Second,
	}

	go func(){
		log.Println("Listening to port", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

}

func todoHandlers(){
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router){
		r.Get("/", fetchItems)
		r.Post("/", createItem)
		r.Put("/{id}", updateItem)
		r.Delete("/{id}", deleteItem) 
	}) 
	return rg
}

func checkErr (err error){
	if err!=nil{
		log.Fatal(err)
	}
}
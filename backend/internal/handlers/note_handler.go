package handlers

/**
	TODO This should be separated into three different components, maybe change it on the future. Because its kinda short for now it will suffice.
**/

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbHost          = "localhost"
	dbPort          = "27017"
	dbName          = "shared-notes"
	notesCollection = "notes"

	connTimeout = 20 * time.Second
	reqTimeout  = 10 * time.Second

	serviceHost = "localhost"
	servicePort = "8000"
)

// Note represents a simple shared note
type Note struct {
	ID      string `json:"_id" bson:"_id"`
	Title   string `json:"_t" bson:"_t"`
	Content string `json:"_c" bson:"_c"`
}

// NoteHandler holds MongoDB client and HTTP server
type NoteHandler struct {
	client  *mongo.Client
	db      *mongo.Database
	server  *http.Server
}

// New creates a new NoteHandler and sets up the server
func New() (*NoteHandler, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://" + dbHost + ":" + dbPort))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	handler := &NoteHandler{client: client, db: db}

	// Set up routes
	r := mux.NewRouter()
	r.HandleFunc("/notes", handler.handleCreateNote).Methods("POST")
	r.HandleFunc("/notes", handler.handleGetNoteNames).Methods("GET")
	r.HandleFunc("/notes/{id}", handler.handleDeleteNote).Methods("DELETE")

	// HTTP server
	srv := &http.Server{
		Handler:      r,
		Addr:         serviceHost + ":" + servicePort,
		WriteTimeout: reqTimeout,
		ReadTimeout:  reqTimeout,
	}
	handler.server = srv

	return handler, nil
}

func (h *NoteHandler) Start() error {
	return h.server.ListenAndServe()
}

// -- Handlers --

func (h *NoteHandler) handleCreateNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	_, err := h.db.Collection(notesCollection).InsertOne(ctx, note)
	if err != nil {
		http.Error(w, "failed to insert", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *NoteHandler) handleGetNoteNames(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	cursor, err := h.db.Collection(notesCollection).Find(ctx, bson.D{}, options.Find().SetProjection(bson.M{"_t": 1}))
	if err != nil {
		http.Error(w, "failed to fetch", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		http.Error(w, "cursor error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *NoteHandler) handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	_, err := h.db.Collection(notesCollection).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

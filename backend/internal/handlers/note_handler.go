package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID      *string `json:"_id" bson:"_id,omitempty"`
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

	// Register routes
	r.HandleFunc("/notes", handler.handleCreateNote).Methods("POST")
	r.HandleFunc("/notes", handler.handleGetNoteNames).Methods("GET")
	r.HandleFunc("/notes/{_id}", handler.handleDeleteNote).Methods("DELETE")
	r.HandleFunc("/notes/{_id}/rename/{_t}", handler.handleRenameNote).Methods("PUT")

	// Apply CORS
	corsWrappedRouter := enableCORS(r)

	// HTTP server
	srv := &http.Server{
		Handler:      corsWrappedRouter,
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
	note := Note{
		Title:   "My note",
		Content: "",
	}

	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	res, err := h.db.Collection(notesCollection).InsertOne(ctx, note)
	if err != nil {
		http.Error(w, "failed to insert", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bson.M{"_id": res.InsertedID})
}


func (h *NoteHandler) handleGetNoteNames(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	cursor, err := h.db.Collection(notesCollection).Find(ctx, bson.D{}, options.Find().SetProjection(bson.M{"_t": 1, "_id": 1}))
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
	id := mux.Vars(r)["_id"]
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	
	_, err = h.db.Collection(notesCollection).DeleteOne(ctx, bson.M{"_id": objectID})

	if err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NoteHandler) handleRenameNote(w http.ResponseWriter, r *http.Request){
	title := mux.Vars(r)["_t"]
	id := mux.Vars(r)["_id"]

	if title == "" {
		http.Error(w, "missing title", http.StatusBadRequest)
		return
	}

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	_ = h.db.Collection(notesCollection).FindOneAndUpdate(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"_t": title}})
}

// CORS middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow frontend
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed with the request
		next.ServeHTTP(w, r)
	})
}

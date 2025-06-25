package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AdityaHarindar/Richpanel-Assignment/model"
	"github.com/AdityaHarindar/Richpanel-Assignment/store"

	"github.com/gorilla/mux"
)

func GetAllHandler(s store.Store, c store.Cache) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		page, _ := strconv.Atoi(request.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))
		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}
		offset := (page - 1) * limit

		// Set JSON header
		writer.Header().Set("Content-Type", "application/json")

		// Attempt cache read
		key := fmt.Sprintf("%d:%d", limit, offset)
		if cacheResponse, ok := c.Get(key); ok {
			writer.Header().Set("X-CACHE", "HIT")
			writer.Write(cacheResponse)
		}

		// Cache miss, ask datastore
		posts, err := s.GetAll(limit, offset)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		} else {
			response, _ := json.Marshal(posts)
			c.Set(key, response)
			writer.Header().Set("X-CACHE", "MISS")
			writer.Write(response)
		}
	}
}

func GetHandler(s store.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Get ID from request URL
		vars := mux.Vars(request)
		key := vars["key"]
		id, e := strconv.Atoi(key)
		if e != nil {
			writer.WriteHeader(http.StatusNotFound)
		}

		// Get post by Id
		post := s.GetByID(id)

		// Encode and Write response
		if post.ID == 0 {
			http.Error(writer, "Post not found!", http.StatusBadRequest)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(writer).Encode(post)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func PostHandler(s store.Store, c store.Cache) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var post model.Post
		err := json.NewDecoder(request.Body).Decode(&post)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		defer request.Body.Close()

		id, err := s.Create(post)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		c.InvalidateAll()

		// Encode and Write response
		post = s.GetByID(id)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(writer).Encode(post)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func PutHandler(s store.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var post model.Post
		err := json.NewDecoder(request.Body).Decode(&post)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		defer request.Body.Close()

		// Get ID from request URL
		vars := mux.Vars(request)
		key := vars["key"]
		id, e := strconv.Atoi(key)
		if e != nil {
			writer.WriteHeader(http.StatusNotFound)
		}

		// post is shadowed
		post, err = s.Update(id, post)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		// Encode and Write response
		post = s.GetByID(id)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(writer).Encode(post)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func DeleteHandler(s store.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Get ID from request URL
		vars := mux.Vars(request)
		key := vars["key"]
		id, e := strconv.Atoi(key)
		if e != nil {
			writer.WriteHeader(http.StatusNotFound)
		}

		if s.Delete(id) {
			//	deleted successfully
			writer.WriteHeader(http.StatusAccepted)
		} else {
			http.Error(writer, "Unable to delete/Not Found", http.StatusNoContent)
		}
	}
}

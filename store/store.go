package store

import (
	"errors"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/AdityaHarindar/Richpanel-Assignment/model"
)

// Store is an interface for the in memory data store, which can be replaced by other data stores/DBs
type Store interface {
	Create(post model.Post) (int, error)
	GetByID(id int) model.Post
	GetAll(limit, offset int) ([]model.Post, error)
	Update(id int, post model.Post) (model.Post, error)
	Delete(id int) bool
}

type DataStore struct {
	ds map[int]model.Post
	mu sync.RWMutex
}

// NewStore returns a pointer to a new in-memory data store
func NewStore() *DataStore {
	ds := make(map[int]model.Post)
	return &DataStore{ds: ds, mu: sync.RWMutex{}}
}

// Create adds a post to the data store and returns post id, error
func (ds *DataStore) Create(p model.Post) (int, error) {
	// Validations
	if len(p.Title) < 1 {
		return 0, errors.New("bad ISBN")
	}
	if len(p.Author) < 1 {
		return 0, errors.New("bad Author")
	}

	// Lock the shared data store
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// upto 1000 posts, try new random again if ID exists in the store already
	id := rand.IntN(1000)
	_, exists := ds.ds[id]
	for {
		if !exists {
			break
		} else {
			id = rand.IntN(1000)
		}
	}

	ds.ds[id] = model.Post{
		ID:        id,
		Title:     p.Title,
		Author:    p.Author,
		Body:      p.Body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return id, nil

}

// GetByID accepts a post ID and returns the post contents
func (ds *DataStore) GetByID(id int) model.Post {
	// checks
	if id < 1 {
		return model.Post{}
	}

	// RLock the shared data store
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Check if post exists
	post, exists := ds.ds[id]
	if !exists {
		return model.Post{}
	}

	return post
}

// GetAll accepts a limit/offset, returns a list of posts
func (ds *DataStore) GetAll(limit, offset int) ([]model.Post, error) {
	// RLock the shared data store
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if len(ds.ds) < 1 {
		return nil, errors.New("no posts found")
	}

	var posts []model.Post
	for _, post := range ds.ds {
		posts = append(posts, post)
	}

	numPosts := len(posts)
	if offset >= numPosts {
		return nil, nil
	}

	end := offset + limit
	if end > numPosts {
		end = numPosts
	}

	return posts[offset:end], nil
}

// Update accepts post ID, updated post contents, and returns the updated post, error
func (ds *DataStore) Update(id int, p model.Post) (model.Post, error) {
	// Validations
	if len(p.Title) < 1 {
		return model.Post{}, errors.New("bad ISBN")
	}
	if len(p.Author) < 1 {
		return model.Post{}, errors.New("bad Author")
	}

	// Lock the shared data store
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// Update the post if found
	if post, exists := ds.ds[id]; !exists {
		return model.Post{}, errors.New("post not found")
	} else {
		p.ID = post.ID
		p.CreatedAt = post.CreatedAt
		p.UpdatedAt = time.Now()

		ds.ds[id] = p
		return post, nil
	}
}

// Delete accepts a post ID and returns true if the deletion is successful
func (ds *DataStore) Delete(id int) bool {
	_, exists := ds.ds[id]
	if !exists {
		return false
	} else {
		// Lock the shared data store
		ds.mu.Lock()
		defer ds.mu.Unlock()

		//	delete
		delete(ds.ds, id)
	}
	return true
}

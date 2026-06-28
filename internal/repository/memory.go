package repository

import (
	"errors"
	"sync"

	"mini_url/internal/models"
)

type MemoryLinkRepository struct {
	links   map[int]*models.Link
	mu      sync.RWMutex
	counter int
}

func NewMemoryLinkRepository() *MemoryLinkRepository {
	return &MemoryLinkRepository{
		links:   make(map[int]*models.Link),
		counter: 1,
	}
}

func (r *MemoryLinkRepository) Create(link *models.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	link.ID = r.counter
	r.links[link.ID] = link
	r.counter++
	return nil
}

func (r *MemoryLinkRepository) FindByID(id int) (*models.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	link, exists := r.links[id]
	if !exists {
		return nil, errors.New("link not found")
	}
	return link, nil
}

func (r *MemoryLinkRepository) FindByShortID(shortID string) (*models.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, link := range r.links {
		if link.ShortID == shortID {
			return link, nil
		}
	}
	return nil, errors.New("link not found")
}

func (r *MemoryLinkRepository) FindByUserID(userID int) ([]*models.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var links []*models.Link
	for _, link := range r.links {
		if link.UserID == userID {
			links = append(links, link)
		}
	}
	return links, nil
}

func (r *MemoryLinkRepository) Update(link *models.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.links[link.ID] = link
	return nil
}

func (r *MemoryLinkRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.links, id)
	return nil
}
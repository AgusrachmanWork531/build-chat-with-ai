package user

import (
	"context"
	"fmt"
	"sync"
)

// InMemoryUserRepository is an in-memory implementation of the UserRepository.
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*User
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*User),
	}
}

// Create saves a new user.
func (r *InMemoryUserRepository) Create(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("user with id %s already exists", user.ID)
	}
	
	r.users[user.ID] = user
	return nil
}

// GetByID retrieves a user by their ID.
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with id %s not found", id)
	}
	return user, nil
}

// GetByEmail retrieves a user by their email.
func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		fmt.Printf("Looking for user with email: %s\n", user.Email) // Debugging line
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with email %s not found", email)
}

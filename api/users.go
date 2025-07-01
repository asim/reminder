package api

import (
	"encoding/json"
	"os"
	"sync"
)

type User struct {
	Email string `json:"email"`
}

type Users struct {
	mu    sync.RWMutex
	Users map[string]User `json:"users"`
}

var usersFile = ReminderPath("users.json")
var users = &Users{Users: make(map[string]User)}

func LoadUsers() error {
	f, err := os.Open(usersFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&users.Users)
}

func SaveUsers() error {
	users.mu.RLock()
	defer users.mu.RUnlock()
	_ = os.MkdirAll(ReminderDir, 0700)
	f, err := os.Create(usersFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(users.Users)
}

func AddUser(email string) error {
	users.mu.Lock()
	defer users.mu.Unlock()
	users.Users[email] = User{Email: email}
	return SaveUsers()
}

func RemoveUser(email string) error {
	users.mu.Lock()
	defer users.mu.Unlock()
	delete(users.Users, email)
	return SaveUsers()
}

func ListUsers() []User {
	users.mu.RLock()
	defer users.mu.RUnlock()
	result := make([]User, 0, len(users.Users))
	for _, u := range users.Users {
		result = append(result, u)
	}
	return result
}

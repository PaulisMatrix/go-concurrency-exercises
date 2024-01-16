//////////////////////////////////////////////////////////////////////
//
// Given is a SessionManager that stores session information in
// memory. The SessionManager itself is working, however, since we
// keep on adding new sessions to the manager our program will
// eventually run out of memory.
//
// Your task is to implement a session cleaner routine that runs
// concurrently in the background and cleans every session that
// hasn't been updated for more than 5 seconds (of course usually
// session times are much longer).
//
// Note that we expect the session to be removed anytime between 5 and
// 7 seconds after the last update. Also, note that you have to be
// very careful in order to prevent race conditions.
//

package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

const EXPIRY_SECS = 5

// SessionManager keeps track of all sessions from creation, updating
// to destroying.
type SessionManager struct {
	sessions map[string]Session
	mu       sync.Mutex
	done     chan bool
}

// Session stores the session's data
type Session struct {
	Data      map[string]interface{}
	ExpiresAt time.Time
}

func (m *SessionManager) GarbageCollector() {
	fmt.Println("running the garbage collector!")
	ticker := time.NewTicker(1 * time.Second)

	// break the for loop explicitly after 20secs to avoid ticker leaking
	// https://pkg.go.dev/time#Tick
	timeout := time.After(20 * time.Second)

	for {
		select {
		case <-timeout:
			fmt.Println("shutting down the garbage collector!")
			ticker.Stop()
			break
		case <-ticker.C:
			m.mu.Lock()

			for sID, sData := range m.sessions {
				if time.Now().After(sData.ExpiresAt) {
					delete(m.sessions, sID)
				}
			}

			m.mu.Unlock()
		}
	}
}

// NewSessionManager creates a new sessionManager
func NewSessionManager() *SessionManager {
	m := &SessionManager{
		sessions: make(map[string]Session),
	}

	// kick off the garbage collector here?
	go m.GarbageCollector()

	return m
}

// CreateSession creates a new session and returns the sessionID
func (m *SessionManager) CreateSession() (string, error) {
	sessionID, err := MakeSessionID()
	if err != nil {
		return "", err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[sessionID] = Session{
		Data:      make(map[string]interface{}),
		ExpiresAt: time.Now().Add(time.Second * EXPIRY_SECS),
	}

	return sessionID, nil
}

// ErrSessionNotFound returned when sessionID not listed in
// SessionManager
var ErrSessionNotFound = errors.New("SessionID does not exists")

// GetSessionData returns data related to session if sessionID is
// found, errors otherwise
func (m *SessionManager) GetSessionData(sessionID string) (map[string]interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session.Data, nil
}

// UpdateSessionData overwrites the old session data with the new one
func (m *SessionManager) UpdateSessionData(sessionID string, data map[string]interface{}) error {
	_, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Hint: you should renew expiry of the session here
	m.sessions[sessionID] = Session{
		Data:      data,
		ExpiresAt: time.Now().Add(time.Second * EXPIRY_SECS),
	}

	return nil
}

func main() {
	// Create new sessionManager and new session
	m := NewSessionManager()
	sID, err := m.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created new session with ID", sID)

	// Update session data
	data := make(map[string]interface{})
	data["website"] = "longhoang.de"

	err = m.UpdateSessionData(sID, data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Update session data, set website to longhoang.de")

	// Retrieve data from manager again
	updatedData, err := m.GetSessionData(sID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get session data:", updatedData)
}

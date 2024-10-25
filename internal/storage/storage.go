package storage

import (
	"sync"

	"github.com/courage173/quiz-api/internal/models"
)

// Storage holds all in-memory data
type Storage struct {
    Questions []models.Question
    Users     []models.UserSubmission
    Mutex     sync.RWMutex
}

// UserSubmission represents a user's submitted answers
type UserSubmission struct {
    UserName string
    Correct  int
}

// NewStorage initializes the storage with sample questions
func NewStorage() *Storage {
    return &Storage{
        Questions: []models.Question{
            {
                ID:   1,
                Text: "What is the capital of France?",
                Options: []models.Option{
                    {ID: 1, Text: "Berlin", IsCorrect: false},
                    {ID: 2, Text: "London", IsCorrect: false},
                    {ID: 3, Text: "Paris", IsCorrect: true},
                    {ID: 4, Text: "Madrid", IsCorrect: false},
                },
            },
            {
                ID:   2,
                Text: "Which planet is known as the Red Planet?",
                Options: []models.Option{
                    {ID: 1, Text: "Earth", IsCorrect: false},
                    {ID: 2, Text: "Mars", IsCorrect: true},
                    {ID: 3, Text: "Jupiter", IsCorrect: false},
                    {ID: 4, Text: "Saturn", IsCorrect: false},
                },
            },
            // Add more questions as needed
        },
        Users: make([]models.UserSubmission, 0),
    }
}

// AddUserSubmission adds a user's submission to storage
func (s *Storage) AddUserSubmission(submission models.UserSubmission) {
    s.Mutex.Lock()
    defer s.Mutex.Unlock()
    s.Users = append(s.Users, submission)
}

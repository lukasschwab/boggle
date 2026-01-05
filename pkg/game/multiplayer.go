package game

import (
	"crypto/rand"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/google/uuid"
	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/dictionary"
)

// Player represents a player in a multiplayer game
type Player struct {
	ID          string // Unique identifier (session ID or generated UUID)
	Username    string // Display name only
	Session     ssh.Session
	Program     *tea.Program
	State       PlayerState
	ScoredWords map[string]bool
	Connected   bool
}

type PlayerState int

const (
	PlayerReady PlayerState = iota
	PlayerPlaying
	PlayerConceded
)

// Room represents a multiplayer game room
type Room struct {
	ID            string
	Players       map[string]*Player // Key is player ID, not username
	Board         boggle.Board
	Dict          dictionary.Map
	Duration      time.Duration
	State         RoomState
	GameStartTime time.Time
	mutex         sync.RWMutex
}

type RoomState int

const (
	RoomWaitingForPlayers RoomState = iota
	RoomInGame
	RoomShowingResults
)

// RoomManager manages all active rooms
type RoomManager struct {
	rooms map[string]*Room
	mutex sync.RWMutex
}

var GlobalRoomManager = &RoomManager{
	rooms: make(map[string]*Room),
}

// CreateRoom creates a new room with a random UUID
func (rm *RoomManager) CreateRoom() *Room {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	roomID := uuid.New().String()
	board := boggle.Shake(boggle.ClassicDice)
	duration := 3 * time.Minute

	baseDict := dictionary.Filtered{
		Underlying: dictionary.EmptyTrie(),
		Filter:     boggle.Playable,
	}
	if err := dictionary.Load(dictionary.CSW19G, baseDict); err != nil {
		return nil
	}
	boardDict := board.AllWords(baseDict)

	room := &Room{
		ID:       roomID,
		Players:  make(map[string]*Player),
		Board:    board,
		Dict:     boardDict,
		Duration: duration,
		State:    RoomInGame,
	}

	rm.rooms[roomID] = room
	return room
}

// GetRoom returns a room by ID
func (rm *RoomManager) GetRoom(roomID string) (*Room, bool) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	room, exists := rm.rooms[roomID]
	return room, exists
}

// CleanupRoom removes a room if it has no players
func (rm *RoomManager) CleanupRoom(roomID string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return
	}

	room.mutex.RLock()
	connectedCount := 0
	for _, player := range room.Players {
		if player.Connected {
			connectedCount++
		}
	}
	room.mutex.RUnlock()

	if connectedCount == 0 {
		delete(rm.rooms, roomID)
	}
}

// AddPlayer adds a player to the room
func (r *Room) AddPlayer(playerID, username string, session ssh.Session, program *tea.Program) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	player := &Player{
		ID:          playerID,
		Username:    username,
		Session:     session,
		Program:     program,
		State:       PlayerReady,
		ScoredWords: make(map[string]bool),
		Connected:   true,
	}

	r.Players[playerID] = player
}

// RemovePlayer removes a player from the room
func (r *Room) RemovePlayer(playerID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if player, exists := r.Players[playerID]; exists {
		player.Connected = false
	}
}

// GetPlayerCount returns the number of connected players
func (r *Room) GetPlayerCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, player := range r.Players {
		if player.Connected {
			count++
		}
	}
	return count
}

// AllPlayersReady checks if all connected players are ready
func (r *Room) AllPlayersReady() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(r.Players) == 0 {
		return false
	}

	for _, player := range r.Players {
		if player.Connected && player.State != PlayerReady {
			return false
		}
	}
	return true
}

// AllPlayersFinished checks if all connected players have finished playing
func (r *Room) AllPlayersFinished() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, player := range r.Players {
		if player.Connected && player.State == PlayerPlaying {
			return false
		}
	}
	return true
}

// StartGame transitions the room to game state
func (r *Room) StartGame() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.State = RoomInGame
	r.GameStartTime = time.Now()

	// Set all players to playing state
	for _, player := range r.Players {
		if player.Connected {
			player.State = PlayerPlaying
			player.ScoredWords = make(map[string]bool)
		}
	}
}

// CalculateResults calculates multiplayer game results
func (r *Room) CalculateResults() MultiplayerResults {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	allWords := r.Dict.Members()
	wordCounts := make(map[string][]string) // word -> list of usernames who found it
	userWords := make(map[string][]string)  // username -> list of words they found

	// Collect all words found by players
	for _, player := range r.Players {
		if !player.Connected {
			continue
		}

		username := player.Username
		userWords[username] = make([]string, 0)
		for word := range player.ScoredWords {
			userWords[username] = append(userWords[username], word)
			wordCounts[word] = append(wordCounts[word], username)
		}
	}

	// Categorize words
	var sharedWords []WordResult
	uniqueWords := make(map[string][]string) // username -> unique words
	var missedWords []string

	// Find shared words (found by multiple players)
	for word, users := range wordCounts {
		if len(users) > 1 {
			sharedWords = append(sharedWords, WordResult{
				Word:  word,
				Users: users,
			})
		}
	}

	// Find unique words for each user
	for username := range userWords {
		uniqueWords[username] = make([]string, 0)
	}
	for word, users := range wordCounts {
		if len(users) == 1 {
			uniqueWords[users[0]] = append(uniqueWords[users[0]], word)
		}
	}

	// Find missed words
	wordFoundMap := make(map[string]bool)
	for word := range wordCounts {
		wordFoundMap[word] = true
	}
	for _, word := range allWords {
		if !wordFoundMap[word] {
			missedWords = append(missedWords, word)
		}
	}

	return MultiplayerResults{
		SharedWords: sharedWords,
		UniqueWords: uniqueWords,
		MissedWords: missedWords,
		UserWords:   userWords,
	}
}

// WordResult represents a word found by multiple players
type WordResult struct {
	Word  string
	Users []string
}

// MultiplayerResults contains the results of a multiplayer game
type MultiplayerResults struct {
	SharedWords []WordResult        // Words found by multiple players
	UniqueWords map[string][]string // Username -> unique words
	MissedWords []string            // Words nobody found
	UserWords   map[string][]string // Username -> all words found
}

// GenerateVegetableName generates a random vegetable name for anonymous users
func GenerateVegetableName() string {
	vegetables := []string{
		"Carrot", "Broccoli", "Spinach", "Tomato", "Cucumber",
		"Pepper", "Onion", "Garlic", "Potato", "Lettuce",
		"Cabbage", "Radish", "Celery", "Asparagus", "Zucchini",
	}

	b := make([]byte, 1)
	rand.Read(b)
	index := int(b[0]) % len(vegetables)
	return vegetables[index]
}

// GetPlayerID generates a unique player ID from the SSH session
func GetPlayerID(session ssh.Session) string {
	// Try to use a unique session identifier
	if session.Context() != nil {
		if sessionID := session.Context().Value("session_id"); sessionID != nil {
			if id, ok := sessionID.(string); ok && id != "" {
				return id
			}
		}
	}

	// Fallback: generate a UUID for this session
	return uuid.New().String()
}

// GetUsername extracts username from SSH session or generates a vegetable name
func GetUsername(session ssh.Session) string {
	if session.User() != "" && session.User() != "root" {
		return session.User()
	}
	return GenerateVegetableName()
}

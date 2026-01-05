package game

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/dictionary"
)

// MultiplayerGameModel extends the existing game model for multiplayer support
type MultiplayerGameModel struct {
	// Embed the existing model to reuse all functionality
	model

	// Multiplayer-specific fields
	room           *Room
	playerID       string
	username       string
	playerState    PlayerState
	conceding      bool
	showingResults bool
	results        MultiplayerResults
}

// NewMultiplayerGameModel creates a new multiplayer game model
func NewMultiplayerGameModel(
	room *Room,
	playerID string,
	username string,
	style Style,
) *MultiplayerGameModel {
	// Create the base game model using the existing constructor
	baseModel := Model(room.Dict, room.Board, room.Duration, style).(model)

	// Modify the keymap for multiplayer (ESC now concedes instead of quitting)
	baseModel.keymap.quit = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("[esc]", "concede"),
	)

	return &MultiplayerGameModel{
		model:          baseModel,
		room:           room,
		playerID:       playerID,
		username:       username,
		playerState:    PlayerPlaying,
		conceding:      false,
		showingResults: false,
	}
}

// Init initializes the multiplayer game model
func (m *MultiplayerGameModel) Init() tea.Cmd {
	return m.model.Init()
}

// Update handles updates for the multiplayer game model
func (m *MultiplayerGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.showingResults {
		return m.updateResults(msg)
	}

	switch msg := msg.(type) {
	case timer.TimeoutMsg:
		// Game time is up, mark player as finished
		return m.handleGameEnd()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit): // ESC key now concedes
			return m.handleConcede()
		case key.Matches(msg, m.keymap.submit):
			// Handle word submission - reuse existing logic
			updatedModel, cmd := m.model.Update(msg)
			m.model = updatedModel.(model)

			// Update player's scored words in the room
			m.updatePlayerWords()

			return m, cmd
		}
	}

	// Pass other messages to the base model
	updatedModel, cmd := m.model.Update(msg)
	m.model = updatedModel.(model)
	return m, cmd
}

func (m *MultiplayerGameModel) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r", "enter": // Mark ready for next game
			return m.handleReady()
		case "q", "ctrl+c": // Quit the room
			return m.handleQuitRoom()
		}
	}
	return m, nil
}

func (m *MultiplayerGameModel) handleConcede() (tea.Model, tea.Cmd) {
	m.conceding = true
	m.playerState = PlayerConceded

	// Update player state in room
	if player, exists := m.room.Players[m.playerID]; exists {
		player.State = PlayerConceded
		player.ScoredWords = make(map[string]bool)
		for word := range m.scoredWords {
			player.ScoredWords[word] = true
		}
	}

	// Check if all players are done
	if m.room.AllPlayersFinished() {
		return m.showResults()
	}

	return m, nil
}

func (m *MultiplayerGameModel) handleGameEnd() (tea.Model, tea.Cmd) {
	m.playerState = PlayerConceded // Time up = conceded

	// Update player state in room
	if player, exists := m.room.Players[m.playerID]; exists {
		player.State = PlayerConceded
		player.ScoredWords = make(map[string]bool)
		for word := range m.scoredWords {
			player.ScoredWords[word] = true
		}
	}

	// Check if all players are done
	if m.room.AllPlayersFinished() {
		return m.showResults()
	}

	return m, nil
}

func (m *MultiplayerGameModel) showResults() (tea.Model, tea.Cmd) {
	m.showingResults = true
	m.results = m.room.CalculateResults()
	m.room.State = RoomShowingResults
	return m, nil
}

func (m *MultiplayerGameModel) handleReady() (tea.Model, tea.Cmd) {
	// Mark player as ready for next game
	if player, exists := m.room.Players[m.playerID]; exists {
		player.State = PlayerReady
	}

	// Check if all players are ready to start next game
	if m.room.AllPlayersReady() {
		// Start new game
		m.room.StartGame()
		m.showingResults = false
		m.playerState = PlayerPlaying
		m.conceding = false

		// Reset the game state
		m.scoredWords = make(map[string]bool)
		m.history = history{}
		m.userInput.Reset()
		m.timer = timer.New(m.room.Duration)

		// Generate new board
		m.room.Board = boggle.Shake(boggle.ClassicDice)
		m.Board = m.room.Board

		// Recalculate dictionary for new board
		baseDict := dictionary.Filtered{
			Underlying: dictionary.EmptyTrie(),
			Filter:     boggle.Playable,
		}
		dictionary.Load(dictionary.CSW19G, baseDict)
		boardDict := m.room.Board.AllWords(baseDict)
		m.room.Dict = boardDict
		m.Dict = boardDict
		m.totalWordCount = len(boardDict.Members())

		return m, m.timer.Init()
	}

	return m, nil
}

func (m *MultiplayerGameModel) handleQuitRoom() (tea.Model, tea.Cmd) {
	// Remove player from room
	m.room.RemovePlayer(m.playerID)

	// Clean up room if no players left
	GlobalRoomManager.CleanupRoom(m.room.ID)

	return m, tea.Quit
}

func (m *MultiplayerGameModel) updatePlayerWords() {
	if player, exists := m.room.Players[m.playerID]; exists {
		player.ScoredWords = make(map[string]bool)
		for word := range m.scoredWords {
			player.ScoredWords[word] = true
		}
	}
}

// View renders the multiplayer game view
func (m *MultiplayerGameModel) View() string {
	if m.showingResults {
		return m.viewResults()
	}

	// Reuse the existing game view but add multiplayer info
	baseView := m.model.View()

	// Add room info header
	roomInfo := fmt.Sprintf("Room: %s • Players: %d • Username: %s",
		m.room.ID, m.room.GetPlayerCount(), m.username)
	roomInfoStyled := m.style.blurredStyle.Render(roomInfo)

	// Add player state info
	var stateInfo string
	switch m.playerState {
	case PlayerPlaying:
		if m.conceding {
			stateInfo = m.style.blurredStyle.Render("You have conceded. Waiting for other players...")
		} else {
			stateInfo = ""
		}
	case PlayerConceded:
		stateInfo = m.style.blurredStyle.Render("You have conceded. Waiting for other players...")
	case PlayerReady:
		stateInfo = m.style.blurredStyle.Render("Ready for next game!")
	}

	result := roomInfoStyled + "\n" + baseView
	if stateInfo != "" {
		result += "\n" + stateInfo
	}

	return result
}

func (m *MultiplayerGameModel) viewResults() string {
	var builder strings.Builder

	// Header
	header := fmt.Sprintf("Game Results - Room: %s", m.room.ID)
	builder.WriteString(m.style.boardStyle.Render(header) + "\n\n")

	playerCount := m.room.GetPlayerCount()
	builder.WriteString(fmt.Sprintf("Players: %d\n\n", playerCount))

	// Shared words (found by multiple players)
	if len(m.results.SharedWords) > 0 {
		builder.WriteString(m.style.promptStyle.Render("Shared Words:") + "\n")
		for _, result := range m.results.SharedWords {
			users := strings.Join(result.Users, ", ")
			builder.WriteString(fmt.Sprintf("  %s (%s)\n", result.Word, users))
		}
		builder.WriteString("\n")
	}

	// Unique words for each player
	builder.WriteString(m.style.promptStyle.Render("Unique Words:") + "\n")
	for username, words := range m.results.UniqueWords {
		if len(words) > 0 {
			builder.WriteString(fmt.Sprintf("  %s: %s\n", username, strings.Join(words, ", ")))
		}
	}
	builder.WriteString("\n")

	// Missed words (show first 10 to avoid overwhelming)
	if len(m.results.MissedWords) > 0 {
		missedToShow := m.results.MissedWords
		if len(missedToShow) > 10 {
			missedToShow = missedToShow[:10]
		}
		builder.WriteString(m.style.blurredStyle.Render("Some missed words:") + "\n")
		builder.WriteString(m.style.blurredStyle.Render(strings.Join(missedToShow, ", ")) + "\n")
		if len(m.results.MissedWords) > 10 {
			builder.WriteString(m.style.blurredStyle.Render(fmt.Sprintf("... and %d more\n", len(m.results.MissedWords)-10)))
		}
		builder.WriteString("\n")
	}

	// Instructions
	builder.WriteString("Press [R] or [ENTER] to ready up for next game\n")
	builder.WriteString("Press [Q] or [CTRL+C] to leave room\n")

	return builder.String()
}

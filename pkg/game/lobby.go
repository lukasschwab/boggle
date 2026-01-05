package game

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

// LobbyState represents the current state of the lobby
type LobbyState int

const (
	LobbyMenu LobbyState = iota
	LobbyJoinRoom
	LobbyInRoom
	LobbyInGame
)

// LobbyModel represents the lobby interface
type LobbyModel struct {
	state        LobbyState
	playerID     string
	username     string
	session      ssh.Session
	style        Style
	selectedItem int
	textInput    textinput.Model
	help         help.Model
	keymap       lobbyKeymap
	currentRoom  *Room
	gameModel    tea.Model
	errorMessage string
}

type lobbyKeymap struct {
	up     key.Binding
	down   key.Binding
	enter  key.Binding
	quit   key.Binding
	back   key.Binding
	submit key.Binding
}

// NewLobbyModel creates a new lobby model
func NewLobbyModel(session ssh.Session, style Style) *LobbyModel {
	playerID := GetPlayerID(session)
	username := GetUsername(session)

	ti := textinput.New()
	ti.Placeholder = "Enter room UUID..."
	ti.Focus()
	ti.PromptStyle = style.promptStyle

	return &LobbyModel{
		state:        LobbyMenu,
		playerID:     playerID,
		username:     username,
		session:      session,
		style:        style,
		selectedItem: 0,
		textInput:    ti,
		help:         help.New(),
		keymap: lobbyKeymap{
			up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "up"),
			),
			down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "down"),
			),
			enter: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("ctrl+c/q", "quit"),
			),
			back: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "back"),
			),
			submit: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "submit"),
			),
		},
	}
}

func (m *LobbyModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *LobbyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case LobbyMenu:
		return m.updateMenu(msg)
	case LobbyJoinRoom:
		return m.updateJoinRoom(msg)
	case LobbyInRoom:
		return m.updateInRoom(msg)
	case LobbyInGame:
		return m.updateInGame(msg)
	}
	return m, nil
}

func (m *LobbyModel) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.up):
			if m.selectedItem > 0 {
				m.selectedItem--
			}
		case key.Matches(msg, m.keymap.down):
			if m.selectedItem < 2 {
				m.selectedItem++
			}
		case key.Matches(msg, m.keymap.enter):
			return m.handleMenuSelection()
		}
	}
	return m, nil
}

func (m *LobbyModel) updateJoinRoom(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.back):
			m.state = LobbyMenu
			m.errorMessage = ""
			m.textInput.Reset()
			return m, nil
		case key.Matches(msg, m.keymap.submit):
			return m.handleJoinRoom()
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *LobbyModel) updateInRoom(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.currentRoom == nil {
		m.state = LobbyMenu
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			// Leave room and return to menu
			m.currentRoom.RemovePlayer(m.playerID)
			GlobalRoomManager.CleanupRoom(m.currentRoom.ID)
			m.currentRoom = nil
			m.state = LobbyMenu
			return m, nil
		case key.Matches(msg, m.keymap.enter):
			// Handle ready/start game
			if m.currentRoom.State == RoomWaitingForPlayers {
				// Mark player as ready
				if player, exists := m.currentRoom.Players[m.playerID]; exists {
					player.State = PlayerReady
				}
				// Check if all players are ready and start game
				if m.currentRoom.AllPlayersReady() {
					m.currentRoom.StartGame()
					// Transition to game
					m.state = LobbyInGame
					m.gameModel = NewMultiplayerGameModel(m.currentRoom, m.playerID, m.username, m.style)
					return m, m.gameModel.Init()
				}
			}
		}
	}

	return m, nil
}

func (m *LobbyModel) updateInGame(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.gameModel, cmd = m.gameModel.Update(msg)

	// Check if game session is finished (player quit room)
	if multiGameModel, ok := m.gameModel.(*MultiplayerGameModel); ok {
		// TODO: Handle different exit conditions from game
		// For now, we'll assume any quit goes back to lobby
		_ = multiGameModel // Placeholder to avoid unused variable
	}

	return m, cmd
}

func (m *LobbyModel) handleMenuSelection() (tea.Model, tea.Cmd) {
	switch m.selectedItem {
	case 0: // Join Room
		m.state = LobbyJoinRoom
		m.textInput.Focus()
		m.errorMessage = ""
	case 1: // Create Room
		return m.createRoom()
	}
	return m, nil
}

func (m *LobbyModel) createRoom() (tea.Model, tea.Cmd) {
	room := GlobalRoomManager.CreateRoom()
	if room == nil {
		m.errorMessage = "Failed to create room"
		return m, nil
	}

	// Add this player to the room (we'll pass nil for program for now)
	room.AddPlayer(m.playerID, m.username, m.session, nil)
	m.currentRoom = room
	m.state = LobbyInRoom

	return m, nil
}

func (m *LobbyModel) handleJoinRoom() (tea.Model, tea.Cmd) {
	roomID := strings.TrimSpace(m.textInput.Value())
	if roomID == "" {
		m.errorMessage = "Please enter a room UUID"
		return m, nil
	}

	room, exists := GlobalRoomManager.GetRoom(roomID)
	if !exists {
		m.errorMessage = "Room not found"
		return m, nil
	}

	// Add this player to the room (we'll pass nil for program for now)
	room.AddPlayer(m.playerID, m.username, m.session, nil)
	m.currentRoom = room
	m.state = LobbyInRoom
	m.textInput.Reset()
	m.errorMessage = ""

	return m, nil
}

func (m *LobbyModel) View() string {
	switch m.state {
	case LobbyMenu:
		return m.viewMenu()
	case LobbyJoinRoom:
		return m.viewJoinRoom()
	case LobbyInRoom:
		return m.viewInRoom()
	case LobbyInGame:
		return m.gameModel.View()
	}
	return "Unknown state"
}

func (m *LobbyModel) viewMenu() string {
	var builder strings.Builder

	// Title
	title := m.style.boardStyle.Render("SSH Boggle Multiplayer")
	builder.WriteString(title + "\n\n")

	// Username display
	usernameDisplay := fmt.Sprintf("Welcome, %s!", m.username)
	builder.WriteString(m.style.blurredStyle.Render(usernameDisplay) + "\n\n")

	// Menu options
	menuItems := []string{
		"Join room",
		"Create room",
	}

	for i, item := range menuItems {
		cursor := "  "
		if i == m.selectedItem {
			cursor = "> "
			item = m.style.promptStyle.Render(item)
		} else {
			item = m.style.blurredStyle.Render(item)
		}
		builder.WriteString(cursor + item + "\n")
	}

	builder.WriteString("\n")
	builder.WriteString(m.help.ShortHelpView([]key.Binding{
		m.keymap.up,
		m.keymap.down,
		m.keymap.enter,
		m.keymap.quit,
	}))

	return builder.String()
}

func (m *LobbyModel) viewJoinRoom() string {
	var builder strings.Builder

	title := m.style.boardStyle.Render("Join Room")
	builder.WriteString(title + "\n\n")

	builder.WriteString("Enter the room UUID to join:\n\n")
	builder.WriteString(m.textInput.View() + "\n\n")

	if m.errorMessage != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
		builder.WriteString(errorStyle.Render(m.errorMessage) + "\n\n")
	}

	builder.WriteString(m.help.ShortHelpView([]key.Binding{
		m.keymap.submit,
		m.keymap.back,
		m.keymap.quit,
	}))

	return builder.String()
}

func (m *LobbyModel) viewInRoom() string {
	if m.currentRoom == nil {
		return "Error: No current room"
	}

	var builder strings.Builder

	// Room info
	title := fmt.Sprintf("Room: %s", m.currentRoom.ID)
	builder.WriteString(m.style.boardStyle.Render(title) + "\n\n")

	playerCount := m.currentRoom.GetPlayerCount()
	builder.WriteString(fmt.Sprintf("Players in room: %d\n", playerCount))

	// List connected players
	m.currentRoom.mutex.RLock()
	builder.WriteString("Connected players:\n")
	for _, player := range m.currentRoom.Players {
		if player.Connected {
			status := ""
			switch player.State {
			case PlayerReady:
				status = " (ready)"
			case PlayerPlaying:
				status = " (playing)"
			case PlayerConceded:
				status = " (finished)"
			}
			builder.WriteString(fmt.Sprintf("  - %s%s\n", player.Username, status))
		}
	}
	m.currentRoom.mutex.RUnlock()
	builder.WriteString("\n")

	// Room state display
	switch m.currentRoom.State {
	case RoomWaitingForPlayers:
		builder.WriteString("Waiting for all players to be ready...\n")
		builder.WriteString("Press ENTER when ready to start the game.\n")
		builder.WriteString("Press Q to leave the room.\n")
	case RoomInGame:
		builder.WriteString("Game in progress...\n")
		builder.WriteString("You will join the next game when it starts.\n")
	case RoomShowingResults:
		builder.WriteString("Showing game results...\n")
		builder.WriteString("You will join the next game when it starts.\n")
	}

	return builder.String()
}

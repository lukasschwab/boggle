package game

import (
       tea "github.com/charmbracelet/bubbletea"
       "github.com/charmbracelet/ssh"
)

// AppModel is the top-level model that coordinates the entire application
type AppModel struct {
       lobby      *LobbyModel
       session    ssh.Session
       playerID   string
       username   string
       style      Style
       currentProgram *tea.Program // Reference to self for room management
}

// NewAppModel creates a new application model
func NewAppModel(session ssh.Session, style Style) *AppModel {
       playerID := GetPlayerID(session)
       username := GetUsername(session)

       return &AppModel{
               lobby:    NewLobbyModel(session, style),
               session:  session,
               playerID: playerID,
               username: username,
               style:    style,
       }
}

// SetProgram sets the tea.Program reference for room management
func (m *AppModel) SetProgram(program *tea.Program) {
       m.currentProgram = program
}

func (m *AppModel) Init() tea.Cmd {
       return m.lobby.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
       updatedLobby, cmd := m.lobby.Update(msg)
       m.lobby = updatedLobby.(*LobbyModel)
       return m, cmd
}

func (m *AppModel) View() string {
       return m.lobby.View()
}

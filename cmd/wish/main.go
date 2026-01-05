package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/lukasschwab/boggle/pkg/game"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

func main() {
	// Key is coupled to fly.toml.
	hostKeyPath, ok := os.LookupEnv("SSH_KEY_PATH")
	if !ok {
		// Default from wish demo app.
		hostKeyPath = ".ssh/id_ed25519"
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	// Create channels to differentiate shutdown reasons
	signalChan := make(chan os.Signal, 1)
	errorChan := make(chan error, 1)

	// Set up signal handling with detailed logging
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port, "pid", os.Getpid())

	// Start server in goroutine
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("SSH server error", "error", err)
			errorChan <- err
		} else {
			log.Debug("SSH server ListenAndServe completed normally")
			errorChan <- nil
		}
	}()

	// Wait for shutdown signal or server error
	var shutdownReason string
	select {
	case sig := <-signalChan:
		shutdownReason = "signal"
		log.Info("Received shutdown signal", "signal", sig.String(), "pid", os.Getpid())
	case err := <-errorChan:
		if err != nil {
			shutdownReason = "server-error"
			log.Error("Shutting down due to server error", "error", err)
		} else {
			shutdownReason = "server-completed"
			log.Info("Server completed normally")
		}
	}

	log.Info("Stopping SSH server", "reason", shutdownReason, "pid", os.Getpid())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	} else {
		log.Info("SSH server stopped successfully")
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	style := game.NewStyle(bubbletea.MakeRenderer(s))

       // Create the new multiplayer app model instead of single-player game
       return game.NewAppModel(s, style), []tea.ProgramOption{}
}

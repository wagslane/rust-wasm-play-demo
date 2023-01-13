package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wagslane/rust-wasm-play-demo/internal/workspacemanager"
)

const (
	numWorkspaces = 1

	// this is a directory that will be created
	// in your home directory and used to build rust code
	workspaceRoot = ".rustwasmplay"
)

func main() {
	workspacePath, err := getRustPath()
	if err != nil {
		log.Fatal(err)
	}

	wm := workspacemanager.New(numWorkspaces, workspacePath, createWorkspace)
	defer func() {
		err := wm.CleanupDisk()
		if err != nil {
			log.Printf("cleanupDisk err: %s", err)
		}
	}()

	cfg := config{
		WorkspaceManager: wm,
	}

	go cfg.startAPI("5000")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Print("Server Stopped")
}

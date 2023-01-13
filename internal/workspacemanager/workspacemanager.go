package workspacemanager

import (
	"log"
	"os"
)

// WorkspaceManager -
type WorkspaceManager struct {
	workingDirs   chan string
	workspacePath string
}

// CheckoutWorkingDir -
func (wm WorkspaceManager) CheckoutWorkingDir() string {
	return <-wm.workingDirs
}

// CheckinWorkingDir -
func (wm WorkspaceManager) CheckinWorkingDir(workingDir string) {
	wm.workingDirs <- workingDir
}

// CleanupDisk -
func (wm WorkspaceManager) CleanupDisk() error {
	log.Println("cleaning up workspace manager disk...")
	err := os.RemoveAll(wm.workspacePath)
	if err != nil {
		return err
	}
	log.Println("cleaned up workspace manager disk")
	return nil
}

// New -
func New(
	numWorkdirs int,
	workspacePath string,
	createWorkspace func(string) (string, error),
) WorkspaceManager {
	pwm := WorkspaceManager{
		workingDirs:   make(chan string, numWorkdirs),
		workspacePath: workspacePath,
	}

	for i := 0; i < numWorkdirs; i++ {
		log.Printf("creating workspace #%v", i)
		workingDir, err := createWorkspace(workspacePath)
		if err != nil {
			log.Printf("error in createWorkspace: %v", err)
			continue
		}
		pwm.workingDirs <- workingDir
		log.Printf("done with workspace #%v", i)
	}

	return pwm
}

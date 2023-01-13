package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	_ "embed"

	"github.com/google/uuid"
)

const projectName = "lib"

//go:embed assets/glue.rs
var wasmGlue string

//go:embed assets/Cargo.toml
var cargoToml string

func (cfg config) compileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Code string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters")
		return
	}

	workingDir := cfg.WorkspaceManager.CheckoutWorkingDir()
	defer cfg.WorkspaceManager.CheckinWorkingDir(workingDir)

	glueCode := addGlue(params.Code)

	rustFilePath := filepath.Join(workingDir, "src", projectName+".rs")
	err = writeFileToDisk(rustFilePath, glueCode)
	if err != nil {
		log.Printf("error in writeFileToDisk: %v, workingDir: %v", err, workingDir)
		respondWithError(w, 500, "Couldn't write rust to disk")
		return
	}

	err = runCmd(time.Second*120, workingDir, "cargo", "build", "--target", "wasm32-unknown-unknown", "--release")
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	dat, err := os.ReadFile(filepath.Join(workingDir, "target", "wasm32-unknown-unknown", "release", projectName+".wasm"))
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	w.Write(dat)
}

func writeFileToDisk(codePath, code string) error {
	// remove old code
	os.Remove(codePath)

	// create the new file
	f, err := os.Create(codePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// write the rest of the code
	dat := []byte(code)
	_, err = f.Write(dat)
	if err != nil {
		return err
	}
	return nil
}

func getRustPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, workspaceRoot), nil
}

func createWorkspace(workspacePath string) (string, error) {
	workingDir := filepath.Join(workspacePath, uuid.New().String())
	err := os.MkdirAll(workingDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	err = runCmd(time.Second*120, workingDir, "cargo", "new", "--lib", projectName)
	if err != nil {
		return "", fmt.Errorf("error in cargo new: %w", err)
	}

	projectDir := filepath.Join(workingDir, projectName)

	cargoTomlPath := filepath.Join(projectDir, "Cargo.toml")
	err = writeFileToDisk(cargoTomlPath, cargoToml)
	if err != nil {
		return "", fmt.Errorf("error in writeFileToDisk: %v, workingDir: %v", err, workingDir)
	}

	err = runCmd(time.Second*120, projectDir, "cargo", "build", "--target", "wasm32-unknown-unknown", "--release")
	if err != nil {
		return "", fmt.Errorf("error in cargo build: %w", err)
	}

	return projectDir, nil
}

func addGlue(code string) string {
	allCode := wasmGlue + code

	finalLines := []string{}
	lines := strings.Split(allCode, "\n")
	for _, line := range lines {

		// add correct entry point
		if strings.Contains(line, "fn") &&
			strings.Contains(line, "main") &&
			strings.Contains(line, "()") {

			finalLines = append(finalLines, "#[no_mangle]")
			finalLines = append(finalLines, `pub extern "C" `+strings.ReplaceAll(line, "main", "lib"))
		} else {
			finalLines = append(finalLines, line)
		}
	}
	return strings.Join(finalLines, "\n")
}

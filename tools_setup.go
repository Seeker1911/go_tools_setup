package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// setupDirs creates a new directory and changes the current working directory to that directory.
func setupDirs(oldDotfileDir string, dotfileDir string, cwd string) {
	fmt.Println("Creating dotfiles_old directory")
	os.MkdirAll(oldDotfileDir, 0777)
	if cwd != dotfileDir {
		fmt.Printf("Changing directory to %s\n", dotfileDir)
		os.Chdir(dotfileDir)

	}
}

// symlinkFiles creates symlink of files from dotfileDir to home directory.
func symlinkFiles(files []string, dotfileDir string, home string) {
	fmt.Printf("Symlinking dotfiles to %s\n", home)
	for _, element := range files {
		os.Symlink(dotfileDir+"/"+element, home+"/"+element)
	}
}

// moveFiles existing dotfiles to dotfilesOld
func moveFiles(files []string, oldDotfileDir string, home string) {
	fmt.Printf("Moving existing dotfiles to %s\n", oldDotfileDir)
	for _, element := range files {
		os.Rename(home+"/"+element, home+"/"+oldDotfileDir+"/"+element)
	}
}

// runPkgMgr installs a Brewfile given a directory it exists in.
// Brew is installed if not found. Working environment, MacOS or LInux, is automatically detected.
func runPkgMgr(dotfileDir string) {
	platform := runtime.GOOS
	fmt.Printf("Runtime is %s\n", platform)
	var shCommand string
	var args []string
	switch {
	case platform == "darwin":
		shCommand = "/usr/local/bin/brew"
		args = append(args, "bundle")

		if _, err := os.Stat(shCommand); os.IsNotExist(err) {
			shCommand = "/usr/bin/ruby"
			args = append(args, "-e", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)")
			fmt.Println("Installing homebrew")
			runCommands(shCommand, args)
			shCommand = "/usr/local/bin/brew"
			args = append(args, "bundle")
		}
	case platform == "linux":
		shCommand = "/home/linuxbrew/.linuxbrew/bin/brew"
		args = append(args, "bundle")

		if _, err := os.Stat(shCommand); os.IsNotExist(err) {
			shCommand = "sh"
			args = append(args, "-c", "$(curl -fsSL https://raw.githubusercontent.com/Linuxbrew/install/master/install.sh)")
			fmt.Println("Installing linuxbrew")
			runCommands(shCommand, args)
			shCommand = "/home/linuxbrew/.linuxbrew/bin/brew"
			args = append(args, "bundle")
		}
	default:
		fmt.Println("Cant find platform for package manager")
	}

	os.Chdir(dotfileDir)
	fmt.Println("Installing Brewfile, this may take a while...")
	runCommands(shCommand, args)
}

// runCommands runs a single command given the command and its arguments.
func runCommands(shCommand string, args []string) {
	cmd := exec.Command(shCommand, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + out.String())
		fmt.Println("Setup complete.")
		return
	}
	fmt.Println("Result: " + out.String())
	fmt.Println("Setup complete.")
}

func main() {
	home := os.Getenv("HOME")
	dotfileDir := fmt.Sprintf("%s/dotfiles", home)
	oldDotfileDir := fmt.Sprintf("%v/dotfiles_old", home)
	var files []string

	testFiles, err := ioutil.ReadDir(dotfileDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Grabbing dotfiles from repo")
	for _, f := range testFiles {
		if strings.HasPrefix(f.Name(), ".") {
			files = append(files, f.Name())
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current working directory is: %s\n", cwd)
	//set up directories
	setupDirs(oldDotfileDir, dotfileDir, cwd)
	//move existing files to dotfiles_old
	moveFiles(files, oldDotfileDir, home)
	//make symlinks and move to home dir.
	symlinkFiles(files, dotfileDir, home)
	//run package manager (i.e. brew or apt)
	runPkgMgr(dotfileDir)
	fmt.Println("You may want to run a few tools:")
	fmt.Println("tmux package manager")
	fmt.Println("vim plugged")
	fmt.Println("check pyenv versions")
}

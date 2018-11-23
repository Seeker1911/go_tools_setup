package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func setupDirs(oldDotfileDir string, dotfileDir string, cwd string) {
	fmt.Println("Creating dotfiles_old directory")
	os.MkdirAll(oldDotfileDir, 0777)
	if cwd != dotfileDir {
		fmt.Printf("changing directory to %s\n", dotfileDir)
		os.Chdir(dotfileDir)

	}
}

func symlinkFiles(files []string, dotfileDir string, home string) {
	fmt.Printf("Symlinking dotfiles to %s\n", home)
	for _, element := range files {
		os.Symlink(dotfileDir+"/"+element, home+"/"+element)
		//fmt.Printf("TEST: Would symlink %s to %s\n",dotfileDir+"/"+element, home+"/"+element )
	}
}

func moveFiles(files []string, oldDotfileDir string, home string) {
	fmt.Printf("Moving existing dotfiles to %s\n", oldDotfileDir)
	for _, element := range files {
		os.Rename(home+"/"+element, home+"/"+oldDotfileDir+"/"+element)
		//fmt.Printf("TEST: would move %s to %s\n",home+"/"+element, home+oldDotfileDir+"/"+element )
	}
}

func runPkgMgr(dotfileDir string) {
	platform := runtime.GOOS
	pkgmgr := ""
	var args []string
	switch {
	case platform == "darwin":
		fmt.Println("platform is macos")
		pkgmgr = "/usr/local/bin/brew"
		args = append(args, "bundle", "check")
	case platform == "linux":
		fmt.Println("platform is linux")
		pkgmgr = "apt"
	default:
		fmt.Println("Cant find platform")
	}
	//run package manager
	os.Chdir(dotfileDir)
	cmd := exec.Command(pkgmgr, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("stderr...")
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("stdout...")
		fmt.Println(fmt.Sprint(err) + ": " + out.String())
		return
	}
	fmt.Println("Result: " + out.String())
	fmt.Println("Setup complete.")
}

func main() {
	home := os.Getenv("HOME")
	dotfileDir := home+"/dotfiles"
	oldDotfileDir := home+"/dotfiles_old"
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory is: %s\n", cwd)
	files := []string{".bashrc", ".tmux.conf", ".vimrc", ".vim", ".ctags", ".gitignore_global", ".Xresources"}


	//set up directories
	setupDirs(oldDotfileDir, dotfileDir, cwd)
	//move existing files to dotfiles_old
	moveFiles(files, oldDotfileDir, home)
	//make symlinks and move to home dir.
	symlinkFiles(files, dotfileDir, home)
	//run package manager (i.e. brew or apt)
	runPkgMgr(dotfileDir)
	//restart bash
	
}
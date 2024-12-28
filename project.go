package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pelletier/go-toml"
)

type ProjectConfig struct {
	Repo     string
	Path     string
	Teardown string
	Setup    string
	Run      string
	Logs     string
	SshKey   string
}

type ProjectHealth struct {
	Sha    string `json:"sha"`
	Status string `json:"status"`
}

// Util methods
func runStep(command []string, workdir string) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = workdir

	fmt.Println("Running: ", command, workdir)

	stdout, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error: ", err.Error(), string(stdout))
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}

func checkout(repo string, commit string, sshKey string) string {
	tmp, _ := os.MkdirTemp(os.TempDir(), "")
	fmt.Println(tmp)

	publicKeys, _ := ssh.NewPublicKeysFromFile("git", sshKey, "")
	localRepo, err := git.PlainClone(tmp, false, &git.CloneOptions{
		URL:  repo,
		Auth: publicKeys,
	})

	refs, _ := localRepo.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			fmt.Println(ref)
		}

		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	w, _ := localRepo.Worktree()

	fmt.Println("Cloning ", commit)
	err = w.Checkout(&git.CheckoutOptions{
		Force: true,
		Hash:  plumbing.NewHash(commit),
	})

	if err != nil {
		fmt.Println("Failed checkout")
		panic(err)
	}

	return tmp
}

func getProjectConfig(key string) ProjectConfig {
	cfg, err := toml.LoadFile("projects.toml")
	fmt.Println("Keys", cfg.Keys())
	if err != nil {
		panic(err)
	}
	var project ProjectConfig
	cfg.Get(key).(*toml.Tree).Unmarshal(&project)
	return project
}

func deployProject(project string, commitSha string) {
	config := getProjectConfig(project)

	fmt.Println("Checking out repo")
	// clone repo into tmp dir
	tmpDir := checkout(config.Repo, commitSha, config.SshKey)

	runStep(strings.Split(config.Setup, " "), tmpDir) // build
	runStep([]string{"rm", "-rf", "./.git"}, tmpDir)  // clean git

	os.MkdirAll(config.Path, fs.ModePerm) // ensure directory exists

	runStep([]string{"/bin/sh", "-c", config.Teardown}, config.Path)             // kill the current running server
	runStep([]string{"/bin/sh", "-c", "cp -r ./* " + config.Path + "/"}, tmpDir) // mv build results
	runStep([]string{"/bin/sh", "-c", config.Run}, config.Path)                  // start the server again
	runStep([]string{"/bin/sh", "-c", "rm -rf ./*"}, tmpDir)                     // clean tmp dir

	os.WriteFile(config.Path+"/sha.log", []byte(commitSha), fs.ModePerm) // dump the commit so we know what version is running
	fmt.Println("Deployed new version of: " + project)
}

func getProjectHealth(project string) ProjectHealth {
	config := getProjectConfig(project)
	var health ProjectHealth
	dat, _ := os.ReadFile(config.Path + "/sha.log")
	health.Sha = string(dat)
	health.Status = "Online"
	return health
}

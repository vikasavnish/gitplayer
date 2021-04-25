package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

var (
	currentPath string
	cleanedPath string
	tracker     string
	commits     []string
	index       int
	prev        int
)

func init() {
	fmt.Println(os.Getwd())
	currentPath, _ := os.Getwd()
	currentUser, _ := user.Current()
	userhome := currentUser.HomeDir
	if currentUser.Username == "root" {
		os.Exit(1)
	}

	cleanedPath = userhome + "/gitlogs/" + strings.ReplaceAll(currentPath, "/", "_") + ".gitlog"
	tracker = userhome + "/gitlogs/" + strings.ReplaceAll(currentPath, "/", "_") + ".yaml"
	_, err := os.Create(cleanedPath)
	_, err = os.Create(tracker)
	if err != nil {
		fmt.Println(err)
	}
	info, err := os.Stat(cleanedPath)
	if os.IsNotExist(err) {
		os.Create(cleanedPath)
	}

	fmt.Println(info.Name())
	cmd := exec.Command("git", "log", "--reverse", "--pretty=oneline")
	out, err := cmd.Output()
	if err != nil {
		os.Remove(cleanedPath)
		os.Remove(tracker)
		log.Fatal(err)

	}
	// fmt.Printf("The date is %s\n", out)
	os.WriteFile(cleanedPath, out, 0644)
	fmt.Println(cleanedPath, tracker)
	_, err = os.Open(cleanedPath)
	if err != nil {
		os.Exit(1)
	}
	commits = readLines(cleanedPath)
}

func readLines(path string) []string {
	file, err := os.Open(cleanedPath)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	return txtlines
}

func diff(first, second int) ([]byte, error) {
	if first < 0 {
		first = 0
	}
	if second < 0 {
		second = 0
	}
	cmd := exec.Command("git", "diff", commits[first][:40], commits[second][:40])
	return cmd.Output()
}
func checkout() {
}

func trigger(action string) {
	switch action {
	case "next", "n", "N":
		prev = index
		index = index + 1
		fmt.Println("next", index, prev)
	case "p", "prev", "back", "-1", "before", "b":
		fmt.Println("back")
		prev = index
		index = index - 1
	case "last", "latest", "final", "end", "e", "l":
		fmt.Println("latest")
		index = len(commits) - 1
		prev = index - 1

	case "first", "start", "f", "s", "reset", "r":
		fmt.Println("start")
		index = 0
		prev = 0

	}
}

func main() {
	var action string

	for {
		fmt.Println("Where to go .. first/last/next/back/previous")
		fmt.Scan(&action)

		fmt.Println(action, prev, index, len(commits))
		if prev <= 0 {
			prev = 0
		}
		if index > len(commits)-1 {
			index = len(commits) - 1
		}
		trigger(action)
		diffOutput, _ := diff(prev, index)
		fmt.Printf("%v", string(diffOutput))
		// proceed to checkout
		if index < 0 {
			index = 0
		}
		fmt.Printf("git checkout %v", commits[index][:40])
		cmdOut, _ := exec.Command("git", "checkout", commits[index][:40]).Output()
		fmt.Printf("%v", cmdOut)

		fmt.Println(action)

	}
}

package main

import (
	"fmt"
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
)


func createNewRepository() {

	fmt.Print("Insert sourceowner: ")
	var sourceowner string
	fmt.Scanln(&sourceowner)

	fmt.Print("Insert sourcerepo: ")
	var sourcerepo string
	fmt.Scanln(&sourcerepo)

	// Imposta il repository di origine che si desidera forcellare
	sourceOwner := sourceowner
	sourceRepo := sourcerepo

	fmt.Print("Insert new repo name: ")
	var reponame string
	fmt.Scanln(&reponame)

	repoName := reponame
	repoDescription := "Repository description"
	private := false

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Personal token as argument")
		return
	}
	accessToken := args[1]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		//&oauth2.Token{AccessToken: ""},
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	baseURL, err := url.Parse("https://api.github.com/")
	if err != nil {
		fmt.Println("Errore during base URL parsing:", err)
		return
	}
	client.BaseURL = baseURL

	url := fmt.Sprintf("repos/%s/%s/forks", sourceOwner, sourceRepo)
	payload := map[string]interface{}{
		"name":        repoName,
		"description": repoDescription,
		"private":     private,
	}

	req, err := client.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		fmt.Println("Errore during request creation:", err)
		return
	}

	_, err = client.Do(ctx, req, nil)
	if err != nil {
		fmt.Println("Errore during fork of the repository:", err)
		return
	}

	fmt.Println("Repository created (fork) with success!")

}


func main() {
	createNewRepository()
}
package main

import (
	"fmt"
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	//"net/url"
	"os"
	"os/exec"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"bufio"
	"path/filepath"
)

// func createNewRepositoryFromFork() {

// 	fmt.Print("Insert sourceowner: ")
// 	var sourceowner string
// 	fmt.Scanln(&sourceowner)

// 	fmt.Print("Insert sourcerepo: ")
// 	var sourcerepo string
// 	fmt.Scanln(&sourcerepo)

// 	// Imposta il repository di origine che si desidera forcellare
// 	sourceOwner := sourceowner
// 	sourceRepo := sourcerepo

// 	fmt.Print("Insert new repo name: ")
// 	var reponame string
// 	fmt.Scanln(&reponame)

// 	repoName := reponame
// 	repoDescription := "Repository description"
// 	private := false

// 	args := os.Args
// 	if len(args) < 2 {
// 		fmt.Println("Personal token as argument")
// 		return
// 	}
// 	accessToken := args[1]

// 	ctx := context.Background()
// 	ts := oauth2.StaticTokenSource(
// 		//&oauth2.Token{AccessToken: ""},
// 		&oauth2.Token{AccessToken: accessToken},
// 	)
// 	tc := oauth2.NewClient(ctx, ts)

// 	client := github.NewClient(tc)

// 	baseURL, err := url.Parse("https://api.github.com/")
// 	if err != nil {
// 		fmt.Println("Errore during base URL parsing:", err)
// 		return
// 	}
// 	client.BaseURL = baseURL

// 	url := fmt.Sprintf("repos/%s/%s/forks", sourceOwner, sourceRepo)
// 	payload := map[string]interface{}{
// 		"name":        repoName,
// 		"description": repoDescription,
// 		"private":     private,
// 	}

// 	req, err := client.NewRequest(http.MethodPost, url, payload)
// 	if err != nil {
// 		fmt.Println("Errore during request creation:", err)
// 		return
// 	}

// 	_, err = client.Do(ctx, req, nil)
// 	if err != nil {
// 		fmt.Println("Errore during fork of the repository:", err)
// 		return
// 	}

// 	fmt.Println("Repository created (fork) with success!")

// }

type RepositoryCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type RepositoryCreateResponse struct {
	Name      string `json:"name"`
	HTMLURL   string `json:"html_url"`
	CloneURL  string `json:"clone_url"`
}

func createRepository(token, repoOwner, repoName, repoDescription string, isPrivate bool) (*RepositoryCreateResponse, error) {
	createURL := fmt.Sprintf("https://api.github.com/user/repos")
	repoCreateRequest := RepositoryCreateRequest{
		Name:        repoName,
		Description: repoDescription,
		Private:     isPrivate,
	}

	payload, err := json.Marshal(repoCreateRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", createURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create repository: %s", body)
	}

	var createResponse RepositoryCreateResponse
	err = json.NewDecoder(resp.Body).Decode(&createResponse)
	if err != nil {
		return nil, err
	}

	return &createResponse, nil
}


func cloneAndPublishRepository(repoName string) {
	// Imposta le informazioni del repository sorgente
	sourceOwner := "attgua"
	sourceRepo := "HelmChart-lens"

	// Imposta le informazioni del nuovo repository
	repoOwner := "attgua"
	//repoName := "NOME_REPO_4"

	// Ottieni il token di accesso personale dall'argomento della riga di comando
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Specificare il token di accesso personale come argomento.")
		return
	}
	accessToken := args[1]

	// Crea una connessione HTTP client personalizzata per includere l'autenticazione
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Esegui la richiesta HTTP GET per ottenere le informazioni sul repository sorgente
	repo, _, err := client.Repositories.Get(ctx, sourceOwner, sourceRepo)
	if err != nil {
		fmt.Println("Errore durante l'ottenimento delle informazioni sul repository sorgente:", err)
		return
	}

	// Esegui il comando "git clone" per clonare il repository sorgente in locale
	cmd := exec.Command("git", "clone", *repo.CloneURL, repoName)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Errore durante l'esecuzione del comando 'git clone':", err)
		return
	}

	// Cambia la directory di lavoro corrente nella directory del repository clonato
	err = os.Chdir(repoName)
	if err != nil {
		fmt.Println("Errore durante il cambio della directory di lavoro:", err)
		return
	}

	// Esegui il comando "git remote" per rimuovere il repository sorgente come remote
	cmd = exec.Command("git", "remote", "remove", "origin")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Errore durante l'esecuzione del comando 'git remote remove origin':", err)
		return
	}

	// Esegui il comando "git remote" per aggiungere il tuo repository come remote
	cmd = exec.Command("git", "remote", "add", "origin", fmt.Sprintf("https://%s@github.com/%s/%s.git",accessToken, repoOwner, repoName))
	err = cmd.Run()
	if err != nil {
		fmt.Println("Errore durante l'esecuzione del comando 'git remote add origin':", err)
		return
	}

	// Esegui il comando "git push" per pubblicare il repository sul tuo account Git
	cmd = exec.Command("git", "push", "-u", "origin", "main")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Errore durante l'esecuzione del comando 'git push':", err)
		return
	}
	fmt.Println("Repository clonato e pubblicato con successo!")

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Errore durante l'ottenimento del percorso della directory corrente:", err)
		return
	}

	// Crea il percorso completo al file "values.yaml" a un livello superiore
	filePath := filepath.Join(currentPath, "/../values.yaml")

	// Stampa il percorso della directory corrente
	fmt.Println("Percorso della directory corrente:", filePath)

	cmd_2 := exec.Command("cp", filePath, "./values-really-custom.yaml")
	err_2 := cmd_2.Run()
	if err_2 != nil {
		fmt.Println("Errore durante la copia", err)
		return
	}

	cmd_3 := exec.Command("git", "add", ".")
	err_3 := cmd_3.Run()
	if err_3 != nil {
		fmt.Println("Errore durante l'esecuzione del comando 'git add':", err)
		return
	}

	// Esegui il comando "git commit" per creare un commit con il file aggiunto
	cmd = exec.Command("git", "commit", "-m", "Aggiunto file")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Errore durante l'esecuzione del comando 'git commit':", err)
		return
	}

	fmt.Println("File aggiunto e commit creato con successo!")

	// Esegui il comando "git push" per pubblicare il repository sul tuo account Git
	cmd = exec.Command("git", "push", "-u", "origin", "main")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Errore durante l'esecuzione del comando 'git push':", err)
		return
	}
	fmt.Println("Repository clonato e pubblicato con successo!")

	currentPath_2, err := os.Getwd()
	if err != nil {
		fmt.Println("Errore durante l'ottenimento del percorso della directory corrente:", err)
		return
	}

	// Rimuovi la directory
	err_new := os.RemoveAll(currentPath_2)
	if err_new != nil {
		fmt.Println("Errore durante l'eliminazione della directory:", err)
		return
	}

	fmt.Println("Directory eliminata con successo!")

}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Specificare il token di accesso personale come argomento.")
		return
	}
	token := args[1]

	// Crea un nuovo scanner per leggere l'input
	scanner := bufio.NewScanner(os.Stdin)

	// Chiedi all'utente di inserire un valore
	fmt.Print("Inserisci il nome del repo da creare: ")

	// Leggi l'input dall'utente
	scanner.Scan()
	input := scanner.Text()

	repoOwner := "attgua"
	repoName := input
	repoDescription := "DESCRIZIONE_REPO"
	private := false

	response, err := createRepository(token, repoOwner, repoName, repoDescription, private)
	if err != nil {
		fmt.Println("Errore durante la creazione del repository:", err)
		return
	}

	fmt.Println("Repository creato con successo:", response.HTMLURL)

	cloneAndPublishRepository(repoName)

}


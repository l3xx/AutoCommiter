package main

import (
	"context"
	"github.com/google/go-github/v24/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"
)

type Config struct {
	PrivateRepo string
	Org         string
	Token       string
}

const fileName = "readme.md"

func main() {
	//init config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var config Config
	config.Token = os.Getenv("TOKEN")
	config.Org = os.Getenv("ORG")
	config.PrivateRepo = os.Getenv("PRIVATE_REPO")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)

	if err != nil {
		log.Fatal(err)
	}
	isExistRepo := false
	for _, repo := range repos {
		if config.PrivateRepo == *repo.Name {
			isExistRepo = true
		}
	}

	if !isExistRepo {
		name := config.PrivateRepo
		newRepository := github.Repository{
			Name:          &name,
			FullName:      &name,
			Description:   &[]string{"Auto commit repo"}[0],
			DefaultBranch: &[]string{"master"}[0],
			MasterBranch:  &[]string{"master"}[0],
		}
		//create repo
		_, _, err := client.Repositories.Create(ctx, "", &newRepository)
		if err != nil {
			log.Fatal(err)
		}
	}

	file := &github.RepositoryContentFileOptions{
		Message: &[]string{"Test message"}[0],
		Content: []byte("Content"),
	}
	_, response, err := client.Repositories.CreateFile(ctx, config.Org, config.PrivateRepo, fileName, file)
	if response.StatusCode != 422 && err != nil {
		log.Fatal(err)
	}

	getFile := &github.RepositoryContentGetOptions{}
	contentFile, _, _, err := client.Repositories.GetContents(ctx, config.Org, config.PrivateRepo, fileName, getFile)
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	file.Content = []byte(t.String())
	file.SHA = contentFile.SHA
	_, _, err = client.Repositories.UpdateFile(ctx, config.Org, config.PrivateRepo, fileName, file)
	if err != nil {
		log.Fatal(err)
	}

}

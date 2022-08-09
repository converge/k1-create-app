package git_client

import (
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/kubefirst/kubefirst-create-app/pkg"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type jsonResponse struct {
	Message struct {
		Name []string `json:"name"`
		Path []string `json:"path"`
	} `json:"message"`
}

func DownloadTemplate(gitHubUser string, gitHubToken string) (*git.Repository, error) {

	err := os.RemoveAll("./tmp")
	if err != nil {
		return nil, err
	}

	// todo: manage tmp creation
	repo, err := git.PlainClone("./tmp", false, &git.CloneOptions{
		URL:      pkg.TemplateURL,
		Progress: os.Stdout,
		Auth: &gitHttp.BasicAuth{
			Username: gitHubUser,
			Password: gitHubToken,
		},
	})
	if err == git.ErrRepositoryAlreadyExists {
		log.Info().Msg("repository already exist, skipping clone...")
		repo, err = git.PlainClone("./tmp", false, &git.CloneOptions{})
		if err != nil {
			return nil, nil
		}
		return repo, nil
	}
	if err != nil {
		return repo, err
	}

	return repo, nil
}

func CreateGitLabRepository(gitLabToken string) error {

	url := "https://gitlab.joao.kubefirst.tech/api/v4/projects"

	payload := strings.NewReader("{\n\t\"name\": \"new-project\",\n\t\"description\": \"New Project\",\n\t\"path\": \"new-project\",\n\t\"initialize_with_readme\": \"false\"\n}")

	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("PRIVATE-TOKEN", gitLabToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var data jsonResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Err(err).Msg("")
	}

	if len(data.Message.Name) > 0 && data.Message.Name[0] == "has already been taken" {
		log.Warn().Msg("GitLab repo. already exist, not creating a new one")
		return nil
	}

	return nil
}

func PushToGitLab(gitLabToken string, repo git.Repository) error {

	err := repo.DeleteRemote("origin")
	if err != nil {
		log.Err(err).Msg("")
		return err
	}
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://gitlab.joao.kubefirst.tech/kubefirst/new-project.git"},
		//URLs: []string{"https://gitlab.com/joao_o/new_project.git"},
	})
	if err != nil {
		log.Err(err).Msg("-2")
		return err
	}

	list, err := repo.Remotes()
	if err != nil {
		log.Err(err).Msg("-1")
		return err
	}

	for _, remotes := range list {
		fmt.Println(remotes)
	}

	workTree, err := repo.Worktree()
	if err != nil {
		log.Err(err).Msg("0")
		return err
	}

	_, err = workTree.Add("go.mod")
	if err != nil {
		fmt.Println(err)
		log.Err(err).Msg("1")
		return err
	}

	commitHash, err := workTree.Commit("go-git commitHash", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Joao",
			Email: "joao@example.org",
			When:  time.Now(),
		},
	})

	_, err = repo.CommitObject(commitHash)
	if err != nil {
		log.Err(err).Msg("2")
		return err
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &gitHttp.BasicAuth{
			Username: "joao_o",
			Password: gitLabToken,
		},
		Force: true,
	})
	if err != nil {
		return err
	}

	return nil
}

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

func DownloadTemplate(gitHubUser string, gitHubToken string) error {

	// todo: manage tmp creation
	_, err := git.PlainClone("./tmp", false, &git.CloneOptions{
		URL:      pkg.TemplateURL,
		Progress: os.Stdout,
		Auth: &gitHttp.BasicAuth{
			Username: gitHubUser,
			Password: gitHubToken,
		},
	})
	if err == git.ErrRepositoryAlreadyExists {
		log.Info().Msg("repository already exist, skipping clone...")
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func CreateGitLabRepository(gitLabToken string) error {

	url := "https://gitlab.com/api/v4/projects"

	payload := strings.NewReader("{\n\t\"name\": \"new_project\",\n\t\"description\": \"New Project\",\n\t\"path\": \"new_project\",\n\t\"initialize_with_readme\": \"true\"\n}")

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

	if data.Message.Name[0] == "has already been taken" {
		log.Warn().Msg("GitLab repo. already exist, not creating a new one")
		return nil
	}

	return nil
}

func SwitchToGitLab(gitLabToken string) (*git.Repository, error) {
	// workarounds :P
	// delete .git
	// create git
	// set remote
	// push
	err := os.RemoveAll("./tmp/.git")
	if err != nil {
		log.Err(err).Msg("")
		return nil, err
	}
	// in memory git repo
	//repo, err := git.Init(memory.NewStorage(), nil)
	//if err != nil {
	//	log.Err(err).Msg("")
	//	return err
	//}

	repo, err := git.PlainClone("./tmp", false, &git.CloneOptions{
		//URL:      pkg.TemplateURL,
		URL:      "https://gitlab.com/joao_o/new_project",
		Progress: os.Stdout,
		Auth: &gitHttp.BasicAuth{
			Username: "joao_o",
			Password: gitLabToken,
		},
	})
	if err != nil {
		log.Err(err).Msg("")
		return nil, err
	}

	err = repo.DeleteRemote("origin")
	if err != nil {
		log.Err(err).Msg("")
		return nil, err
	}
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://gitlab.com/joao_o/new_project.git"},
	})
	if err != nil {
		log.Err(err).Msg("-2")
		return nil, err
	}

	list, err := repo.Remotes()
	if err != nil {
		log.Err(err).Msg("-1")
		return nil, err
	}

	for _, remotes := range list {
		fmt.Println(remotes)
	}

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})

	return repo, nil
}

func PushToGitLab(gitLabToken string, repo git.Repository) error {

	workTree, err := repo.Worktree()
	if err != nil {
		log.Err(err).Msg("0")
		return err
	}

	//_, err = workTree.Add("go.mod")
	_, err = workTree.Add(".gitlab-ci.yml")
	if err != nil {
		log.Err(err).Msg("1")
		return err
	}

	commitHash, err := workTree.Commit("example go-git commitHash", &git.CommitOptions{
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
		Auth: &gitHttp.BasicAuth{
			Username: "joao_o",
			Password: gitLabToken,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

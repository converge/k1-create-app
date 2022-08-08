package main

import (
	"flag"
	gitClient "github.com/kubefirst/kubefirst-create-app/internal/git-client"
	"github.com/kubefirst/kubefirst-create-app/internal/template"
	"github.com/rs/zerolog/log"
)

func main() {

	createAppFlag := flag.Bool("create-app", false, "Create Kubefirst Go Application")
	languageFlag := flag.String("language", "", "Set application programming language to be created")
	gitHubUser := flag.String("gh-user", "", "GitHub user")
	gitHubToken := flag.String("gh-token", "", "GitHub token with repo read access")
	gitLabToken := flag.String("gitlab-token", "", "GitLab token with creating new repo. permission")
	flag.Parse()

	if !*createAppFlag {
		log.Warn().Msg("create app not enable, exiting...")
		return
	}

	if *languageFlag != "go" {
		log.Warn().Msg("Go is the best, no other language is necessary!")
		return
	}

	if len(*gitHubUser) == 0 || len(*gitHubToken) == 0 {
		log.Warn().Msg("missing GitHub credentials")
		return
	}

	if len(*gitLabToken) == 0 {
		log.Warn().Msg("missing GitLab token")
		return
	}

	// download template
	err := gitClient.DownloadTemplate(*gitHubUser, *gitHubToken)
	if err != nil {
		log.Err(err).Msg("")
		return
	}

	// apply changes to:
	projectBase := template.ProjectBase{
		ProjectName:  "kubefirst-create-app",
		ProjectOwner: "converge",
	}
	err = template.ApplyGoModChange(projectBase)
	if err != nil {
		log.Err(err).Msg("-4")
		return
	}

	// todo: chart.yaml.template
	err = gitClient.CreateGitLabRepository(*gitLabToken)
	if err != nil {
		log.Err(err).Msg("")
		return
	}

	repo, err := gitClient.SwitchToGitLab(*gitLabToken)
	if err != nil {
		log.Err(err).Msg("")
		return
	}

	err = gitClient.PushToGitLab(*gitLabToken, *repo)
	if err != nil {
		log.Err(err).Msg("")
		return
	}

}

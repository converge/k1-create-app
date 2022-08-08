package template

import (
	"html/template"
	"os"
)

type ProjectBase struct {
	ProjectName  string
	ProjectOwner string
}

func ApplyGoModChange(projectBase ProjectBase) error {
	tpl, err := template.ParseFiles("./tmp/go.mod.template")
	if err != nil {
		return err
	}

	localFile, err := os.OpenFile("./tmp/go.mod", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	err = tpl.Execute(localFile, projectBase)
	if err != nil {
		return err
	}

	return nil
}

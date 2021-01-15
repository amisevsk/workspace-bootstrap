package main

import (
	"flag"
	"log"

	"github.com/amisevsk/workspace-bootstrap/devworkspace-sync/funcs"
)

const (
	syncDevfileArg              = "sync-devfile"
	syncDevWorkspaceTemplateArg = "sync-template"
)

func main() {
	var syncDevfile, syncDevWorkspaceTemplate bool

	flag.BoolVar(&syncDevfile, syncDevfileArg, false, "Sync devfile from repo to current DevWorkspace")
	flag.BoolVar(&syncDevWorkspaceTemplate, syncDevWorkspaceTemplateArg, false, "Sync DevWorkspaceTemplate from repo to current DevWorkspace")

	flag.Parse()

	if syncDevWorkspaceTemplate {
		err := funcs.SyncDevWorkspaceTemplate()
		if err != nil {
			log.Fatal(err)
		}
	}
	if syncDevfile {
		err := funcs.SyncDevfile()
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Operation completed successfully")
}

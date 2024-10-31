package main

import (
	"encoding/json"

	"net/http"
)

type APIError struct {
	Success bool  `json:"success"`
	Reason  error `json:"reason"`
}

type APISuccess struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type DeployRequest struct {
	Project string `json:"project"`
	Commit  string `json:"commit"`
}

func handleDeploy(w http.ResponseWriter, r *http.Request) {
	var req DeployRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go deployProject(req.Project, req.Commit)

	response := APISuccess{Success: true, Message: "Deploy queued for project"}
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	project := r.FormValue("project")
	health := getProjectHealth(project)
	json.NewEncoder(w).Encode(health)
}

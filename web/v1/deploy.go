package v1

import (
	"encoding/json"
	"log"
	"net/http"
)

type Deployment struct {
	Cluster     string `json:"cluster"`
	Environment string `json:"env"`
	Namespace   string `json:"namespace"`
	Project     string `json:"project"`
	Version     string `json:"version"`
}

func (d Deployment) IsValid() bool {
	return len(d.Cluster) > 0 && len(d.Environment) > 0 &&
		len(d.Namespace) > 0 && len(d.Project) > 0 && len(d.Version) > 0
}

func DeployHandler(w http.ResponseWriter, r *http.Request) {
	var d Deployment
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, `{msg: "bad json"}`, http.StatusBadRequest)
		log.Println("[deploy]", "bad json")
		return
	}

	if !d.IsValid() {
		http.Error(w, `{msg: "invalid params"}`, http.StatusBadRequest)
		log.Println("[deploy]", "invalid params")
		return
	}

	log.Println("[deploy]", "deployment started:", d)
}

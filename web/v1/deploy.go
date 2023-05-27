package v1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"greenfield-deploy/pkg/github"
	"greenfield-deploy/pkg/k8s"
	"io"
	"log"
	"net/http"
	"strings"
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

	log.Printf("[deploy] deployment started: %+v", d)
	cc, err := github.Content(d.Project, d.Environment)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[deploy] content: %+v", cc)
	for _, c := range cc {
		r, err := github.DownloadContent(c)
		if err != nil {
			log.Printf("error on download content: %v\n", err)
			continue
		}

		var (
			kind, version string
			vm            strings.Builder
		)
		s := bufio.NewScanner(r)
		for s.Scan() {
			var l = s.Text()
			switch {
			case strings.Contains(l, "apiVersion:"):
				version = l
				fallthrough
			case strings.Contains(l, "kind:"):
				kind = l
				fallthrough
			case !strings.Contains(l, "image:"):
				vm.WriteString(l)
				vm.WriteString("\n")
				continue
			}
			vm.WriteString(strings.ReplaceAll(l, ":latest", fmt.Sprintf(":%s", d.Version)))
			vm.WriteString("\n")
		}
		log.Println("found manifest", "name", c.GetName(), "kind", kind, "api version", version)
		manifest := io.NopCloser(strings.NewReader(vm.String()))
		if err := k8s.DeployToNamespace(k8s.NewKubernetesConfigLocal(), d.Namespace, manifest, false); err != nil {
			log.Fatal(err)
		}
	}
	w.WriteHeader(http.StatusOK)
}

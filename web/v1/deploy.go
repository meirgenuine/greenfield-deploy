package v1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"greenfield-deploy/pkg/github"
	"greenfield-deploy/pkg/k8s"
	"greenfield-deploy/pkg/notification"
	"io"
	"log"
	"net/http"
	"strings"
)

type DeployRequest struct {
	Username string `json:"username"`
	ChatID   int64  `json:"chatID"`
	Deployment
}

func (r DeployRequest) String() string {
	return fmt.Sprintf("%+v", r.Deployment)
}

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

type deploymentService struct {
	messenger notification.Messenger
}

func NewHandler(m notification.Messenger) *deploymentService {
	return &deploymentService{
		messenger: m,
	}
}

func (h deploymentService) DeployHandler(w http.ResponseWriter, r *http.Request) {
	var d DeployRequest
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, `bad json`, http.StatusBadRequest)
		log.Println("[deploy]", "bad json")
		return
	}

	if !d.IsValid() {
		http.Error(w, `invalid params`, http.StatusBadRequest)
		log.Println("[deploy]", "invalid params")
		return
	}

	log.Printf("[deploy] deployment started: %+v", d)

	w.WriteHeader(http.StatusOK)

	go h.deploy(&d)
}

func (h deploymentService) deploy(d *DeployRequest) {
	var reqErr error
	midErr := make([]error, 0, 2)

	defer func() {
		if reqErr != nil {
			h.Notify(d, fmt.Sprintf("Error occurred: %s", reqErr))
		} else if len(midErr) > 0 {
			h.Notify(d, fmt.Sprintf("Deployed with error: %s", midErr[0]))
		} else if r := recover(); r != nil {
			h.Notify(d, fmt.Sprintf("Server error occurred: %v", r))
		} else {
			h.Notify(d, "Successfully deployed")
		}
	}()
	cc, err := github.Content(d.Project, d.Environment)
	if err != nil {
		log.Println("[deploy]", "github", err)
		reqErr = fmt.Errorf("github: content: %w", err)
		return
	}
	log.Printf("[deploy] content: %+v", cc)
	for _, c := range cc {
		r, err := github.DownloadContent(c)
		if err != nil {
			log.Printf("error on download content: %v\n", err)
			midErr = append(midErr, fmt.Errorf("github: download: %w", err))
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
		if err := k8s.Deploy(k8s.NewKubernetesConfigLocal(), d.Namespace, manifest); err != nil {
			log.Println("[deploy]", "deploy to namespace", err)
			midErr = append(midErr, fmt.Errorf("k8s: deploy: %w", err))
			continue
		}
	}
}

func (h deploymentService) Notify(r *DeployRequest, message string) {
	h.messenger.Notify(
		notification.User{
			ChatID: r.ChatID,
		},
		notification.Notification{
			Message: fmt.Sprintf(
				"%s:```\nProject: %s\nVersion: %s\nCluster: %s\nNamespace: %s\nEnv: %s```",
				message,
				r.Project,
				r.Version,
				r.Cluster,
				r.Namespace,
				r.Environment,
			),
		},
	)
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	web "greenfield-deploy/web/v1"
	"io"
	"log"
	"net/http"
)

func (b Bot) StartHandler() string {
	return getListCommands()
}

func (b Bot) DeployHandler(u User, args ...string) string {
	if len(args) < 6 {
		return "invalid args"
	}

	d := web.DeployRequest{
		Username: u.Name,
		ChatID:   u.ChatID,
		Deployment: web.Deployment{
			Project:     args[1],
			Version:     args[2],
			Cluster:     args[3],
			Namespace:   args[4],
			Environment: args[5],
		},
	}
	body, err := json.Marshal(d)
	if err != nil {
		log.Printf("deploy failed, args: %v, resp: %s\n", args, err)
		return err.Error()
	}
	resp, err := b.sendRequest(body)
	if err != nil {
		log.Printf("deploy failed, args: %v, resp: %s\n", args, err)
		return err.Error()
	}
	return resp
}

func (b *Bot) sendRequest(body []byte) (string, error) {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/deploy", b.conf.DeploymentServiceURL),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

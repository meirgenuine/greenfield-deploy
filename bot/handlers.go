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

type Handler struct {
	client *http.Client
	conf   *Config
}

func NewHandler(conf *Config) *Handler {
	return &Handler{
		// todo config client
		client: &http.Client{},
		conf:   conf,
	}
}

func (h Handler) Start() string {
	return getListCommands()
}

func (h Handler) Deploy(args ...string) string {
	if len(args) < 6 {
		return "invalid args"
	}

	d := web.Deployment{
		Project:     args[1],
		Version:     args[2],
		Cluster:     args[3],
		Namespace:   args[4],
		Environment: args[5],
	}
	body, err := json.Marshal(d)
	if err != nil {
		log.Printf("deploy failed, args: %v, resp: %s\n", args, err)
		return err.Error()
	}
	resp, err := h.sendRequest(body)
	if err != nil {
		log.Printf("deploy failed, args: %v, resp: %s\n", args, err)
		return err.Error()
	}
	return resp
}

func (h Handler) sendRequest(body []byte) (string, error) {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/deploy", h.conf.API.URL),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	resp, err := h.client.Do(req)
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

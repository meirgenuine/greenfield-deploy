package main

import (
	"bytes"
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
	if len(args) < 2 {
		return "invalid args"
	}

	resp, err := h.sendRequest("")
	if err != nil {
		resp = err.Error()
		log.Printf("deploy failed, args: %v, resp: %s\n", args, resp)
	}
	return resp
}

func (h Handler) sendRequest(message string) (string, error) {
	req, err := http.NewRequest(
		"POST",
		h.conf.API.URL,
		bytes.NewBufferString(message),
	)
	if err != nil {
		return "", err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

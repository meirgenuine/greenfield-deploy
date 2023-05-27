package main

import (
	"bytes"
	"greenfield-deploy/pkg/k8s"
	"io"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "./config")
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	k8scfg := k8s.KubernetesConfig{
		Clientset: clientset,
	}

	f, err := os.ReadFile("./busybox.yaml")
	if err != nil {
		panic(err)
	}

	c := io.NopCloser(bytes.NewReader(f))
	err = k8s.DeployToNamespace(&k8scfg, "prod", c, false)
	if err != nil {
		panic(err)
	}
}

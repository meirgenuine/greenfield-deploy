package k8s

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesConfig struct {
	sync.RWMutex
	Clientset *kubernetes.Clientset

	IsLocal   bool
	CaFile    string `mapstructure:"ca_file"`
	Host      string `mapstructure:"host"`
	TokenFile string `mapstructure:"token_file"`
}

func NewKubernetesConfigLocal() *KubernetesConfig {
	return &KubernetesConfig{
		IsLocal: true,
	}
}

func NewKubernetesConfig(host, caFile, tokenFile string) *KubernetesConfig {
	return &KubernetesConfig{
		CaFile:    caFile,
		Host:      host,
		TokenFile: tokenFile,
	}
}

func ClientSet(cfg *KubernetesConfig) (*kubernetes.Clientset, error) {
	if cfg == nil {
		return nil, errors.New("k8s configuration is not defined for the cluster")
	}

	cfg.RLock()
	if cfg.Clientset != nil {
		defer cfg.RUnlock()
		return cfg.Clientset, nil
	}
	cfg.RUnlock()

	cfg.Lock()
	defer cfg.Unlock()
	if cfg.Clientset != nil {
		return cfg.Clientset, nil
	}

	var err error
	if cfg.IsLocal {
		config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
		if err != nil {
			return nil, err
		}
		cfg.Clientset, err = kubernetes.NewForConfig(config)
		return cfg.Clientset, err
	}

	cfg.Clientset, err = kubernetes.NewForConfig(&rest.Config{
		BearerTokenFile: cfg.TokenFile,
		Host:            cfg.Host,
		TLSClientConfig: rest.TLSClientConfig{
			CAFile: cfg.CaFile,
		},
	})
	return cfg.Clientset, err
}

func Deploy(cfg *KubernetesConfig, n string, r io.ReadCloser) error {
	mm, err := Manifests(r)
	if err != nil {
		return err
	}

	for _, m := range mm {
		err = MustUpdate(cfg, n, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func waitJob(ctx context.Context, status watch.EventType, selector, n string, cs *kubernetes.Clientset, out chan error) {
	w, err := cs.BatchV1().Jobs(n).Watch(ctx, v1.ListOptions{
		FieldSelector: selector,
	})
	if err != nil {
		out <- err
		return
	}
	defer w.Stop()
	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				out <- fmt.Errorf("watcher closed unexpectedly")
				return
			}
			if event.Type == status {
				out <- nil
				return
			}
		case <-ctx.Done():
			out <- nil
			return
		}
	}
}

func MustUpdate(cfg *KubernetesConfig, n string, m runtime.Object) error {
	var (
		cs, err = ClientSet(cfg)
		ctx     = context.Background()
	)
	if err != nil {
		return err
	}
	switch m.(type) {
	case *appsv1.Deployment:
		u := cs.AppsV1().Deployments(n).Update
		c := cs.AppsV1().Deployments(n).Create
		_, err := u(ctx, m.(*appsv1.Deployment), metav1.UpdateOptions{})
		if k8s_errors.IsNotFound(err) {
			err = nil
			_, err = c(ctx, m.(*appsv1.Deployment), metav1.CreateOptions{})
		}
		if err != nil {
			return err
		}
	case *batchv1.CronJob:
		u := cs.BatchV1().CronJobs(n).Update
		c := cs.BatchV1().CronJobs(n).Create
		_, err := u(ctx, m.(*batchv1.CronJob), metav1.UpdateOptions{})
		if k8s_errors.IsNotFound(err) {
			err = nil
			_, err = c(ctx, m.(*batchv1.CronJob), metav1.CreateOptions{})
		}
		if err != nil {
			return err
		}
	case *batchv1.Job:
		deleteCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		d := cs.BatchV1().Jobs(n).Delete
		c := cs.BatchV1().Jobs(n).Create
		errChan := make(chan error)

		// Delete fully job
		go waitJob(deleteCtx, watch.Deleted, fmt.Sprintf("metadata.name=%s", m.(*batchv1.Job).Name), n,
			cs, errChan)
		dp := metav1.DeletePropagationForeground
		err := d(deleteCtx, m.(*batchv1.Job).Name, metav1.DeleteOptions{
			PropagationPolicy: &dp,
		})
		if err != nil {
			if k8s_errors.IsNotFound(err) {
				cancel()
			} else {
				return err
			}
		}

		select {
		case err = <-errChan:
			break
		case <-time.After(time.Minute):
			err = errors.New("deadline exceed for waiting job's delete event")
		}
		close(errChan)

		if err != nil && !k8s_errors.IsNotFound(err) {
			return err
		}

		// Create new job
		_, err = c(ctx, m.(*batchv1.Job), metav1.CreateOptions{})
		if err != nil {
			return err
		}
	case *corev1.Pod:
		u := cs.CoreV1().Pods(n).Update
		c := cs.CoreV1().Pods(n).Create
		_, err := u(ctx, m.(*corev1.Pod), metav1.UpdateOptions{})
		if k8s_errors.IsNotFound(err) {
			err = nil
			_, err = c(ctx, m.(*corev1.Pod), metav1.CreateOptions{})
		}
		if err != nil {
			return err
		}
	case *corev1.Service:
		u := cs.CoreV1().Services(n).Update
		c := cs.CoreV1().Services(n).Create
		_, err := u(ctx, m.(*corev1.Service), metav1.UpdateOptions{})
		if k8s_errors.IsNotFound(err) {
			err = nil
			_, err = c(ctx, m.(*corev1.Service), metav1.CreateOptions{})
		}
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("couldn't deploy manifest for resource: %s", m.GetObjectKind().GroupVersionKind().Kind)
	}
	return nil
}

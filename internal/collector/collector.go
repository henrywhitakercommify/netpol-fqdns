// Package collector
package collector

import (
	"context"
	"fmt"
	"os"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Rule struct {
	Name      string   `json:"name"      csv:"name"`
	Namespace string   `json:"namespace" csv:"namespace"`
	FQDNs     []string `json:"fqdns"     csv:"fqdns"`
}

func Collect(ctx context.Context, kubeconfig string) ([]Rule, error) {
	path, err := resolveConfigPath(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("resolve kubeconfig path: %w", err)
	}
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, fmt.Errorf("build client config: %w", err)
	}
	_, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("init k8s client: %w", err)
	}
	return nil, fmt.Errorf("not implemeted yet")
}

func resolveConfigPath(path string) (string, error) {
	if strings.Contains(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = strings.ReplaceAll(path, "~", home)
	}

	return path, nil
}

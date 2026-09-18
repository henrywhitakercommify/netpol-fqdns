// Package collector
package collector

import (
	"context"
	"fmt"
	"os"
	"strings"

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	ciliumclient "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type Rule struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	FQDNs     []string `json:"fqdns"`
}

func Collect(ctx context.Context, kubeconfig string) ([]Rule, error) {
	client, err := newClient(kubeconfig)
	if err != nil {
		return nil, err
	}

	policies, err := listCiliumNetworkPolicies(ctx, client)
	if err != nil {
		return nil, err
	}

	out := []Rule{}
	for _, p := range policies {
		rule := Rule{
			Name:      p.Name,
			Namespace: p.Namespace,
		}
		for _, e := range p.Spec.Egress {
			for _, f := range e.ToFQDNs {
				if f.MatchName != "" {
					rule.FQDNs = append(rule.FQDNs, f.MatchName)
				}
				if f.MatchPattern != "" {
					rule.FQDNs = append(rule.FQDNs, f.MatchPattern)
				}
			}
		}
		if len(rule.FQDNs) > 0 {
			out = append(out, rule)
		}
	}

	return out, nil
}

// listCiliumNetworkPolicies returns every CiliumNetworkPolicy in the cluster.
func listCiliumNetworkPolicies(
	ctx context.Context,
	client ciliumclient.Interface,
) ([]ciliumv2.CiliumNetworkPolicy, error) {
	policies := []ciliumv2.CiliumNetworkPolicy{}

	opts := metav1.ListOptions{Limit: 500}
	for {
		list, err := client.CiliumV2().
			CiliumNetworkPolicies(metav1.NamespaceAll).
			List(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("list ciliumnetworkpolicies: %w", err)
		}

		policies = append(policies, list.Items...)

		opts.Continue = list.GetContinue()
		if opts.Continue == "" {
			return policies, nil
		}
	}
}

func newClient(kubeconfig string) (ciliumclient.Interface, error) {
	path, err := resolveConfigPath(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("resolve kubeconfig path: %w", err)
	}
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, fmt.Errorf("build client config: %w", err)
	}
	client, err := ciliumclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("init k8s client: %w", err)
	}
	return client, nil
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

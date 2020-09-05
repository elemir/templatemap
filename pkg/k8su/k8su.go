package k8su

import (
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func RunAs(clientset *kubernetes.Clientset, config *rest.Config, namespace string, serviceAccount string) (*kubernetes.Clientset, error) {
	sa, err := clientset.CoreV1().ServiceAccounts(namespace).Get(serviceAccount, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot find service account '%s/%s': %w", namespace, serviceAccount, err)
	}

	if len(sa.Secrets) == 0 {
		return nil, fmt.Errorf("cannot find any token for serviceAccount '%s/%s'", namespace, serviceAccount)
	}

	if sa.Secrets[0].Namespace != "" {
		namespace = sa.Secrets[0].Namespace
	}

	secret, err := clientset.CoreV1().Secrets(namespace).Get(sa.Secrets[0].Name, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot get token for serviceAccount '%s/%s': %w", namespace, sa.Secrets[0].Name, err)
	}

	config = rest.AnonymousClientConfig(config)
	config.BearerToken = string(secret.Data["token"])

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create clientset from config: %w", err)
	}

	return clientset, nil
}

package template

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Template struct {
	Pod  Metadata
	Node Metadata
}

func NewTemplate(clientset *kubernetes.Clientset, namespace string, podName string) (*Template, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(podName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	node, err := clientset.CoreV1().Nodes().Get(pod.Spec.NodeName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &Template{
		Pod: Metadata{
			Annotations: pod.ObjectMeta.GetAnnotations(),
			Labels:      pod.ObjectMeta.GetLabels(),
			Name:        pod.ObjectMeta.GetName(),
			Namespace:   pod.ObjectMeta.GetNamespace(),
		},
		Node: Metadata{
			Annotations: node.ObjectMeta.GetAnnotations(),
			Labels:      node.ObjectMeta.GetLabels(),
			Name:        node.ObjectMeta.GetName(),
			Namespace:   node.ObjectMeta.GetNamespace(),
		},
	}, nil
}

func (t *Template) GenerateFile(dir string, filename string, data string) error {
	fp := filepath.Join(dir, filename)

	if !filepath.IsAbs(fp) {
		return fmt.Errorf("Cannot use nonabsolute path")
	}

	file, err := os.Create(fp)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filename).Parse(data)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, t)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

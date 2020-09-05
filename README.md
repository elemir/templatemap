# TemplateMap

TemplateMap is a ephemeral CSI driver for kubernetes that allow ConfigMap gotemplating

## Installation

You can use prebuilt image from my dockerhub with `make deploy` command

## Usage

TemplateMap CSI plugin support several modes of mounting ConfigMap:

1. Mounting whole ConfigMap as a directory
2. Mounting only key file from ConfigMap with subPath

Every file is templated by standard go templating. You can access to pod and node metadata with eponymous fields and use its name, namespace, labels and annotations. For security issues templatemap use ServiceAccount from the pod or via special serviceAccountName attribute of ephemeral storage.

Example of pod manifest:

```.yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.19.2-alpine
    ports:
    - containerPort: 80
    volumeMounts:
    - name: test
      mountPath: /usr/share/nginx/html/index.html
      subPath: index.html
  volumes:
  - name: test
    csi:
      driver: csi-templatemap
      volumeAttributes:
        configMapName: test-configmap
        serviceAccountName: test-sa
```

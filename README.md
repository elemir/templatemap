# TemplateMap

TemplateMap is a ephemeral CSI driver for kubernetes that allow ConfigMap gotemplating

## Installation

You can use prebuilt image from my dockerhub with `make deploy` command

## Usage

TemplateMap CSI plugin support several modes of mounting ConfigMap:

1. Mounting ConfigMap as a directory
2. Mounting only key file from ConfigMap with subPath

Every file is templated by standard go templating. You can access to pod and node metadata with eponymous fields and use its name, namespace, labels and annotations.

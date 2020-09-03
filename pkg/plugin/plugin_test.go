package plugin

import (
	"fmt"
	"os"
	"testing"

	"github.com/kubernetes-csi/csi-test/v4/pkg/sanity"
	"github.com/sirupsen/logrus"
)

func TestPlugin(t *testing.T) {
	endpoint := "unix:///tmp/csi.sock"
	log := logrus.New()
	log.Level = logrus.DebugLevel

	config := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))

	plugin, err := NewPlugin("test-node", endpoint, log)
	if err != nil {
		t.Errorf("cannot initialize plugin: %s", err)
		return
	}

	plugin.SetKubeConfig(config)

	end := make(chan struct{})

	go func() {
		err := plugin.Run()
		if err != nil {
			t.Errorf("cannot run plugin: %s", err)
		}
		end <- struct{}{}
	}()

	go func() {
		config := sanity.NewTestConfig()
		config.Address = endpoint

		sanity.Test(t, config)
		end <- struct{}{}
	}()

	_ = <-end
}

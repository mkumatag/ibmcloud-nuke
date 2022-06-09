package app

import (
	"github.com/mkumatag/ibmcloud-nuke/internal/pkg"
	"github.com/mkumatag/ibmcloud-nuke/internal/resources"
	"gopkg.in/yaml.v2"
	"os"
)

func Run(config string) error {
	content, err := os.ReadFile(config)
	if err != nil {
		return err
	}
	conf := &pkg.Config{}
	if err := yaml.Unmarshal(content, conf); err != nil {
		return err
	}
	//spew.Dump(conf)
	for _, r := range conf.Resources {
		for _, f := range resources.ResourceFuncs {
			n, err := f(*conf, r)
			if err != nil {
				return err
			}
			if err := n.NukeIt(); err != nil {
				return err
			}
		}

	}
	return nil
}

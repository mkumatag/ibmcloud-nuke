package utils

import (
	"github.com/mkumatag/ibmcloud-nuke/internal/pkg"
)

func GetRegion(c pkg.Config, r pkg.Resource) (string, error) {
	if c.Global.Region != "" {
		return c.Global.Region, nil
	} else if r.Region != "" {
		return r.Region, nil
	}
	return "", pkg.ErrorMissingRegion
}

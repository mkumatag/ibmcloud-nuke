package resources

import (
	"github.com/mkumatag/ibmcloud-nuke/internal/pkg"
)

//var Resources []pkg.Resource

var ResourceFuncs []ResourceFunc

type ResourceFunc = func(pkg.Config, pkg.Resource) (pkg.Nuke, error)

//func RegisterResource(r pkg.Resource) {
//	Resources = append(Resources, r)
//}

func RegisterResourceFunc(rf ResourceFunc) {
	ResourceFuncs = append(ResourceFuncs, rf)
}

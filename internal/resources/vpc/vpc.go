package vpc

import (
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/mkumatag/ibmcloud-nuke/internal/pkg"
	"github.com/mkumatag/ibmcloud-nuke/internal/pkg/utils"
	"github.com/mkumatag/ibmcloud-nuke/internal/resources"
	"log"
)

func init() {
	resources.RegisterResourceFunc(New)
}

const rtype = "vpc"

type VPC struct {
	auth     core.Authenticator
	nukes    []pkg.Nuke
	service  *vpcv1.VpcV1
	resource pkg.Resource
}

func (v *VPC) Validate() error {
	//TODO implement me
	panic("implement me")
}

func (v *VPC) NukeIt() error {
	log.Println("Nuking VPC resource")
	vpc, err := v.getVPC()
	if err != nil {
		return err
	}
	if vpc == nil {
		return pkg.ErrorResourceNotFound
	}
	log.Printf("%s(vpc) found with crn: %s", *vpc.Name, *vpc.ID)
	// Step 1: nuke the subnets
	{
		subnet := subnet{
			service: v.service,
			vpc:     vpc,
		}
		if err := subnet.Nuke(); err != nil {
			return err
		}
	}

	// Step 2: nuke the public gateways
	{
		gateway := publicGateway{
			service: v.service,
			vpc:     vpc,
		}
		if err := gateway.Nuke(); err != nil {
			return err
		}
	}

	// Step 3: nuke the vpc
	{
		deleteVPCOpt := &vpcv1.DeleteVPCOptions{
			ID: vpc.ID,
		}
		log.Printf("deleting the %s(vpc) with ID: %s", *vpc.Name, *vpc.ID)
		if _, err := v.service.DeleteVPC(deleteVPCOpt); err != nil {
			return err
		}
	}
	return nil
}

func (v *VPC) getVPC() (*vpcv1.VPC, error) {
	var found *vpcv1.VPC
	f := func(start string) (isDone bool, nextUrl string, err error) {
		listOptions := &vpcv1.ListVpcsOptions{}
		if start != "" {
			listOptions.Start = &start
		}
		vpcs, _, err := v.service.ListVpcs(listOptions)
		if err != nil {
			return
		}
		if vpcs == nil || len(vpcs.Vpcs) == 0 {
			log.Println("no vpcs found!")
			return
		}
		for _, vpc := range vpcs.Vpcs {
			if *vpc.Name == v.resource.Filter.Name {
				found = &vpc
				return
			}
		}
		// For paging over next set of resources getting the start token and passing it for next iteration
		if vpcs.Next != nil && *vpcs.Next.Href != "" {
			nextUrl = *vpcs.Next.Href
			return
		}

		isDone = true
		return
	}
	if err := utils.PagingHelper(f); err != nil {
		return nil, err
	}
	return found, nil
}

func (v *VPC) nukePublicGateways(vpc *vpcv1.VPC) {

}

func New(c pkg.Config, r pkg.Resource) (pkg.Nuke, error) {
	if r.Type != rtype {
		return nil, pkg.ErrorResourceTypeMatch
	}
	auth, err := pkg.GetAuthenticator()
	if err != nil {
		return nil, err
	}

	service, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: auth,
	})
	region, err := utils.GetRegion(c, r)
	if err != nil {
		return nil, err
	}
	if err := service.SetServiceURL(fmt.Sprintf("https://%s.iaas.cloud.ibm.com/v1", region)); err != nil {
		return nil, err
	}
	return &VPC{
		auth:     auth,
		service:  service,
		resource: r,
	}, nil
}

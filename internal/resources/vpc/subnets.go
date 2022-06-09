package vpc

import (
	"fmt"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/mkumatag/ibmcloud-nuke/internal/pkg/utils"
	"log"
)

type subnet struct {
	vpc     *vpcv1.VPC
	service *vpcv1.VpcV1
}

func (l *subnet) Nuke() error {
	if l.vpc == nil {
		return fmt.Errorf("vpc can't be empty")
	}

	f := func(start string) (isDone bool, nextUrl string, err error) {
		listOptions := &vpcv1.ListSubnetsOptions{}
		if start != "" {
			listOptions.Start = &start
		}
		subnets, _, err := l.service.ListSubnets(listOptions)
		if err != nil {
			return
		}
		if subnets == nil || len(subnets.Subnets) == 0 {
			log.Println("no subnets found!")
			return
		}
		for _, subnet := range subnets.Subnets {
			if *subnet.VPC.CRN == *l.vpc.CRN {
				log.Printf("%s(subnet) found with crn: %s", *subnet.Name, *subnet.CRN)
				err = l.deleteSubnet(subnet)
				if err != nil {
					return
				}
			}
		}
		// For paging over next set of resources getting the start token and passing it for next iteration
		if subnets.Next != nil && *subnets.Next.Href != "" {
			nextUrl = *subnets.Next.Href
			return
		}

		isDone = true
		return
	}
	return utils.PagingHelper(f)
}

func (l *subnet) deleteSubnet(s vpcv1.Subnet) error {
	// Step 1: delete the loadbalancers
	{
		lb := loadBalancer{
			vpc:     l.vpc,
			service: l.service,
			subnet:  &s,
		}
		if err := lb.Nuke(); err != nil {
			return err
		}
	}
	// Step 2: delete the subnet
	{
		log.Printf("deleting the %s(subnet) with ID: %s", *s.Name, *s.ID)
		subnetDeleteOpt := &vpcv1.DeleteSubnetOptions{
			ID: s.ID,
		}
		if _, err := l.service.DeleteSubnet(subnetDeleteOpt); err != nil {
			return err
		}
	}

	return nil
}

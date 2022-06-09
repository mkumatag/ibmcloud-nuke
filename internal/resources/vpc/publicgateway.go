package vpc

import (
	"fmt"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/mkumatag/ibmcloud-nuke/internal/pkg/utils"
	"log"
)

type publicGateway struct {
	vpc     *vpcv1.VPC
	service *vpcv1.VpcV1
}

func (g *publicGateway) Nuke() error {
	if g.vpc == nil {
		return fmt.Errorf("vpc can't be empty")
	}

	f := func(start string) (isDone bool, nextUrl string, err error) {
		listOptions := &vpcv1.ListPublicGatewaysOptions{}
		if start != "" {
			listOptions.Start = &start
		}
		gateways, _, err := g.service.ListPublicGateways(listOptions)
		if err != nil {
			return
		}
		if gateways == nil || len(gateways.PublicGateways) == 0 {
			log.Println("no public gateways found!")
			return
		}
		for _, gateway := range gateways.PublicGateways {
			if *gateway.VPC.CRN == *g.vpc.CRN {
				log.Printf("%s(gateway) found with ID: %s", *gateway.Name, *gateway.ID)
				err = g.deletePublicGateway(gateway)
				if err != nil {
					return
				}
			}
		}
		// For paging over next set of resources getting the start token and passing it for next iteration
		if gateways.Next != nil && *gateways.Next.Href != "" {
			nextUrl = *gateways.Next.Href
			return
		}

		isDone = true
		return
	}
	return utils.PagingHelper(f)
}

func (g *publicGateway) deletePublicGateway(gateway vpcv1.PublicGateway) error {
	log.Printf("deleting the %s(public gateway) with ID: %s", *gateway.Name, *gateway.ID)
	gatewayDeleteOpt := &vpcv1.DeletePublicGatewayOptions{
		ID: gateway.ID,
	}
	if _, err := g.service.DeletePublicGateway(gatewayDeleteOpt); err != nil {
		return err
	}
	return nil
}

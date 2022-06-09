package vpc

import (
	"fmt"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/mkumatag/ibmcloud-nuke/internal/pkg/utils"
	"k8s.io/apimachinery/pkg/util/wait"
	"log"
	"time"
)

type loadBalancer struct {
	vpc     *vpcv1.VPC
	service *vpcv1.VpcV1
	subnet  *vpcv1.Subnet
}

func (l *loadBalancer) Nuke() error {
	if l.vpc == nil {
		return fmt.Errorf("vpc can't be empty")
	}
	if l.subnet == nil {
		return fmt.Errorf("subnet can't be empty")
	}

	f := func(start string) (isDone bool, nextUrl string, err error) {
		listOptions := &vpcv1.ListLoadBalancersOptions{}
		if start != "" {
			listOptions.Start = &start
		}
		lbs, _, err := l.service.ListLoadBalancers(listOptions)
		if err != nil {
			return
		}
		if lbs == nil || len(lbs.LoadBalancers) == 0 {
			log.Println("no lbs found!")
			return
		}
		for _, lb := range lbs.LoadBalancers {
			for _, subnet := range lb.Subnets {
				if *subnet.CRN == *l.subnet.CRN {
					log.Printf("%s(lb) found with crn: %s", *lb.Name, *lb.CRN)
					err = l.deleteLoadbalancer(lb)
					if err != nil {
						return
					}
				}
			}
		}
		// For paging over next set of resources getting the start token and passing it for next iteration
		if lbs.Next != nil && *lbs.Next.Href != "" {
			nextUrl = *lbs.Next.Href
			return
		}

		isDone = true
		return
	}
	return utils.PagingHelper(f)
}

func (l *loadBalancer) deleteLoadbalancer(lb vpcv1.LoadBalancer) error {
	// Steps 1: Delete all the listeners
	if err := l.deleteListeners(lb); err != nil {
		return err
	}

	// Step 2: Delete all the backend pools
	for _, pool := range lb.Pools {
		if err := l.deletePool(lb, pool); err != nil {
			return err
		}
	}

	// Step 3: Delete the load balancer
	log.Printf("deleting the %s(LB) with ID: %s", *lb.Name, *lb.ID)
	deleteLBOpt := &vpcv1.DeleteLoadBalancerOptions{
		ID: lb.ID,
	}
	if err := l.waitForLBStatus(lb, "active"); err != nil {
		return err
	}
	if _, err := l.service.DeleteLoadBalancer(deleteLBOpt); err != nil {
		return err
	}

	// TODO: add a code to wait till LB disappears
	return nil
}

func (l *loadBalancer) deleteListeners(lb vpcv1.LoadBalancer) error {
	opt := &vpcv1.ListLoadBalancerListenersOptions{
		LoadBalancerID: lb.ID,
	}
	listeners, _, err := l.service.ListLoadBalancerListeners(opt)
	if err != nil {
		return err
	}
	if listeners == nil || len(listeners.Listeners) == 0 {
		log.Println("no pool listeners found!")
		return nil
	}
	for _, listener := range listeners.Listeners {
		log.Printf("removing listener with ID: %s", *listener.ID)
		if err := l.waitForLBStatus(lb, "active"); err != nil {
			return err
		}
		listenerOpt := &vpcv1.DeleteLoadBalancerListenerOptions{
			LoadBalancerID: lb.ID,
			ID:             listener.ID,
		}
		_, err := l.service.DeleteLoadBalancerListener(listenerOpt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *loadBalancer) deletePool(lb vpcv1.LoadBalancer, pool vpcv1.LoadBalancerPoolReference) error {
	opt := &vpcv1.ListLoadBalancerPoolMembersOptions{
		LoadBalancerID: lb.ID,
		PoolID:         pool.ID,
	}
	members, _, err := l.service.ListLoadBalancerPoolMembers(opt)
	if err != nil {
		return err
	}
	if members == nil || len(members.Members) == 0 {
		log.Println("no pool members found!")
		return nil
	}
	for _, member := range members.Members {
		if err := l.waitForLBStatus(lb, "active"); err != nil {
			return err
		}
		log.Printf("removing LB pool memeber with ID: %s", *member.ID)
		delopt := &vpcv1.DeleteLoadBalancerPoolMemberOptions{
			LoadBalancerID: lb.ID,
			PoolID:         pool.ID,
			ID:             member.ID,
		}
		if _, err := l.service.DeleteLoadBalancerPoolMember(delopt); err != nil {
			return err
		}
	}

	log.Printf("deleting the %s(LB pool) with ID: %s", *pool.Name, *pool.ID)
	if err := l.waitForLBStatus(lb, "active"); err != nil {
		return err
	}
	deleteLBopt := &vpcv1.DeleteLoadBalancerPoolOptions{
		LoadBalancerID: lb.ID,
		ID:             pool.ID,
	}
	if _, err := l.service.DeleteLoadBalancerPool(deleteLBopt); err != nil {
		return err
	}
	return nil
}

func (l *loadBalancer) waitForLBStatus(lb vpcv1.LoadBalancer, status string) error {
	f := func() (cond bool, err error) {
		opt := &vpcv1.GetLoadBalancerOptions{
			ID: lb.ID,
		}
		lb, _, err := l.service.GetLoadBalancer(opt)
		log.Printf("waiting for %s(lb) to be %s and current operating status is: %s", *lb.Name, status, *lb.ProvisioningStatus)

		if err != nil {
			return
		}

		if *lb.ProvisioningStatus == status {
			cond = true
			return
		}
		return
	}

	return wait.PollImmediate(20*time.Second, 5*time.Minute, f)
}

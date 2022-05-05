package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type (
	// Route53ResolverRule is the resource type
	Route53ResolverRule struct {
		svc        *route53resolver.Route53Resolver
		id         *string
		name       *string
		domainName *string
		vpcIds     []*string
	}
)

func init() {
	register("Route53ResolverRule", ListRoute53ResolverRules)
}

// ListRoute53ResolverRules produces the resources to be nuked.
func ListRoute53ResolverRules(sess *session.Session) ([]Resource, error) {
	svc := route53resolver.New(sess)

	vpcAssociations, err := resolverRulesToVpcIDs(svc)
	if err != nil {
		return nil, err
	}

	var resources []Resource

	params := &route53resolver.ListResolverRulesInput{}
	for {
		resp, err := svc.ListResolverRules(params)

		if err != nil {
			return nil, err
		}

		for _, rule := range resp.ResolverRules {
			resources = append(resources, &Route53ResolverRule{
				svc:        svc,
				id:         rule.Id,
				name:       rule.Name,
				domainName: rule.DomainName,
				vpcIds:     vpcAssociations[*rule.Id],
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

// Associate all the vpcIDs to their resolver rule ID to be disassociated before deleting the rule.
func resolverRulesToVpcIDs(svc *route53resolver.Route53Resolver) (map[string][]*string, error) {
	vpcAssociations := map[string][]*string{}

	params := &route53resolver.ListResolverRuleAssociationsInput{}

	for {
		resp, err := svc.ListResolverRuleAssociations(params)

		if err != nil {
			return nil, err
		}

		for _, ruleAssociation := range resp.ResolverRuleAssociations {
			vpcID := ruleAssociation.VPCId
			if vpcID != nil {
				resolverRuleID := *ruleAssociation.ResolverRuleId

				if _, ok := vpcAssociations[resolverRuleID]; !ok {
					vpcAssociations[resolverRuleID] = []*string{vpcID}
				} else {
					vpcAssociations[resolverRuleID] = append(vpcAssociations[resolverRuleID], vpcID)
				}
			}
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return vpcAssociations, nil
}

// Filter removes resources automatically from being nuked
func (r *Route53ResolverRule) Filter() error {
	if *r.domainName == "." {
		return fmt.Errorf(`Filtering DomainName "."`)
	}

	return nil
}

// Remove implements Resource
func (r *Route53ResolverRule) Remove() error {
	for _, vpcID := range r.vpcIds {
		_, err := r.svc.DisassociateResolverRule(&route53resolver.DisassociateResolverRuleInput{
			ResolverRuleId: r.id,
			VPCId:          vpcID,
		})

		if err != nil {
			return err
		}
	}

	_, err := r.svc.DeleteResolverRule(&route53resolver.DeleteResolverRuleInput{
		ResolverRuleId: r.id,
	})

	return err
}

// Properties provides debugging output
func (r *Route53ResolverRule) Properties() types.Properties {
	return types.NewProperties().
		Set("ID", r.id).
		Set("Name", r.name)
}

// String implements Stringer
func (r *Route53ResolverRule) String() string {
	return fmt.Sprintf("%s (%s)", *r.id, *r.name)
}

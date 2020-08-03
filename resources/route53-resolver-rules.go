package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/rebuy-de/aws-nuke/pkg/types"
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

	association struct {
		id    *string
		vpcID *string
	}
)

func init() {
	register("Route53ResolverRules", ListRout53ResolverRules)
}

// ListRout53ResolverRules produces the resources to be nuked.
func ListRout53ResolverRules(sess *session.Session) ([]Resource, error) {
	svc := route53resolver.New(sess)

	var resources []Resource
	output, err := svc.ListResolverRules(&route53resolver.ListResolverRulesInput{})

	if err != nil {
		return resources, err
	}

	associationsOutput, err := svc.ListResolverRuleAssociations(&route53resolver.ListResolverRuleAssociationsInput{})

	if err != nil {
		return resources, err
	}

	vpcAssociations := map[string][]*string{}
	for _, ruleAssociation := range associationsOutput.ResolverRuleAssociations {
		vpcId := ruleAssociation.VPCId
		if vpcId != nil {
			resolverRuleID := *ruleAssociation.ResolverRuleId

			if _, ok := vpcAssociations[resolverRuleID]; !ok {
				vpcAssociations[resolverRuleID] = []*string{vpcId}
			} else {
				vpcAssociations[resolverRuleID] = append(vpcAssociations[resolverRuleID], vpcId)
			}
		}
	}

	for _, rule := range output.ResolverRules {
		resources = append(resources, &Route53ResolverRule{
			svc:        svc,
			id:         rule.Id,
			name:       rule.Name,
			domainName: rule.DomainName,
			vpcIds:     vpcAssociations[*rule.Id],
		})
	}

	return resources, nil
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
		Set("Id", r.id).
		Set("Name", r.name)
}

// String implements Stringer
func (r *Route53ResolverRule) String() string {
	return fmt.Sprintf("%s (%s)", *r.id, *r.name)
}

package promrulesviewer

import (
	"context"
	"fmt"
	"os"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/prometheus/promql/parser"
)

func FetchPromRules(promAddress string) (*v1.RulesResult, error) {
	client, err := api.NewClient(api.Config{
		// Address: "https://prometheus.demo.do.prometheus.io/rules",
		Address: promAddress,
	})

	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	rules, err := v1api.Rules(context.Background())

	if err != nil {
		return nil, err
	}

	return &rules, nil
}

func extractQueryMatchers(query string) []string {
	expr, _ := parser.ParseExpr(query)

	labelMatchers := parser.ExtractSelectors(expr)

	queryMatchersMap := map[string]bool{}
	result := []string{}

	for _, lms := range labelMatchers {
		for _, lm := range lms {
			if lm.Name == "__name__" {
				queryMatchersMap[lm.Value] = true
			}
		}
	}

	for k := range queryMatchersMap {
		result = append(result, k)
	}

	return result
}

// IndexByRule indexes the raw rules result by rules and associated expressions
func IndexByDirectExprAssociation(rules *v1.RulesResult) (RulesWithAssociatedExpressions, error) {
	rulesWithAssociatedExpressions := RulesWithAssociatedExpressions{}

	for _, group := range rules.Groups {
		for _, r := range group.Rules {
			var associatedExpressions []string

			switch v := r.(type) {
			case v1.RecordingRule:
				associatedExpressions = extractQueryMatchers(r.(v1.RecordingRule).Query)
			case v1.AlertingRule:
				associatedExpressions = extractQueryMatchers(r.(v1.AlertingRule).Query)
			default:
				fmt.Printf("unknown rule type %s", v)
				return nil, fmt.Errorf("unknown rule type %s", v)
			}

			rulesWithAssociatedExpressions[&RuleAndGroup{
				Rule:  r,
				Group: group,
			}] = associatedExpressions
		}
	}

	return rulesWithAssociatedExpressions, nil
}

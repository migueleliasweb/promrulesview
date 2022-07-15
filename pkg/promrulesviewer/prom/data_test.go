package promrulesviewer

import (
	"testing"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

var rules = v1.RulesResult{
	Groups: []v1.RuleGroup{
		{
			Name:  "Group1",
			Rules: v1.Rules{},
		},
	},
}

func TestIndexByDirectExprAssociation(t *testing.T) {
	IndexByDirectExprAssociation(&rules)
}

package promrulesviewer

import v1 "github.com/prometheus/client_golang/api/prometheus/v1"

type RuleAndGroup struct {
	Rule  interface{} // v1.RecordingRule or v1.AlertingRule
	Group v1.RuleGroup
}

// RulesWithAssociatesRules models a single RuleAndGroup to its associates rules
//
// TODO: consider specializing []*RuleAndGroup to a more memory efficient struct
// with a single group and a slice of Rules
type RulesWithAssociatedExpressions map[*RuleAndGroup][]string

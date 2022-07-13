package flowchart

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/prometheus/promql/parser"
)

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

func generateData() map[string][]string {
	client, err := api.NewClient(api.Config{
		Address: "https://prometheus.demo.do.prometheus.io/rules",
	})

	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rules, _ := v1api.Rules(ctx)

	// need to add group name as metrics can repeat in different groups
	ruleNameWithRelatedQueryNames := map[string][]string{}

	for _, group := range rules.Groups {
		// groupName := group.Name
		for _, r := range group.Rules {
			switch v := r.(type) {
			case v1.RecordingRule:
				rule := r.(v1.RecordingRule)
				ruleNameWithRelatedQueryNames[rule.Name] = extractQueryMatchers(rule.Query)
			case v1.AlertingRule:
				rule := r.(v1.AlertingRule)
				ruleNameWithRelatedQueryNames[rule.Name] = extractQueryMatchers(rule.Query)
			default:
				fmt.Printf("unknown rule type %s", v)
			}
		}
	}

	return ruleNameWithRelatedQueryNames
}

func outputMetricRecursive(data map[string][]string, metricName string) string {
	relatedMetrics, found := data[metricName]

	if !found {
		return ""
	}

	var b bytes.Buffer

	for _, relatedMetric := range relatedMetrics {
		b.WriteString(fmt.Sprintf("    %s-->%s\n",
			relatedMetric,
			metricName,
		))

		b.WriteString(outputMetricRecursive(data, relatedMetric))
	}

	return b.String()
}

func TestFakeMain(t *testing.T) {
	data := generateData()

	var b bytes.Buffer

	b.WriteString("flowchart TB\n")

	topMetricName := "Foo"

	b.WriteString(outputMetricRecursive(data, topMetricName))

	fmt.Println(b.String())
}

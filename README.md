# PROMRULESVIEWER

A tool to help Prometheus rules (rules can be Alerts & RecordingRules) visualisation

## Note

- Recording rules can be composited by other recording rules and normal metrics
- Alerts can be composited by recording rules and normal metrics
- Metrics cannot be composited by any rules nor other metrics **directly\*** (that's when you should use recording rules)

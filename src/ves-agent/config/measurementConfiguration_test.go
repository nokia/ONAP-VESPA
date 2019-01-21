package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricRuleHasLabel(t *testing.T) {
	rule := MetricRule{
		Expr:      "foobar",
		Target:    "target",
		VMIDLabel: "id",
		Labels: []Label{
			{Name: "label-1", Expr: "foobar1"},
			{Name: "label-2", Expr: "foobar2"},
		},
	}
	assert.True(t, rule.hasLabel("label-1"))
	assert.True(t, rule.hasLabel("label-2"))
	assert.False(t, rule.hasLabel("label-3"))
}

func TestMetricRuleWithDefaults(t *testing.T) {
	defMetric := MetricRule{
		Target:    "defautlTarget",
		VMIDLabel: "id",
		Labels: []Label{
			{Name: "label-2", Expr: "foobar2"},
		},
	}
	rule := MetricRule{
		Expr:   "foobar",
		Target: "target",
		Labels: []Label{
			{Name: "label-1", Expr: "foobar1"},
		},
	}

	newRule := rule.WithDefaults(&defMetric)

	assert.Len(t, defMetric.Labels, 1)
	assert.Len(t, rule.Labels, 1)
	assert.Len(t, newRule.Labels, 2)
	assert.True(t, newRule.hasLabel("label-1"))
	assert.True(t, newRule.hasLabel("label-2"))
	assert.Equal(t, "foobar", newRule.Expr)
	assert.Equal(t, "target", newRule.Target)
	assert.Equal(t, "id", newRule.VMIDLabel)

	rule.Target = ""
	newRule = rule.WithDefaults(&defMetric)
	assert.Equal(t, "defautlTarget", newRule.Target)

	newRule = rule.WithDefaults(nil)
	assert.Equal(t, rule, newRule)

}

/*
	Copyright 2019 Nokia

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

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

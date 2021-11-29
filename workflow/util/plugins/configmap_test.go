package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/argoproj/argo-workflows/v3/pkg/plugins/spec"
	"github.com/argoproj/argo-workflows/v3/workflow/common"
)

func TestToConfigMap(t *testing.T) {
	cm, err := ToConfigMap(&spec.Plugin{
		TypeMeta: metav1.TypeMeta{
			Kind: "ExecutorPlugin",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-plug",
			Annotations: map[string]string{
				"my-anno": "my-value",
			},
			Labels: map[string]string{
				"my-label": "my-value",
			},
		},
		Spec: spec.PluginSpec{
			Address: "http://localhost:1234",
		},
	})
	if assert.NoError(t, err) {
		assert.Equal(t, "my-plug-executor-plugin", cm.Name)
		assert.Len(t, cm.Annotations, 1)
		assert.Equal(t, map[string]string{
			"my-label":                             "my-value",
			"workflows.argoproj.io/configmap-type": "ExecutorPlugin",
		}, cm.Labels)
	}
}

func TestFromConfigMap(t *testing.T) {
	p, err := FromConfigMap(&apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-plug-executor-plugin",
			Annotations: map[string]string{
				"my-anno": "my-value",
			},
			Labels: map[string]string{
				common.LabelKeyConfigMapType: "ExecutorPlugin",
				"my-label":                   "my-value",
			},
		},
		Data: map[string]string{
			"address":   "http://my-addr",
			"container": "{'name': 'my-name'}",
		},
	})
	if assert.NoError(t, err) {
		assert.Equal(t, "ExecutorPlugin", p.Kind)
		assert.Equal(t, "my-plug", p.Name)
		assert.Len(t, p.Annotations, 1)
		assert.Len(t, p.Labels, 1)
		assert.Equal(t, "http://my-addr", p.Spec.Address)
		assert.Equal(t, apiv1.Container{Name: "my-name"}, p.Spec.Container)
	}
}
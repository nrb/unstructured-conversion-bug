package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// TestFromUnstructuredIntToFloatBug tests for a bug where runtime.DefaultUnstructuredConverter.FromUnstructured can't take a whole number into a float.
// This test should fail when the bug is fixed upstream, letting us know we can remove the IsUnstructuredCRDReady function.
func TestFromUnstructuredIntToFloatBug(t *testing.T) {
	b := []byte(`
{
	"apiVersion": "apiextensions.k8s.io/v1beta1",
	"kind": "CustomResourceDefinition",
	"metadata": {
	  "name": "foos.example.foo.com"
	},
	"spec": {
	  "group": "example.foo.com",
	  "version": "v1alpha1",
	  "scope": "Namespaced",
	  "names": {
		"plural": "foos",
		"singular": "foo",
		"kind": "Foo"
	  },
	  "validation": {
		"openAPIV3Schema": {
		  "required": [
			"spec"
		  ],
		  "properties": {
			"spec": {
			  "required": [
				"bar"
			  ],
			  "properties": {
				"bar": {
				  "type": "integer",
				  "minimum": 1
				}
			  }
			}
		  }
		}
	  }
	}
  }
`)

	var obj unstructured.Unstructured
	err := json.Unmarshal(b, &obj)
	require.NoError(t, err)

	var newCRD apiextv1beta1.CustomResourceDefinition
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &newCRD)
	// If there's no error, then the upstream issue is fixed, and we need to remove our workarounds.
	require.NoError(t, err)
}

package printer

import (
	"fmt"
	"os"

	"github.com/aofekiko/kubectl-undo/pkg/logger"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	log = logger.NewLogger()
)

func PrintUnstructured(outputFormat string, resource *unstructured.Unstructured) {
	var serializer *json.Serializer
	switch outputFormat {
	case "yaml":
		serializer = json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	case "json":
		serializer = json.NewSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)
	case "none":
		return
	}
	err := serializer.Encode(resource, os.Stdout)
	if err != nil {
		log.Info(fmt.Sprintf("failed to encode objects: %v\n", err))
	}
}

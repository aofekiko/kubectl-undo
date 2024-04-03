package printer

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/aofekiko/kubectl-undo/pkg/logger"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	log = logger.NewLogger()
)

func PrintUnstructured(outputFormat string, resource *unstructured.Unstructured) (string, error) {
	var serializer *json.Serializer
	switch outputFormat {
	case "yaml":
		serializer = json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	case "json":
		serializer = json.NewSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)
	case "none":
		return "", errors.New("output format of none cannot be printed")
	}

	var buffer bytes.Buffer
	err := serializer.Encode(resource, &buffer)
	if err != nil {
		log.Info(fmt.Sprintf("failed to encode objects: %v\n", err))
	}
	return buffer.String(), nil
}

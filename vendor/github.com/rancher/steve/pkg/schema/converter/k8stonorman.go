package converter

import (
	"fmt"
	"strings"

	"github.com/rancher/steve/pkg/schemaserver/types"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/apiextensions.k8s.io/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

func GVKToSchemaID(gvk schema.GroupVersionKind) string {
	if gvk.Group == "" {
		return strings.ToLower(fmt.Sprintf("core.%s.%s", gvk.Version, gvk.Kind))
	}
	return strings.ToLower(fmt.Sprintf("%s.%s.%s", gvk.Group, gvk.Version, gvk.Kind))
}

func GVRToPluralName(gvr schema.GroupVersionResource) string {
	if gvr.Group == "" {
		return fmt.Sprintf("core.%s.%s", gvr.Version, gvr.Resource)
	}
	return fmt.Sprintf("%s.%s.%s", gvr.Group, gvr.Version, gvr.Resource)
}

func ToSchemas(crd v1beta1.CustomResourceDefinitionClient, client discovery.DiscoveryInterface) (map[string]*types.APISchema, error) {
	result := map[string]*types.APISchema{}

	if err := AddOpenAPI(client, result); err != nil {
		return nil, err
	}

	if err := AddDiscovery(client, result); err != nil {
		return nil, err
	}

	if err := AddCustomResources(crd, result); err != nil {
		return nil, err
	}

	return result, nil
}

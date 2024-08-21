package ingress2gateway

import (
	"fmt"
	"testing"
)

func TestToGatewayResource(t *testing.T) {
	ipath := "../../tests/httproute-canary-header.yaml"
	gwResource, err := ToGatewayResource(ipath)
	if err != nil {
		t.Fatalf("failed to convert ingress to gateway resource: %v", err)
	}
	// t.Logf("🧳 Ingress: %s translate to Gateway resource successfully", ipath)
	fmt.Println("output: ", gwResource)
}

package megaport

import (
	"context"
	"fmt"
	"testing"
)

const (
	TEST_LOCATION_ID_A = 19 // 	Interactive 437 Williamstown
)

func TestListPorts(t *testing.T) {
	setup()
	defer teardown()
	ctx := context.Background()
	portRes, err := megaportClient.PortService.ListPorts(ctx)
	if err != nil {
		t.Fatalf("error listing ports: %s", err.Error())
	}
	fmt.Println("port response", portRes)
}

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testDroneUser string = os.Getenv("DRONE_USER")
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"drone": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DRONE_SERVER"); v == "" {
		t.Fatal("DRONE_SERVER must be set for acceptance tests")
	}
	if v := os.Getenv("DRONE_TOKEN"); v == "" {
		t.Fatal("DRONE_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("DRONE_USER"); v == "" {
		t.Fatal("DRONE_USER must be set for acceptance tests")
	}
}

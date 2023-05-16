package pkg

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	config := terraform.NewResourceConfigRaw(map[string]interface{}{})
	testAccProvider.Configure(context.Background(), config)
	testAccProviders = map[string]*schema.Provider{
		"redis": testAccProvider,
	}

	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"redis": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("REDIS_ADDR"); v == "" {
		t.Fatal("REDIS_ADDR must be set for acceptance tests")
	}
}

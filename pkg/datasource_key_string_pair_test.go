package pkg

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testCheckDataSourceRedisKeyStringPairExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource '%s' not found in state", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource '%s' ID is not set", resourceName)
		}

		return nil
	}
}

func renderDataSourceKeyStringPair(key string) string {
	return fmt.Sprintf(`
data "redis_key_string_pair" "test" {
  key = "%s"
}
	`, key)
}

func renderKeyStringPair(key, value string) string {
	return fmt.Sprintf(`
resource "redis_key_string_pair" "foo" {
  key   = "%s"
  value = "%s"
}
	`, key, value)
}

func TestStringDataSource(t *testing.T) {
	key := "foo"
	expectedValue := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: renderKeyStringPair(key, expectedValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyStringPairExists("redis_key_string_pair.foo"),
					resource.TestCheckResourceAttr("redis_key_string_pair.foo", "key", key),
				),
			},
			{
				Config: renderDataSourceKeyStringPair(key),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataSourceRedisKeyStringPairExists("data.redis_key_string_pair.test"),
					resource.TestCheckResourceAttr("data.redis_key_string_pair.test", "key", key),
					resource.TestCheckResourceAttr("data.redis_key_string_pair.test", "value", expectedValue),
				),
			},
		},
	})
}

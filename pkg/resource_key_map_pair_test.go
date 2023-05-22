package pkg

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// <Common>
func testAccDeleteKeyMapPair(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "redis_key_map_pair" {
			continue
		}

		client := testAccProvider.Meta().(*Client).goRedisClient()

		exists, err := client.Exists(context.Background(), rs.Primary.ID).Result()
		if err != nil {
			return fmt.Errorf("Redis Error: %s", err.Error())
		}
		if exists == 1 {
			return fmt.Errorf("Key still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckkeyMapPairExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		client := testAccProvider.Meta().(*Client).goRedisClient()

		exists, err := client.Exists(context.Background(), rs.Primary.ID).Result()
		if err != nil {
			return fmt.Errorf("Redis Error: %s", err.Error())
		}
		if exists != 1 {
			return fmt.Errorf("Key Does Not Exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

// </Common>
// <Without Expiry>
type testKeyMapPairConfig struct {
	key   string
	value map[string]string
}

func (c *testKeyMapPairConfig) render() string {
	modifiedStrings := make([]string, len(c.value))
	for key, val := range c.value {
		newVal := `"` + key + `": "` + val + `"`
		modifiedStrings = append(modifiedStrings, newVal)
	}
	val := strings.Join(modifiedStrings, "\n")

	result := fmt.Sprintf(`
resource "redis_key_map_pair" "foo" {
  key   = "%s"
  value = {%s}
}
		`, c.key, val)
	return result
}

func TestKeyMapPair(t *testing.T) {
	val := map[string]string{
		acctest.RandString(5):  acctest.RandString(5),
		acctest.RandString(1):  acctest.RandString(1),
		acctest.RandString(10): acctest.RandString(10),
	}

	c := testKeyMapPairConfig{
		key:   "foo",
		value: val,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccDeleteKeyMapPair,
		Steps: []resource.TestStep{
			{
				Config: c.render(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckkeyMapPairExists("redis_key_map_pair.foo"),
					resource.TestCheckResourceAttr("redis_key_map_pair.foo", "key", c.key),
					resource.TestCheckResourceAttr("redis_key_map_pair.foo", "value.%", strconv.Itoa(len(c.value))),
				),
			},
		},
	})
}

// </Without Expiry>
// <With Expiry>
type testKeyMapPairWithExpiryConfig struct {
	key    string
	value  map[string]string
	expiry int // seconds
}

func (c *testKeyMapPairWithExpiryConfig) render() string {
	modifiedStrings := make([]string, len(c.value))
	for key, val := range c.value {
		newVal := `"` + key + `": "` + val + `"`
		modifiedStrings = append(modifiedStrings, newVal)
	}
	val := strings.Join(modifiedStrings, "\n")

	return fmt.Sprintf(`
resource "redis_key_map_pair" "foo" {
  key   = "%s"
  expiry = "%ds"
  value = {%s}
}
		`, c.key, c.expiry, val)
}

func testAccCheckkeyMapPairExpiry(n string, t int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		client := testAccProvider.Meta().(*Client).goRedisClient()

		pttl, err := client.PTTL(context.Background(), rs.Primary.ID).Result()
		if err != nil {
			return fmt.Errorf("Redis Error: %s", err.Error())
		}

		_t := int(pttl.Seconds())

		if _t > t || t < 0 {
			return fmt.Errorf("PTTL Error")
		}
		return nil
	}
}

func TestKeyMapPairWithExpiry(t *testing.T) {
	val := map[string]string{
		acctest.RandString(5):  acctest.RandString(5),
		acctest.RandString(1):  acctest.RandString(1),
		acctest.RandString(10): acctest.RandString(10),
	}

	c := testKeyMapPairWithExpiryConfig{
		key:    "foo",
		value:  val,
		expiry: 200,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: c.render(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckkeyMapPairExists("redis_key_map_pair.foo"),
					resource.TestCheckResourceAttr("redis_key_map_pair.foo", "key", c.key),
					testAccCheckkeyMapPairExpiry("redis_key_map_pair.foo", c.expiry),
					resource.TestCheckResourceAttr("redis_key_map_pair.foo", "value.%", strconv.Itoa(len(c.value))),
				),
			},
		},
	})
}

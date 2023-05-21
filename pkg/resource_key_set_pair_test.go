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
func testAccDeleteKeySetPair(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "redis_key_set_pair" {
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

func testAccCheckkeySetPairExists(n string) resource.TestCheckFunc {
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
type testKeySetPairConfig struct {
	key   string
	value []string
}

func (c *testKeySetPairConfig) render() string {
	modifiedStrings := make([]string, len(c.value))
	for i, str := range c.value {
		modifiedStrings[i] = `"` + str + `"`
	}
	val := strings.Join(modifiedStrings, ",")

	result := fmt.Sprintf(`
resource "redis_key_set_pair" "foo" {
  key   = "%s"
  value = [%s]
}
		`, c.key, val)
	return result
}

func TestKeySetPair(t *testing.T) {
	c := testKeySetPairConfig{
		key:   "foo",
		value: []string{acctest.RandString(5), acctest.RandString(1), acctest.RandString(10)},
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccDeleteKeySetPair,
		Steps: []resource.TestStep{
			{
				Config: c.render(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckkeySetPairExists("redis_key_list_pair.foo"),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "key", c.key),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "value.0", c.value[0]),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "value.#", strconv.Itoa(len(c.value))),
				),
			},
		},
	})
}

// </Without Expiry>
// <With Expiry>
type testKeySetPairWithExpiryConfig struct {
	key    string
	value  []string
	expiry int // seconds
}

func (c *testKeySetPairWithExpiryConfig) render() string {
	modifiedStrings := make([]string, len(c.value))
	for i, str := range c.value {
		modifiedStrings[i] = `"` + str + `"`
	}
	val := strings.Join(modifiedStrings, ",")

	return fmt.Sprintf(`
resource "redis_key_set_pair" "foo"{
  key 	 = "%s"
  value  = [%s]
  expiry = "%ds"
}
		`, c.key, val, c.expiry)
}

func testAccCheckkeySetPairExpiry(n string, t int) resource.TestCheckFunc {
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

func TestKeySetPairWithExpiry(t *testing.T) {
	c := testKeySetPairWithExpiryConfig{
		key:    "foo",
		value:  []string{acctest.RandString(5), acctest.RandString(1), acctest.RandString(10)},
		expiry: 200,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: c.render(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckkeySetPairExists("redis_key_list_pair.foo"),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "key", c.key),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "value.0", c.value[0]),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "value.#", strconv.Itoa(len(c.value))),
					testAccCheckkeySetPairExpiry("redis_key_list_pair.foo", c.expiry),
				),
			},
		},
	})
}

// </With Expiry>

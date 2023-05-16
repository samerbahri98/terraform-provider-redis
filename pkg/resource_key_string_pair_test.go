package pkg

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// <Common>

func testAccDeleteKeyStringPair(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "redis_key_string_pair" {
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

func testAccCheckKeyStringPairExists(n string) resource.TestCheckFunc {
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

type testKeyValuePairConfig struct {
	key   string
	value string
}

func (c *testKeyValuePairConfig) render() string {
	return fmt.Sprintf(`
		resource "redis_key_string_pair" "foo" {
			key = "%s"
			value = "%s"
		}
		`, c.key, c.value)
}

func TestKeyValuePair(t *testing.T) {
	c := testKeyValuePairConfig{
		key:   "foo",
		value: acctest.RandString(5),
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccDeleteKeyStringPair,
		Steps: []resource.TestStep{
			{
				Config: c.render(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyStringPairExists("redis_key_string_pair.foo"),
					resource.TestCheckResourceAttr("redis_key_string_pair.foo", "key", c.key),
					resource.TestCheckResourceAttr("redis_key_string_pair.foo", "value", c.value),
				),
			},
		},
	})
}

// </Without Expiry>
// <With Expiry>

type testKeyValuePairWithExpiryConfig struct {
	key    string
	value  string
	expiry int // seconds
}

func (c *testKeyValuePairWithExpiryConfig) render() string {
	return fmt.Sprintf(`
		resource "redis_key_string_pair" "foo"{
			key = "%s"
			value = "%s"
			expiry = "%ds"
		}
		`, c.key, c.value, c.expiry)
}

func testAccCheckKeyStringPairExpiry(n string, t int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
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

func TestKeyValuePairWithExpiry(t *testing.T) {
	c := testKeyValuePairWithExpiryConfig{
		key:    "foo",
		value:  acctest.RandString(5),
		expiry: 200,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: c.render(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyStringPairExists("redis_key_string_pair.foo"),
					resource.TestCheckResourceAttr("redis_key_string_pair.foo", "key", c.key),
					resource.TestCheckResourceAttr("redis_key_string_pair.foo", "value", c.value),
					testAccCheckKeyStringPairExpiry("redis_key_string_pair.foo", c.expiry),
				),
			},
		},
	})
}

// </With Expiry>

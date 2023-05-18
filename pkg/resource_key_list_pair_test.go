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
func testAccDeleteKeyListPair(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "redis_key_list_pair" {
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

func testAccCheckkeyListPairExists(n string) resource.TestCheckFunc {
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
type testKeyListPairConfig struct {
	key   string
	value []string
}

func (c *testKeyListPairConfig) render() string {
	modifiedStrings := make([]string, len(c.value))
	for i, str := range c.value {
		modifiedStrings[i] = `"` + str + `"`
	}
	val := strings.Join(modifiedStrings, ",")

	result := fmt.Sprintf(`
resource "redis_key_list_pair" "foo" {
  key   = "%s"
  value = [%s]
}
		`, c.key, val)
	return result
}

func TestKeyListPair(t *testing.T) {
	c := testKeyListPairConfig{
		key:   "foo",
		value: []string{acctest.RandString(5), acctest.RandString(1), acctest.RandString(10)},
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccDeleteKeyListPair,
		Steps: []resource.TestStep{
			{
				Config: c.render(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckkeyListPairExists("redis_key_list_pair.foo"),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "key", c.key),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "value.0", c.value[0]),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "value.#", strconv.Itoa(len(c.value))),
				),
			},
		},
	})
}

// </Without Expiry>
<<<<<<< HEAD
// <With Expiry>
type testKeyListPairWithExpiryConfig struct {
	key    string
	value  []string
	expiry int // seconds
}

func (c *testKeyListPairWithExpiryConfig) render() string {
	modifiedStrings := make([]string, len(c.value))
	for i, str := range c.value {
		modifiedStrings[i] = `"` + str + `"`
	}
	val := strings.Join(modifiedStrings, ",")

	return fmt.Sprintf(`
resource "redis_key_list_pair" "foo"{
  key 	 = "%s"
  value  = [%s]
  expiry = "%ds"
}
		`, c.key, val, c.expiry)
}

func testAccCheckkeyListPairExpiry(n string, t int) resource.TestCheckFunc {
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

func TestKeyListPairWithExpiry(t *testing.T) {
	c := testKeyListPairWithExpiryConfig{
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
					testAccCheckkeyListPairExists("redis_key_list_pair.foo"),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "key", c.key),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "value.0", c.value[0]),
					resource.TestCheckResourceAttr("redis_key_list_pair.foo", "value.#", strconv.Itoa(len(c.value))),
					testAccCheckkeyListPairExpiry("redis_key_list_pair.foo", c.expiry),
				),
			},
		},
	})
}

// </With Expiry>
=======
>>>>>>> 759c354 (Added: create and read key list pair)

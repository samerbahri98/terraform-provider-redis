package pkg

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("REDIS_ADDR", nil),
				Description: "Redis Address (And Port)",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REDIS_PASS", nil),
				Description: "Redis Password",
			},
			"db": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Redis Database",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"redis_key_string_pair": resourceKeyStringPair(),
			"redis_key_map_pair":    resourceKeyMapPair(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Addr:     d.Get("address").(string),
		Password: d.Get("password").(string),
		DB:       d.Get("db").(int),
	}

	return config.Client()
}

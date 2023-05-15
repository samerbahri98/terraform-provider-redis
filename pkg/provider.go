package pkg

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func Provider() *schema.Provider {
	return &schema.Provider{
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

package pkg

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceKeyStringPair() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceKeyStringPairRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// "expiry": {
			// 	Type:     schema.TypeString,
			// 	Default:  "0s",
			// 	Optional: true,
			// },
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceKeyStringPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	key := d.Get("key").(string)

	client := meta.(*Client).goRedisClient()

	value, _ := client.Get(ctx, key).Result()

	d.Set("value", value)
	d.SetId(key)

	return nil
}

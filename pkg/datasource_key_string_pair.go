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
			"expiry": {
				Type:     schema.TypeString,
				Default:  "0s",
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func datasourceKeyStringPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	if client.Exists(ctx, d.Get("key").(string)).Val() == 0 {
		return diag.Errorf("Key does not exist")
	}

	v := client.Get(ctx, d.Id())
	if v.Val() != d.Get("value").(string) {
		return diag.Errorf("Redis Error")
	}

	return nil
}

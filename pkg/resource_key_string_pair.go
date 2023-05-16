package pkg

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeyStringPair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeyStringPairCreate,
		ReadContext:   resourceKeyStringPairRead,
		UpdateContext: resourceKeyStringPairCreate,
		DeleteContext: resourceKeyStringPairDelete,
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

func resourceKeyStringPairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).goRedisClient()
	duration, _ := time.ParseDuration(d.Get("expiry").(string))
	key := d.Get("key").(string)
	d.SetId(key)
	client.Set(ctx, key, d.Get("value").(string), duration)

	return resourceKeyStringPairRead(ctx, d, meta)
}

func resourceKeyStringPairDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).goRedisClient()
	client.Del(ctx, d.Get("key").(string))
	return nil
}

func resourceKeyStringPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	k := d.Id()
	v := client.Get(ctx, k)
	if v.Val() != d.Get("value").(string) {
		return diag.Errorf("Redis Error")
	}

	return nil
}

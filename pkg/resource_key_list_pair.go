package pkg

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeyListPair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeyListPairCreate,
		ReadContext:   resourceKeyListPairRead,
		UpdateContext: resourceKeyListPairUpdate,
		DeleteContext: resourceKeyListPairDelete,
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
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	duration, _ := time.ParseDuration(d.Get("expiry").(string))
	key := d.Get("key").(string)

	d.SetId(key)
	client.Expire(ctx, key, duration)

	v := d.Get("value").([]interface{})
	stringSlice := make([]interface{}, len(v))
	for i, value := range v {
		stringSlice[i] = value.(string)
	}
	client.RPush(ctx, key, stringSlice)

	return resourceKeyListPairRead(ctx, d, meta)
}

func resourceKeyListPairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	if client.Exists(ctx, d.Get("key").(string)).Val() != 0 {
		return diag.Errorf("Key exists, use update instead")
	}

	return create(ctx, d, meta)
}

func resourceKeyListPairUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	if client.Exists(ctx, d.Get("key").(string)).Val() == 0 {
		return diag.Errorf("Key does not exist")
	}

	resourceKeyListPairDelete(ctx, d, meta)
	return create(ctx, d, meta)
}

func resourceKeyListPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	key := d.Id()
	val := d.Get("value").([]interface{})

	v, _ := client.LRange(ctx, key, 0, -1).Result()
	for i, value := range val {
		if v[i] != value.(string) {
			return diag.Errorf("Redis Error")
		}
	}

	return nil
}

func resourceKeyListPairDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	client.Del(ctx, d.Get("key").(string))
	return nil
}

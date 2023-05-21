package pkg

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeySetPair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeySetPairCreate,
		ReadContext:   resourceKeySetPairRead,
		UpdateContext: resourceKeySetPairUpdate,
		DeleteContext: resourceKeySetPairDelete,
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
				Type:     schema.TypeSet,
				Required: true,
			},
		},
	}
}

func createSet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	client.SAdd(ctx, key, stringSlice)

	return resourceKeySetPairRead(ctx, d, meta)
}

func resourceKeySetPairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	if client.Exists(ctx, d.Get("key").(string)).Val() != 0 {
		return diag.Errorf("Key exists, use update instead")
	}

	return createSet(ctx, d, meta)
}

func resourceKeySetPairUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	if client.Exists(ctx, d.Get("key").(string)).Val() == 0 {
		return diag.Errorf("Key does not exist")
	}

	resourceKeySetPairDelete(ctx, d, meta)
	return createSet(ctx, d, meta)
}

func resourceKeySetPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	key := d.Id()
	val := d.Get("value").([]interface{})

	v, _ := client.SMembers(ctx, key).Result()
	for i, value := range val {
		if v[i] != value.(string) {
			return diag.Errorf("Redis Error")
		}
	}

	return nil
}

func resourceKeySetPairDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	client.Del(ctx, d.Get("key").(string))
	return nil
}

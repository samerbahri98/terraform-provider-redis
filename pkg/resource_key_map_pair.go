package pkg

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeyMapPair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeyMapPairCreate,
		ReadContext:   resourceKeyMapPairRead,
		UpdateContext: resourceKeyMapPairUpdate,
		DeleteContext: resourceKeyMapPairDelete,
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
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func createMap(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	duration, _ := time.ParseDuration(d.Get("expiry").(string))
	key := d.Get("key").(string)

	d.SetId(key)
	client.Expire(ctx, key, duration)

	v := d.Get("value").([]interface{})
	// stringSlice := make([]interface{}, len(v))
	// for i, value := range v {
	// 	stringSlice[i] = value.(string)
	// }
	client.HSet(ctx, key, v)

	return resourceKeyMapPairRead(ctx, d, meta)
}

func resourceKeyMapPairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	if client.Exists(ctx, d.Get("key").(string)).Val() != 0 {
		return diag.Errorf("Key exists, use update instead")
	}

	return createMap(ctx, d, meta)
}

func resourceKeyMapPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	key := d.Id()
	val := d.Get("value").([]interface{})

	v, _ := client.HGetAll(ctx, key).Result()
	if len(v) != len(val) {
		return diag.Errorf("Redis Error")
	}

	return nil
}

func resourceKeyMapPairUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	if client.Exists(ctx, d.Get("key").(string)).Val() == 0 {
		return diag.Errorf("Key does not exist")
	}

	resourceKeyMapPairDelete(ctx, d, meta)
	return createMap(ctx, d, meta)
}

func resourceKeyMapPairDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).goRedisClient()
	client.Del(ctx, d.Get("key").(string))
	return nil
}

package pkg

import (
	"context"

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

func resourceKeySetPairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceKeySetPairUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceKeySetPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceKeySetPairDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

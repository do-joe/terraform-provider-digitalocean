package database

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanDatabaseClusterMetricsCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanDatabaseClusterMetricsCredentialsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceDigitalOceanDatabaseClusterMetricsCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	// Note: The GetMetricsCredentials method doesn't take a cluster ID parameter
	// as it returns global metrics credentials for the account
	creds, _, err := client.Databases.GetMetricsCredentials(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving database metrics credentials: %s", err)
	}

	d.SetId(clusterID)
	d.Set("username", creds.BasicAuthUsername)
	d.Set("password", creds.BasicAuthPassword)

	return nil
}

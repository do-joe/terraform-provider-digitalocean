package database

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseLogsinkOpensearch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseLogsinkOpensearchCreate,
		ReadContext:   resourceDigitalOceanDatabaseLogsinkOpensearchRead,
		UpdateContext: resourceDigitalOceanDatabaseLogsinkOpensearchUpdate,
		DeleteContext: resourceDigitalOceanDatabaseLogsinkOpensearchDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanDatabaseLogsinkOpensearchImport,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "UUID of the source database cluster that will forward logs",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Display name for the logsink",
			},
			"endpoint": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateHTTPSEndpoint,
				Description:  "HTTPS URL to OpenSearch (https://host:port)",
			},
			"index_prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIndexPrefix,
				Description:  "Prefix for OpenSearch indices",
			},
			"index_days_max": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIndexDaysMax,
				Description:  "Maximum number of days to retain indices (>= 1)",
			},
			"ca_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "CA certificate for TLS verification (PEM format)",
			},
			"timeout_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateLogsinkTimeout,
				Description:  "Request timeout for log deliveries in seconds (>= 1)",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Composite ID of the logsink resource",
			},
			"logsink_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API sink_id returned by DigitalOcean",
			},
		},

		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("name", func(_ context.Context, old, new, meta interface{}) bool {
				// Force recreation if name changes
				return old.(string) != new.(string)
			}),
		),
	}
}

func resourceDigitalOceanDatabaseLogsinkOpensearchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	req := buildCreateLogsinkRequest(d, "opensearch")

	log.Printf("[DEBUG] Database logsink opensearch create configuration: %#v", req)
	logsink, _, err := client.Databases.CreateLogsink(context.Background(), clusterID, req)
	if err != nil {
		return diag.Errorf("Error creating database logsink opensearch: %s", err)
	}

	d.SetId(createLogsinkID(clusterID, logsink.ID))
	log.Printf("[INFO] Database logsink opensearch ID: %s", logsink.ID)

	// Post-create read for consistency
	return resourceDigitalOceanDatabaseLogsinkOpensearchRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkOpensearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID, logsinkID := splitLogsinkID(d.Id())

	if clusterID == "" || logsinkID == "" {
		return diag.Errorf("Invalid logsink ID format: %s", d.Id())
	}

	logsink, resp, err := client.Databases.GetLogsink(context.Background(), clusterID, logsinkID)
	if err != nil {
		// If the logsink is somehow already destroyed, mark as successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving database logsink opensearch: %s", err)
	}

	d.Set("cluster_id", clusterID)

	err = setLogsinkResourceData(d, logsink, "opensearch")
	if err != nil {
		return diag.Errorf("Error setting logsink resource data: %s", err)
	}

	return nil
}

func resourceDigitalOceanDatabaseLogsinkOpensearchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID, logsinkID := splitLogsinkID(d.Id())

	if clusterID == "" || logsinkID == "" {
		return diag.Errorf("Invalid logsink ID format: %s", d.Id())
	}

	req := buildUpdateLogsinkRequest(d, "opensearch")

	log.Printf("[DEBUG] Database logsink opensearch update configuration: %#v", req)
	_, err := client.Databases.UpdateLogsink(context.Background(), clusterID, logsinkID, req)
	if err != nil {
		return diag.Errorf("Error updating database logsink opensearch: %s", err)
	}

	// Re-read the resource to refresh state
	return resourceDigitalOceanDatabaseLogsinkOpensearchRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkOpensearchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID, logsinkID := splitLogsinkID(d.Id())

	if clusterID == "" || logsinkID == "" {
		return diag.Errorf("Invalid logsink ID format: %s", d.Id())
	}

	log.Printf("[INFO] Deleting database logsink opensearch: %s", d.Id())
	_, err := client.Databases.DeleteLogsink(context.Background(), clusterID, logsinkID)
	if err != nil {
		// Treat 404 as success (already removed)
		if godoErr, ok := err.(*godo.ErrorResponse); ok && godoErr.Response.StatusCode == 404 {
			log.Printf("[INFO] Database logsink opensearch %s was already deleted", d.Id())
		} else {
			return diag.Errorf("Error deleting database logsink opensearch: %s", err)
		}
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanDatabaseLogsinkOpensearchImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Validate the import ID format
	clusterID, logsinkID := splitLogsinkID(d.Id())
	if clusterID == "" || logsinkID == "" {
		return nil, fmt.Errorf("must use the format 'cluster_id,logsink_id' for import (e.g. 'deadbeef-dead-4aa5-beef-deadbeef347d,01234567-89ab-cdef-0123-456789abcdef')")
	}

	// The Read function will handle populating all fields from the API
	return []*schema.ResourceData{d}, nil
}

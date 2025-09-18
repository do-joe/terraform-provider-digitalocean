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

func ResourceDigitalOceanDatabaseLogsinkRsyslog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseLogsinkRsyslogCreate,
		ReadContext:   resourceDigitalOceanDatabaseLogsinkRsyslogRead,
		UpdateContext: resourceDigitalOceanDatabaseLogsinkRsyslogUpdate,
		DeleteContext: resourceDigitalOceanDatabaseLogsinkRsyslogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanDatabaseLogsinkRsyslogImport,
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
			"server": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Hostname or IP address of the rsyslog server",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateLogsinkPort,
				Description:  "Port number for the rsyslog server (1-65535)",
			},
			"tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable TLS encryption for rsyslog connection",
			},
			"format": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "rfc5424",
				ValidateFunc: validateRsyslogFormat,
				Description:  "Log format: rfc5424, rfc3164, or custom",
			},
			"logline": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom logline template (required when format is 'custom')",
			},
			"structured_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Structured data for rsyslog",
			},
			"ca_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "CA certificate for TLS verification (PEM format)",
			},
			"client_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client certificate for mTLS (PEM format)",
			},
			"client_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Client private key for mTLS (PEM format)",
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
			func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return validateLogsinkCustomDiff(diff, "rsyslog")
			},
		),
	}
}

func resourceDigitalOceanDatabaseLogsinkRsyslogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	req := buildCreateLogsinkRequest(d, "rsyslog")

	log.Printf("[DEBUG] Database logsink rsyslog create configuration: %#v", req)
	logsink, _, err := client.Databases.CreateLogsink(context.Background(), clusterID, req)
	if err != nil {
		return diag.Errorf("Error creating database logsink rsyslog: %s", err)
	}

	log.Printf("[DEBUG] API Response logsink: %#v", logsink)
	log.Printf("[DEBUG] Logsink ID: '%s'", logsink.ID)
	log.Printf("[DEBUG] Logsink Name: '%s'", logsink.Name)
	log.Printf("[DEBUG] Logsink Type: '%s'", logsink.Type)

	d.SetId(createLogsinkID(clusterID, logsink.ID))
	log.Printf("[INFO] Database logsink rsyslog ID: %s", logsink.ID)

	// Post-create read for consistency
	return resourceDigitalOceanDatabaseLogsinkRsyslogRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkRsyslogRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("Error retrieving database logsink rsyslog: %s", err)
	}

	d.Set("cluster_id", clusterID)

	err = setLogsinkResourceData(d, logsink, "rsyslog")
	if err != nil {
		return diag.Errorf("Error setting logsink resource data: %s", err)
	}

	return nil
}

func resourceDigitalOceanDatabaseLogsinkRsyslogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID, logsinkID := splitLogsinkID(d.Id())

	if clusterID == "" || logsinkID == "" {
		return diag.Errorf("Invalid logsink ID format: %s", d.Id())
	}

	req := buildUpdateLogsinkRequest(d, "rsyslog")

	log.Printf("[DEBUG] Database logsink rsyslog update configuration: %#v", req)
	_, err := client.Databases.UpdateLogsink(context.Background(), clusterID, logsinkID, req)
	if err != nil {
		return diag.Errorf("Error updating database logsink rsyslog: %s", err)
	}

	// Re-read the resource to refresh state
	return resourceDigitalOceanDatabaseLogsinkRsyslogRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkRsyslogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID, logsinkID := splitLogsinkID(d.Id())

	if clusterID == "" || logsinkID == "" {
		return diag.Errorf("Invalid logsink ID format: %s", d.Id())
	}

	log.Printf("[INFO] Deleting database logsink rsyslog: %s", d.Id())
	_, err := client.Databases.DeleteLogsink(context.Background(), clusterID, logsinkID)
	if err != nil {
		// Treat 404 as success (already removed)
		if godoErr, ok := err.(*godo.ErrorResponse); ok && godoErr.Response.StatusCode == 404 {
			log.Printf("[INFO] Database logsink rsyslog %s was already deleted", d.Id())
		} else {
			return diag.Errorf("Error deleting database logsink rsyslog: %s", err)
		}
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanDatabaseLogsinkRsyslogImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Validate the import ID format
	clusterID, logsinkID := splitLogsinkID(d.Id())
	if clusterID == "" || logsinkID == "" {
		return nil, fmt.Errorf("must use the format 'cluster_id,logsink_id' for import (e.g. 'deadbeef-dead-4aa5-beef-deadbeef347d,01234567-89ab-cdef-0123-456789abcdef')")
	}

	// The Read function will handle populating all fields from the API
	return []*schema.ResourceData{d}, nil
}

package database

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// createLogsinkID creates a composite ID for logsink resources
// Format: <cluster_id>,<logsink_id>
func createLogsinkID(clusterID string, logsinkID string) string {
	return fmt.Sprintf("%s,%s", clusterID, logsinkID)
}

// splitLogsinkID splits a composite logsink ID into cluster ID and logsink ID
// Splits on the first comma to handle edge cases
func splitLogsinkID(id string) (string, string) {
	parts := strings.SplitN(id, ",", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// expandLogsinkConfigRsyslog converts Terraform schema data to godo.DatabaseLogsinkConfig for rsyslog
func expandLogsinkConfigRsyslog(d *schema.ResourceData) *godo.DatabaseLogsinkConfig {
	config := &godo.DatabaseLogsinkConfig{}

	if v, ok := d.GetOk("server"); ok {
		config.Server = v.(string)
	}
	if v, ok := d.GetOk("port"); ok {
		config.Port = v.(int)
	}
	if v, ok := d.GetOk("tls"); ok {
		config.TLS = v.(bool)
	}
	if v, ok := d.GetOk("format"); ok {
		config.Format = v.(string)
	}
	if v, ok := d.GetOk("logline"); ok {
		config.Logline = v.(string)
	}
	if v, ok := d.GetOk("structured_data"); ok {
		config.SD = v.(string)
	}
	if v, ok := d.GetOk("ca_cert"); ok {
		config.CA = trimPEMString(v.(string))
	}
	if v, ok := d.GetOk("client_cert"); ok {
		config.Cert = trimPEMString(v.(string))
	}
	if v, ok := d.GetOk("client_key"); ok {
		config.Key = trimPEMString(v.(string))
	}
	if v, ok := d.GetOk("timeout_seconds"); ok {
		config.Timeout = float32(v.(int))
	}

	return config
}

// expandLogsinkConfigElasticsearch converts Terraform schema data to godo.DatabaseLogsinkConfig for elasticsearch
func expandLogsinkConfigElasticsearch(d *schema.ResourceData) *godo.DatabaseLogsinkConfig {
	config := &godo.DatabaseLogsinkConfig{}

	if v, ok := d.GetOk("endpoint"); ok {
		config.URL = v.(string)
	}
	if v, ok := d.GetOk("index_prefix"); ok {
		config.IndexPrefix = v.(string)
	}
	if v, ok := d.GetOk("index_days_max"); ok {
		config.IndexDaysMax = v.(int)
	}
	if v, ok := d.GetOk("ca_cert"); ok {
		config.CA = trimPEMString(v.(string))
	}
	if v, ok := d.GetOk("timeout_seconds"); ok {
		config.Timeout = float32(v.(int))
	}

	return config
}

// expandLogsinkConfigOpensearch converts Terraform schema data to godo.DatabaseLogsinkConfig for opensearch
func expandLogsinkConfigOpensearch(d *schema.ResourceData) *godo.DatabaseLogsinkConfig {
	config := &godo.DatabaseLogsinkConfig{}

	if v, ok := d.GetOk("endpoint"); ok {
		config.URL = v.(string)
	}
	if v, ok := d.GetOk("index_prefix"); ok {
		config.IndexPrefix = v.(string)
	}
	if v, ok := d.GetOk("index_days_max"); ok {
		config.IndexDaysMax = v.(int)
	}
	if v, ok := d.GetOk("ca_cert"); ok {
		config.CA = trimPEMString(v.(string))
	}
	if v, ok := d.GetOk("timeout_seconds"); ok {
		config.Timeout = float32(v.(int))
	}

	return config
}

// flattenLogsinkConfigRsyslog converts godo.DatabaseLogsinkConfig to Terraform schema data for rsyslog
func flattenLogsinkConfigRsyslog(d *schema.ResourceData, config *godo.DatabaseLogsinkConfig) error {
	if config == nil {
		return nil
	}

	if config.Server != "" {
		d.Set("server", config.Server)
	}
	if config.Port != 0 {
		d.Set("port", config.Port)
	}
	d.Set("tls", config.TLS)
	if config.Format != "" {
		d.Set("format", config.Format)
	}
	if config.Logline != "" {
		d.Set("logline", config.Logline)
	}
	if config.SD != "" {
		d.Set("structured_data", config.SD)
	}
	if config.CA != "" {
		d.Set("ca_cert", trimPEMString(config.CA))
	}
	if config.Cert != "" {
		d.Set("client_cert", trimPEMString(config.Cert))
	}
	// Preserve sensitive client_key from prior state if not returned by API
	if config.Key != "" {
		d.Set("client_key", trimPEMString(config.Key))
	}
	if config.Timeout != 0 {
		d.Set("timeout_seconds", int(config.Timeout))
	}

	return nil
}

// flattenLogsinkConfigElasticsearch converts godo.DatabaseLogsinkConfig to Terraform schema data for elasticsearch
func flattenLogsinkConfigElasticsearch(d *schema.ResourceData, config *godo.DatabaseLogsinkConfig) error {
	if config == nil {
		return nil
	}

	if config.URL != "" {
		d.Set("endpoint", config.URL)
	}
	if config.IndexPrefix != "" {
		d.Set("index_prefix", config.IndexPrefix)
	}
	if config.IndexDaysMax != 0 {
		d.Set("index_days_max", config.IndexDaysMax)
	}
	if config.CA != "" {
		d.Set("ca_cert", trimPEMString(config.CA))
	}
	if config.Timeout != 0 {
		d.Set("timeout_seconds", int(config.Timeout))
	}

	return nil
}

// flattenLogsinkConfigOpensearch converts godo.DatabaseLogsinkConfig to Terraform schema data for opensearch
func flattenLogsinkConfigOpensearch(d *schema.ResourceData, config *godo.DatabaseLogsinkConfig) error {
	if config == nil {
		return nil
	}

	if config.URL != "" {
		d.Set("endpoint", config.URL)
	}
	if config.IndexPrefix != "" {
		d.Set("index_prefix", config.IndexPrefix)
	}
	if config.IndexDaysMax != 0 {
		d.Set("index_days_max", config.IndexDaysMax)
	}
	if config.CA != "" {
		d.Set("ca_cert", trimPEMString(config.CA))
	}
	if config.Timeout != 0 {
		d.Set("timeout_seconds", int(config.Timeout))
	}

	return nil
}

// validateLogsinkTimeout validates timeout is >= 1 second
func validateLogsinkTimeout(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(int)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be an integer", key))
		return
	}

	if v < 1 {
		errs = append(errs, fmt.Errorf("%q must be >= 1", key))
	}

	return
}

// validateLogsinkPort validates port is in range 1-65535
func validateLogsinkPort(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(int)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be an integer", key))
		return
	}

	if v < 1 || v > 65535 {
		errs = append(errs, fmt.Errorf("%q must be between 1 and 65535", key))
	}

	return
}

// validateHTTPSEndpoint validates that URL uses HTTPS scheme
func validateHTTPSEndpoint(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be a string", key))
		return
	}

	u, err := url.Parse(v)
	if err != nil {
		errs = append(errs, fmt.Errorf("%q must be a valid URL: %s", key, err))
		return
	}

	if u.Scheme != "https" {
		errs = append(errs, fmt.Errorf("%q must use HTTPS scheme", key))
	}

	return
}

// validateRsyslogFormat validates format is one of the allowed values
func validateRsyslogFormat(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be a string", key))
		return
	}

	validFormats := []string{"rfc5424", "rfc3164", "custom"}
	for _, format := range validFormats {
		if v == format {
			return
		}
	}

	errs = append(errs, fmt.Errorf("%q must be one of: %s", key, strings.Join(validFormats, ", ")))
	return
}

// validateIndexPrefix validates index_prefix is not empty
func validateIndexPrefix(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be a string", key))
		return
	}

	if strings.TrimSpace(v) == "" {
		errs = append(errs, fmt.Errorf("%q cannot be empty", key))
	}

	return
}

// validateIndexDaysMax validates index_days_max is >= 1
func validateIndexDaysMax(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(int)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be an integer", key))
		return
	}

	if v < 1 {
		errs = append(errs, fmt.Errorf("%q must be >= 1", key))
	}

	return
}

// trimPEMString trims whitespace from PEM strings while preserving content
func trimPEMString(s string) string {
	return strings.TrimSpace(s)
}

// validateLogsinkCustomDiff validates cross-field dependencies for logsink resources
func validateLogsinkCustomDiff(diff *schema.ResourceDiff, sinkType string) error {
	if sinkType == "rsyslog" {
		// If format is custom, require logline
		format := diff.Get("format").(string)
		logline := diff.Get("logline").(string)

		if format == "custom" && strings.TrimSpace(logline) == "" {
			return fmt.Errorf("logline is required when format is 'custom'")
		}

		// If any TLS cert fields are set, require tls = true
		tls := diff.Get("tls").(bool)
		caCert := diff.Get("ca_cert").(string)
		clientCert := diff.Get("client_cert").(string)
		clientKey := diff.Get("client_key").(string)

		if !tls && (caCert != "" || clientCert != "" || clientKey != "") {
			return fmt.Errorf("tls must be true when ca_cert, client_cert, or client_key is set")
		}

		// If client_cert or client_key is set, require both
		if (clientCert != "" || clientKey != "") && (clientCert == "" || clientKey == "") {
			return fmt.Errorf("both client_cert and client_key must be set for mTLS")
		}
	}

	return nil
}

// getLogsinkSinkType returns the sink_type value based on resource name
func getLogsinkSinkType(resourceName string) string {
	switch {
	case strings.Contains(resourceName, "_rsyslog"):
		return "rsyslog"
	case strings.Contains(resourceName, "_elasticsearch"):
		return "elasticsearch"
	case strings.Contains(resourceName, "_opensearch"):
		return "opensearch"
	default:
		return ""
	}
}

// buildCreateLogsinkRequest builds a godo.DatabaseCreateLogsinkRequest from resource data
func buildCreateLogsinkRequest(d *schema.ResourceData, sinkType string) *godo.DatabaseCreateLogsinkRequest {
	var config *godo.DatabaseLogsinkConfig

	switch sinkType {
	case "rsyslog":
		config = expandLogsinkConfigRsyslog(d)
	case "elasticsearch":
		config = expandLogsinkConfigElasticsearch(d)
	case "opensearch":
		config = expandLogsinkConfigOpensearch(d)
	}

	return &godo.DatabaseCreateLogsinkRequest{
		Name:   d.Get("name").(string),
		Type:   sinkType,
		Config: config,
	}
}

// buildUpdateLogsinkRequest builds a godo.DatabaseUpdateLogsinkRequest from resource data
func buildUpdateLogsinkRequest(d *schema.ResourceData, sinkType string) *godo.DatabaseUpdateLogsinkRequest {
	var config *godo.DatabaseLogsinkConfig

	switch sinkType {
	case "rsyslog":
		config = expandLogsinkConfigRsyslog(d)
	case "elasticsearch":
		config = expandLogsinkConfigElasticsearch(d)
	case "opensearch":
		config = expandLogsinkConfigOpensearch(d)
	}

	return &godo.DatabaseUpdateLogsinkRequest{
		Config: config,
	}
}

// setLogsinkResourceData sets the resource data from a godo.DatabaseLogsink response
func setLogsinkResourceData(d *schema.ResourceData, logsink *godo.DatabaseLogsink, sinkType string) error {
	d.Set("name", logsink.Name)
	d.Set("logsink_id", logsink.ID)

	switch sinkType {
	case "rsyslog":
		return flattenLogsinkConfigRsyslog(d, logsink.Config)
	case "elasticsearch":
		return flattenLogsinkConfigElasticsearch(d, logsink.Config)
	case "opensearch":
		return flattenLogsinkConfigOpensearch(d, logsink.Config)
	}

	return nil
}
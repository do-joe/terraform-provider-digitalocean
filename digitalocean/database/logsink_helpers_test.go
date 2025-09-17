package database

import (
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCreateLogsinkID(t *testing.T) {
	tests := []struct {
		name      string
		clusterID string
		logsinkID string
		expected  string
	}{
		{
			name:      "normal case",
			clusterID: "abc123",
			logsinkID: "def456",
			expected:  "abc123,def456",
		},
		{
			name:      "with special characters",
			clusterID: "cluster-uuid-123",
			logsinkID: "logsink_id_456",
			expected:  "cluster-uuid-123,logsink_id_456",
		},
		{
			name:      "empty strings",
			clusterID: "",
			logsinkID: "",
			expected:  ",",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createLogsinkID(tt.clusterID, tt.logsinkID)
			if result != tt.expected {
				t.Errorf("createLogsinkID() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSplitLogsinkID(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		expectedClusterID string
		expectedLogsinkID string
	}{
		{
			name:              "normal case",
			id:                "abc123,def456",
			expectedClusterID: "abc123",
			expectedLogsinkID: "def456",
		},
		{
			name:              "with special characters",
			id:                "cluster-uuid-123,logsink_id_456",
			expectedClusterID: "cluster-uuid-123",
			expectedLogsinkID: "logsink_id_456",
		},
		{
			name:              "multiple commas - split on first",
			id:                "abc123,def456,ghi789",
			expectedClusterID: "abc123",
			expectedLogsinkID: "def456,ghi789",
		},
		{
			name:              "no comma",
			id:                "abc123def456",
			expectedClusterID: "",
			expectedLogsinkID: "",
		},
		{
			name:              "empty string",
			id:                "",
			expectedClusterID: "",
			expectedLogsinkID: "",
		},
		{
			name:              "only comma",
			id:                ",",
			expectedClusterID: "",
			expectedLogsinkID: "",
		},
		{
			name:              "trailing comma",
			id:                "abc123,",
			expectedClusterID: "abc123",
			expectedLogsinkID: "",
		},
		{
			name:              "leading comma",
			id:                ",def456",
			expectedClusterID: "",
			expectedLogsinkID: "def456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusterID, logsinkID := splitLogsinkID(tt.id)
			if clusterID != tt.expectedClusterID {
				t.Errorf("splitLogsinkID() clusterID = %v, want %v", clusterID, tt.expectedClusterID)
			}
			if logsinkID != tt.expectedLogsinkID {
				t.Errorf("splitLogsinkID() logsinkID = %v, want %v", logsinkID, tt.expectedLogsinkID)
			}
		})
	}
}

func TestExpandLogsinkConfigRsyslog(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"server": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"tls": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"logline": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"structured_data": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ca_cert": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"client_cert": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"client_key": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"timeout_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}, map[string]interface{}{
		"server":          "test.example.com",
		"port":            514,
		"tls":             true,
		"format":          "rfc5424",
		"logline":         "%timestamp% %HOSTNAME%",
		"structured_data": "[test@123]",
		"ca_cert":         "  \n-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----\n  ",
		"client_cert":     "  client-cert  ",
		"client_key":      "  client-key  ",
		"timeout_seconds": 30,
	})

	config := expandLogsinkConfigRsyslog(d)

	if config.Server != "test.example.com" {
		t.Errorf("Expected Server to be 'test.example.com', got %s", config.Server)
	}
	if config.Port != 514 {
		t.Errorf("Expected Port to be 514, got %d", config.Port)
	}
	if config.TLS != true {
		t.Errorf("Expected TLS to be true, got %v", config.TLS)
	}
	if config.Format != "rfc5424" {
		t.Errorf("Expected Format to be 'rfc5424', got %s", config.Format)
	}
	if config.Logline != "%timestamp% %HOSTNAME%" {
		t.Errorf("Expected Logline to be '%%timestamp%% %%HOSTNAME%%', got %s", config.Logline)
	}
	if config.SD != "[test@123]" {
		t.Errorf("Expected SD to be '[test@123]', got %s", config.SD)
	}
	if config.CA != "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----" {
		t.Errorf("Expected CA to be trimmed PEM, got %s", config.CA)
	}
	if config.Cert != "client-cert" {
		t.Errorf("Expected Cert to be 'client-cert', got %s", config.Cert)
	}
	if config.Key != "client-key" {
		t.Errorf("Expected Key to be 'client-key', got %s", config.Key)
	}
	if config.Timeout != 30.0 {
		t.Errorf("Expected Timeout to be 30.0, got %f", config.Timeout)
	}
}

func TestExpandLogsinkConfigElasticsearch(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"endpoint": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"index_prefix": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"index_days_max": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"ca_cert": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"timeout_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}, map[string]interface{}{
		"endpoint":        "https://es.example.com:9200",
		"index_prefix":    "test-logs",
		"index_days_max":  14,
		"ca_cert":         "  test-ca-cert  ",
		"timeout_seconds": 60,
	})

	config := expandLogsinkConfigElasticsearch(d)

	if config.URL != "https://es.example.com:9200" {
		t.Errorf("Expected URL to be 'https://es.example.com:9200', got %s", config.URL)
	}
	if config.IndexPrefix != "test-logs" {
		t.Errorf("Expected IndexPrefix to be 'test-logs', got %s", config.IndexPrefix)
	}
	if config.IndexDaysMax != 14 {
		t.Errorf("Expected IndexDaysMax to be 14, got %d", config.IndexDaysMax)
	}
	if config.CA != "test-ca-cert" {
		t.Errorf("Expected CA to be 'test-ca-cert', got %s", config.CA)
	}
	if config.Timeout != 60.0 {
		t.Errorf("Expected Timeout to be 60.0, got %f", config.Timeout)
	}
}

func TestValidateLogsinkTimeout(t *testing.T) {
	tests := []struct {
		name      string
		val       interface{}
		expectErr bool
	}{
		{
			name:      "valid timeout",
			val:       30,
			expectErr: false,
		},
		{
			name:      "minimum valid timeout",
			val:       1,
			expectErr: false,
		},
		{
			name:      "zero timeout",
			val:       0,
			expectErr: true,
		},
		{
			name:      "negative timeout",
			val:       -1,
			expectErr: true,
		},
		{
			name:      "non-integer",
			val:       "30",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errs := validateLogsinkTimeout(tt.val, "timeout_seconds")
			hasErr := len(errs) > 0
			if hasErr != tt.expectErr {
				t.Errorf("validateLogsinkTimeout() error = %v, expectErr %v", hasErr, tt.expectErr)
			}
		})
	}
}

func TestValidateLogsinkPort(t *testing.T) {
	tests := []struct {
		name      string
		val       interface{}
		expectErr bool
	}{
		{
			name:      "valid port",
			val:       514,
			expectErr: false,
		},
		{
			name:      "minimum port",
			val:       1,
			expectErr: false,
		},
		{
			name:      "maximum port",
			val:       65535,
			expectErr: false,
		},
		{
			name:      "zero port",
			val:       0,
			expectErr: true,
		},
		{
			name:      "too high port",
			val:       65536,
			expectErr: true,
		},
		{
			name:      "negative port",
			val:       -1,
			expectErr: true,
		},
		{
			name:      "non-integer",
			val:       "514",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errs := validateLogsinkPort(tt.val, "port")
			hasErr := len(errs) > 0
			if hasErr != tt.expectErr {
				t.Errorf("validateLogsinkPort() error = %v, expectErr %v", hasErr, tt.expectErr)
			}
		})
	}
}

func TestValidateHTTPSEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		val       interface{}
		expectErr bool
	}{
		{
			name:      "valid HTTPS URL",
			val:       "https://es.example.com:9200",
			expectErr: false,
		},
		{
			name:      "HTTP URL",
			val:       "http://es.example.com:9200",
			expectErr: true,
		},
		{
			name:      "invalid URL",
			val:       "not-a-url",
			expectErr: true,
		},
		{
			name:      "non-string",
			val:       123,
			expectErr: true,
		},
		{
			name:      "empty string",
			val:       "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errs := validateHTTPSEndpoint(tt.val, "endpoint")
			hasErr := len(errs) > 0
			if hasErr != tt.expectErr {
				t.Errorf("validateHTTPSEndpoint() error = %v, expectErr %v", hasErr, tt.expectErr)
			}
		})
	}
}

func TestValidateRsyslogFormat(t *testing.T) {
	tests := []struct {
		name      string
		val       interface{}
		expectErr bool
	}{
		{
			name:      "rfc5424 format",
			val:       "rfc5424",
			expectErr: false,
		},
		{
			name:      "rfc3164 format",
			val:       "rfc3164",
			expectErr: false,
		},
		{
			name:      "custom format",
			val:       "custom",
			expectErr: false,
		},
		{
			name:      "invalid format",
			val:       "invalid",
			expectErr: true,
		},
		{
			name:      "non-string",
			val:       123,
			expectErr: true,
		},
		{
			name:      "empty string",
			val:       "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errs := validateRsyslogFormat(tt.val, "format")
			hasErr := len(errs) > 0
			if hasErr != tt.expectErr {
				t.Errorf("validateRsyslogFormat() error = %v, expectErr %v", hasErr, tt.expectErr)
			}
		})
	}
}

func TestValidateIndexPrefix(t *testing.T) {
	tests := []struct {
		name      string
		val       interface{}
		expectErr bool
	}{
		{
			name:      "valid prefix",
			val:       "logs",
			expectErr: false,
		},
		{
			name:      "prefix with dashes",
			val:       "test-logs",
			expectErr: false,
		},
		{
			name:      "empty string",
			val:       "",
			expectErr: true,
		},
		{
			name:      "only spaces",
			val:       "   ",
			expectErr: true,
		},
		{
			name:      "non-string",
			val:       123,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errs := validateIndexPrefix(tt.val, "index_prefix")
			hasErr := len(errs) > 0
			if hasErr != tt.expectErr {
				t.Errorf("validateIndexPrefix() error = %v, expectErr %v", hasErr, tt.expectErr)
			}
		})
	}
}

func TestValidateIndexDaysMax(t *testing.T) {
	tests := []struct {
		name      string
		val       interface{}
		expectErr bool
	}{
		{
			name:      "valid days",
			val:       7,
			expectErr: false,
		},
		{
			name:      "minimum days",
			val:       1,
			expectErr: false,
		},
		{
			name:      "zero days",
			val:       0,
			expectErr: true,
		},
		{
			name:      "negative days",
			val:       -1,
			expectErr: true,
		},
		{
			name:      "non-integer",
			val:       "7",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errs := validateIndexDaysMax(tt.val, "index_days_max")
			hasErr := len(errs) > 0
			if hasErr != tt.expectErr {
				t.Errorf("validateIndexDaysMax() error = %v, expectErr %v", hasErr, tt.expectErr)
			}
		})
	}
}

func TestTrimPEMString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal PEM",
			input:    "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
			expected: "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
		},
		{
			name:     "PEM with leading/trailing whitespace",
			input:    "  \n-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----\n  ",
			expected: "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   \n   ",
			expected: "",
		},
		{
			name:     "simple string",
			input:    "  hello world  ",
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimPEMString(tt.input)
			if result != tt.expected {
				t.Errorf("trimPEMString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetLogsinkSinkType(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		expected     string
	}{
		{
			name:         "rsyslog resource",
			resourceName: "digitalocean_database_logsink_rsyslog",
			expected:     "rsyslog",
		},
		{
			name:         "elasticsearch resource",
			resourceName: "digitalocean_database_logsink_elasticsearch",
			expected:     "elasticsearch",
		},
		{
			name:         "opensearch resource",
			resourceName: "digitalocean_database_logsink_opensearch",
			expected:     "opensearch",
		},
		{
			name:         "unknown resource",
			resourceName: "digitalocean_database_logsink_unknown",
			expected:     "",
		},
		{
			name:         "non-logsink resource",
			resourceName: "digitalocean_database_cluster",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLogsinkSinkType(tt.resourceName)
			if result != tt.expected {
				t.Errorf("getLogsinkSinkType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBuildCreateLogsinkRequest(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"server": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}, map[string]interface{}{
		"name":   "test-logsink",
		"server": "test.example.com",
		"port":   514,
	})

	req := buildCreateLogsinkRequest(d, "rsyslog")

	if req.Name != "test-logsink" {
		t.Errorf("Expected Name to be 'test-logsink', got %s", req.Name)
	}
	if req.Type != "rsyslog" {
		t.Errorf("Expected Type to be 'rsyslog', got %s", req.Type)
	}
	if req.Config == nil {
		t.Error("Expected Config to be non-nil")
	}
	if req.Config.Server != "test.example.com" {
		t.Errorf("Expected Config.Server to be 'test.example.com', got %s", req.Config.Server)
	}
}

func TestSetLogsinkResourceData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"logsink_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"server": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}, map[string]interface{}{})

	logsink := &godo.DatabaseLogsink{
		ID:   "logsink-123",
		Name: "test-logsink",
		Type: "rsyslog",
		Config: &godo.DatabaseLogsinkConfig{
			Server: "test.example.com",
			Port:   514,
		},
	}

	err := setLogsinkResourceData(d, logsink, "rsyslog")
	if err != nil {
		t.Errorf("setLogsinkResourceData() error = %v", err)
	}

	if d.Get("name").(string) != "test-logsink" {
		t.Errorf("Expected name to be 'test-logsink', got %s", d.Get("name").(string))
	}
	if d.Get("logsink_id").(string) != "logsink-123" {
		t.Errorf("Expected logsink_id to be 'logsink-123', got %s", d.Get("logsink_id").(string))
	}
}

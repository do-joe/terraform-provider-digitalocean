package database_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDigitalOceanDatabaseLogsinkOpensearch_Basic tests creating a basic opensearch logsink
// with required fields (endpoint, index_prefix, index_days_max). Expected: successful creation.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_Basic(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigBasic, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_opensearch.test", &logsink),
					testAccCheckDigitalOceanDatabaseLogsinkAttributes(&logsink, logsinkName, "opensearch"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "name", logsinkName),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "endpoint", "https://opensearch.example.com:9200"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_prefix", "db-logs"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_days_max", "7"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_opensearch.test", "cluster_id"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_opensearch.test", "logsink_id"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_opensearch.test", "id"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkOpensearch_Update tests updating an opensearch logsink
// configuration (index_prefix, index_days_max, timeout_seconds). Expected: successful update.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_Update(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigBasic, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_opensearch.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_prefix", "db-logs"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_days_max", "7"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigUpdated, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_opensearch.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_prefix", "updated-logs"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_days_max", "14"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "timeout_seconds", "30"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkOpensearch_WithCA tests creating an opensearch logsink
// with CA certificate for TLS verification. Expected: successful creation with certificate.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_WithCA(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigWithCA, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_opensearch.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "ca_cert", "-----BEGIN CERTIFICATE-----\nMIIBkTCB+wIJANOxiCFJwTkMMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNVBAMMCWxv\nY2FsaG9zdDAeFw0yMzEwMTAwMDAwMDBaFw0yNDEwMDkwMDAwMDBaMBQxEjAQBgNV\nBAMMCWxvY2FsaG9zdDBcMA0GCSqGSIb3DQEBAQUAAksAMEgCQQC7k3M1Y7s+7k3M\n1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M\n1Y7s+7k3AgMBAAEwDQYJKoZIhvcNAQELBQADQQA7k3M1Y7s+7k3M1Y7s+7k3M1Y7\ns+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M\n-----END CERTIFICATE-----"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkOpensearch_MongoDB tests creating an opensearch logsink
// with a MongoDB cluster (engine compatibility test). Expected: successful creation.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_MongoDB(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigMongoDB, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_opensearch.test", &logsink),
					testAccCheckDigitalOceanDatabaseLogsinkAttributes(&logsink, logsinkName, "opensearch"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "name", logsinkName),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "endpoint", "https://opensearch.example.com:9200"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_prefix", "mongo-logs"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkOpensearch_InvalidIndexDaysMax tests validation for invalid index_days_max.
// Uses value 0 which is below minimum (must be >= 1). Expected: validation error.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_InvalidIndexDaysMax(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigInvalidIndexDays, clusterName, logsinkName),
				ExpectError: regexp.MustCompile("must be >= 1"),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkOpensearch_EmptyIndexPrefix tests validation for empty index_prefix.
// Uses empty string for index_prefix which is not allowed. Expected: validation error.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_EmptyIndexPrefix(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigEmptyIndexPrefix, clusterName, logsinkName),
				ExpectError: regexp.MustCompile("index_prefix cannot be empty"),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkOpensearch_MalformedEndpoint tests validation for malformed endpoint URLs.
// Uses invalid URL format that fails scheme validation. Expected: validation error.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_MalformedEndpoint(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigMalformedEndpoint, clusterName, logsinkName),
				ExpectError: regexp.MustCompile("must use HTTPS scheme"),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigBasic = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id      = digitalocean_database_cluster.test.id
  name            = "%s"
  endpoint        = "https://opensearch.example.com:9200"
  index_prefix    = "db-logs"
  index_days_max  = 7
}`

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigUpdated = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id        = digitalocean_database_cluster.test.id
  name              = "%s"
  endpoint          = "https://opensearch.example.com:9200"
  index_prefix      = "updated-logs"
  index_days_max    = 14
  timeout_seconds   = 30
}`

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigWithCA = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id     = digitalocean_database_cluster.test.id
  name           = "%s"
  endpoint       = "https://opensearch.example.com:9200"
  index_prefix   = "db-logs"
  index_days_max = 7
  ca_cert        = "-----BEGIN CERTIFICATE-----\nMIIBkTCB+wIJANOxiCFJwTkMMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNVBAMMCWxv\nY2FsaG9zdDAeFw0yMzEwMTAwMDAwMDBaFw0yNDEwMDkwMDAwMDBaMBQxEjAQBgNV\nBAMMCWxvY2FsaG9zdDBcMA0GCSqGSIb3DQEBAQUAAksAMEgCQQC7k3M1Y7s+7k3M\n1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M\n1Y7s+7k3AgMBAAEwDQYJKoZIhvcNAQELBQADQQA7k3M1Y7s+7k3M1Y7s+7k3M1Y7\ns+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M\n-----END CERTIFICATE-----"
}`

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigMongoDB = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "mongodb"
  version    = "6"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id     = digitalocean_database_cluster.test.id
  name           = "%s"
  endpoint       = "https://opensearch.example.com:9200"
  index_prefix   = "mongo-logs"
  index_days_max = 7
}`

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigInvalidIndexDays = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id     = digitalocean_database_cluster.test.id
  name           = "%s"
  endpoint       = "https://opensearch.example.com:9200"
  index_prefix   = "db-logs"
  index_days_max = 0
}`

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigEmptyIndexPrefix = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id     = digitalocean_database_cluster.test.id
  name           = "%s"
  endpoint       = "https://opensearch.example.com:9200"
  index_prefix   = ""
  index_days_max = 7
}`

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigMalformedEndpoint = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id     = digitalocean_database_cluster.test.id
  name           = "%s"
  endpoint       = "not-a-valid-url"
  index_prefix   = "db-logs"
  index_days_max = 7
}`

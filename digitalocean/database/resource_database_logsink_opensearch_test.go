package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

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

func TestAccDigitalOceanDatabaseLogsinkOpensearch_ElasticsearchCompatibility(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigElasticsearch, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_opensearch.test", &logsink),
					testAccCheckDigitalOceanDatabaseLogsinkAttributes(&logsink, logsinkName, "opensearch"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "name", logsinkName),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "endpoint", "https://elasticsearch.example.com:9200"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_prefix", "es-logs"),
				),
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

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigElasticsearch = `
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
  endpoint       = "https://elasticsearch.example.com:9200"
  index_prefix   = "es-logs"
  index_days_max = 7
}`

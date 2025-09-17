package database_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseLogsinkRsyslog_Basic(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigBasic, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_rsyslog.test", &logsink),
					testAccCheckDigitalOceanDatabaseLogsinkAttributes(&logsink, logsinkName, "rsyslog"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "name", logsinkName),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "server", "192.168.1.100"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "port", "514"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "tls", "false"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "format", "rfc5424"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_rsyslog.test", "cluster_id"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_rsyslog.test", "logsink_id"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_rsyslog.test", "id"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseLogsinkRsyslog_Update(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigBasic, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_rsyslog.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "port", "514"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "tls", "false"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "format", "rfc5424"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigUpdated, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_rsyslog.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "port", "1514"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "tls", "true"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "format", "rfc3164"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "structured_data", "[example@41058 iut=\"3\"]"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseLogsinkRsyslog_CustomFormat(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigCustom, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_rsyslog.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "format", "custom"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "logline", "%%timestamp%% %%HOSTNAME%% %%app-name%% %%procid%% %%msgid%% %%msg%%"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseLogsinkRsyslog_TLS(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigTLS, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_rsyslog.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "tls", "true"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "ca_cert", "-----BEGIN CERTIFICATE-----\nMIIBkTCB+wIJANOxiCFJwTkMMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNVBAMMCWxv\nY2FsaG9zdDAeFw0yMzEwMTAwMDAwMDBaFw0yNDEwMDkwMDAwMDBaMBQxEjAQBgNV\nBAMMCWxvY2FsaG9zdDBcMA0GCSqGSIb3DQEBAQUAAksAMEgCQQC7k3M1Y7s+7k3M\n1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M\n1Y7s+7k3AgMBAAEwDQYJKoZIhvcNAQELBQADQQA7k3M1Y7s+7k3M1Y7s+7k3M1Y7\ns+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M\n-----END CERTIFICATE-----"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseLogsinkRsyslog_MongoDB_ShouldFail(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigMongoDB, clusterName, logsinkName),
				ExpectError: regexp.MustCompile("rsyslog sink type is not supported for MongoDB"),
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseLogsinkExists(resource string, logsink *godo.DatabaseLogsink) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		// Parse composite ID
		clusterID := rs.Primary.Attributes["cluster_id"]
		logsinkID := rs.Primary.Attributes["logsink_id"]

		foundLogsink, _, err := client.Databases.GetLogsink(context.Background(), clusterID, logsinkID)
		if err != nil {
			return err
		}

		if foundLogsink.ID != logsinkID {
			return fmt.Errorf("Record not found")
		}

		*logsink = *foundLogsink

		return nil
	}
}

func testAccCheckDigitalOceanDatabaseLogsinkAttributes(logsink *godo.DatabaseLogsink, name, sinkType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if logsink.Name != name {
			return fmt.Errorf("Bad name: %s", logsink.Name)
		}

		if logsink.Type != sinkType {
			return fmt.Errorf("Bad sink type: %s", logsink.Type)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDatabaseLogsinkDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_database_logsink_rsyslog" && rs.Type != "digitalocean_database_logsink_opensearch" {
			continue
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		logsinkID := rs.Primary.Attributes["logsink_id"]

		_, _, err := client.Databases.GetLogsink(context.Background(), clusterID, logsinkID)
		if err == nil {
			return fmt.Errorf("DatabaseLogsink still exists")
		}
	}

	return nil
}

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigBasic = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 514
  tls        = false
  format     = "rfc5424"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigUpdated = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id       = digitalocean_database_cluster.test.id
  name             = "%s"
  server           = "192.168.1.100"
  port             = 1514
  tls              = true
  format           = "rfc3164"
  structured_data  = "[example@41058 iut=\"3\"]"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigCustom = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 514
  format     = "custom"
  logline    = "%%timestamp%% %%HOSTNAME%% %%app-name%% %%procid%% %%msgid%% %%msg%%"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigTLS = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 514
  tls        = true
  ca_cert    = "-----BEGIN CERTIFICATE-----\nMIIBkTCB+wIJANOxiCFJwTkMMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNVBAMMCWxv\nY2FsaG9zdDAeFw0yMzEwMTAwMDAwMDBaFw0yNDEwMDkwMDAwMDBaMBQxEjAQBgNV\nBAMMCWxvY2FsaG9zdDBcMA0GCSqGSIb3DQEBAQUAAksAMEgCQQC7k3M1Y7s+7k3M\n1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M\n1Y7s+7k3AgMBAAEwDQYJKoZIhvcNAQELBQADQQA7k3M1Y7s+7k3M1Y7s+7k3M1Y7\ns+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M1Y7s+7k3M\n-----END CERTIFICATE-----"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigMongoDB = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "mongodb"
  version    = "6"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 514
  format     = "rfc5424"
}`

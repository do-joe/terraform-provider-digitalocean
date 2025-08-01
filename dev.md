# Terraform Provider Enhancement: DigitalOcean Database Metrics Support

## Objective
Enable Terraform users to access PostgreSQL metrics endpoint information and credentials from a DigitalOcean database cluster.

---

## 1. Update `digitalocean_database_cluster` Resource & Data Source

### Goal
Include `metrics_endpoint` attribute in both the resource and data source schemas to expose the Prometheus-compatible metrics URL.

### Implementation
- In `resource_database_cluster.go` and `datasource_database_cluster.go`:
    - Add new schema field:
      ```go
      "metrics_endpoint": {
          Type:     schema.TypeString,
          Computed: true,
      }
      ```
    - In read functions, populate the field using the `ServiceAddress` struct:
      ```go
      if cluster.Metrics != nil {
          addr := cluster.Metrics
          d.Set("metrics_endpoint", fmt.Sprintf("https://%s:%d/metrics", addr.Host, addr.Port))
      }
      ```

### Notes
- Uses `cluster.Metrics` of type `*godo.ServiceAddress` which includes `Host` and `Port`.
- This avoids hardcoding port 9273 and supports potential changes or custom setups.
- Godo references:
  - https://pkg.go.dev/github.com/digitalocean/godo#Database

---

## 2. Add `digitalocean_database_cluster_metrics_credentials` Data Source

### Goal
Add a new data source to fetch metrics username and password using godo’s `GetMetricsCredentials` API call.

### Schema
- `cluster_id` (string, required)
- `username` (string, computed)
- `password` (string, computed, sensitive)

### File
`datasource_database_cluster_metrics_credentials.go`

### Example Schema Block
```go
"username": {
    Type:     schema.TypeString,
    Computed: true,
},
"password": {
    Type:      schema.TypeString,
    Computed:  true,
    Sensitive: true,
},
```

### Data Source Read Function
```go
func dataSourceDatabaseClusterMetricsCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*CombinedConfig).godoClient
    clusterID := d.Get("cluster_id").(string)
    creds, _, err := client.Databases.GetMetricsCredentials(ctx, clusterID)
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId(clusterID)
    d.Set("username", creds.Username)
    d.Set("password", creds.Password)
    return nil
}
```

### Example Usage in HCL
```hcl
data "digitalocean_database_cluster_metrics_credentials" "example" {
  cluster_id = digitalocean_database_cluster.example.id
}
```

### Note
- The CA certificate is intentionally excluded from this data source, as it can already be retrieved using the existing [`digitalocean_database_ca`](https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/data-sources/database_ca) data source.
- Godo references:
  - https://pkg.go.dev/github.com/digitalocean/godo#DatabasesServiceOp.GetMetricsCredentials

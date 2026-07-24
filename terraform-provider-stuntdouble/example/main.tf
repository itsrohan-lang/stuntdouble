terraform {
  required_providers {
    stuntdouble = {
      source = "itsrohan-lang/stuntdouble"
      version = "1.0.0"
    }
  }
}

provider "stuntdouble" {
  control_plane_url = "http://localhost:4439"
}

# Define a global zero-trust policy for the entire organization
resource "stuntdouble_policy" "global_backend" {
  org_id         = "ent_global_engineering"
  blocked_ports  = [5432, 27017, 3306, 6379, 443]
  allowed_agents = ["claude", "cursor"]
  strict_egress  = true
}

# The Terraform provider will push this JSON directly to the Go Control Plane API

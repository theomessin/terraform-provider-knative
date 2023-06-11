# ⛰️ Terraform Knative Provider

[![Tests](https://github.com/theomessin/terraform-provider-knative/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/theomessin/terraform-provider-knative/actions/workflows/test.yml)

[Knative Serving](https://knative.dev/docs/serving/) builds on Kubernetes to support deploying and serving of applications and functions as serverless containers. This Terraform Provider lets you deploy and manage your Knative Services using Terraform.
Maybe one day it'll support Knative Eventing, or even handle installing Knative for you.

## Getting Started

Simply declare the provider as a requirement:

```tf
terraform {
  required_version = ">= 0.13"

  required_providers {
    knative = {
      source  = "theomessin/knative"
      version = ">= 0.1.1"
    }
  }
}
```

Now use the provider. For example, read the URL of an existing Service:

```tf
data "knative_service" "this" {
  name = "app"
}

output "app" {
  url = data.knative_service.this.status.url
}
```

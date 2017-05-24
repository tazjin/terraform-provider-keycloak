Terraform Keycloak Provider
===========================

This project implements a [Terraform provider][] for declaratively configuring
API resources in [Keycloak][].

## Status

This provider can currently manage Keycloak `client` resources. Not all fields of
this resource are supported at the moment.

## Installation

Installation is simple:

1. Install and configure Go
2. `go get github.com/tazjin/terraform-provider-keycloak`

You must also [configure][] the provider in Terraform. In your `~/.terraformrc` add

```
providers {
  keycloak = "/path/to/gopath/bin/terraform-provider-keycloak"
}
```

## Setup instructions

The Keycloak instance to manage needs to be configured with a client that has
permission to change the resources in Keycloak.

If you want to create and manage realms directly you should grant this client
the `admin` role.

The provider needs to be configured with credentials to access the API:

```
provider "keycloak" {
  # These parameters are required:
  client_id     = "dingus"
  client_secret = "Oox7luexoofeuquaosh5ti3aequie7sh"
  api_base      = "https://keycloak.my-company.acme"
  
  # These parameters are optional:
  realm = "my-company"  # defaults to 'master'
}
```

[configure]: https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin

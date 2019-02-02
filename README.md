Terraform Keycloak Provider
===========================

[![Build Status](https://travis-ci.org/tazjin/terraform-provider-keycloak.svg?branch=master)](https://travis-ci.org/tazjin/terraform-provider-keycloak)

This project implements a [Terraform provider][] for declaratively configuring
API resources in [Keycloak][].

## Status

This provider can currently manage Keycloak `client` resources and user-role mappings.

Not all fields of those resources are supported at the moment.

## Installation

Grab a binary release for your operating system from the [releases][] page and drop it into
`~/.terraform.d/plugins`.

Run `terraform init` to initialise the new provider in the folder containing your configuration
files and `terraform providers` to check that it has been loaded correctly.

**Note**: The targeted version of Terraform is currently **v0.11.11**.

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
Note the following steps will need to be completed as part of the provider setup: 
1. The client ("dingus" in above example) will have to be created under the chosen realm
2. "Service Accounts Enabled" need to be enabled under client settings
3. Under "Service Account Roles" the create-client, create-group, manage-clients, manage-groups, manage-users roles will need to be assigned under "realm-management

Groups can be created using the keycloak_group resource:
```
resource "keycloak_group" "group1" {
  name       = "<group_name>"
  realm      = "<realm_name>"
}

```

Users can be created using the keycloak_user resource:
```
resource "keycloak_user" "user1" {
  realm      = "<realm_name>"
  username   = "user1"
  firstname  = "user"
  lastname   = "cameron"
  email      = "jcameron@abc.com"
}
```

User group mapping can be created using the keycloak_user_group_mapping resource. You have to reference group and user ids as listed below. 
```
resource "keycloak_user_group_mapping" "group1_map" {
  group_id = "${keycloak_group.group1.id}"
  user_ids   = ["${keycloak_user.user1.id}", ]
  realm      = "<realm_name>"
}

```
To import a user or group use the following command:
```
terraform import <keycloak_resource>.<resource_name> <realm_name>.<resource_id>
terraform import keycloak_group.group2 Jenkins.310f73af-3b70-4e4a-9a6f-a3f4de8c8f
```

## Building from source

For "vanilla"-builds do this:

1. Install and configure Go v1.11.x or later
2. `$ make install`

The last step will build the provider for your machine and copy the binary into `~/.terraform.d/plugins` ready to be used.

## Publish new releases

All maintainers have the ability to publish new releases of this provider
by pushing a new git tag to the repository after new changes has been merged.

As an example for publishing a new releases named `v1.0.1`:

```
$ git tag v1.0.1
$ git push origin --tags
```

The above will make Travis CI build the project as usual, then upload the resulting binaries
up to github.com, making them visible on our [releases][] page.

[Terraform provider]: https://www.terraform.io/docs/plugins/provider.html
[Keycloak]: http://www.keycloak.org/
[configure]: https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin
[releases]: https://github.com/tazjin/terraform-provider-keycloak/releases

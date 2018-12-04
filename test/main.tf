/*
resource "keycloak_client" "client_test" {
  realm     = "Jenkins"
  client_id = "test"
  redirect_uris = ["http://127.0.0.1:8080/"]
  service_accounts_enabled = true
}
*/

resource "keycloak_group" "group1" {
  name       = "group1"
  realm      = "Jenkins"
}

resource "keycloak_group" "group2" {
  name       = "group2"
  realm      = "Jenkins"
}

resource "keycloak_group" "group11" {
  name       = "group11"
  realm      = "Jenkins"
}
resource "keycloak_group" "group21" {
  name       = "group21"
  realm      = "Jenkins"
}
resource "keycloak_user" "martin" {
  realm      = "Jenkins"
  username   = "martin"
  firstname  = "Martin"
  lastname   = "Patel"
  email      = "mpatel@abc.com"
}
resource "keycloak_user" "martin1" {
  realm      = "Jenkins"
  username   = "martin1"
  firstname  = "martin1"
  lastname   = "patel"
  email      = "mpatel1@abc.com"
}

resource "keycloak_user" "josh" {
  realm      = "Jenkins"
  username   = "josh"
  firstname  = "josh"
  lastname   = "cameron"
  email      = "jcameron@abc.com"
}

resource "keycloak_user" "josh1" {
  realm      = "Jenkins"
  username   = "josh1"
  firstname  = "josh"
  lastname   = "cameron"
  email      = "jcameron1@abc.com"
  initial_required_actions = ["UPDATE_PASSWORD"]
}

resource "keycloak_user_group_mapping" "group1_map" {
  group_id = "${keycloak_group.group1.id}"
  user_ids   = ["${keycloak_user.martin1.id}", ]
  realm      = "Jenkins"
}

resource "keycloak_user_group_mapping" "group2_map" {
  group_id = "${keycloak_group.group2.id}"
  user_ids   = ["${keycloak_user.josh1.id}"]
  realm      = "Jenkins"
}

provider "keycloak" {
  # These parameters are required:
  client_id     = "terraform"
  client_secret = "71c07a1e-f114-42b9-ae88-f5c627bab599"
  api_base      = "http://127.0.0.1:8080"
  
  # These parameters are optional:
  realm = "Jenkins"  # defaults to 'master'
}


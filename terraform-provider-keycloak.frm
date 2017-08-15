inputs:
  "/": # https://github.com/tklx/base
    type: "tar"
    hash: "9nkvYhmJHaeK_Agc3Lm5rg444dSLWDp0Pri-KilHiX3A9Pt4TaQ7RxOj5qMSs6XT"
    silo: "https://github.com/tklx/base/releases/download/0.1.1/rootfs.tar.xz"
  "/opt":
    type: "tar"
    hash: "gi0Kpb-VH3TK0UBX6YmpuKsrMAUlxicPrY2YvXPo9sBQm_NsD_hKrn7pmc95zrmM"
    silo: "https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz"

  # Terraform provider dependencies
  "/go/src/github.com/hashicorp/terraform":
    type: "git"
    # This hash is tag v0.10.0
    hash: "2041053ee9444fa8175a298093b55a89586a1823"
    silo: "https://github.com/hashicorp/terraform"
action:
  policy: governor
  command:
    - "/bin/sh"
    - "-e"
    - "-c"
    - |
      export PATH="/opt/go/bin:$PATH"
      export GOROOT=/opt/go
      export GOPATH=/go
      echo 'nameserver 8.8.8.8' > /etc/resolv.conf
      apt-get update && apt-get install -y git ca-certificates
      mkdir -p /go/src/github.com/tazjin
      git clone --single-branch --branch master https://github.com/tazjin/terraform-provider-keycloak /go/src/github.com/tazjin/terraform-provider-keycloak
      cd /go/src/github.com/tazjin/terraform-provider-keycloak
      ./build-release.sh build
outputs:
  "release":
    type: "dir"
    mount: "/go/src/github.com/tazjin/terraform-provider-keycloak/release"
    silo: "file://release"

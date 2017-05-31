package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
)

func TestValidateProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("Validation error: %s", err)
	}
}

package provider

import "github.com/hashicorp/terraform/helper/schema"

func realm(d *schema.ResourceData) string {
	return d.Get("realm").(string)
}

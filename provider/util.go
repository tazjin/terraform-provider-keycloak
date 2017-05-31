package provider

import "github.com/hashicorp/terraform/helper/schema"

func realm(d *schema.ResourceData) string {
	return d.Get("realm").(string)
}

func getOptionalBool(d *schema.ResourceData, key string) *bool {
	if v, present := d.GetOk(key); present {
		b := v.(bool)
		return &b
	}
	return nil
}

func setOptionalBool(d *schema.ResourceData, key string, b *bool) {
	if b != nil {
		d.Set(key, *b)
	}
}

func getOptionalInt(d *schema.ResourceData, key string) *int {
	if v, present := d.GetOk(key); present {
		i := v.(int)
		return &i
	}
	return nil
}

func setOptionalInt(d *schema.ResourceData, key string, i *int) {
	if i != nil {
		d.Set(key, *i)
	}
}

func getStringSlice(d *schema.ResourceData, key string) []string {
	var stringSlice []string = []string{}
	untyped, present := d.GetOk(key)

	if !present {
		return stringSlice
	}

	for _, value := range untyped.([]interface{}) {
		stringSlice = append(stringSlice, value.(string))
	}

	return stringSlice
}

package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func realm(d *schema.ResourceData) string {
	return d.Get("realm").(string)
}

func clientId(d *schema.ResourceData) string {
	return d.Get("client_id").(string)
}

func containsSameElements(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for _, a_elem := range a {
		if !contains(b, a_elem) {
			return false
		}
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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

// This function is used when importing realm-specific resources. The realm must be specified by the user when
// importing by using a `${realm}.${resource_id}` syntax.
func splitRealmId(raw string) (string, string, error) {
	split := strings.Split(raw, ".")

	if len(split) != 2 {
		return "", "", fmt.Errorf("Import ID must be specified as '${realm}.${resource_id}'")
	}

	return split[0], split[1], nil
}

func getMandatoryStringList(d *schema.ResourceData, key string) []string {
	stringList := []string{}

	for _, stringVal := range d.Get(key).([]interface{}) {
		stringList = append(stringList, stringVal.(string))
	}
	return stringList
}

func getOptionalStringList(d *schema.ResourceData, key string) []string {
	stringList := []string{}
	rawList, present := d.GetOk(key)
	if present {
		for _, stringVal := range rawList.([]interface{}) {
			stringList = append(stringList, stringVal.(string))
		}
	}
	return stringList

}

func getOptionalStringMap(d *schema.ResourceData, key string) map[string]string {
	stringMap := map[string]string{}
	rawMap, present := d.GetOk(key)
	if present {
		for stringKey, stringVal := range rawMap.(map[string]interface{}) {
			stringMap[stringKey] = stringVal.(string)
		}
	}
	return stringMap

}

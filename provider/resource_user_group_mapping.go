package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tazjin/terraform-provider-keycloak/keycloak"
)

func resourceUserGroupMapping() *schema.Resource {
	return &schema.Resource{
		// API methods
		Read:   schema.ReadFunc(resourceUserGroupMappingRead),
		Create: schema.CreateFunc(resourceUserGroupMappingCreate),
		Update: schema.UpdateFunc(resourceUserGroupMappingUpdate),
		Delete: schema.DeleteFunc(resourceUserGroupMappingDelete),

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"realm": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserGroupMappingRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)

	ug, err := c.GetUsersInGroup(d.Id(), realm(d))
	if err != nil {
		return err
	}
	userGroupMappingToResourceData(ug, d)

	return nil
}

func resourceUserGroupMappingCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)

	ug := resourceDataToUserGroupMapping(d)
	err := c.AddUsersToGroup(ug.UserIds, ug.GroupId, realm(d))
	if err != nil {
		return err
	}
	d.SetId(ug.GroupId)
	return resourceUserGroupMappingRead(d, m)
}

func resourceUserGroupMappingUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)

	newUg := resourceDataToUserGroupMapping(d)
	currentUg, err := c.GetUsersInGroup(newUg.GroupId, newUg.Realm)
	if err != nil {
		return err
	}
	err = c.RemoveUsersFromGroup(currentUg.UserIds, newUg.GroupId, realm(d))
	if err != nil {
		return err
	}
	err = c.AddUsersToGroup(newUg.UserIds, newUg.GroupId, newUg.Realm)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserGroupMappingDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)

	group := resourceDataToUserGroupMapping(d)
	err := c.RemoveUsersFromGroup(group.UserIds, group.GroupId, realm(d))
	if err != nil {
		return err
	}

	return nil
}

func userGroupMappingToResourceData(ug *keycloak.UserGroupMap, d *schema.ResourceData) {
	d.SetId(ug.GroupId)
	d.Set("user_ids", ug.UserIds)
	d.Set("group_id", ug.GroupId)
	d.Set("realm", ug.Realm)
}

func resourceDataToUserGroupMapping(d *schema.ResourceData) keycloak.UserGroupMap {
	return keycloak.UserGroupMap{
		UserIds: getMandatoryStringList(d, "user_ids"),
		GroupId: d.Get("group_id").(string),
		Realm:   d.Get("realm").(string),
	}
}

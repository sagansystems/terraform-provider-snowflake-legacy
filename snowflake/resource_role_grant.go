package snowflake

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRoleGrant() *schema.Resource {
	return &schema.Resource{
		Create: createRoleGrant,
		Read:   readRoleGrant,
		Delete: deleteRoleGrant,

		Schema: map[string]*schema.Schema{
			"role": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the role to grant",
				ForceNew:    true,
			},
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user to which this role should be granted",
				ForceNew:    true,
			},
		},
	}
}

func createRoleGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	role := d.Get("role").(string)
	user := d.Get("user").(string)

	d.SetId(roleGrantIDFromParams(role, user))

	stmtSQL := fmt.Sprintf("GRANT ROLE \"%s\" TO USER \"%s\"", role, user)

	log.Println("Executing statement:", stmtSQL)

	if _, err := db.Exec(stmtSQL); err != nil {
		return err
	}

	return readRoleGrant(d, meta)
}

func readRoleGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	role, user := paramsFromRoleGrantID(d.Id())

	stmtSQL := fmt.Sprintf("SHOW GRANTS TO USER \"%s\"", user)

	log.Println("Executing statement:", stmtSQL)

	rows, err := db.Query(stmtSQL)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var dbCreatedAt string
		var dbRole string
		var dbGrantedTo string
		var dbGranteeName string
		var dbGrantedBy string
		if err := rows.Scan(&dbCreatedAt, &dbRole, &dbGrantedTo, &dbGranteeName, &dbGrantedBy); err != nil {
			return err
		}
		if role == dbRole {
			d.Set("role", dbRole)
			d.Set("user", dbGranteeName)
			return nil
		}
	}

	return fmt.Errorf("The grant of role %s to user %s does not exist.", role, user)
}

func deleteRoleGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	role, user := paramsFromRoleGrantID(d.Id())

	stmtSQL := fmt.Sprintf("REVOKE ROLE \"%s\" FROM USER \"%s\"", role, user)

	log.Println("Executing statement:", stmtSQL)

	_, err := db.Exec(stmtSQL)
	if err == nil {
		d.SetId("")
	}

	return err
}

func paramsFromRoleGrantID(id string) (role, user string) {
	splits := strings.Split(id, "-")
	return splits[0], splits[1]
}

func roleGrantIDFromParams(role, user string) string {
	return fmt.Sprintf("%s-%s", role, user)
}

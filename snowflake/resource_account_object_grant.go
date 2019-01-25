package snowflake

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAccountObjectGrant() *schema.Resource {
	return &schema.Resource{
		Create: createAccountObjectGrant,
		Update: nil,
		Read:   readAccountObjectGrant,
		Delete: deleteAccountObjectGrant,

		Schema: map[string]*schema.Schema{
			"object_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"object_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"privileges": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"grant_option": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
		},
	}
}

func createAccountObjectGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	stmtSQL := fmt.Sprintf("GRANT %s ON %s \"%s\" TO ROLE \"%s\"",
		getPrivilegesString(d.Get("privileges").(*schema.Set)),
		d.Get("object_type").(string),
		d.Get("object_name").(string),
		d.Get("role").(string))

	if d.Get("grant_option").(bool) {
		stmtSQL += " WITH GRANT OPTION"
	}

	log.Println("Executing statement:", stmtSQL)
	_, err := db.Exec(stmtSQL)
	if err != nil {
		return err
	}

	id := generateGrantID(
		d.Get("object_type").(string),
		d.Get("object_name").(string),
		d.Get("role").(string))
	d.SetId(id)

	return readAccountObjectGrant(d, meta)
}

func readAccountObjectGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	objectType, objectName, role := getParamsFromGrantID(d.Id())

	stmtSQL := fmt.Sprintf("SHOW GRANTS ON %s \"%s\"",
		objectType,
		objectName)

	log.Println("Executing statement:", stmtSQL)

	rows, err := db.Query(stmtSQL)
	if err != nil {
		return err
	}

	defer rows.Close()

	var (
		createdOn         string
		privilege         string
		grantedOn         string
		name              string
		grantedTo         string
		granteeName       string
		grantOption       bool
		grantedBy         string
		privileges        []interface{}
		objectGrantOption bool
	)

	for rows.Next() {
		if err := rows.Scan(&createdOn, &privilege, &grantedOn, &name, &grantedTo, &granteeName, &grantOption, &grantedBy); err != nil {
			return err
		}

		if grantedTo == "ROLE" && granteeName == role {
			privileges = append(privileges, privilege)
			objectGrantOption = grantOption
		}
	}

	if len(privileges) > 0 {
		d.Set("objectType", objectType)
		d.Set("objectName", objectName)
		d.Set("role", role)
		d.Set("privileges", schema.NewSet(schema.HashString, privileges))
		d.Set("grant_option", objectGrantOption)
		return nil
	}

	return fmt.Errorf("The grant of role %s on %s %s does not exist.", role, objectType, objectName)
}

func deleteAccountObjectGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	objectType, objectName, role := getParamsFromGrantID(d.Id())

	stmtSQL := fmt.Sprintf("REVOKE ALL PRIVILEGES ON %s \"%s\" FROM ROLE \"%s\"",
		objectType,
		objectName,
		role)

	log.Println("Executing statement:", stmtSQL)
	_, err := db.Exec(stmtSQL)
	if err == nil {
		d.SetId("")
	}
	return err
}

func getPrivilegesString(priviligesSet *schema.Set) string {
	if priviligesSet.Contains("ALL") {
		return "ALL"
	}

	var privilegesList []string
	for _, v := range priviligesSet.List() {
		privilegesList = append(privilegesList, v.(string))
	}

	return strings.Join(privilegesList, ",")
}

func generateGrantID(objectType, objectName, role string) string {
	return fmt.Sprintf("%s-%s-%s", objectType, objectName, role)
}

func getParamsFromGrantID(id string) (objectType, objectName, role string) {
	params := strings.Split(id, "-")
	return params[0], params[1], params[2]
}

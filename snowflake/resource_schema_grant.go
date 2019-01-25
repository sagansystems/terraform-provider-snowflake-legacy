package snowflake

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSchemaGrant() *schema.Resource {
	return &schema.Resource{
		Create: createSchemaGrant,
		Update: nil,
		Read:   readSchemaGrant,
		Delete: deleteSchemaGrant,

		Schema: map[string]*schema.Schema{
			"schema": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"database": &schema.Schema{
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

func createSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	stmtSQL := fmt.Sprintf("GRANT %s ON %s TO ROLE \"%s\"",
		getPrivilegesString(d.Get("privileges").(*schema.Set)),
		generateSchemalalal(d.Get("schema").(string), d.Get("database").(string)),
		d.Get("role").(string))

	if d.Get("grant_option").(bool) {
		stmtSQL += " WITH GRANT OPTION"
	}

	log.Println("Executing statement:", stmtSQL)
	_, err := db.Exec(stmtSQL)
	if err != nil {
		return err
	}

	id := generateAccountObjectGrantID(
		d.Get("database").(string),
		d.Get("schema").(string),
		d.Get("role").(string))
	d.SetId(id)

	return readSchemaGrant(d, meta)
}

func readSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	database, schema, role := getParamsFromSchemaGrantID(d.Id())
	return nil

}

func deleteSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	database, schema, role := getParamsFromSchemaGrantID(d.Id())

	stmtSQL := fmt.Sprintf("REVOKE ALL PRIVILEGES ON %s FROM ROLE \"%s\"",
		generateSchemalalal(schema, database),
		role)

	log.Println("Executing statement:", stmtSQL)
	_, err := db.Exec(stmtSQL)
	if err == nil {
		d.SetId("")
	}
	return err
}

func generateSchemalalal(schema, database string) string {
	if schema == "ALL" {
		return fmt.Sprintf("ALL SCHEMAS IN DATABASE \"%s\"", database)
	}
	return fmt.Sprintf("SCHEMA \"%s.%s\"", database, schema)
}

func generateSchemaGrantID(database, schema, role string) string {
	return fmt.Sprintf("%s-%s-%s", database, schema, role)
}

func getParamsFromSchemaGrantID(id string) (database, schema, role string) {
	params := strings.Split(id, "-")
	return params[0], params[1], params[2]
}

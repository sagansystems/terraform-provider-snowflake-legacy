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
	schemaName := d.Get("schema").(string)
	databaseName := d.Get("database").(string)
	role := d.Get("role").(string)

	stmtSQL := fmt.Sprintf("GRANT %s ON %s TO ROLE \"%s\"",
		privilegesSetToString(d.Get("privileges").(*schema.Set)),
		generateRecipientSchemaString(schemaName, databaseName),
		role)

	if d.Get("grant_option").(bool) {
		stmtSQL += " WITH GRANT OPTION"
	}

	log.Println("Executing statement:", stmtSQL)
	_, err := db.Exec(stmtSQL)
	if err != nil {
		return err
	}

	id := generateSchemaGrantID(databaseName, schemaName, role)
	d.SetId(id)

	return readSchemaGrant(d, meta)
}

func readSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	databaseName, schemaName, role := getParamsFromSchemaGrantID(d.Id())

	stmtSQL := fmt.Sprintf("SHOW GRANTS TO ROLE \"%s\"", role)

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

		if grantedOn == "SCHEMA" && validateSchemaName(name, databaseName, schemaName) {
			privileges = append(privileges, privilege)
			objectGrantOption = grantOption
		}
	}

	if len(privileges) > 0 {
		d.Set("schema", schemaName)
		d.Set("database", databaseName)
		d.Set("role", role)
		d.Set("privileges", schema.NewSet(schema.HashString, privileges))
		d.Set("grant_option", objectGrantOption)
		return nil
	}

	return fmt.Errorf("The grant of role %s on %s does not exist.", role, generateRecipientSchemaString(schemaName, databaseName))
}

func deleteSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	databaseName, schemaName, role := getParamsFromSchemaGrantID(d.Id())

	stmtSQL := fmt.Sprintf("REVOKE ALL PRIVILEGES ON %s FROM ROLE \"%s\"",
		generateRecipientSchemaString(schemaName, databaseName),
		role)

	log.Println("Executing statement:", stmtSQL)
	_, err := db.Exec(stmtSQL)
	if err == nil {
		d.SetId("")
	}
	return err
}

func validateSchemaName(nameToValidate, databaseName, schemaName string) bool {
	databaseToValidate := strings.Split(nameToValidate, ".")[0]
	schemaToValidate := strings.Split(nameToValidate, ".")[1]

	if schemaName == "ALL" {
		return databaseToValidate == databaseName
	}
	return databaseToValidate == databaseName && schemaToValidate == schemaName
}

func generateRecipientSchemaString(schema, database string) string {
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

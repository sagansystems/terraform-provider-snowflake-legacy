package snowflake

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSchemaObjectGrant() *schema.Resource {
	return &schema.Resource{
		Create: createSchemaObjectGrant,
		Update: nil,
		Read:   readSchemaObjectGrant,
		Delete: deleteSchemaObjectGrant,

		Schema: map[string]*schema.Schema{
			"object_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"object_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"database": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"schema": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"future": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
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

func createSchemaObjectGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	txn, err := db.Begin()

	defer func() {
		_ = txn.Rollback()
	}()

	var (
		objectName   = d.Get("object_name").(string)
		objectType   = d.Get("object_type").(string)
		databaseName = d.Get("database").(string)
		schemaName   = d.Get("schema").(string)
		future       = d.Get("future").(bool)
		role         = d.Get("role").(string)
	)

	stmtSQL := fmt.Sprintf("USE DATABASE \"%s\"", databaseName)
	_, err = txn.Exec(stmtSQL)
	if err != nil {
		return err
	}

	stmtSQL = fmt.Sprintf("USE SCHEMA \"%s\"", schemaName)
	_, err = txn.Exec(stmtSQL)
	if err != nil {
		return err
	}

	stmtSQL = fmt.Sprintf("GRANT %s ON %s TO ROLE \"%s\"",
		privilegesSetToString(d.Get("privileges").(*schema.Set)),
		generateRecipientSchemaObjectString(objectType, objectName, schemaName, future),
		role)

	if d.Get("grant_option").(bool) {
		stmtSQL += " WITH GRANT OPTION"
	}

	log.Println("Executing statement:", stmtSQL)
	_, err = txn.Exec(stmtSQL)
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	id := generateSchemaObjectGrantID(objectType, objectName, databaseName, schemaName, role, future)
	d.SetId(id)

	return readSchemaObjectGrant(d, meta)
}

func readSchemaObjectGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	txn, err := db.Begin()
	objectType, objectName, databaseName, schemaName, role, future := getParamsFromSchemaObjectGrantID(d.Id())

	defer func() {
		_ = txn.Rollback()
	}()

	stmtSQL := fmt.Sprintf("USE DATABASE \"%s\"", databaseName)
	_, err = txn.Exec(stmtSQL)
	if err != nil {
		return err
	}

	stmtSQL = fmt.Sprintf("USE SCHEMA \"%s\"", schemaName)
	_, err = txn.Exec(stmtSQL)
	if err != nil {
		return err
	}

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

	if future {
		stmtSQL = fmt.Sprintf("SHOW FUTURE GRANTS IN SCHEMA \"%s\"", schemaName)

		log.Println("Executing statement:", stmtSQL)
		rows, err := txn.Query(stmtSQL)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&createdOn, &privilege, &grantedOn, &name, &grantedTo, &granteeName, &grantOption); err != nil {
				return err
			}

			if grantedTo == "ROLE" && granteeName == role && grantedOn == objectType {
				privileges = append(privileges, privilege)
				objectGrantOption = grantOption
			}
		}
	} else {
		stmtSQL = fmt.Sprintf("SHOW GRANTS TO ROLE \"%s\"", role)

		log.Println("Executing statement:", stmtSQL)
		rows, err := txn.Query(stmtSQL)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&createdOn, &privilege, &grantedOn, &name, &grantedTo, &granteeName, &grantOption, &grantedBy); err != nil {
				return err
			}

			if grantedOn == objectType && validateSchemaObjectName(name, databaseName, schemaName, objectName) {
				privileges = append(privileges, privilege)
				objectGrantOption = grantOption
			}
		}
	}

	if len(privileges) > 0 {
		d.Set("object_type", objectType)
		d.Set("object_name", objectName)
		d.Set("database", databaseName)
		d.Set("schema", schemaName)
		d.Set("future", future)
		d.Set("role", role)
		d.Set("privileges", schema.NewSet(schema.HashString, privileges))
		d.Set("grant_option", objectGrantOption)
		return nil
	}

	return fmt.Errorf("The grant of role %s on %s does not exist.", role, generateRecipientSchemaObjectString(objectType, objectName, schemaName, future))
}

func deleteSchemaObjectGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	txn, err := db.Begin()
	objectType, objectName, databaseName, schemaName, role, future := getParamsFromSchemaObjectGrantID(d.Id())

	defer func() {
		_ = txn.Rollback()
	}()

	stmtSQL := fmt.Sprintf("USE DATABASE \"%s\"", databaseName)
	_, err = txn.Exec(stmtSQL)
	if err != nil {
		return err
	}

	stmtSQL = fmt.Sprintf("USE SCHEMA \"%s\"", schemaName)
	_, err = txn.Exec(stmtSQL)
	if err != nil {
		return err
	}

	stmtSQL = fmt.Sprintf("REVOKE ALL PRIVILEGES ON %s FROM ROLE \"%s\"",
		generateRecipientSchemaObjectString(objectType, objectName, schemaName, future),
		role)

	log.Println("Executing statement:", stmtSQL)
	_, err = txn.Exec(stmtSQL)
	if err != nil {
		return nil
	}

	err = txn.Commit()
	if err == nil {
		d.SetId("")
	}

	return err
}

func validateSchemaObjectName(nameToValidate, databaseName, schemaName, objectName string) bool {
	databaseToValidate := strings.Split(nameToValidate, ".")[0]
	schemaToValidate := strings.Split(nameToValidate, ".")[1]
	objectNameToValidate := strings.Split(nameToValidate, ".")[2]

	if len(objectName) == 0 {
		return databaseToValidate == databaseName && schemaToValidate == schemaName
	}
	return databaseToValidate == databaseName && schemaToValidate == schemaName && objectNameToValidate == objectName
}

func generateRecipientSchemaObjectString(objectType, objectName, schema string, future bool) string {
	if future {
		return fmt.Sprintf("FUTURE %sS IN SCHEMA\"%s\"", objectType, schema)
	}

	if len(objectName) > 0 {
		return fmt.Sprintf("%s \"%s\"", objectType, objectName)
	}

	return fmt.Sprintf("ALL %sS IN SCHEMA \"%s\"", objectType, schema)
}

func generateSchemaObjectGrantID(objectType, objectName, database, schema, role string, future bool) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s-%s", objectType, objectName, strconv.FormatBool(future), database, schema, role)
}

func getParamsFromSchemaObjectGrantID(id string) (objectType, objectName, database, schema, role string, future bool) {
	params := strings.Split(id, "-")
	future, _ = strconv.ParseBool(params[2])
	return params[0], params[1], params[3], params[4], params[5], future
}

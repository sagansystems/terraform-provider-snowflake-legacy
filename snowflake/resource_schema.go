package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Create: createSchema,
		Read:   readSchema,
		Update: updateSchema,
		Delete: deleteSchema,

		Schema: map[string]*schema.Schema{
			"database": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the database in which to place schema",
			},

			"schema": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the schema to create",
			},
		},
	}
}

func createSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	database := d.Get("database").(string)
	schema := d.Get("schema").(string)

	d.SetId(schemaIDFromParams(database, schema))

	stmtSQL := fmt.Sprintf("CREATE SCHEMA \"%s\".\"%s\"", database, schema)

	log.Println("Executing statement:", stmtSQL)

	if _, err := db.Exec(stmtSQL); err != nil {
		return err
	}

	return readSchema(d, meta)
}

func readSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	database, schema := paramsFromSchemaID(d.Id())

	stmtSQL := fmt.Sprintf("SHOW TERSE SCHEMAS LIKE '%s' IN DATABASE \"%s\"", schema, database)

	log.Println("Executing statement:", stmtSQL)

	rows, err := db.Query(stmtSQL)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var dbCreatedAt string
		var dbName string
		var dbKind sql.NullString
		var dbDBName string
		var dbSchemaName sql.NullString
		if err := rows.Scan(&dbCreatedAt, &dbName, &dbKind, &dbDBName, &dbSchemaName); err != nil {
			return err
		}
		// The SHOW TERSE SCHEMAS LIKE will return case-insensitive matches but schema names are case-sensitive so we
		// need to make sure we have what we expect.
		if database == dbDBName && schema == dbName {
			d.Set("schema", dbName)
			d.Set("database", dbDBName)
			return nil
		}
	}

	// the terraform thing to do if the resource does not exist is set id to the empty string
	d.SetId("")
	return nil
}

func deleteSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	database, schema := paramsFromSchemaID(d.Id())

	stmtSQL := fmt.Sprintf("DROP SCHEMA \"%s\".\"%s\"", database, schema)

	log.Println("Executing statement:", stmtSQL)

	if _, err := db.Exec(stmtSQL); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func updateSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	oldDB, oldSchema := paramsFromSchemaID(d.Id())

	var newDB string
	if d.HasChange("database") {
		_, newDBI := d.GetChange("database")
		newDB = newDBI.(string)
	} else {
		newDB = oldDB
	}

	var newSchema string
	if d.HasChange("schema") {
		_, newSchemaI := d.GetChange("schema")
		newSchema = newSchemaI.(string)
	} else {
		newSchema = oldSchema
	}

	stmtSQL := fmt.Sprintf("ALTER SCHEMA \"%s\".\"%s\" RENAME TO \"%s\".\"%s\"", oldDB, oldSchema, newDB, newSchema)
	d.SetId(schemaIDFromParams(newDB, newSchema))

	if _, err := db.Exec(stmtSQL); err != nil {
		return err
	}

	return nil
}

func paramsFromSchemaID(id string) (database, schema string) {
	splits := strings.Split(id, "-")
	return splits[0], splits[1]
}

func schemaIDFromParams(database, schema string) string {
	return fmt.Sprintf("%s-%s", database, schema)
}

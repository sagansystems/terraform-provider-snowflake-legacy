package snowflake

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTable() *schema.Resource {
	return &schema.Resource{
		Create: createTable,
		Update: updateTable,
		Read:   readTable,
		Delete: deleteTable,

		Schema: map[string]*schema.Schema{
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Name of the database to use",
			},
			"schema": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the schema to use",
			},
			"table": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the table to use",
			},
			"columns": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name and type of columns",
			},
		},
	}
}

func createTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	table := d.Get("table").(string)
	columns := d.Get("columns").(string)

	stmtSQL := fmt.Sprintf(`CREATE TABLE "%s"."%s"."%s" (%s)`,
		database, schema, table, columns,
	)
	d.SetId(tableIDFromParams(database, schema, table))

	log.Println("Executing statement:", stmtSQL)

	if _, err := db.Exec(stmtSQL); err != nil {
		return err
	}

	return readTable(d, meta)
}

func updateTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	table := d.Get("table").(string)
	columns := d.Get("columns").(string)

	stmtSQL := fmt.Sprintf(`REPLACE TABLE "%s"."%s"."%s" (%s)`,
		database, schema, table, columns,
	)
	d.SetId(tableIDFromParams(database, schema, table))

	log.Println("Executing statement:", stmtSQL)

	if _, err := db.Exec(stmtSQL); err != nil {
		return err
	}

	return readTable(d, meta)
}

func readTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	database, schema, table := paramsFromTableID(d.Id())

	stmtSQL := fmt.Sprintf("SHOW TABLES LIKE '%s' IN DATABASE \"%s\" SCHEMA \"%s\"", table, database, schema)

	log.Println("Executing statement:", stmtSQL)

	rows, err := db.Query(stmtSQL)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var tableCreatedAt string
		var tableName string
		var dbName string
		var schemaName string
		if err := rows.Scan(&tableCreatedAt, &tableName, &dbName, &schemaName); err != nil {
			return err
		}
		if database == dbName && schema == schemaName && table == tableName {
			d.Set("database", dbName)
			d.Set("schema", schemaName)
			d.Set("table", tableName)
			return nil
		}
	}

	return fmt.Errorf("the table %s.%s.%s does not exist", database, schema, table)
}

func deleteTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	database, schema, table := paramsFromTableID(d.Id())

	stmtSQL := fmt.Sprintf("DROP TABLE \"%s\".\"%s\".\"%s\"", database, schema, table)

	log.Println("Executing statement:", stmtSQL)

	if _, err := db.Exec(stmtSQL); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func paramsFromTableID(id string) (database, schema, table string) {
	splits := strings.Split(id, "-")
	return splits[0], splits[1], splits[2]
}

func tableIDFromParams(database, schema, table string) string {
	return fmt.Sprintf("%s-%s-%s", database, schema, table)
}

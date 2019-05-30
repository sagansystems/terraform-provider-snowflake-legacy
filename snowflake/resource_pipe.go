package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePipe() *schema.Resource {
	return &schema.Resource{
		Create: createPipe,
		Update: updatePipe,
		Read:   readPipe,
		Delete: deletePipe,

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
			"stage": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the stage to use",
			},
			"table": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the table to copy into",
			},
			"pipe": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the pipe to create",
			},
			"auto_ingest": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "specifies whether auto ingest is on",
				Default:     false,
			},
		},
	}
}

func createPipe(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	stage := d.Get("stage").(string)
	table := d.Get("table").(string)
	pipe := d.Get("pipe").(string)
	autoIngest := d.Get("auto_ingest").(bool)

	stmtSQL := fmt.Sprintf(`
		CREATE PIPE "%s"."%s"."%s" auto_ingest=%t AS
			COPY INTO "%s"."%s"."%s"
			FROM @"%s"."%s"."%s"
			file_format = (type = 'JSON')
			`,
		database, schema, pipe, autoIngest,
		database, schema, table,
		database, schema, stage,
	)
	d.SetId(pipeIDFromParams(database, schema, pipe))

	log.Println("Executing statement:", stmtSQL)

	if _, err := db.Exec(stmtSQL); err != nil {
		return err
	}

	return readPipe(d, meta)
}

func updatePipe(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	stage := d.Get("stage").(string)
	table := d.Get("table").(string)
	pipe := d.Get("pipe").(string)
	autoIngest := d.Get("autoIngest").(bool)

	stmtSQL := fmt.Sprintf(`
		REPLACE PIPE "%s"."%s"."%s" auto_ingest=%t AS
			COPY INTO "%s"."%s"."%s"
			FROM @"%s"."%s"."%s"
			file_format = (type = 'JSON')
			`,
		database, schema, pipe, autoIngest,
		database, schema, table,
		database, schema, stage,
	)
	d.SetId(pipeIDFromParams(database, schema, pipe))

	log.Println("Executing statement:", stmtSQL)

	if _, err := db.Exec(stmtSQL); err != nil {
		return err
	}

	return readPipe(d, meta)
}

func readPipe(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	database, schema, pipe := paramsFromPipeID(d.Id())

	stmtSQL := fmt.Sprintf("SHOW PIPES LIKE '%s' IN DATABASE \"%s\" SCHEMA \"%s\"", pipe, database, schema)

	log.Println("Executing statement:", stmtSQL)

	rows, err := db.Query(stmtSQL)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var pipeCreatedAt string
		var pipeName string
		var dbName string
		var schemaName string
		var definition sql.NullString
		var owner sql.NullString
		var notificationChannel string
		if err := rows.Scan(&pipeCreatedAt, &pipeName, &dbName, &schemaName, &definition, &owner, &notificationChannel); err != nil {
			return err
		}
		// The SHOW TERSE SCHEMAS LIKE will return case-insensitive matches but schema names are case-sensitive so we
		// need to make sure we have what we expect.
		if database == dbName && schema == schemaName && pipe == pipeName {
			d.Set("database", dbName)
			d.Set("schema", schemaName)
			d.Set("pipe", pipeName)
			d.Set("notificationChannel", notificationChannel)
			return nil
		}
	}

	return fmt.Errorf("the pipe %s.%s.%s does not exist.", database, schema, pipe)
}

func deletePipe(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	database, schema, pipe := paramsFromPipeID(d.Id())

	stmtSQL := fmt.Sprintf("DROP PIPE \"%s\".\"%s\".\"%s\"", database, schema, pipe)

	log.Println("Executing statement:", stmtSQL)

	if _, err := db.Exec(stmtSQL); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func paramsFromPipeID(id string) (database, schema, pipe string) {
	splits := strings.Split(id, "-")
	return splits[0], splits[1], splits[2]
}

func pipeIDFromParams(database, schema, pipe string) string {
	return fmt.Sprintf("%s-%s-%s", database, schema, pipe)
}

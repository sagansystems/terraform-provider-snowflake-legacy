package snowflake

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: createRole,
		Update: updateRole,
		Read:   readRole,
		Delete: deleteRole,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func createRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	stmtSQL := fmt.Sprintf("CREATE ROLE \"%s\"", d.Get("name").(string))

	if _, ok := d.GetOk("comment"); ok {
		stmtSQL = stmtSQL + fmt.Sprintf(" COMMENT = \"%s\"", d.Get("comment").(string))
	}

	log.Println("Executing statement:", stmtSQL)
	_, err := db.Exec(stmtSQL)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%s", d.Get("name").(string))
	d.SetId(name)

	return nil
}

func updateRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	var newName interface{}
	var oldName interface{}
	if d.HasChange("name") {
		oldName, newName = d.GetChange("name")
	} else {
		oldName = d.Get("name")
		newName = nil
	}

	var newComment interface{}
	if d.HasChange("comment") {
		_, newComment = d.GetChange("comment")
	} else {
		newComment = nil
	}

	var queries = make([]string, 0)
	if newComment != nil {
		if newComment.(string) == "" {
			queries = append(queries, fmt.Sprintf("ALTER ROLE \"%s\" UNSET COMMENT",
				oldName.(string)))
		} else {
			queries = append(queries, fmt.Sprintf("ALTER ROLE \"%s\" SET COMMENT = \"%s\"",
				oldName.(string),
				newComment.(string)))
		}
	}

	if newName != nil {
		queries = append(queries, fmt.Sprintf("ALTER ROLE \"%s\" RENAME TO \"%s\"",
			oldName.(string),
			newName.(string)))
	}

	// Execute queries in one transaction.
	if len(queries) > 0 {
		txn, err := db.Begin()
		if err != nil {
			return err
		}

		defer func() {
			_ = txn.Rollback()
		}()

		for _, query := range queries {
			log.Println("Executing statement:", query)

			_, err := txn.Exec(query)

			if err != nil {
				return err
			}
		}

		return txn.Commit()
	}

	return nil
}

func readRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	stmtSQL := fmt.Sprintf("SHOW ROLES LIKE '%s'", d.Get("name").(string))

	log.Println("Executing statement:", stmtSQL)

	rows, err := db.Query(stmtSQL)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() && rows.Err() == nil {
		d.SetId("")
	}
	return rows.Err()
}

func deleteRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	name := d.Get("name").(string)
	stmtSQL := fmt.Sprintf("DROP ROLE  %s ", name)

	_, err := db.Exec(stmtSQL)
	if err == nil {
		d.SetId("")
	}
	return err
}

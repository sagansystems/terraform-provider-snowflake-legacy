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

	name := d.Get("name").(string)
	d.SetId(name)

	return nil
}

func updateRole(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("comment") {
		return nil
	}

	var stmtSQL string
	db := meta.(*providerConfiguration).DB
	_, newComment := d.GetChange("comment")

	if newComment.(string) == "" {
		stmtSQL = fmt.Sprintf("ALTER ROLE \"%s\" UNSET COMMENT", d.Id())
	} else {
		stmtSQL = fmt.Sprintf("ALTER ROLE \"%s\" SET COMMENT = \"%s\"",
			d.Id(),
			newComment.(string))
	}

	log.Println("Executing statement:", stmtSQL)
	_, err := db.Exec(stmtSQL)
	if err != nil {
		return err
	}

	return nil
}

func readRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	stmtSQL := fmt.Sprintf("SHOW ROLES LIKE '%s'", d.Id())

	log.Println("Executing statement:", stmtSQL)

	rows, err := db.Query(stmtSQL)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			createdOn      string
			name           string
			isDefault      string
			isCurrent      string
			isInherited    string
			assignedTo     string
			grantedToRoles string
			grantedRoles   string
			owner          string
			comment        string
		)

		if err := rows.Scan(&createdOn, &name, &isDefault, &isCurrent, &isInherited, assignedTo, &assignedTo, &grantedToRoles, &grantedRoles, &owner, &comment); err != nil {
			return err
		}

		if name == d.Id() {
			d.Set("name", name)
			d.Set("comment", comment)
			return nil
		}
	}

	return fmt.Errorf("The role %s does not exist", d.Id())
}

func deleteRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	stmtSQL := fmt.Sprintf("DROP ROLE \"%s\"", d.Id())

	_, err := db.Exec(stmtSQL)
	if err == nil {
		d.SetId("")
	}
	return err
}

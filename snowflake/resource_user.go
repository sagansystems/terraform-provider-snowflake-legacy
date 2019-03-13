package snowflake

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: CreateUser,
		Update: UpdateUser,
		Read:   ReadUser,
		Delete: DeleteUser,

		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plaintext_password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				StateFunc: hashSum,
			},
			"rsa_public_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"password": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"plaintext_password"},
				Sensitive:     true,
				Deprecated:    "Please use plaintext_password instead",
			},
			"default_role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func CreateUser(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	stmtSQL := fmt.Sprintf("CREATE USER \"%s\"", d.Get("user").(string))

	var password string
	if v, ok := d.GetOk("plaintext_password"); ok {
		password = v.(string)
	} else {
		password = d.Get("password").(string)
	}

	if password != "" {
		stmtSQL = stmtSQL + fmt.Sprintf(" PASSWORD = \"%s\"", password)
	}

	if v, ok := d.GetOk("rsa_public_key"); ok {
		stmtSQL = stmtSQL + fmt.Sprintf(" RSA_PUBLIC_KEY = \"%s\"", v.(string))
	}

	if v, ok := d.GetOk("default_role"); ok {
		stmtSQL = stmtSQL + fmt.Sprintf(" DEFAULT_ROLE = \"%s\"", v.(string))
	}

	log.Println("Executing statement:", stmtSQL)
	_, err := db.Exec(stmtSQL)
	if err != nil {
		return err
	}

	user := fmt.Sprintf("%s", d.Get("user").(string))
	d.SetId(user)

	return nil
}

func UpdateUser(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(*providerConfiguration)

	var newpw interface{}
	if d.HasChange("plaintext_password") {
		_, newpw = d.GetChange("plaintext_password")
	} else if d.HasChange("password") {
		_, newpw = d.GetChange("password")
	} else {
		newpw = nil
	}

	var newdefrole interface{}
	if d.HasChange("default_role") {
		_, newdefrole = d.GetChange("default_role")
	} else {
		newdefrole = nil
	}

	var newRSAPublicKey interface{}
	if d.HasChange("rsa_public_key") {
		_, newRSAPublicKey = d.GetChange("rsa_public_key")
	} else {
		newRSAPublicKey = nil
	}

	if newpw != nil || newdefrole != nil || newRSAPublicKey != nil {
		stmtSQL := fmt.Sprintf("ALTER USER \"%s\" SET ", d.Get("user").(string))

		if newpw != nil {
			stmtSQL = stmtSQL + fmt.Sprintf(" PASSWORD = \"%s\"", newpw.(string))
		}

		if newRSAPublicKey != nil {
			stmtSQL = stmtSQL + fmt.Sprintf(" RSA_PUBLIC_KEY = \"%s\"", newRSAPublicKey.(string))
		}

		if newdefrole != nil {
			stmtSQL = stmtSQL + fmt.Sprintf(" DEFAULT_ROLE = \"%s\"", newdefrole.(string))
		}

		log.Println("Executing query:", stmtSQL)
		_, err := conf.DB.Exec(stmtSQL)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadUser(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	stmtSQL := fmt.Sprintf("SHOW USERS LIKE '%s'", d.Get("user").(string))

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

func DeleteUser(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	stmtSQL := fmt.Sprintf("DROP USER \"%s\"", d.Get("user").(string))

	log.Println("Executing statement:", stmtSQL)

	_, err := db.Exec(stmtSQL)
	if err == nil {
		d.SetId("")
	}
	return err
}

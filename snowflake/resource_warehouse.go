package snowflake

import (
	"database/sql"
	"fmt"
	"log"

	"bytes"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	whNameAttr    = "name"
	whCommentAttr = "comment"

	//Properties
	whSizeAttr           = "warehouse_size"
	whMaxClusterCount    = "max_cluster_count"
	whMinClusterCount    = "min_cluster_count"
	whAutoSuspend        = "auto_suspend"
	whAutoResume         = "auto_resume"
	whInitiallySuspended = "initially_suspended"
)

func resourceWarehouse() *schema.Resource {
	return &schema.Resource{
		Create: createWarehouse,
		Update: updateWarehouse,
		Read:   readWarehouse,
		Delete: deleteWarehouse,

		Schema: map[string]*schema.Schema{
			whNameAttr: {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         false,
				Description:      "Identifier for the Snowflake warehouse;must be unique for your account ",
				DiffSuppressFunc: ignoreCase,
			},
			whCommentAttr: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    false,
				Description: "Specifies a comment for the warehouse.",
			},
			whSizeAttr: {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "XSMALL",
				ForceNew:         false,
				Description:      "Specifies the size of virtual warehouse to create.",
				DiffSuppressFunc: ignoreCase,
			},
			whMaxClusterCount: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				ForceNew:    false,
				Description: "Specifies the maximum number of server clusters for the warehouse.",
			},
			whMinClusterCount: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				ForceNew:    false,
				Description: "Specifies the minimum number of server clusters for a multi-cluster warehouse. ",
			},
			whAutoSuspend: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				ForceNew:    false,
				Description: "Specifies the number of seconds of inactivity after which a warehouse is automatically suspended.",
			},
			whAutoResume: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    false,
				Description: "Specifies whether to automatically resume a warehouse when it is accessed by a SQL statement, ",
			},
			whInitiallySuspended: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    false,
				Description: "Specifies whether the warehouse is created initially in suspended state.",
			},
		},
	}
}

func createWarehouse(d *schema.ResourceData, meta interface{}) error {
	whName := d.Get(whNameAttr).(string)
	db := meta.(*providerConfiguration).DB
	b := bytes.NewBufferString("CREATE  WAREHOUSE IF NOT EXISTS ")
	fmt.Fprint(b, whName)
	fmt.Fprintf(b, " WITH ")
	for _, attr := range []string{whMaxClusterCount, whMinClusterCount, whAutoSuspend, whAutoResume, whInitiallySuspended} {
		fmt.Fprintf(b, " %s=%v ", attr, d.Get(attr))
	}
	// Wrap string values in quotes
	for _, attr := range []string{whSizeAttr, whCommentAttr} {
		fmt.Fprintf(b, " %s='%v' ", attr, d.Get(attr))
	}

	sql := b.String()
	if _, err := db.Exec(sql); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error creating warehouse sql(%s) \n %q: {{err}}", sql, whName), err)
	}
	d.SetId(whName)
	return readWarehouse(d, meta)
}

func updateWarehouse(d *schema.ResourceData, meta interface{}) error {
	whName := d.Get(whNameAttr).(string)
	db := meta.(*providerConfiguration).DB
	b := bytes.NewBufferString("ALTER WAREHOUSE IF EXISTS ")
	fmt.Fprint(b, whName)
	fmt.Fprintf(b, " SET ")
	for _, attr := range []string{whMaxClusterCount, whMinClusterCount, whAutoSuspend, whAutoResume} {
		fmt.Fprintf(b, " %s=%v ", attr, d.Get(attr))
	}
	// Wrap string values in quotes
	for _, attr := range []string{whSizeAttr, whCommentAttr} {
		fmt.Fprintf(b, " %s='%v' ", attr, d.Get(attr))
	}

	sql := b.String()
	if _, err := db.Exec(sql); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error Altering warehouse %q: {{err}}", whName), err)
	}
	d.SetId(whName)
	return readWarehouse(d, meta)
}

func readWarehouse(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	warehouseName := d.Id()
	stmtSQL := fmt.Sprintf("show warehouses like '%s'", warehouseName)

	fmt.Printf(" Read Warehouse Executing query: %s \n", stmtSQL)
	log.Println("Executing query:", stmtSQL)

	var name, size, comment string
	var minClusterCount, maxClusterCount int
	var autoResume bool
	var state, instanceType, startedClusters, running, queued sql.NullString
	var isDefault, isCurrent, autoSuspend, available, provisioning, quiescing, other sql.NullString
	var createdOn, resumedOn, updatedOn, owner, resourceMonitor sql.NullString
	var actives, pendings, failed, suspended, uuid, scalingPolicy sql.NullString

	err := db.QueryRow(stmtSQL).Scan(
		&name, &state, &instanceType, &size, &minClusterCount, &maxClusterCount, &startedClusters, &running, &queued,
		&isDefault, &isCurrent, &autoSuspend, &autoResume, &available, &provisioning, &quiescing, &other,
		&createdOn, &resumedOn, &updatedOn, &owner, &comment, &resourceMonitor,
		&actives, &pendings, &failed, &suspended, &uuid, &scalingPolicy,
	)
	if err != nil {
		return fmt.Errorf("Error during show create warehouse: %s", err)
	}

	//Properties
	d.Set(whNameAttr, name)
	d.Set(whSizeAttr, size)
	d.Set(whMinClusterCount, minClusterCount)
	d.Set(whMaxClusterCount, maxClusterCount)
	d.Set(whAutoSuspend, autoSuspend)
	d.Set(whAutoResume, autoResume)
	d.Set(whCommentAttr, comment)
	return nil
}

func deleteWarehouse(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	whName := d.Get(whNameAttr).(string)
	sql := fmt.Sprintf("DROP WAREHOUSE  %s ", whName)
	if _, err := db.Exec(sql); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error droping warehouse %q: {{err}}", whName), err)
	}
	return nil
}

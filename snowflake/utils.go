package snowflake

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func hashSum(contents interface{}) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(contents.(string))))
}

func generateGrantCreationSQLStmt(privileges *schema.Set, object, role string, grantOption bool) string {
	stmtSQL := fmt.Sprintf("GRANT %s ON \"%s\" TO ROLE \"%s\"",
		privilegesSetToString(privileges),
		object,
		role)

	if grantOption {
		stmtSQL += " WITH GRANT OPTION"
	}

	return stmtSQL
}

func privilegesSetToString(priviligesSet *schema.Set) string {
	if priviligesSet.Contains("ALL") {
		return "ALL"
	}

	var privilegesList []string
	for _, v := range priviligesSet.List() {
		privilegesList = append(privilegesList, v.(string))
	}

	return strings.Join(privilegesList, ",")
}

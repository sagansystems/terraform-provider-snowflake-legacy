resource "snowflake_warehouse" "warehouse_terraform" {
  name           = "dev_wh"
  warehouse_size = "SMALL"
  auto_resume    = false
  auto_suspend   = 600
  comment        = "terraform development warehouse"
}

resource "snowflake_database" "database_terraform" {
  name    = "dev_db"
  comment = "terraform development database"
}
resource "snowflake_schema" "schema_terraform" {
  schema   = "dev_schema"
  database = "dev_db"
}

resource "snowflake_table" "table_terraform" {
  database = "dev_db"
  schema   = "dev_schema"
  table    = "dev_table"
  columns  = "jsontext variant"
}

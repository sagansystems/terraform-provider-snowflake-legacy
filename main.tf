resource "snowflake_warehouse" "warehouse_terraform" {
  name           = "DEV_WH"
  warehouse_size = "SMALL"
  auto_resume    = false
  auto_suspend   = 600
  comment        = "terraform development warehouse"
}

resource "snowflake_database" "database_terraform" {
  name    = "DEV_DB"
  comment = "terraform development database"
}

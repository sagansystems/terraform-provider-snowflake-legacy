Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-$PROVIDER_NAME`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-$PROVIDER_NAME
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-$PROVIDER_NAME
$ make build
```

Using the provider
----------------------
## Fill in for each provider
```
go get github.com/snowflakedb/gosnowflake
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-$PROVIDER_NAME
...
```


Provider Configuration
--------------------------

Firstly, to use the provider you will need to create a user within Snowflake that can execute the resource requests made by Terraform

```sh
$ export SF_USER
$ export SF_PASSWORD
$ export SF_REGION
$ export SF_ACCOUNT
```

### Snowflake Warehouse Management
```
resource "snowflake_warehouse" "warehouse_terraform" {
      name           = "dev_wh"
      warehouse_size = "SMALL"
      auto_resume    = false
      auto_suspend   = 600
      comment        = "terraform development warehouse"
}
```

##### Properties
| Property | Description | Type | Required |
| ------ | ------ | ------ | ------ |
| `name` | Name of the Snowflake warehouse | String | TRUE |
| `max_concurrency_level` | Max concurrent SQL statements that can run on warehouse | Integer | FALSE |
| `statement_queued_timeout_in_seconds` | Time, in seconds, an SQL statement can be queued before being cancelled | Integer | FALSE |
| `statement_timeout_in_seconds` | Time, in seconds, after which an SQL statement will be terminated | Integer | FALSE |
| `warehouse_size` | Size of the warehouse | String | FALSE |
| `max_cluster_count` | Min number of warehouses | Integer | FALSE |
| `min_cluster_count` | Max number of warehouses | Integer | FALSE |
| `auto_resume` | Should warehouse should auto resume | Boolean | FALSE |
| `auto_suspend` | Number of seconds after which the warehouse should suspend | Integer | FALSE |
| `initially_suspended` | Should warehouse start off suspended  | Boolean | FALSE |
| `comment` | Additional comments | String | FALSE |

### Snowflake Database Management
```
resource "snowflake_database" "database_terraform" {
      name    = "dev_db"
      comment = "terraform development database"
}
```

##### Properties
| Property | Description | Type | Required |
| ------ | ------ | ------   | ------ |
| `name` | Name of the Snowflake database | String | TRUE |
| `comment` | Additional comments | String | FALSE |

### Snowflake Schema Management
```
resource "snowflake_schema" "default" {
      database = "dev_db"
      schema   = "default"
}
```

##### Properties
| Property | Description | Type | Required |
| ------ | ------ | ------   | ------ |
| `database` | Database in which schema should be created| String | TRUE |
| `schema` | Name of the schema | String | TRUE |


### Snowflake User Management
```
resource "snowflake_user" "tf_test_user" {
  user               = "terraform.test"
  plaintext_password = "12345QWERTYqwerty"
  rsa_public_key     = "MIIBIjANBgkqhkiG9w...AQAB"
  default_role       = "READONLY"
}
```

##### Properties
| Property | Description | Type | Required |
| ------ | ------ | ------ | ------ |
| `user` | The username of the user | String | TRUE |
| `plaintext_password` | Password of the user. Ensure that passwords conform to the complexity requirements by Snowflake | String | FALSE |
| `rsa_public_key` | RSA public key to associate with the user. | String | FALSE |
| `default_role` | Default role the user assumes. Defaults to `null` | String | FALSE |

### Snowflake Role Management
```
resource "snowflake_role" "tf_test_role" {
  name    = "EXAMPLE_ROLE"
  comment = "example role"
}
```

##### Properties
| Property | Description | Type | Required |
| ------ | ------ | ------ | ------ |
| `name` | The name of the role | String | TRUE |
| `comment` | Additional comments | String | FALSE |

### Snowflake Role Grant Management
```
resource "snowflake_role_grant" "tf_test_role_grant" {
  role = "ACCOUNTADMIN"
  user = "tf_test_user"
}
```

##### Properties
| Property | Description | Type | Required |
| ------ | ------ | ------ | ------ |
| `role` | The role to grant | String | TRUE |
| `user` | The user to which to grant the role| String | TRUE |

### Snowflake Grant Management
Please note that that mixing grant management for a role with and without terraform is strongly discouraged, since it may result in terraform deleting grants which were added outside of terraform when performing updates or deletes.

### Snowflake Account Object Grant Management
```
resource "snowflake_account_object_grant" "tf_test_grant" {
  object_type  = "DATABASE"
  object_name  = "EXAMPLE_NAME"
  privileges   = ["MODIFY"]
  role         = "EXAMPLE_ROLE"
  grant_option = false
}
```

##### Properties
| Property | Description | Type | Required |
| ------ | ------ | ------ | ------ |
| `object_type` | Type of the object: DATABASE, WAREHOUSE, RESOURCE MONITOR | String | TRUE |
| `object_name` | The name of the object | String | TRUE |
| `privileges` | Privileges to grant (["ALL"] for all privileges) | String set | FALSE |
| `role` | The role to which the privileges are granted | String | TRUE |
| `grant_option` | Allows the recipient role to grant the privileges to other roles | Boolean | FALSE |

### Snowflake Schema Grant Management
```
resource "snowflake_schema_grant" "tf_test_grant" {
  database     = "DATABASE"
  schema       = "EXAMPLE_SCHEMA"
  privileges   = ["MODIFY"]
  role         = "EXAMPLE_ROLE"
  grant_option = false
}
```

##### Properties
| Property | Description | Type | Required |
| ------ | ------ | ------ | ------ |
| `schema` | The name of the schema ("ALL" if changes should be applied to all) | String | TRUE |
| `database` | The name of the database | String | TRUE |
| `privileges` | Privileges to grant (["ALL"] for all privileges) | String set | FALSE |
| `role` | The role to which the privileges are granted | String | TRUE |
| `grant_option` | Allows the recipient role to grant the privileges to other roles | Boolean | FALSE |
# dbac ([dee-BEK] the coffee) Database Command Line Interface (CLI) Documentation
This document provides an overview of the functionalities provided by the Database Command Line Interface (CLI) tool, which supports interactions with MySQL and PostgreSQL databases. The CLI allows users to manage databases, users, and permissions through various commands.

<img src="assets/dbac-logo.png" width="200">

## Table of Contents
- [dbac (\[dee-BEK\] the coffee) Database Command Line Interface (CLI) Documentation](#dbac-dee-bek-the-coffee-database-command-line-interface-cli-documentation)
  - [Table of Contents](#table-of-contents)
  - [Installation - Build from Source](#installation---build-from-source)
    - [Prerequisites](#prerequisites)
    - [Build Instructions](#build-instructions)
  - [Usage](#usage)
    - [General Commands](#general-commands)
    - [Profile Commands](#profile-commands)
    - [Database Commands](#database-commands)
    - [User Commands](#user-commands)
    - [Permission Commands](#permission-commands)
    - [Query Commands](#query-commands)
    - [Batch Command Functions](#batch-command-functions)
      - [YAML Configuration Example](#yaml-configuration-example)
  - [Examples](#examples)
  - [Functions Overview](#functions-overview)
    - [MySQL Functions](#mysql-functions)
    - [PostgreSQL Functions](#postgresql-functions)
  - [Contributing](#contributing)
  - [License](#license)
## Installation
### Homebrew
You can install dbac using Homebrew:
```sh
brew install bestcloudforme/homebrew-tap/dbac
```
### Build from Source
To build the CLI tool from source, you will need to have `go` installed on your machine. The source code also uses `make` to simplify the build and installation process.
#### Prerequisites
- Install [Go](https://golang.org/dl/) (version 1.14 or higher)
- Install [Make](https://www.gnu.org/software/make/)
#### Build Instructions
1. Clone the repository:
```bash
git clone https://github.com/bestcloudforme/dbac.git
cd dbac
```
2. Build the CLI:
```bash
make build
```
3. Install the Binary:
If you want to use `dbac` system-wide, you can move the binary to `/usr/local/bin`. This step may require superuser privileges.
```bash
sudo mv .bin/dbac /usr/local/bin
```
## Usage
The CLI tool provides various commands for managing databases, users, and permissions. Below is a detailed guide on how to use these commands.
### General Commands
- `dbac init`: Initialize the CLI configuration.
### Profile Commands
- `dbac profile current`: Display the current profile.
- `dbac profile list`: List all profiles.
- `dbac profile switch [PROFILE-NAME]`: Switch to a specified profile.
- `dbac profile add --db-type [DB-TYPE] --host [HOST] --user [USER] --port [PORT] --password [PASSWORD] --database [DATABASE] --profile-name [NAME]`: Add a new profile with specified details.
- `dbac profile add --file [FILE]`: Add a new profile from a file.
- `dbac profile delete --profile-name [PROFILE-NAME]`: Delete a specified profile.
### Database Commands
- `dbac database ping`: Ping the current database to check the connection.
- `dbac database list-databases`: List all databases.
- `dbac database list-tables`: List all tables in the current database.
- `dbac database create-database --database [DATABASE]`: Create a new database.
- `dbac database delete-database --database [DATABASE]`: Delete an existing database.
- `dbac database dump --filepath [FILEPATH] --table [TABLE] [--allTables] --database [DATABASE]`: Dump the specified database to the specified path.
### User Commands
- `dbac database list-user`: List all users.
- `dbac database create-user --username [USERNAME] --user-password [PASSWORD]`: Create a new database user.
- `dbac database delete-user --username [USERNAME]`: Delete an existing database user.
- `dbac database change-password --username [USERNAME] --new-password [NEW_PASSWORD]`: Change the password of an existing user.
### Permission Commands
- `dbac database grant-database --username [USERNAME] --permission [PERMISSION] --database [DATABASE]`: Grant database-level permissions to a user.
- `dbac database revoke-database --username [USERNAME] --permission [PERMISSION] --database [DATABASE]`: Revoke database-level permissions from a user.
- `dbac database grant-table --username [USERNAME] --permission [PERMISSION] --table [TABLE]`: Grant table-level permissions to a user.
- `dbac database revoke-table --username [USERNAME] --permission [PERMISSION] --table [TABLE]`: Revoke table-level permissions from a user.
### Query Commands
- `dbac database exec --query [QUERY]`: Execute a SQL query.
- `dbac database exec --file [FILE]`: Execute SQL commands from a file.
  - With file `some.sql`
```sql
CREATE DATABASE test5;
CREATE DATABASE test6;
CREATE DATABASE test7;
CREATE DATABASE test8;
```
### Batch Command Functions
This section describes the batch command functions integrated in the CLI. These functions allow for batch processing of database operations such as user and database management, permission handling, and executing SQL commands based on YAML configuration files.
- `dbac batch --file [FILE]`: Execute `dbac` commands from a file.
#### YAML Configuration Example
An example YAML configuration for a batch operation might look like this:
```yaml
steps:
  - type: create-user
    description: Create a new database user
    username: john_doe
    password: securepassword123
    profile: default
  - type: grant-database
    description: Grant all privileges
    username: john_doe
    permission: ALL
    database: test_db
    profile: default
```
This batch system simplifies complex database operations, ensuring they are executed in a controlled and sequential manner based on the defined steps in the YAML configuration.
## Examples
1. **Ping a database**:
```bash
dbac database ping
```
2. **List all databases**:
```bash
dbac database list-databases
```
3. **Create a new database**:
```bash
dbac database create-database --database testdb
```
4. **Create a new user**:
```bash
dbac database create-user --username testuser --user-password testpass
```
5. **Grant permissions to a user**:
```bash
dbac database grant-database --username testuser --permission ALL PRIVILEGES --database testdb
```
6. **Batch operations**:
```bash
dbac batch --file="batch.yaml"
```
## Functions Overview
### MySQL Functions
- **NewConnection**: Establishes a new connection to the MySQL database.
- **Ping**: Checks the MySQL database connection.
- **Close**: Closes the MySQL database connection.
- **Exec**: Executes a given SQL query on the MySQL database.
- **FileExec**: Executes SQL commands from a file on the MySQL database.
- **ListUsers**: Lists all users in the MySQL database.
- **CreateUser**: Creates a new user in the MySQL database.
- **DeleteUser**: Deletes a user from the MySQL database.
- **ChangeUserPassword**: Changes the password of a MySQL user.
- **GrantPermissions**: Grants database-level permissions to a MySQL user.
- **RevokePermissions**: Revokes database-level permissions from a MySQL user.
- **GrantTablePermissions**: Grants table-level permissions to a MySQL user.
- **RevokeTablePermissions**: Revokes table-level permissions from a MySQL user.
- **Dump**: Creates a dump of the MySQL database, allowing data to be saved in a specified path.
### PostgreSQL Functions
- **NewConnection**: Establishes a new connection to the PostgreSQL database.
- **Ping**: Checks the PostgreSQL database connection.
- **Close**: Closes the PostgreSQL database connection.
- **Exec**: Executes a given SQL query on the PostgreSQL database.
- **FileExec**: Executes SQL commands from a file on the PostgreSQL database.
- **ListUsers**: Lists all users in the PostgreSQL database.
- **CreateUser**: Creates a new user in the PostgreSQL database.
- **DeleteUser**: Deletes a user from the PostgreSQL database.
- **ChangeUserPassword**: Changes the password of a PostgreSQL user.
- **GrantPermissions**: Grants database-level permissions to a PostgreSQL user.
- **RevokePermissions**: Revokes database-level permissions from a PostgreSQL user.
- **GrantTablePermissions**: Grants table-level permissions to a PostgreSQL user.
- **RevokeTablePermissions**: Revokes table-level permissions from a PostgreSQL user.
- **ListDatabases**: Lists all databases in the PostgreSQL server.
- **ListTables**: Lists all tables in the current PostgreSQL database.
- **CreateDatabase**: Creates a new database in the PostgreSQL server.
- **DeleteDatabase**: Deletes a database from the PostgreSQL server.
- **Dump**: Creates a dump of the PostgreSQL database, enabling the backup of data to a specified path.
## Contributing
We welcome contributions! Please fork the repository and submit pull requests.
## License
This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.

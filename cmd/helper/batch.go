package helper

import (
	"dbac/cmd/mysql"
	"dbac/cmd/psql"
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Batch struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
}

type Steps struct {
	Steps []Batch `yaml:"steps"`
}

type CreateUser struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Profile     string `yaml:"profile"`
}

type CreateUserSteps struct {
	Steps []CreateUser `yaml:"steps"`
}

type GrantDatabase struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Username    string `yaml:"username"`
	Permission  string `yaml:"permission"`
	Database    string `yaml:"database"`
	Profile     string `yaml:"profile"`
}

type GrantDatabaseSteps struct {
	Steps []GrantDatabase `yaml:"steps"`
}

type ChangePassword struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Username    string `yaml:"username"`
	NewPassword string `yaml:"newPassword"`
	Profile     string `yaml:"profile"`
}

type ChangePasswordSteps struct {
	Steps []ChangePassword `yaml:"steps"`
}

type RevokeDatabase struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Username    string `yaml:"username"`
	Permission  string `yaml:"permission"`
	Profile     string `yaml:"profile"`
	Database    string `yaml:"database"`
}

type RevokeDatabaseSteps struct {
	Steps []RevokeDatabase `yaml:"steps"`
}

type DeleteUser struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Username    string `yaml:"username"`
	Profile     string `yaml:"profile"`
}

type DeleteUserSteps struct {
	Steps []DeleteUser `yaml:"steps"`
}

type DeleteDatabase struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Database    string `yaml:"database"`
	Profile     string `yaml:"profile"`
}

type DeleteDatabaseSteps struct {
	Steps []DeleteDatabase `yaml:"steps"`
}

type ListDatabases struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Profile     string `yaml:"profile"`
}

type ListDatabasesSteps struct {
	Steps []ListDatabases `yaml:"steps"`
}

type ListTables struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Profile     string `yaml:"profile"`
}

type ListTablesSteps struct {
	Steps []ListTables `yaml:"steps"`
}

type ListUsers struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Profile     string `yaml:"profile"`
}

type ListUsersSteps struct {
	Steps []ListUsers `yaml:"steps"`
}

type GrantTable struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Username    string `yaml:"username"`
	Permission  string `yaml:"permission"`
	Table       string `yaml:"table"`
	Profile     string `yaml:"profile"`
}

type GrantTableSteps struct {
	Steps []GrantTable `yaml:"steps"`
}

type RevokeTable struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Username    string `yaml:"username"`
	Permission  string `yaml:"permission"`
	Profile     string `yaml:"profile"`
	Table       string `yaml:"table"`
}

type RevokeTableSteps struct {
	Steps []RevokeTable `yaml:"steps"`
}

type CreateDatabase struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Database    string `yaml:"database"`
	Profile     string `yaml:"profile"`
}

type CreateDatabaseSteps struct {
	Steps []CreateDatabase `yaml:"steps"`
}

type ExecCommand struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Query       string `yaml:"query"`
	File        string `yaml:"file"`
	Profile     string `yaml:"profile"`
}

type ExecCommandSteps struct {
	Steps []ExecCommand `yaml:"steps"`
}

func RunBatch(file string) {
	var stepper Steps
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &stepper)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for i := 0; i < len(stepper.Steps); i++ {
		App(stepper.Steps[i].Type, i, file)
	}
}

func App(param string, step int, file string) {
	switch param {
	case "create-user":
		var currentProfile Profile
		var stepper CreateUserSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}

		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)

			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.CreateUser(stepper.Steps[step].Username, stepper.Steps[step].Password)
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.CreateUser(stepper.Steps[step].Username, stepper.Steps[step].Password)
			mysql.Close()
		}

	case "create-database":
		var currentProfile Profile
		var stepper CreateDatabaseSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)

			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.CreateDatabase(stepper.Steps[step].Database)
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.CreateDatabase(stepper.Steps[step].Database)
			mysql.Close()
		}

	case "grant-database":
		var currentProfile Profile
		var stepper GrantDatabaseSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.GrantPermissions(stepper.Steps[step].Database, stepper.Steps[step].Username, stepper.Steps[step].Permission)
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.GrantPermissions(stepper.Steps[step].Database, stepper.Steps[step].Username, stepper.Steps[step].Permission)
			mysql.Close()
		}

	case "change-password":
		var currentProfile Profile
		var stepper ChangePasswordSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.ChangeUserPassword(stepper.Steps[step].Username, stepper.Steps[step].NewPassword)
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.ChangeUserPassword(stepper.Steps[step].Username, stepper.Steps[step].NewPassword)
			mysql.Close()
		}

	case "revoke-database":
		var currentProfile Profile
		var stepper RevokeDatabaseSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.RevokePermissions(stepper.Steps[step].Database, stepper.Steps[step].Username, stepper.Steps[step].Permission)
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.RevokePermissions(stepper.Steps[step].Database, stepper.Steps[step].Username, stepper.Steps[step].Permission)
			mysql.Close()
		}

	case "delete-user":
		var currentProfile Profile
		var stepper DeleteUserSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.DeleteUser(stepper.Steps[step].Username)
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.DeleteUser(stepper.Steps[step].Username)
			mysql.Close()
		}

	case "delete-database":
		var currentProfile Profile
		var stepper DeleteDatabaseSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.DeleteDatabase(stepper.Steps[step].Database)
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.DeleteDatabase(stepper.Steps[step].Database)
			mysql.Close()
		}

	case "exec":
		var currentProfile Profile
		var stepper ExecCommandSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}
		if stepper.Steps[step].File == "" {
			if currentProfile.DbType == "psql" {
				db_port, _ := strconv.Atoi(currentProfile.Port)
				psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
				psql.Exec(stepper.Steps[step].Query)
				psql.Close()
			} else if currentProfile.DbType == "mysql" {
				mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
				mysql.Exec(stepper.Steps[step].Query)
				mysql.Close()
			}
		} else {
			if currentProfile.DbType == "psql" {
				db_port, _ := strconv.Atoi(currentProfile.Port)
				psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
				psql.FileExec(stepper.Steps[step].File)
				psql.Close()
			} else if currentProfile.DbType == "mysql" {
				fmt.Println(stepper.Steps[step].File)
				mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
				mysql.FileExec(stepper.Steps[step].File)
				mysql.Close()
			}
		}

	case "revoke-table":
		var currentProfile Profile
		var stepper RevokeTableSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.RevokeTablePermissions(stepper.Steps[step].Table, stepper.Steps[step].Username, stepper.Steps[step].Permission)
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.RevokeTablePermissions(currentProfile.Database, stepper.Steps[step].Table, stepper.Steps[step].Username, stepper.Steps[step].Permission)
			mysql.Close()
		}

	case "grant-table":
		var currentProfile Profile
		var stepper GrantTableSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.GrantTablePermissions(stepper.Steps[step].Table, stepper.Steps[step].Username, stepper.Steps[step].Permission)
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.GrantTablePermissions(currentProfile.Database, stepper.Steps[step].Table, stepper.Steps[step].Username, stepper.Steps[step].Permission)
			mysql.Close()
		}

	case "list-databases":
		var currentProfile Profile
		var stepper ListDatabasesSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.ListDatabases()
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.ListDatabases()
			mysql.Close()
		}
	case "list-tables":
		var currentProfile Profile
		var stepper ListTablesSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.ListTables()
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.ListTables()
			mysql.Close()
		}
	case "list-users":
		var currentProfile Profile
		var stepper ListUsersSteps
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &stepper)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		if stepper.Steps[step].Profile == "" {
			currentProfileName, _ := GetCurrentProfileName()
			currentProfile = ReadProfile(currentProfileName)
		} else {
			currentProfile = ReadProfile(stepper.Steps[step].Profile)
		}

		if currentProfile.DbType == "psql" {
			db_port, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, db_port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.ListUsers()
			psql.Close()
		} else if currentProfile.DbType == "mysql" {
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.ListUsers()
			mysql.Close()
		}
	}
}

package functions

import (
	"dbac/cmd/helper"
	"dbac/cmd/mysql"
	"dbac/cmd/psql"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func AddPingCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "ping",
		Short:   "Ping your database",
		Long:    `This command pings your database to ensure that you can reach your database instance correctly.`,
		Example: "dbac db ping",
		Args:    cobra.NoArgs,
		Run:     runPingDatabase,
	}
	subcommand.AddCommand(cmd)
}

func runPingDatabase(cmd *cobra.Command, args []string) {
	currentProfileName, err := helper.GetCurrentProfileName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current profile name: %v\n", err)
		os.Exit(1)
	}
	profile := helper.ReadProfile(currentProfileName)
	switch profile.DbType {
	case "postgres":
		dbPort, _ := strconv.Atoi(profile.Port)
		psql.NewConnection(profile.Host, dbPort, profile.User, profile.Password, profile.Database)
		psql.Ping()
		psql.Close()
	case "mysql":
		mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
		mysql.Ping()
		mysql.Close()
	}
}

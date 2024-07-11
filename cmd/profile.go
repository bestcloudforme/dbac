package cmd

import (
	functions "dbac/cmd/functions/profile"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage database connection profiles",
	Long:  `Manage your database connection profiles, including creating, listing, switching, and deleting profiles.`,
	Example: `
# Add a new database profile
dbac profile add myprofile --db-type mysql --user admin --pass secure123 --host db.example.com --port 3306 --database sample_db

# List all database profiles
dbac profile list

# Switch to a specific database profile
dbac profile switch myprofile

# Show the current active database profile
dbac profile current

# Delete a database profile
dbac profile delete myprofile
`,
}

func init() {
	rootCmd.AddCommand(profileCmd)

	functions.AddProfileCommand(profileCmd)
	functions.ListProfileCommand(profileCmd)
	functions.SwitchProfileCommand(profileCmd)
	functions.CurrentProfileCommand(profileCmd)
	functions.DeleteProfileCommand(profileCmd)
}

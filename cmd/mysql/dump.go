package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/JamesStewy/go-mysqldump"
)

func Dump2(path string, database string) error {
	dumpDir := path
	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", database)
	dumper, err := mysqldump.Register(DbConnection, dumpDir, dumpFilenameFormat)
	if err != nil {
		fmt.Printf("Error registering database: %v\n", err)
		return err
	}

	defer dumper.Close()

	resultFilename, err := dumper.Dump()
	if err != nil {
		fmt.Printf("Error dumping database: %v\n", err)
		return err
	}

	fmt.Printf("Dump successful, file saved to: %s\n", resultFilename)
	return nil
}

func Dump(path string, database string, filename string, table string, allTables bool) error {
	file, err := os.Create(path + filename)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "--\n-- dbac MySQL database dump\n--\n")
	fmt.Fprintf(file, "-- Database: %s\n-- ------------------------------------------------------\n", database)
	writeSessionVariables(file)

	if table != "" {
		dumpTable(DbConnection, table, file)
	} else if allTables {
		dumpAllTables(DbConnection, file, database)
	} else {
		log.Fatal("No table specified and --all-tables flag is not set. Nothing to dump.")
	}

	writeSessionRestore(file)
	fmt.Fprintf(file, "\n-- Dump completed on %s\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

func dumpTable(db *sql.DB, tableName string, file *os.File) {

	fmt.Fprintf(file, "DROP TABLE IF EXISTS `%s`;\n", tableName)
	fmt.Fprintf(file, "/*!40101 SET @saved_cs_client     = @@character_set_client */;\n")
	fmt.Fprintf(file, "/*!50503 SET character_set_client = utf8mb4 */;\n")
	dumpCreateTable(db, tableName, file)

	fmt.Fprintf(file, "/*!40101 SET character_set_client = @saved_cs_client */;\n\n")
	dumpTableData(db, tableName, file)
}

func dumpCreateTable(db *sql.DB, tableName string, file *os.File) {
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		var tblName, createStmt string
		if err := rows.Scan(&tblName, &createStmt); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(file, createStmt+";\n")
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func dumpTableData(db *sql.DB, tableName string, file *os.File) {
	query := fmt.Sprintf("SELECT * FROM `%s`", tableName)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(file, "LOCK TABLES `%s` WRITE;\n", tableName)
	fmt.Fprintf(file, "/*!40000 ALTER TABLE `%s` DISABLE KEYS */;\n", tableName)

	first := true
	for rows.Next() {
		if first {
			columnList := "`" + strings.Join(columns, "`, `") + "`"
			fmt.Fprintf(file, "INSERT INTO `%s` (%s) VALUES\n", tableName, columnList)
			first = false
		} else {
			fmt.Fprintf(file, ",\n")
		}

		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}

		var valueText []string
		for _, value := range values {
			if value == nil {
				valueText = append(valueText, "NULL")
			} else {
				valueText = append(valueText, fmt.Sprintf("'%s'", escapeSQLString(string(value))))
			}
		}
		fmt.Fprintf(file, "(%s)", strings.Join(valueText, ", "))
	}

	if !first {
		fmt.Fprintf(file, ";\n")
	}

	fmt.Fprintf(file, "/*!40000 ALTER TABLE `%s` ENABLE KEYS */;\n", tableName)
	fmt.Fprintf(file, "UNLOCK TABLES;\n")
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func dumpAllTables(db *sql.DB, file *os.File, database string) {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = ?", database)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		dumpTable(db, tableName, file)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func writeSessionVariables(file *os.File) {
	fmt.Fprintf(file, "/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;\n")
	fmt.Fprintf(file, "/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;\n")
	fmt.Fprintf(file, "/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;\n")
	fmt.Fprintf(file, "/*!50503 SET NAMES utf8mb4 */;\n")
	fmt.Fprintf(file, "/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;\n")
	fmt.Fprintf(file, "/*!40103 SET TIME_ZONE='+00:00' */;\n")
	fmt.Fprintf(file, "/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;\n")
	fmt.Fprintf(file, "/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;\n")
	fmt.Fprintf(file, "/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;\n")
	fmt.Fprintf(file, "/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;\n")
}

func writeSessionRestore(file *os.File) {
	fmt.Fprintf(file, "/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;\n")
	fmt.Fprintf(file, "/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;\n")
	fmt.Fprintf(file, "/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;\n")
	fmt.Fprintf(file, "/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;\n")
	fmt.Fprintf(file, "/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;\n")
	fmt.Fprintf(file, "/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;\n")
	fmt.Fprintf(file, "/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;\n")
}

func escapeSQLString(value string) string {
	return strings.ReplaceAll(value, "'", "\\'")
}

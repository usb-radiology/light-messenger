package lmdatabase

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" // mysql driver ..
	"github.com/usb-radiology/light-messenger/src/configuration"
)

// GetDB ..
func GetDB(initConfig *configuration.Configuration) (*sql.DB, error) {
	conn := initConfig.Database.Username + ":" + initConfig.Database.Password + "@tcp(" + initConfig.Database.Host + ":" + strconv.Itoa(initConfig.Database.Port) + ")/" + initConfig.Database.DBName
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal("Error opening database connection", err)
	}
	return db, nil
}

// ReadStatementsFromSQL ..
func ReadStatementsFromSQL(sqlFilePath string) (*[]string, error) {
	// sqlFilePath :=

	fileContents, errFileRead := ioutil.ReadFile(sqlFilePath)
	if errFileRead != nil {
		return nil, errFileRead
	}

	strStatements := strings.Split(string(fileContents), ";\n")
	return &strStatements, nil
}

// ExecStatements ..
func ExecStatements(db *sql.DB, sqlStatements []string) (*[]sql.Result, error) {

	results := make([]sql.Result, 0)

	for _, statement := range sqlStatements {

		trimedStatement := strings.Trim(statement, " \n")

		// skip empty statements
		if len(trimedStatement) > 0 {
			execResult, err := db.Exec(trimedStatement)
			if err != nil {
				// if err.Error() != "Error 1065: Query was empty" { // skip empty line errors (or alternatively skip empty line statements)
				return nil, err
				// }
			}
			results = append(results, execResult)
		}

	}

	return &results, nil
}

// ExecStatement ..
func ExecStatement(db *sql.DB, statement string) (sql.Result, error) {

	trimedStatement := strings.Trim(statement, " \n")

	// skip empty statements
	if len(trimedStatement) > 0 {
		execResult, err := db.Exec(trimedStatement)
		if err != nil {
			// if err.Error() != "Error 1065: Query was empty" { // skip empty line errors (or alternatively skip empty line statements)
			return nil, err
			// }
		}
		return execResult, nil
	}

	return nil, nil
}

package database

import (
	"database/sql"
	"fmt"
)

type ShareCodesRecord struct {
	ShareCode        string
	ShareDescription string
	PollStart        string
	PollEnd          string
	PollInterval     int
}

type SharesDB struct {
	*sql.DB
}

func (db *SharesDB) ConnectToDatabase(host string, port int, user string, password string, dbname string) (SharesDB, error) {

	var sharesDB SharesDB

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newdb, err := sql.Open("postgres", psqlInfo)

	sharesDB.DB = newdb

	if err != nil {
		return sharesDB, err
	}

	return sharesDB, nil
}

func (db *SharesDB) RetrieveShares() []ShareCodesRecord {

	sqlStatement := "SELECT * FROM read_share_codes()"

	rows, err := db.DB.Query(sqlStatement)

	if err != nil {
		fmt.Println(err)
	}

	var shareCodesTable []ShareCodesRecord

	for rows.Next() {

		var shareCode, shareDescription string
		var pollStart, pollEnd string
		var pollInterval int

		err := rows.Scan(&shareCode, &shareDescription, &pollStart, &pollEnd, &pollInterval)

		if err != nil {
			fmt.Println(err)
		}

		var shareCodesRecord ShareCodesRecord

		shareCodesRecord.ShareCode = shareCode
		shareCodesRecord.ShareDescription = shareDescription
		shareCodesRecord.PollStart = pollStart
		shareCodesRecord.PollEnd = pollEnd
		shareCodesRecord.PollInterval = pollInterval

		shareCodesTable = append(shareCodesTable, shareCodesRecord)
	}

	defer rows.Close()

	return shareCodesTable
}

func (db *SharesDB) AddShareCode(shareCode string, shareDescription string, pollStart string, pollEnd string, pollInterval int) (int64, error) {

	stmt, err := db.Prepare("SELECT * FROM create_share_code($1,$2,$3,$4,$5)")

	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(shareCode, shareDescription, pollStart, pollEnd, pollInterval)

	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil

}

package database

import (
	"database/sql"
	"fmt"
	"time"
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

func (db *SharesDB) RemoveShareCode(shareCode string) (int64, error) {

	var returnAmount int

	fmt.Printf("In remove and share code = %s\n", shareCode)

	stmt, err := db.Prepare("SELECT * FROM remove_share_code($1)")

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	result, err := stmt.Query(shareCode)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	for result.Next() {

		err := result.Scan(&returnAmount)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Rows pulled = %d.\n", returnAmount)

	}

	//rowsAffected, _ := result.RowsAffected()
	//return rowsAffected, nil
	return 1, nil

}

func (db *SharesDB) ListShareCodes() (int64, error) {

	//var returnAmount int

	fmt.Printf("In list share code = %s\n", "WWWW")

	stmt, err := db.Prepare("SELECT * FROM read_share_codes()")

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	//result, err := stmt.Query("WWWW")
	result, err := stmt.Query()

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	for result.Next() {

		var share_code, share_description, poll_start, poll_end string
		var poll_interval int

		err := result.Scan(&share_code, &share_description, &poll_start, &poll_end, &poll_interval)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Rows pulled = %s, %s. %s, %s, %d\n", share_code, share_description, poll_start, poll_end, poll_interval)

	}

	return 1, nil
}

func (db *SharesDB) LogSharePrice(code string, midPrice float64, bidPrice float64, askPrice float64) {

	//	statement := `INSERT INTO SharePrices
	//				(TimeDate,ShareCodeID,MidPrice,BidPrice,AskPrice)
	//				VALUES (?,?,?,?,?)`

	fmt.Println("*aaa*")

	stmt, err := db.Prepare("SELECT * FROM create_share_price_record($1,$2,$3,$4,$5)")
	//stmt, err := db.Prepare("SELECT * FROM create_share_price_record()")
	//stmt, err := db.Prepare("SELECT * FROM read_share_codes()")

	/*
		stmt, err := db.Prepare(`	INSERT INTO SharePrices
												(TimeDate,ShareCodeID,MidPrice,BidPrice,AskPrice)
									VALUES 		(?,?,?,?,?)`)
	*/
	fmt.Println("*aab*")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("*aac*")
	t := time.Now()

	fmt.Println("aaa")
	fmt.Printf("%s\t%s\t%f\t%f\t%f\n", code, t.Format("2006-01-02 15:04:05"), midPrice, bidPrice, askPrice)
	_, err = stmt.Exec(code, t.Format("2006-01-02 15:04:05"), midPrice, bidPrice, askPrice)

	fmt.Println("bbb")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\t\t%f\t%f\t%f\n", "TEP", midPrice, bidPrice, askPrice)

	//fmt.Println(result.RowsAffected())

}

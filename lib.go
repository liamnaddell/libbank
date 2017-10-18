package libbank

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"syscall"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "bankmngr"
	dbname = "bankmngr"
)

type Name struct {
	Firstname string
	Lastname  string
}

type BinReturn struct {
	Comment      string
	Emailaddress string
	Firstname    string
	Lastname     string
	Balance      int
}

/*Plan:
We need a few functions:
	one:
		CreateUser(emailaddress, firstname, lastname, dob, hash(password))
		insert into useraccounts values ('bigdaddytony@hotmail.com', 'Tony', 'FrostedFlakes(tm)', '2001-09-28', 'iambetter');
	two:
		CreateBin(userid, Comment, balance) binid ommitted
	three:
		AddSecurityQuestionToUser(userid, questionid, answer)
*/

func Connect() *sql.DB {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	password = strings.TrimSpace(password)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password='%s' dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()

	fmt.Println("\nSuccessfully connected!")
	return db
}

func GetGenericColumn(db *sql.DB, columnName string, tableName string) ([]string, error) {
	rows, err := db.Query("select " + columnName + " from " + tableName)
	if err != nil {
		return nil, err
	}
	var final []string
	var temp string
	for rows.Next() {
		err = rows.Scan(&temp)
		final = append(final, temp)
	}
	return final, nil
}

func AddSecurityQuestionToUser(db *sql.DB, userid int, questionid int, answer int) error {
	stmt, err := db.Prepare("insert into securityquestionstousers(questionid, answer, userid) values($1,$2,$3)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(questionid, answer, userid)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Rows changed \n%d\n", affect)
	return nil
}

func CreateBin(db *sql.DB, userid int, comment string, balance int) error {
	stmt, err := db.Prepare("insert into bins(userid,balance,comment) values($1,$2,$3)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	res, err := stmt.Exec(userid, balance, comment)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Rows changed \n%d\n", affect)
	return nil
}

func CreateUser(db *sql.DB, emailaddress string, firstname string, lastname string, dob string, password string) error {
	stmt, err := db.Prepare("insert into useraccounts(emailaddress,firstname,lastname,dob,password) values($1,$2,$3,$4,$5)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(emailaddress, firstname, lastname, dob, password)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Rows changed \n%d\n", affect)
	return nil
}

func UpdateBins(db *sql.DB, comment string, binid int) error {
	stmt, err := db.Prepare("update bins set comment=$1 where binid=$2")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(comment, binid)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Rows changed \n%d\n", affect)
	return nil
}

func JoinBins(db *sql.DB) ([]BinReturn, error) {
	//select bins.comment,useraccounts.emailaddress,useraccounts.firstname,useraccounts.lastname,bins.balance from bins inner join useraccounts on bins.userid=useraccounts.userid;
	rows, err := db.Query("select bins.comment,useraccounts.emailaddress,useraccounts.firstname,useraccounts.lastname,bins.balance from bins inner join useraccounts on bins.userid=useraccounts.userid")

	var comment, emailaddress, firstname, lastname string
	var balance int
	if err != nil {
		return []BinReturn{}, err
	}
	var finishes []BinReturn
	for rows.Next() {
		err = rows.Scan(&comment, &emailaddress, &firstname, &lastname, &balance)
		finishes = append(finishes, BinReturn{comment, emailaddress, firstname, lastname, balance})
	}
	return finishes, nil
}

func UpdateOrInsert(db *sql.DB, comment string, balance, binid, userid int) error {
	//insert into bins(comment,balance,binid,userid) values('My name jeff', 494, 5, 1) on conflict (binid) do update set comment='Hello world';
	stmt, err := db.Prepare("insert into bins(comment,balance,binid,userid) values($1, $2, $3, $4) on conflict (binid) do update set comment=$1, balance=$2, userid=$4")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(comment, balance, binid, userid)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Rows changed \n%d\n", affect)
	return nil
}
func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

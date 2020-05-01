package main

import (
	"database/sql"
	"fmt"
	"time"

	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	User     = "root"
	Password = "lrc2529130"
	DBName   = "ass3"
)

type Library struct {
	db *sqlx.DB
}

type book struct {
	title  string
	ISBN   string
	author string
}

func mustExecute(db *sqlx.DB, SQLs []string) {
	for _, s := range SQLs {
		_, err := db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}

func (lib *Library) CreateDB() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", User, Password))
	if err != nil {
		panic(err)
	}
	lib.db = db
	mustExecute(lib.db, []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS ass3"),
		fmt.Sprintf("CREATE DATABASE ass3"),
		fmt.Sprintf("USE ass3"),
	})
	fmt.Println("create DB:ass3 success")
}

func (lib *Library) ConnectDB() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
	// s := fmt.Sprintf("DROP DATABASE IF EXISTS ass3_%s "+
	// 	"CREATE DATABASE ass3_%s"+
	// 	"USE ass3_%s", User, User, User)
	// _, err = lib.db.Exec(s)
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println("connect DB:ass3 success")
}

// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {
	s := `CREATE TABLE book
	(title VARCHAR(32) NOT NULL,
	ISBN VARCHAR(32) NOT NULL,
	author VARCHAR(32) NOT NULL,
	condp bool NOT NULL,
	reason VARCHAR(32) NOT NULL,
	PRIMARY KEY(ISBN)
	);`
	_, err := lib.db.Exec(s)
	if err != nil {
		return err
	}
	s = `CREATE TABLE student
	(account VARCHAR(32) NOT NULL,
	ability bool NOT NULL,
	PRIMARY KEY(account)
	);`
	_, err = lib.db.Exec(s)
	if err != nil {
		return err
	}
	s = `CREATE TABLE lend
	(ret bool NOT NULL,
	getdate datetime NOT NULL,
	delay_cnt INT NOT NULL,
	book_id VARCHAR(32) NOT NULL,
	student_id VARCHAR(32) NOT NULL,
	FOREIGN KEY (book_id) REFERENCES book(ISBN),
	FOREIGN KEY (student_id) REFERENCES student(account)
	);`
	_, err = lib.db.Exec(s)
	if err != nil {
		panic(err)
	}
	fmt.Println("Createtables success")
	return nil
}

// AddBook add books into the library
func (lib *Library) AddBooks(admin string, books []book) error {
	fmt.Printf("Add books by %s\n", admin)
	if admin == User {
		if books != nil {
			fmt.Println("add books auto:")
			for _, abook := range books {
				var a = abook.author
				var i = abook.ISBN
				var t = abook.title
				lib.AddBook(admin, t, a, i)
			}
		} else {
			fmt.Println("add books on hand:")
			condition := 1
			var (
				title  string
				author string
				ISBN   string
			)
			for {
				if condition == 1 {
					fmt.Println("please input the title:")
					fmt.Scanln(&title)
					fmt.Println("please input the author:")
					fmt.Scanln(&author)
					fmt.Println("please input the ISBN:")
					fmt.Scanln(&ISBN)
					lib.AddBook(admin, title, author, ISBN)
					fmt.Println("input 1 to continue,input 0 to quit")
					fmt.Scanln(&condition)
				} else {
					break
				}
			}
		}
	} else {
		fmt.Println("you aren't the administrator")
	}
	fmt.Println("add books success")
	return nil
}

// AddBook add a book into the library

func (lib *Library) AddBook(admin, title, author, ISBN string) error {

	if admin == User {
		var iSBN sql.NullString
		err := lib.db.QueryRow("select ISBN from book where ISBN = ?", ISBN).Scan(&iSBN)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("Add the book %s %s %s\n", ISBN, title, author)
				s := fmt.Sprintf(
					`INSERT INTO book(title,ISBN,author,condp,reason)
				VALUES("%s","%s","%s",1,"");`, title, ISBN, author)
				_, err := lib.db.Exec(s)
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {
			fmt.Printf("the book:%s already exits\n,addbook invalid", ISBN)
		}
	} else {
		fmt.Println("you aren't the administrator")
	}
	return nil
}

//remove a book from library
func (lib *Library) RemoveBook(ISBN, reason string) error {
	var iSBN string
	err := lib.db.QueryRow("select ISBN from book where ISBN = ? AND condp=1", ISBN).Scan(&iSBN)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("the book:%s doesn't exit,remove invalid\n", ISBN)
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("remove the book %s for the reason: %s\n", ISBN, reason)
		s := fmt.Sprintf("UPDATE book SET condp=0,reason='%s' WHERE ISBN='%s' ", reason, ISBN)
		_, err = lib.db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

//没考虑已经在
//add student account
func (lib *Library) AddAccount(account string) error {
	fmt.Printf(" Add account %s\n", account)
	s1 := fmt.Sprintf(
		`INSERT INTO student(account,ability)
	 VALUES("%s",1);`, account)
	_, err := lib.db.Exec(s1)
	if err != nil {
		panic(err)
	}
	return nil
}

//query books by title/author/ISBN
func (lib *Library) QueryByTitle(s string) error {
	s1 := fmt.Sprintf(`
		SELECT *
		FROM book
		WHERE title='%s'
	`, s)
	rows, err := lib.db.Query(s1)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf(" Query By title:%s\n", s)
		fmt.Println("title", "ISBN", "author", "condp", "reason")
		var author, title, ISBN, reason string
		var con bool
		for rows.Next() {
			err = rows.Scan(&title, &ISBN, &author, &con, &reason)
			if err != nil {
				panic(err)
			}
			fmt.Println(title, ISBN, author, con, reason)
		}
	}
	defer rows.Close()
	return nil
}
func (lib *Library) QueryByAuthor(s string) error {
	s1 := fmt.Sprintf(`
		SELECT *
		FROM book
		WHERE author='%s'
	`, s)
	rows, err := lib.db.Query(s1)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf(" Query By author:%s\n", s)
		fmt.Println("title", "ISBN", "author", "condp", "reason")
		var author, title, ISBN, reason string
		var con bool
		for rows.Next() {
			err = rows.Scan(&title, &ISBN, &author, &con, &reason)
			if err != nil {
				panic(err)
			}
			fmt.Println(title, ISBN, author, con, reason)
		}
	}
	defer rows.Close()
	return nil
}
func (lib *Library) QueryByISBN(s string) error {
	s1 := fmt.Sprintf(`
		SELECT *
		FROM book
		WHERE ISBN='%s'
	`, s)
	rows, err := lib.db.Query(s1)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf(" Query By ISBN:%s\n", s)
		fmt.Println("title", "ISBN", "author", "condp", "reason")
		var author, title, ISBN, reason string
		var con bool
		for rows.Next() {
			err = rows.Scan(&title, &ISBN, &author, &con, &reason)
			if err != nil {
				panic(err)
			}
			fmt.Println(title, ISBN, author, con, reason)
		}
	}
	defer rows.Close()
	return nil
}

//borrow a book
func (lib *Library) BorrowBook(account, ISBN string) error {
	s := fmt.Sprintf(`
		SELECT ability
		FROM student
		WHERE account='%s' 
	`, account)
	var able bool
	err := lib.db.QueryRow(s).Scan(&able)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("the student:%s doesn't exit,borrow invalid\n", account)
		} else {
			panic(err)
		}
	} else {
		var exist bool
		err := lib.db.QueryRow("SELECT condp FROM book WHERE ISBN=? AND condp=1", ISBN).Scan(&exist)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("the book:%s doesn't in the library,borrow invalid\n", ISBN)
			} else {
				panic(err)
			}
		} else {
			if able {
				var id, sid string
				err := lib.db.QueryRow("SELECT book_id,student_id FROM lend WHERE book_id=? AND ret=0", ISBN).Scan(&id, &sid)
				if err != nil {
					if err == sql.ErrNoRows {
						fmt.Printf("the student:%s borrow the book: %s, please return in 30 days\n", account, ISBN)
						now := time.Now()
						//day, _ := time.ParseDuration("24h")
						//day, _ := time.ParseDuration("24h")
						getdate := now.Format("2006-01-02 15:04:05")
						//deadline := now.Add(day * 30).Format("2006-01-02 15:04:05")
						//getdate := time.parse("2006-01-02 15:04:05",date) time.parse time.Time
						//now.Format("2006-01-02 15:04:05")  string
						s := fmt.Sprintf(
							`INSERT INTO lend(ret,getdate,delay_cnt,book_id,student_id)
						VALUES(0,'%s',0,'%s','%s');`, getdate, ISBN, account)
						_, err := lib.db.Exec(s)
						if err != nil {
							panic(err)
						}
						//suspend student's account
						s = fmt.Sprintf(`
							SELECT COUNT(*)
							FROM lend
							WHERE student_id='%s' AND ret=0
						`, account)
						rows, err := lib.db.Query(s)
						rows.Next()
						var cnt int
						rows.Scan(&cnt)
						if cnt == 3 {
							s = fmt.Sprintf(`
							UPDATE student
							SET ability=0
							WHERE account='%s' 
							`, account)
							_, err = lib.db.Exec(s)
							if err != nil {
								panic(err)
							}
						}
					} else {
						panic(err)
					}
				} else {
					fmt.Printf("the book:%s has been borrowed by student:%s, borrow invalid\n", id, sid)
				}

			} else {
				fmt.Printf("the student:%s shouldn't borrow more than three books, borrow invalid\n", account)
			}
		}

	}
	return nil
}

//query the borrow history
func (lib *Library) QueryHistory(account string) error {
	s := fmt.Sprintf(`
		SELECT *
		FROM lend
		WHERE student_id="%s"
	`, account)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	fmt.Printf("the borrow history of student %s:\n", account)
	fmt.Println("ret", "getdate", "delay_cnt", "book_id", "student_id")
	var getdate, book_id, student_id string
	var delay_cnt int
	var ret bool
	for rows.Next() {
		err = rows.Scan(&ret, &getdate, &delay_cnt, &book_id, &student_id)
		if err != nil {
			panic(err)
		}
		fmt.Println(ret, getdate, delay_cnt, book_id, student_id)
	}
	return nil
}

//query the books a student has borrowed and not returned yet
func (lib *Library) QueryBookCon(account string) error {
	s := fmt.Sprintf(`
		SELECT book_id
		FROM lend
		WHERE student_id='%s' AND ret=0
	`, account)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	fmt.Println("ISBN")
	var ISBN string
	fmt.Printf("the books student %s has borrowed and not returned yet: \n", account)
	for rows.Next() {
		err = rows.Scan(&ISBN)
		if err != nil {
			panic(err)
		}
		fmt.Println(ISBN)
	}
	return nil
}

//check the deadline of returning a borrowed book
func (lib *Library) CheckDeadline(ISBN string) error {
	s := fmt.Sprintf(`
		SELECT getdate,delay_cnt
		FROM lend
		WHERE book_id='%s' AND ret=0
	`, ISBN)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	} else {
		var data string
		var cnt int
		rows.Next()
		err = rows.Scan(&data, &cnt)
		getdate, _ := time.Parse("2006-01-02 15:04:05", data)
		day, _ := time.ParseDuration("24h")
		deadline := getdate.Add(day * 30).Format("2006-01-02 15:04:05")
		if cnt == 0 {
		} else if cnt == 1 {
			deadline = getdate.Add(day * 37).Format("2006-01-02 15:04:05")
		} else if cnt == 2 {
			deadline = getdate.Add(day * 44).Format("2006-01-02 15:04:05")
		} else {
			deadline = getdate.Add(day * 51).Format("2006-01-02 15:04:05")
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("the deadline of the borrowed book %s is %s\n", ISBN, deadline)
	}
	return nil
}

//extend the deadline of returning a book
func (lib *Library) ExtendDeadline(ISBN string) error {
	fmt.Printf("extend 7 days to return the book %s\n", ISBN)
	s := fmt.Sprintf(`
	SELECT delay_cnt
	FROM lend
	WHERE book_id='%s' AND ret=0
	`, ISBN)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	rows.Next()
	var cnt int
	err = rows.Scan(&cnt)
	if err != nil {
		panic(err)
	}
	if cnt == 3 {
		fmt.Printf("the book %s couldn't been delayed 3 times\n", ISBN)
	} else {
		// delay := (day * 7).Format("2006-01-02 15:04:05")
		s := fmt.Sprintf(`
			UPDATE lend
			SET delay_cnt=delay_cnt+1
			WHERE book_id=%s AND ret=0 AND delay_cnt<3
		`, ISBN)
		_, err := lib.db.Exec(s)
		if err != nil {
			panic(err)
		}
		fmt.Println("extend success")
	}
	return nil
}

//check if a student has any overdue books that needs to be returned
func (lib *Library) CheckDue(account string) error {
	now := time.Now()
	ndate := now.Format("2006-01-02 15:04:05")
	nowdate, _ := time.Parse("2006-01-02 15:04:05", ndate)
	s := fmt.Sprintf(`
		SELECT book_id,getdate,delay_cnt
		FROM lend
		WHERE student_id='%s' AND ret=0
	`, account)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	cnt := 0
	var ISBN, date string
	var delay_cnt int
	fmt.Printf("the overdue books the student %s has:\n", account)
	for rows.Next() {
		cnt = cnt + 1
		rows.Scan(&ISBN, &date, &delay_cnt)
		getdate, _ := time.Parse("2006-01-02 15:04:05", date)
		day, _ := time.ParseDuration("24h")
		deadline := getdate.Add(day * 30).Format("2006-01-02 15:04:05")
		if delay_cnt == 0 {
		} else if delay_cnt == 1 {
			deadline = getdate.Add(day * 37).Format("2006-01-02 15:04:05")
		} else if delay_cnt == 2 {
			deadline = getdate.Add(day * 44).Format("2006-01-02 15:04:05")
		} else {
			deadline = getdate.Add(day * 51).Format("2006-01-02 15:04:05")
		}
		Dead, _ := time.Parse("2006-01-02 15:04:05", deadline)
		if Dead.Before(nowdate) {
			fmt.Println(ISBN)
		}
	}
	if cnt == 0 {
		fmt.Printf("the student %s has no overdue book", account)
	}
	return nil
}

//return a book to the library by a student account
func (lib *Library) ReturnBook(ISBN, account string) error {

	s := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM lend
		WHERE book_id='%s' AND ret=0 AND student_id='%s'
	`, ISBN, account)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	var cnt int
	rows.Next()
	err = rows.Scan(&cnt)
	if err != nil {
		panic(err)
	}
	if cnt == 0 {
		fmt.Printf("the student %s doesn't borrow the book %s,return invalid\n", account, ISBN)
	} else {
		fmt.Printf("the student %s return the book %s\n", account, ISBN)
		s := fmt.Sprintf(`
			UPDATE lend
			SET ret=1
			WHERE book_id='%s' AND ret=0 AND student_id='%s'
		`, ISBN, account)
		_, err := lib.db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// etc...

// func main() {
// 	fmt.Println("Welcome to the Library Management System!")
// 	var lib Library
// 	lib.CreateDB()
// 	lib.ConnectDB()
// 	lib.CreateTables()
// 	//	lib.AddBook("root", "T1", "A1", "1")
// 	// books := []book{
// 	// 	book{
// 	// 		"B1",
// 	// 		"1",
// 	// 		"A1",
// 	// 	},
// 	// }
// 	// lib.AddBooks("root", books)
// 	// var book1 []book
// 	// lib.AddBooks("root", book1)
// }

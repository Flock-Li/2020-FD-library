package main

import (
	"testing"
)

func TestRemovebook(t *testing.T) {
	err := lib.RemoveBook("50", "the book is lost")
	if err != nil {
		t.Errorf("can't remove book")
	}
	err = lib.RemoveBook("2", "the book is lost")
	if err != nil {
		t.Errorf("can't remove book")
	}
	t.Logf("removebook success")
}

func TestAddAccount(t *testing.T) {
	lib.AddAccount("s1")
	err := lib.AddAccount("s2")
	if err != nil {
		t.Errorf("can't Add Account")
	}
	t.Logf("AddAccount success")
}

func TestQueryByTitle(t *testing.T) {
	err := lib.QueryByTitle("B2")
	if err != nil {
		t.Errorf("can't Query By Title")
	}
	t.Logf("QueryByTitle success")
}

func TestQueryByAuthor(t *testing.T) {
	err := lib.QueryByAuthor("A1")
	if err != nil {
		t.Errorf("can't Query By Author")
	}
	t.Logf("QueryByAuthor success")
}

func TestQueryByISBN(t *testing.T) {
	err := lib.QueryByISBN("2")
	if err != nil {
		t.Errorf("can't Query By ISBN:2")
	}
	t.Logf("QueryByISBN success")
}

func TestBorrowBook(t *testing.T) {
	lib.BorrowBook("s3", "4")
	lib.BorrowBook("s1", "2")
	lib.BorrowBook("s1", "4")
	lib.BorrowBook("s1", "5")
	lib.BorrowBook("s1", "6")
	lib.BorrowBook("s2", "4")
	err := lib.BorrowBook("s1", "2")
	if err != nil {
		t.Errorf("can't Borrow Book")
	}
	t.Logf(" BorrowBook success")
}

func TestQueryHistory(t *testing.T) {
	err := lib.QueryHistory("s1")
	if err != nil {
		t.Errorf("can't Query History")
	}
	t.Logf("QueryHistory success")
}

func TestQueryBookCon(t *testing.T) {
	err := lib.QueryBookCon("s1")
	if err != nil {
		t.Errorf("can't Query BookCon")
	}
	t.Logf("QueryBookCon success")
}

func TestCheckDeadline(t *testing.T) {
	err := lib.CheckDeadline("4")
	if err != nil {
		t.Errorf("can't Check Deadline")
	}
	t.Logf("CheckDeadline success")
}

func TestExtendDeadline(t *testing.T) {
	lib.ExtendDeadline("4")
	lib.ExtendDeadline("4")
	lib.ExtendDeadline("4")
	err := lib.ExtendDeadline("4")
	if err != nil {
		t.Errorf("can't Extend Deadline")
	}
	t.Logf("ExtendDeadline success")
}

func TestCheckDue(t *testing.T) {
	err := lib.CheckDue("s1")
	if err != nil {
		t.Errorf("can't Check Due")
	}
	t.Logf("CheckDue success")
}

func TestReturnBook(t *testing.T) {
	err := lib.ReturnBook("4", "s1")
	err = lib.ReturnBook("2", "s1")
	if err != nil {
		t.Errorf("can't Return Book")
	}
	t.Logf("ReturnBook success")
}

//TestAddbooks and TestCreateTables

var lib = Library{}

func TestMain(m *testing.M) {
	books := []book{
		book{"B1", "1", "A1"},
		book{"B2", "2", "A2"},
		book{"B3", "3", "A3"},
		book{"B4", "4", "A4"},
		book{"B5", "5", "A5"},
		book{"B6", "6", "A6"},
		book{"B7", "7", "A7"},
		book{"B8", "8", "A8"},
		book{"B9", "9", "A9"},
		book{"B10", "10", "A10"},
	}
	lib.CreateDB()
	lib.ConnectDB()
	lib.CreateTables()
	lib.AddBooks("root", books)
	lib.AddBook("root", "B10", "A10", "10")
	m.Run()
}

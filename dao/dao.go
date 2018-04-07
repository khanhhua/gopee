package dao

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const CLEARDB_DATABASE_URL = "b804b4e297b1f3:3575f604@tcp(us-cdbr-iron-east-05.cleardb.net:3306)/heroku_5b0142bce7e674e"

type Client struct {
	ID                 int64
	ClientKey          string
	ClientDomain       string
	DropboxAccountID   string
	DropboxAccessToken string
}

func CreateClient(client Client) (ret Client, err error) {
	db, dberr := sql.Open("mysql", CLEARDB_DATABASE_URL)
	if dberr != nil {
		err = dberr
		return
	}
	defer db.Close()
	// Insert into table clients and return the latest ID
	var result sql.Result
	result, dberr = db.Exec("INSERT INTO clients (client_key, domain, dropbox_account_id, dropbox_access_token) VALUES (?, ?, ?, ?)",
		client.ClientKey,
		client.ClientDomain,
		client.DropboxAccountID,
		client.DropboxAccessToken)
	if dberr != nil {
		err = dberr
		return
	}
	client.ID, dberr = result.LastInsertId()
	if dberr != nil {
		err = dberr
		return
	}

	ret = client
	return
}

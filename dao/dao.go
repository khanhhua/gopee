package dao

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Client struct {
	ID                 int64
	ClientKey          string
	ClientDomain       string
	DropboxAccountID   string
	DropboxAccessToken string
}

type FuncSpec struct {
	ID                 int64
	DropboxAccountID   string
	DropboxAccessToken string
	FnName             string
	XlsxFile           string
	InputMappings      map[string]string
	OutputMappings     map[string]string
}

// CreateClient Persists a new client
func CreateClient(client Client) (ret Client, err error) {
	CLEARDB_DATABASE_URL := os.Getenv("CLEARDB_DATABASE_URL")

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

// GetFuncByName Gets function by name
func GetFuncByName(clientKey string, fnName string) (ret FuncSpec, err error) {
	CLEARDB_DATABASE_URL := os.Getenv("CLEARDB_DATABASE_URL")

	db, dberr := sql.Open("mysql", CLEARDB_DATABASE_URL)
	if dberr != nil {
		err = dberr
		return
	}
	defer db.Close()

	var row *sql.Row
	row = db.QueryRow(`SELECT xlsx_file, dropbox_account_id, dropbox_access_token,
		input_mappings, output_mappings
			FROM funs
			JOIN clients ON funs.client_id = clients.id
			WHERE clients.client_key = ? AND funs.fn_name = ?`,
		clientKey, fnName)

	var inputMappingsRaw, outputMappingsRaw string
	if err = row.Scan(&ret.XlsxFile, &ret.DropboxAccountID, &ret.DropboxAccessToken,
		&inputMappingsRaw, &outputMappingsRaw); err != nil {
		fmt.Printf("Could not retrieve. Reason: %v", err)
		return
	}
	ret.FnName = fnName

	ret.InputMappings = make(map[string]string)
	for _, pair := range strings.Split(inputMappingsRaw, ";") {
		splat := strings.Split(pair, "=")
		ret.InputMappings[splat[0]] = splat[1]
	}

	ret.OutputMappings = make(map[string]string)
	for _, pair := range strings.Split(outputMappingsRaw, ";") {
		splat := strings.Split(pair, "=")
		ret.OutputMappings[splat[0]] = splat[1]
	}

	return
}

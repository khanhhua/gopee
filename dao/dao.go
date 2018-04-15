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
	row := db.QueryRow(`SELECT id FROM clients where client_key = ?`, client.ClientKey)
	if dberr = row.Scan(&client.ID); dberr == nil { // Do update
		result, dberr = db.Exec("UPDATE clients SET domain = ?, dropbox_account_id = ?, dropbox_access_token = ? WHERE id = ?",
			client.ClientDomain,
			client.DropboxAccountID,
			client.DropboxAccessToken,
			client.ID)
		if dberr != nil {
			err = dberr
			return
		}
	} else { // Do insert
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

// CreateFunc Register a new function named fnName
func CreateFunc(clientKey string, fnName string, xlsxFile string,
	inputMappings map[string]string, outputMappings map[string]string) (id int64, err error) {
	CLEARDB_DATABASE_URL := os.Getenv("CLEARDB_DATABASE_URL")

	db, dberr := sql.Open("mysql", CLEARDB_DATABASE_URL)
	if dberr != nil {
		err = dberr
		return
	}
	defer db.Close()

	im := make([]string, 0)
	for key, value := range inputMappings {
		im = append(im, fmt.Sprintf("%s=%s", key, value))
	}
	serializedIM := strings.Join(im, ";")

	om := make([]string, 0)
	for key, value := range outputMappings {
		om = append(om, fmt.Sprintf("%s=%s", key, value))
	}
	serializedOM := strings.Join(om, ";")

	var result sql.Result

	result, dberr = db.Exec(`INSERT INTO
		funs   (client_id, fn_name, xlsx_file, input_mappings, output_mappings)
	 	VALUES ((SELECT id FROM clients WHERE clients.client_key = ?), ?, ?, ?, ?)`,
		clientKey, fnName, xlsxFile, serializedIM, serializedOM)
	if dberr != nil {
		err = dberr
		return
	}

	id, dberr = result.LastInsertId()
	if dberr != nil {
		err = dberr
		return
	}
	return
}

// FindFuns Finds all the funs this client has
func FindFuns(clientKey string) (result []FuncSpec, err error) {
	CLEARDB_DATABASE_URL := os.Getenv("CLEARDB_DATABASE_URL")

	db, dberr := sql.Open("mysql", CLEARDB_DATABASE_URL)
	if dberr != nil {
		err = dberr
		return
	}
	defer db.Close()

	var rows *sql.Rows
	rows, dberr = db.Query(`
		SELECT id, fn_name, xlsx_file, input_mappings, output_mappings
		FROM funs
		WHERE client_id = (SELECT id FROM clients WHERE client_key = ?)`, clientKey)
	if dberr != nil {
		err = dberr
		return
	}

	for rows.Next() {
		row := FuncSpec{}
		var rawInputMappings, rawOutputMappings string
		dberr = rows.Scan(&row.ID, &row.FnName, &row.XlsxFile, &rawInputMappings, &rawOutputMappings)
		if dberr != nil {
			err = dberr
			return
		}

		row.InputMappings = deserializeRawMapping(rawInputMappings)
		row.OutputMappings = deserializeRawMapping(rawOutputMappings)

		result = append(result, row)
	}

	return
}

func GetFunc(clientKey string, id int64) (result FuncSpec, err error) {
	CLEARDB_DATABASE_URL := os.Getenv("CLEARDB_DATABASE_URL")

	db, dberr := sql.Open("mysql", CLEARDB_DATABASE_URL)
	if dberr != nil {
		err = dberr
		return
	}
	defer db.Close()

	var row *sql.Row
	row = db.QueryRow(`
		SELECT id, fn_name, xlsx_file, input_mappings, output_mappings
		FROM funs
		WHERE client_id = (SELECT id FROM clients WHERE client_key = ?)
			AND id = ?`, clientKey, id)
	var rawInputMappings, rawOutputMappings string
	if dberr = row.Scan(&result.ID, &result.FnName, &result.XlsxFile,
		&rawInputMappings, &rawOutputMappings); dberr != nil {
		err = dberr
		return
	}

	result.InputMappings = deserializeRawMapping(rawInputMappings)
	result.OutputMappings = deserializeRawMapping(rawOutputMappings)

	return
}

func UpdateFunc(clientKey string, fun FuncSpec) error {
	CLEARDB_DATABASE_URL := os.Getenv("CLEARDB_DATABASE_URL")

	db, dberr := sql.Open("mysql", CLEARDB_DATABASE_URL)
	if dberr != nil {
		return dberr
	}
	defer db.Close()

	if _, dberror := db.Exec(`
		UPDATE funs
		SET fn_name=?, xlsx_file=?, input_mappings=?, output_mappings=?
		WHERE client_id=(SELECT id FROM clients WHERE client_key = ?) 
		 AND id = ?
	`, fun.FnName, fun.XlsxFile, serializeMapping(fun.InputMappings),
		serializeMapping(fun.OutputMappings), clientKey, fun.ID); dberror != nil {
		return dberr
	}

	return nil
}

func serializeMapping(mappings map[string]string) string {
	im := make([]string, 0)
	for key, value := range mappings {
		im = append(im, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(im, ";")
}

func deserializeRawMapping(raw string) map[string]string {
	ret := make(map[string]string)
	if len(raw) == 0 {
		return ret
	}
	for _, pair := range strings.Split(raw, ";") {
		splat := strings.Split(pair, "=")
		ret[splat[0]] = splat[1]
	}

	return ret
}

package dao

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DAO struct {
	io.Closer
	db *sql.DB
}

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

func New(url string) (instance *DAO, err error) {
	db, dberr := sql.Open("mysql", url)
	if dberr != nil {
		err = dberr
		return
	}

	instance = &DAO{
		db: db,
	}
	return
}

func (instance *DAO) Close() error {
	return instance.db.Close()
}

// CreateClient Persists a new client
func (instance *DAO) CreateClient(client Client) (ret Client, err error) {
	db := instance.db
	var dberr error
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
func (instance *DAO) GetFuncByName(clientKey string, fnName string) (ret FuncSpec, err error) {
	db := instance.db

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

	ret.InputMappings = deserializeRawMapping(inputMappingsRaw)
	ret.OutputMappings = deserializeRawMapping(outputMappingsRaw)

	return
}

// CreateFunc Register a new function named fnName
func (instance *DAO) CreateFunc(clientKey string, fnName string, xlsxFile string,
	inputMappings map[string]string, outputMappings map[string]string) (id int64, err error) {
	db := instance.db
	var dberr error

	serializedIM := serializeMapping(inputMappings)
	serializedOM := serializeMapping(outputMappings)

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
func (instance *DAO) FindFuns(clientKey string) (result []FuncSpec, err error) {
	db := instance.db
	var dberr error

	var rows *sql.Rows
	rows, dberr = db.Query(`
		SELECT id, fn_name, xlsx_file, input_mappings, output_mappings
		FROM funs
		WHERE client_id = (SELECT id FROM clients WHERE client_key = ?)`, clientKey)
	if dberr != nil {
		err = dberr
		return
	}
	defer rows.Close()

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

func (instance *DAO) GetFunc(clientKey string, id int64) (result FuncSpec, err error) {
	db := instance.db
	var dberr error

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

func (instance *DAO) UpdateFunc(clientKey string, fun FuncSpec) error {
	db := instance.db
	var dberr error

	if _, dberr = db.Exec(`
		UPDATE funs
		SET fn_name=?, xlsx_file=?, input_mappings=?, output_mappings=?
		WHERE client_id=(SELECT id FROM clients WHERE client_key = ?)
		 AND id = ?
	`, fun.FnName, fun.XlsxFile, serializeMapping(fun.InputMappings),
		serializeMapping(fun.OutputMappings), clientKey, fun.ID); dberr != nil {
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

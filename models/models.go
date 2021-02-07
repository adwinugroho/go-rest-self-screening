package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/adwinugroho/go-rest-self-screening/config"
	"github.com/arangodb/go-driver"
)

type (
	// Health assessment
	HealthAssessment struct {
		ID    string               `json:"_key"`
		Data  DataHealthAssessment `json:"data"`
		Audit *Audit               `json:"audit,omitempty"`
	}

	DataHealthAssessment struct {
		Status string      `json:"status"`
		Detail interface{} `json:"detail"`
		Date   string      `json:"date,omitempty"`
	}

	Audit struct {
		Key          string `json:"_key,omitempty"`
		CurrNo       int    `json:"curr_no"`
		Inputter     string `json:"inputter"`
		InputterName string `json:"inputterName,omitempty"`
		Datetime     string `json:"datetime"`
		Inputtime    string `json:"inputtime,omitempty"`
		LogReason    string `json:"log_reason,omitempty"`
	}

	BindVars struct {
		KeyID  string `json:"_key,omitempty"`
		Offset int    `json:"offset"`
		Limit  int    `json:"limit"`
		Page   int    `json:"page"`
		Filter Filter `json:"filter"`
		Search Search `json:"search,omitempty"`
	}

	Filter struct {
		DateAge []string
		Age     []int  `json:"age,omitempty"`
		Gender  string `json:"gender,omitempty"`
		Covid   string `json:"covid,omitempty"`
		PCR     string `json:"pcr,omitempty"`
		Rapid   string `json:"rapid,omitempty"`
	}

	Search struct {
		Text string `json:"text"`
	}

	// struct db for get connection instance
	DB struct {
		DBLive driver.Database
		DBLog  driver.Database
	}
)

// create a function we call in main to get connection
func NewConnection(conn *config.Connection) *DB {
	return &DB{
		DBLive: conn.DBLive,
		DBLog:  conn.DBLog,
	}
}

// func add data to DB (insert)
func (db *DB) AddData(model HealthAssessment) (*HealthAssessment, error) {
	// variable parent context
	ctx := context.Background()
	// declare model we used
	data := HealthAssessment{}
	// open connection to DB collection
	col, err := db.DBLive.Collection(ctx, "health_assessment")
	if err != nil {
		log.Printf("[models.go:AddData, db.Collection] Error open connection to collection, cause: %+v\n", err)
		return nil, err
	}
	//WithReturnNew is used to configure a context to make create, update & replace document
	driverCtx := driver.WithReturnNew(ctx, &data)
	// call func CreateDocument and passing driverCtx and parameter
	meta, err := col.CreateDocument(driverCtx, data)
	if err != nil {
		log.Printf("[models.go:AddData, col.CreateDoc] Error while creating document, cause: %+v\n", err)
		return nil, err
	}
	// passing _key in db in our model (data)
	data.ID = meta.Key
	// print _key
	fmt.Printf("Created document with key '%s', revision '%s'\n", meta.Key, meta.Rev)
	// return deference model (data) and error == nil
	return &data, nil
}

func (db *DB) DeleteByKey(id string) (*[]HealthAssessment, error) {
	ctx := context.Background()
	var result []HealthAssessment
	var bindVars = map[string]interface{}{
		"keyid": id,
	}
	var query = "FOR x IN health_assessment FILTER x._key == @keyid REMOVE x IN health_assessment LET removed = OLD RETURN removed"
	cursor, err := db.DBLive.Query(ctx, query, bindVars)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {
		var data HealthAssessment
		_, err := cursor.ReadDocument(ctx, &data)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, errors.New("error read document")
		}
		result = append(result, data)
	}
	return &result, nil
}

func (db *DB) GetDataByKey(id string) (*HealthAssessment, error) {
	ctx := context.Background()
	data := HealthAssessment{}
	col, err := db.DBLive.Collection(ctx, "health_assessment")
	if err != nil {
		log.Printf("[models.go:GetDataByKey, db.Collection] Error open connection to collection, cause: %+v\n", err)
		return nil, err
	}
	_, err = col.ReadDocument(ctx, id, &data)
	if err != nil {
		log.Printf("[models.go:GetDataByKey, col.ReadDoc] Error reading document, cause: %+v\n", err)
		return nil, err
	}
	return &data, nil
}

func (db *DB) GetListAllData(vars BindVars) ([]HealthAssessment, int64, error) {
	var query, queryCount, search string
	var filter bool
	var ctx = context.Background()
	var datas []HealthAssessment
	var bindVars = map[string]interface{}{}
	//Search by
	if vars.Search.Text != "" {
		vars.Search.Text = "%" + vars.Search.Text + "%"
		search = fmt.Sprintf(`LIKE(x.data.longUrl, "%s", true)`, vars.Search.Text)
		filter = true
	}

	if filter {
		query = fmt.Sprintf(`FOR x IN health_assessment %s %s LIMIT %v, %v RETURN x`, "FILTER", search, vars.Offset, vars.Limit)
		queryCount = fmt.Sprintf(`FOR x IN health_assessment %s %s LIMIT %v, %v RETURN x`, "FILTER", search, vars.Offset, vars.Limit)
		//fmt.Printf("query: %+v\n", query)
	} else {
		query = fmt.Sprintf(`FOR x IN health_assessment LIMIT %v, %v RETURN x`, vars.Offset, vars.Limit)
		queryCount = `FOR x IN health_assessment RETURN x`
	}

	cursor, err := db.DBLive.Query(ctx, query, bindVars)
	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close()
	for {
		data := new(HealthAssessment)
		meta, err := cursor.ReadDocument(ctx, &data)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Printf("Error reading document: %+v \n", err)
		}
		data.ID = meta.Key
		datas = append(datas, *data)
	}
	if len(datas) == 0 {
		return nil, 0, nil
	}

	ctxCount := driver.WithQueryCount(context.Background())
	cursorCount, err := db.DBLive.Query(ctxCount, queryCount, bindVars)
	if err != nil {
		fmt.Printf("[model:ListAllData]: Error execute query count [%s], cause: %+v \n", queryCount, err)
		return nil, 0, err
	}
	defer cursorCount.Close()
	fmt.Printf("query count func:%+v\n", cursorCount.Count())
	return datas, cursorCount.Count(), nil
}

//Add log to DB
func (db *DB) SaveLog(model *HealthAssessment) (*string, error) {
	var ctx = context.Background()
	model.ID = fmt.Sprintf("%s-%s", model.ID, strconv.Itoa(model.Audit.CurrNo))
	col, err := db.DBLog.Collection(ctx, "health_assessment_log")
	if err != nil {
		return nil, err
	}

	meta, err := col.CreateDocument(ctx, &model)
	if err != nil {
		return nil, err
	}

	return &meta.Key, nil
}

//Update data
func (db *DB) UpdateData(model *HealthAssessment) (*HealthAssessment, error) {
	ctx := context.Background()
	data := HealthAssessment{}
	col, err := db.DBLive.Collection(ctx, "health_assessment")
	if err != nil {
		return nil, err
	}

	driverCtx := driver.WithReturnNew(ctx, &data)
	meta, err := col.ReplaceDocument(driverCtx, model.ID, model)
	fmt.Printf("Doc Revision: %s \n", meta.Rev)
	if err != nil {
		return nil, err
	}

	data.ID = meta.Key
	return &data, nil
}

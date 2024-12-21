package crudder

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"net/url"
	"protchain/internal/value"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/uptrace/bun"
)

type PaginationMeta struct {
	Total     int `json:"total"`
	Skipped   int `json:"skipped"`
	PerPage   int `json:"per_page"`
	Page      int `json:"page"`
	PageCount int `json:"page_count"`
}

type DALFilters struct {
	Like       map[string]string
	Exact      map[string]any
	RawWhere   []string
	RawWhereOr []string
	// Raw query overrides every other filter and executes just the query
	RawQuery string
	RawCount string
	RegEx    map[string]string
	Or       map[string]interface{}
	InInt    map[string]interface{}
	InOr     map[string]any
	IsNot    map[string]bool
	Relation map[string][]func(query *bun.SelectQuery) *bun.SelectQuery
	Sorter   DALSorters
	Columns  []string
}

type Join struct {
	Table     string
	Condition string
	// the expression for each column to be mapped to the struct
	// e.g. "table_name.column_name as table_name__column_name"
	Expression string
}

type DALSorters struct {
	Asc  []string
	Desc []string
}

func (ds *DALSorters) ToQuery() []string {
	sorter := make([]string, 0)
	for _, col := range ds.Asc {
		sorter = append(sorter, fmt.Sprintf("%s ASC", col))
	}
	for _, col := range ds.Desc {
		sorter = append(sorter, fmt.Sprintf("%s DESC", col))
	}

	return sorter
}

// Setter defines information to be updated and how they should be
// updated.
type Setter struct {
	// Default identifies columns that are updated in this manner key = value
	Default map[string]any

	// Inc identifies columns that are updated in this manner key = key + value
	Inc map[string]float64

	// Dec identifies columns that are updated in this manner key = key - value
	Dec map[string]float64
}

func (s *Setter) HasUpdate() bool {
	for k := range s.Dec {
		if k != "" {
			return true
		}
	}

	for k := range s.Inc {
		if k != "" {
			return true
		}
	}

	for k := range s.Default {
		if k != "" {
			return true
		}
	}

	return false
}

// Crudder ...
type Crudder struct {
	Filter    *DALFilters
	Sorter    *DALSorters
	Paginator *PaginationMeta
	Setter    *Setter
	DataModel any
	DB        bun.IDB
}

// GenerateSetter generates a map[string]any for database update. This function accepts a database
// model, a list of struct fields that can be updated and returns fields with not nil value
func GenerateSetter(object any, updateFields []string) map[string]any {
	setter := make(map[string]any)

	objectValue := reflect.ValueOf(object).Elem()
	objectType := reflect.TypeOf(object).Elem()

	for _, field := range updateFields {

		sField, found := objectType.FieldByName(field)
		if !found {
			continue
		}

		fieldValue := objectValue.FieldByName(field)
		fieldData := fieldValue.Interface()

		if fieldValue.Kind() != reflect.Bool && (!fieldValue.IsValid() || fieldValue.IsZero()) {
			continue
		}
		if timeV, isTime := fieldValue.Interface().(time.Time); isTime {
			if timeV.IsZero() {
				continue
			}
		}

		setterKey, tagSet := sField.Tag.Lookup("bun")
		if !tagSet {
			continue
		}
		setterKey = strings.Split(setterKey, ",")[0]

		setter[setterKey] = fieldData

	}

	return setter
}

// DefaultCrudder ...
func DefaultCrudder(dataModel any, db bun.IDB, query ...url.Values) Crudder {
	return Crudder{
		Filter:    defaultFilter(),
		Sorter:    defaultSorter(),
		Paginator: defaultPaginator(query...),
		Setter:    defaultSetter(),
		DataModel: dataModel,
		DB:        db,
	}
}

// addSelectQuery ...
func (c Crudder) addSelectQuery(query *bun.SelectQuery) {
	addQueryFilters(query, *c.Filter)
	return
}

// addUpdateQuery ...
func (c Crudder) addUpdateQuery(query *bun.UpdateQuery) {
	for column, value := range c.Filter.Exact {
		query.Where("? = ?", bun.Ident(column), value)
	}

	for _, rQuery := range c.Filter.RawWhere {
		query.Where(rQuery)
	}
}

func (c Crudder) addUpdateSetters(query *bun.UpdateQuery) {
	for column, value := range c.Setter.Default {
		query.Set("? = ?", bun.Ident(column), value)
	}

	for column, value := range c.Setter.Dec {
		query.Set("? = ? - ?", bun.Ident(column), bun.Ident(column), value)
	}

	for column, value := range c.Setter.Inc {
		query.Set("? = ? + ?", bun.Ident(column), bun.Ident(column), value)
	}
}

func (c Crudder) addDeleteQuery(query *bun.DeleteQuery) {
	for column, value := range c.Filter.Exact {
		query.Where("? = ?", bun.Ident(column), value)
	}
}

// addPaginator ...
func (c Crudder) addPaginator(query *bun.SelectQuery) {
	addPaginator(query, c.Paginator, *c.Filter, c.DataModel)
	return
}

// Exists ...
func (c Crudder) Exists() (bool, error) {
	baseQuery := c.DB.NewSelect().
		Model(c.DataModel)
	c.addSelectQuery(baseQuery)

	return baseQuery.Exists(context.TODO())
}

func (c Crudder) Count() (int, error) {
	baseQuery := c.DB.NewSelect().
		Model(c.DataModel)

	c.addSelectQuery(baseQuery)

	return baseQuery.Count(context.Background())
}

// Insert ...
func (c Crudder) Insert() (sql.Result, error) {
	return c.DB.NewInsert().
		Model(c.DataModel).
		Ignore().
		Returning("*").Exec(context.TODO())
}

// Fetch ...
func (c Crudder) Fetch(forQuery ...string) error {
	baseQuery := c.DB.NewSelect().
		Model(c.DataModel)

	c.addSelectQuery(baseQuery)
	c.addPaginator(baseQuery)

	if len(forQuery) > 0 && forQuery[0] != "" {
		baseQuery.For(forQuery[0])
	}

	return baseQuery.Scan(context.Background())
}

// SelectForUpdate finds the row, locks it in place with the tx, and returns it such that no call outside the transaction can affect the row until the associated tx has either been rolled back or committed
func (c Crudder) SelectForUpdate(updateOf ...string) error {
	updateQ := "UPDATE"
	if len(updateOf) > 0 && updateOf[0] != "" {
		updateQ += " of " + updateOf[0]
	}

	baseQuery := c.DB.NewSelect().
		Model(c.DataModel).
		For(updateQ)

	c.addSelectQuery(baseQuery)

	return baseQuery.Scan(context.Background())
}

func (c Crudder) Update() (sql.Result, error) {
	baseQuery := c.DB.NewUpdate().
		Model(c.DataModel).Returning("*")

	if !c.Setter.HasUpdate() {
		return nil, errors.New("invalid update defined")
	}

	c.addUpdateQuery(baseQuery)
	c.addUpdateSetters(baseQuery)
	return baseQuery.Exec(context.Background())
}

func (c Crudder) Delete() (sql.Result, error) {
	baseQuery := c.DB.NewDelete().
		Model(c.DataModel).Returning("*")
	c.addDeleteQuery(baseQuery)

	return baseQuery.Exec(context.Background())
}

func defaultFilter() *DALFilters {
	return &DALFilters{
		Like:     make(map[string]string),
		Exact:    make(map[string]any),
		RegEx:    make(map[string]string),
		Or:       make(map[string]any),
		InInt:    make(map[string]any),
		IsNot:    make(map[string]bool),
		Relation: make(map[string][]func(query *bun.SelectQuery) *bun.SelectQuery),
	}
}

func defaultSorter() *DALSorters {
	return &DALSorters{
		Asc:  make([]string, 0),
		Desc: make([]string, 0),
	}
}

func defaultSetter() *Setter {
	return &Setter{
		Default: make(map[string]any),
		Inc:     make(map[string]float64),
		Dec:     make(map[string]float64),
	}
}

func defaultPaginator(queryV ...url.Values) *PaginationMeta {
	query := url.Values{}
	if len(queryV) > 0 && queryV[0] != nil {
		query = queryV[0]
	}

	var (
		pS  = query.Get("page")
		ppS = query.Get("per-page")
	)

	if query.Get("paginate") != "" {
		paginate, _ := strconv.ParseBool(query.Get("paginate"))
		if !paginate {
			return nil
		}
	}

	paginator := &PaginationMeta{
		Page:    1,
		PerPage: 10,
	}

	pInt, err := strconv.ParseInt(pS, 10, 64)
	if err == nil {
		paginator.Page = int(pInt)
	}

	ppInt, err := strconv.ParseInt(ppS, 10, 64)
	if err == nil {
		paginator.PerPage = int(ppInt)
	}

	return paginator
}

func addQueryFilters(baseQuery *bun.SelectQuery, filter DALFilters) {
	for key, value := range filter.Exact {
		baseQuery.Where("? = ?", bun.Ident(key), value)
	}

	for relation, fn := range filter.Relation {
		baseQuery.Relation(relation, fn...)
	}

	for column, value := range filter.InInt {
		baseQuery.Where("? IN (?)", bun.Ident(column), bun.In(value))
	}

	baseQuery.WhereGroup("AND", func(query *bun.SelectQuery) *bun.SelectQuery {
		for key, value := range filter.Like {
			query.WhereOr("? ilike ?", bun.Ident(key), value)
		}
		return query
	})

	baseQuery.WhereGroup("AND", func(query *bun.SelectQuery) *bun.SelectQuery {
		for key, value := range filter.InOr {
			query.WhereOr("? IN (?)", bun.Ident(key), bun.In(value))
		}
		return query
	})

	baseQuery.WhereGroup("AND", func(query *bun.SelectQuery) *bun.SelectQuery {
		for _, raw := range filter.RawWhereOr {
			query.WhereOr(raw)
		}
		return query
	})

	for column, value := range filter.IsNot {
		baseQuery.Where("? is not ?", bun.Ident(column), value)
	}

	for column, value := range filter.RegEx {
		baseQuery.Where("? ~ ?", column, value)
	}

	if len(filter.Sorter.Asc) > 0 || len(filter.Sorter.Desc) > 0 {
		baseQuery.Order(filter.Sorter.ToQuery()...)
	}

	for _, raw := range filter.RawWhere {
		baseQuery.Where(raw)
	}

	for _, col := range filter.Columns {
		baseQuery.Column(col)
	}
}

// addPaginator processes the pagination information provided by the client and applies it to the
// baseQuery
func addPaginator(baseQuery *bun.SelectQuery, paginator *PaginationMeta, filter DALFilters, dataModel interface{}) (string, string, error) {
	if paginator == nil || paginator.Page == 0 || paginator.PerPage == 0 {
		return value.Success, value.Success, nil
	}
	var err error
	paginator.Total, err = Count(filter, dataModel, baseQuery.DB())
	if err != nil {
		return value.NotAllowed, "Unable to count records", err
	}

	paginator.PageCount = int(math.Ceil(float64(paginator.Total) / float64(paginator.PerPage)))
	paginator.Skipped = (paginator.Page - 1) * paginator.PerPage
	if paginator.PageCount < paginator.Page {
		paginator.Page = paginator.PageCount
	}

	baseQuery.Limit(paginator.PerPage).Offset(paginator.Skipped)

	return value.Success, "Pagination information added successfully", nil
}

func Count(filter DALFilters, dataModel interface{}, db bun.IDB) (int, error) {
	if filter.RawCount != "" {
		var count int
		err := db.NewRaw(filter.RawCount).Scan(context.Background(), &count)
		return count, err
	}

	baseQuery := db.NewSelect().
		Model(dataModel)

	addQueryFilters(baseQuery, filter)

	return baseQuery.Count(context.Background())
}

func CreateNewFilter() DALFilters {
	return DALFilters{
		Like:   map[string]string{},
		Exact:  map[string]interface{}{},
		RegEx:  map[string]string{},
		Or:     map[string]interface{}{},
		InInt:  map[string]interface{}{},
		Sorter: DALSorters{},
	}
}

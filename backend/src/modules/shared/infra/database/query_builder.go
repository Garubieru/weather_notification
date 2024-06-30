package infra_database

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	table   string
	columns []string
	where   string
	joins   []string
	orderBy *string
}

func (queryBuilder *QueryBuilder) SetColumns(columns []string) *QueryBuilder {
	queryBuilder.columns = columns
	return queryBuilder
}

func (queryBuilder *QueryBuilder) Join(join string) *QueryBuilder {
	queryBuilder.joins = append(queryBuilder.joins, join)

	return queryBuilder
}

func (queryBuilder *QueryBuilder) OrderBy(orderBy string) *QueryBuilder {
	queryBuilder.orderBy = &orderBy

	return queryBuilder
}

func (queryBuilder *QueryBuilder) Select() string {
	output := fmt.Sprintf("SELECT %s FROM %s ", strings.Join(queryBuilder.columns, ", "), queryBuilder.table)

	for _, join := range queryBuilder.joins {
		output += join + " "
	}

	output += fmt.Sprintf("WHERE %s", queryBuilder.where)

	if queryBuilder.orderBy != nil {
		output += fmt.Sprintf(" ORDER BY %s", *queryBuilder.orderBy)
	}

	return output
}

func (queryBuilder *QueryBuilder) Insert(count int) string {
	output := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", queryBuilder.table, strings.Join(queryBuilder.columns, ", "))

	for i := 0; i < count; i++ {
		repated := strings.Repeat("?,", len(queryBuilder.columns))
		output += fmt.Sprintf("(%s),", repated[:len(repated)-1])
	}

	return output[:len(output)-1]
}

func (queryBuilder *QueryBuilder) Update() string {
	output := fmt.Sprintf("UPDATE %s SET ", queryBuilder.table)

	for _, column := range queryBuilder.columns {
		output += fmt.Sprintf("%s = ?, ", column)
	}

	output = output[:len(output)-2]

	output += fmt.Sprintf(" WHERE %s", queryBuilder.where)

	return output
}

func (queryBuilder *QueryBuilder) Where(query string) *QueryBuilder {
	queryBuilder.where = query
	return queryBuilder
}

func (queryBuilder *QueryBuilder) Delete() string {
	output := fmt.Sprintf("DELETE FROM %s WHERE %s", queryBuilder.table, queryBuilder.where)

	return output
}

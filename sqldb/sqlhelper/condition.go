package sqlhelper

import "fmt"

func AddConditionClause(column string, value any, query *string, args *[]any) {
	if query == nil {
		return
	}

	*query += fmt.Sprintf(` AND %s = ?`, column)
	*args = append(*args, value)
}

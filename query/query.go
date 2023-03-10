package query

type Query struct {
	queryText string
	_type     string
	results   string
}

func (q *Query) SetQueryText(text string) {
	q.queryText = text
}

func (q *Query) SetType(t string) {
	q._type = t
}

func (q *Query) GetQueryText() string {
	return q.queryText
}

func (q *Query) GetType() string {
	return q._type
}

func (q *Query) Execute() string {
	q.results = "TODO: Implement query execution."
	return q.results
}

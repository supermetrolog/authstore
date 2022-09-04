package mysql

type MysqlQueryError struct {
	Code        string
	Err         error
	QueryString string
}

func NewMysqlQueryError(err error, queryString, code string) *MysqlQueryError {
	return &MysqlQueryError{
		Code:        code,
		Err:         err,
		QueryString: queryString,
	}
}

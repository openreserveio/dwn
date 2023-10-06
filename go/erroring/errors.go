package erroring

type PostgresError struct {
	Msg string
}

func (e *PostgresError) Error() string {
	return e.Msg
}

type RecordError struct {
	Msg string
}

func (e *RecordError) Error() string {
	return e.Msg
}

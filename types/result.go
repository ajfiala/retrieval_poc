package types 

type Result struct {
	Data interface{}
	Error error
	Success bool
}

type ResultChannel chan Result

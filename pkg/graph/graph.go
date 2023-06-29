package graph

type Function struct {
	pack          string
	name          string
	content       string
	funcSignature string
	calls         []Function
}

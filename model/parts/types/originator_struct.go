package types

type OriginatorStruct struct {
	RequestCount int
}

func (o *OriginatorStruct) AddRequest() {
	o.RequestCount++
}

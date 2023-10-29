package declarations

type Group struct {
	Id      uint
	Batch   string
	Names   string
	Project string
	Group1  string
	Group2  string
	Group3  string
	Group4  string
	Group5  string
	Group6  string
	Group7  string
}

type CreateBatch struct {
	Batch string `json:"batch"`
	Names string `json:"names"`
}
package dtos

type WorkTables struct {
	Count    int64  `json:"count"`
	JoinType string `json:"join_type"`
	Env      string `json:"env"`
}

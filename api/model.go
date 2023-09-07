package api

type Telegraph struct {
	Title    string `json:"title"`
	Brief    string `json:"brief"`
	Content  string `json:"content"`
	CTime    int64  `json:"ctime"`
	Subjects []struct {
		SubjectName string `json:"subject_name"`
	} `json:"subjects"`
}

type RefreshTelegraphListResponse struct {
	L map[string]*Telegraph `json:"l"`
}

type RollTelegraphListResponse struct {
	Data struct {
		RollData []*Telegraph `json:"roll_data"`
	} `json:"data"`
}

package models

type TestCase struct {
	ProblemId string `json:"problemId"`
	Serial    int    `json:"serial"`
	In        string `json:"in"`
	Out       string `json:"out"`
}

type TestCaseHeader struct {
	Serial     int    `json:"serial"`
	Name       string `json:"name"`
	InputSize  int    `json:"inputSize"`
	OutputSize int    `json:"outputSize"`
	Score      int    `json:"score"`
}

package models

import "encoding/json"

type User struct {
	ID                         string          `json:"id"`
	Name                       string          `json:"name"`
	Affiliation                *string         `json:"affiliation"`
	RegisterDate               int64           `json:"registerDate"`
	LastSubmitDate             int64           `json:"lastSubmitDate"`
	Policy                     string          `json:"policy"`
	Country                    string          `json:"country"`
	BirthYear                  int             `json:"birthYear"`
	DisplayLanguage            string          `json:"displayLanguage"`
	DefaultProgrammingLanguage string          `json:"defaultProgrammingLanguage"`
	Status                     json.RawMessage `json:"status"`
	URL                        *string         `json:"url"`
}

type UserStatus struct {
	Submissions  int `json:"submissions"`
	Solved       int `json:"solved"`
	Accepted     int `json:"accepted"`
	WrongAnswer  int `json:"wrongAnswer"`
	TimeLimit    int `json:"timeLimit"`
	MemoryLimit  int `json:"memoryLimit"`
	OutputLimit  int `json:"outputLimit"`
	CompileLimit int `json:"compileLimit"`
	CompileError int `json:"compileError"`
	RuntimeError int `json:"runtimeError"`
}

func (user User) GetStatus() (*UserStatus, error) {
	var status UserStatus
	err := json.Unmarshal(user.Status, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

package models

type SubmissionStatus int

const (
	CompileError = iota
	WrongAnswer
	TimeLimit
	MemoryLimit
	Accepted
	Waiting
	OutputLimit
	RuntimeError
	PresentationError
	Running
)

func (status SubmissionStatus) String() string {
	switch status {
	case CompileError:
		return "Compile Error"
	case WrongAnswer:
		return "Wrong Answer"
	case TimeLimit:
		return "Time Limit"
	case MemoryLimit:
		return "Memory Limit"
	case Accepted:
		return "Accepted"
	case OutputLimit:
		return "Output Limit"
	case RuntimeError:
		return "Runtime Error"
	case PresentationError:
		return "Presentation Error"
	case Running:
		return "Running"
	default:
		return "Unknown Status"
	}
}

type SubmissionRecord struct {
	JudgeId        int              `json:"judgeId"`
	JudgeType      int              `json:"judgeType"`
	UserId         string           `json:"userId"`
	ProblemId      string           `json:"problemId"`
	SubmissionDate int64            `json:"submissionDate"`
	Language       string           `json:"language"`
	Status         SubmissionStatus `json:"status"`
	CpuTime        int              `json:"cpuTime"`
	Memory         int              `json:"memory"`
	CodeSize       int              `json:"codeSize"`
	Accuracy       *string          `json:"accuracy"`
	JudgeDate      int64            `json:"judgeDate"`
	Score          int              `json:"score"`
	ProblemTitle   *string          `json:"problemTitle"`
	Token          *string          `json:"token"`
}

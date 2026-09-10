package pkg

import "fmt"

type Severity string

const (
	SeverityInfo  Severity = "INFO"
	SeverityWarn  Severity = "WARN"
	SeverityError Severity = "ERROR"
	SeverityFatal Severity = "FATAL"
)

type CoreBankingError struct {
	Code        string
	TraceID     string
	Severity    Severity
	Message     string
	IsRetryable bool
	Err         error
}

func (e *CoreBankingError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s (Trace: %s): %s - %v", e.Severity, e.Code, e.TraceID, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s (Trace: %s): %s", e.Severity, e.Code, e.TraceID, e.Message)
}

func (e *CoreBankingError) Unwrap() error {
	return e.Err
}

// Helper func untuk generate error
func NewError(code, traceID string, severity Severity, msg string, retryable bool, err error) *CoreBankingError {
	return &CoreBankingError{
		Code:        code,
		TraceID:     traceID,
		Severity:    severity,
		Message:     msg,
		IsRetryable: retryable,
		Err:         err,
	}
}

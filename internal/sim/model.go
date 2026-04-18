package sim

import "time"

type Severity int

const (
	Info Severity = iota
	Low
	Medium
	High
	Critical
)

type LogEntry struct {
	Timestamp   time.Time
	Severity    Severity
	Hostname    string
	User        string
	SrcIP       string
	DstIP       string
	ProcessName string
	Category    string
	Message     string
}

type AlertStatus string

const (
	StatusNew           AlertStatus = "New"
	StatusInvestigating AlertStatus = "Investigating"
	StatusContained     AlertStatus = "Contained"
	StatusResolved      AlertStatus = "Resolved"
)

type Alert struct {
	ID          string
	Timestamp   time.Time
	Severity    Severity
	RuleName    string
	Host        string
	User        string
	Explanation string
	Status      AlertStatus
}

type Metric struct {
	Timestamp time.Time
	Name      string
	Value     float64
}

type EDRState string

const (
	EDRActive   EDRState = "Active"
	EDROffline  EDRState = "Offline"
	EDRWarning  EDRState = "Warning"
)

type Asset struct {
	Hostname           string
	Criticality        string
	RiskScore          int
	EDRState           EDRState
	LastSeen           time.Time
	ActiveAlerts       int
	Isolated           bool
	SuspiciousProcs    int
}

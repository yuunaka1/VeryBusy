package sim

import (
	"fmt"
	"math/rand"
	"time"
)

type Engine struct {
	rnd         *rand.Rand
	seed        int64
	theme       string
	
	// internal state
	tickCount   int
	Assets      []*Asset
	Alerts      []Alert
	Logs        []LogEntry
	Metrics     map[string][]Metric
	
	lastTime    time.Time
}

func NewEngine(seed int64, theme string) *Engine {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	e := &Engine{
		rnd:      rand.New(rand.NewSource(seed)),
		seed:     seed,
		theme:    theme,
		Metrics:  make(map[string][]Metric),
		lastTime: time.Now(),
	}
	e.initAssets()
	e.initAlerts()
	return e
}

func (e *Engine) initAssets() {
	// generate some fake assets
	for i := 0; i < 15; i++ {
		e.Assets = append(e.Assets, &Asset{
			Hostname:    fmt.Sprintf("WKSTN-%04d", e.rnd.Intn(9999)),
			Criticality: "Medium",
			RiskScore:   e.rnd.Intn(10),
			EDRState:    EDRActive,
			LastSeen:    time.Now().Add(-time.Duration(e.rnd.Intn(60)) * time.Second),
		})
	}
	e.Assets = append(e.Assets, &Asset{
		Hostname:    "SRV-DC-01",
		Criticality: "Critical",
		RiskScore:   0,
		EDRState:    EDRActive,
		LastSeen:    time.Now(),
	})
	e.Assets = append(e.Assets, &Asset{
		Hostname:    "SRV-DB-PROD",
		Criticality: "Critical",
		RiskScore:   0,
		EDRState:    EDRActive,
		LastSeen:    time.Now(),
	})
}

func (e *Engine) Tick() {
	e.tickCount++
	now := time.Now()
	
	e.generateAlerts(now)
	e.updateAssets(now)
	e.updateMetrics(now)
	
	e.lastTime = now
}

type logTemplate struct {
	Category    string
	ProcessName string
	Message     string
}

var infoTemplates = []logTemplate{
	{"Auth", "lsass.exe", "Successful interactive logon"},
	{"Auth", "lsass.exe", "Kerberos authentication ticket requested"},
	{"Auth", "sshd", "Accepted publickey session"},
	{"Network", "svchost.exe", "Network connection accepted on port 443"},
	{"Network", "svchost.exe", "Session disconnected due to timeout"},
	{"Network", "firewalld", "Connection dropped by default policy"},
	{"OS", "services.exe", "Service state changed from STOPPED to RUNNING"},
	{"OS", "taskeng.exe", "Scheduled task completed successfully"},
	{"OS", "gpupdate.exe", "Group policy object updated"},
	{"OS", "wininit.exe", "System time synchronized with NTP server"},
	{"App", "docker.exe", "Container stopped gracefully"},
	{"App", "nginx", "HTTP request proxied to internal upstream"},
	{"DB", "postgres", "Connection authorized: user=app_service"},
}

var suspiciousTemplates = []logTemplate{
	{"Process", "powershell.exe", "Suspicious encoded PowerShell command executed"},
	{"Network", "curl.exe", "Outbound connection to known malicious IP (ThreatIntel)"},
	{"Auth", "sshd", "Multiple failed login attempts detected (Brute-force)"},
	{"File", "explorer.exe", "Mass file modification detected (Possible Ransomware)"},
	{"OS", "vssadmin.exe", "Volume Shadow Copies deletion attempt"},
	{"Registry", "reg.exe", "Run keys modified for persistence"},
}

func (e *Engine) GenerateLogs(now time.Time) {
	// Add 1-3 noise logs per tick
	count := 1 + e.rnd.Intn(3)
	for i := 0; i < count; i++ {
		host := e.Assets[e.rnd.Intn(len(e.Assets))].Hostname
		user := fmt.Sprintf("user%d", e.rnd.Intn(50))
		
		tmpl := infoTemplates[e.rnd.Intn(len(infoTemplates))]
		
		log := LogEntry{
			Timestamp:   now,
			Severity:    Info,
			Hostname:    host,
			User:        user,
			ProcessName: tmpl.ProcessName,
			Category:    tmpl.Category,
			Message:     tmpl.Message,
		}
		e.Logs = append(e.Logs, log)
	}
	
	// Randomly add a medium/high log
	if e.rnd.Float32() < 0.05 {
		host := e.Assets[e.rnd.Intn(len(e.Assets))].Hostname
		user := fmt.Sprintf("user%d", e.rnd.Intn(50))
		
		tmpl := suspiciousTemplates[e.rnd.Intn(len(suspiciousTemplates))]
		
		log := LogEntry{
			Timestamp:   now,
			Severity:    Medium,
			Hostname:    host,
			User:        user,
			ProcessName: tmpl.ProcessName,
			Category:    tmpl.Category,
			Message:     tmpl.Message,
		}
		e.Logs = append(e.Logs, log)
	}
	
	// Keep log slice bounded
	if len(e.Logs) > 100 {
		e.Logs = e.Logs[len(e.Logs)-100:]
	}
}

type alertTemplate struct {
	RuleName    string
	Explanation string
	Severity    Severity
}

var alertTemplates = []alertTemplate{
	{"Suspicious Encoded Command", "A powershell command was executed with -enc flag indicating potential obfuscation.", High},
	{"Multiple Failed Logins", "User account experienced multiple authentication failures within a short time.", Medium},
	{"Malicious IP Communication", "Outbound network connection detected to a known bad IP address.", Critical},
	{"Possible Ransomware Activity", "Mass file renaming and high volume of file accesses detected.", Critical},
	{"Privilege Escalation Attempt", "Unexpected process spawned as SYSTEM user.", High},
}

func (e *Engine) initAlerts() {
	count := 3 + e.rnd.Intn(3) // 3 to 5 initial alerts
	now := time.Now()
	
	for i := 0; i < count; i++ {
		host := e.Assets[e.rnd.Intn(len(e.Assets))].Hostname
		user := fmt.Sprintf("user%d", e.rnd.Intn(50))
		tmpl := alertTemplates[e.rnd.Intn(len(alertTemplates))]

		alert := Alert{
			ID:          fmt.Sprintf("ALT-%05d", e.rnd.Intn(99999)),
			Timestamp:   now.Add(-time.Duration(e.rnd.Intn(3600)) * time.Second), // past hour
			Severity:    tmpl.Severity,
			RuleName:    tmpl.RuleName,
			Host:        host,
			User:        user,
			Explanation: tmpl.Explanation,
			Status:      StatusNew,
		}
		
		r := e.rnd.Float32()
		if r < 0.2 {
			alert.Status = StatusInvestigating
		} else if r < 0.4 {
			alert.Status = StatusContained
		}
		
		e.Alerts = append(e.Alerts, alert)
		
		for _, a := range e.Assets {
			if a.Hostname == host {
				a.RiskScore += 15
				a.ActiveAlerts++
			}
		}
	}
}

func (e *Engine) generateAlerts(now time.Time) {
	if e.rnd.Float32() < 0.02 { // 2% chance per tick for an alert
		host := e.Assets[e.rnd.Intn(len(e.Assets))].Hostname
		user := fmt.Sprintf("user%d", e.rnd.Intn(50))
		tmpl := alertTemplates[e.rnd.Intn(len(alertTemplates))]

		alert := Alert{
			ID:          fmt.Sprintf("ALT-%05d", e.rnd.Intn(99999)),
			Timestamp:   now,
			Severity:    tmpl.Severity,
			RuleName:    tmpl.RuleName,
			Host:        host,
			User:        user,
			Explanation: tmpl.Explanation,
			Status:      StatusNew,
		}
		e.Alerts = append(e.Alerts, alert)
		
		// Map it to asset
		for _, a := range e.Assets {
			if a.Hostname == host {
				a.RiskScore += 15
				a.ActiveAlerts++
			}
		}
	}
	
	if len(e.Alerts) > 50 {
		e.Alerts = e.Alerts[len(e.Alerts)-50:]
	}
}

func (e *Engine) updateAssets(now time.Time) {
	for _, a := range e.Assets {
		a.LastSeen = now.Add(-time.Duration(e.rnd.Intn(10)) * time.Second)
		if a.RiskScore > 0 && e.rnd.Float32() < 0.05 {
			a.RiskScore-- // naturally decay risk score
		}
	}
}

func (e *Engine) updateMetrics(now time.Time) {
	mName := "Network Traffic (Anomaly Score)"
	baseValue := 50.0 + e.rnd.Float64()*10.0
	// Spike randomly
	if e.rnd.Float32() < 0.1 {
		baseValue += 30.0 + e.rnd.Float64()*20.0
	}
	
	e.Metrics[mName] = append(e.Metrics[mName], Metric{Timestamp: now, Name: mName, Value: baseValue})
	if len(e.Metrics[mName]) > 60 {
		e.Metrics[mName] = e.Metrics[mName][len(e.Metrics[mName])-60:]
	}
}

func (e *Engine) Theme() string {
	return e.theme
}

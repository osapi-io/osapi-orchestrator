package orchestrator

// Values these tests build fixtures from. They live here rather than
// inline because each is used in four or more places, and something
// that means the same thing in every one of them should say so once.
const (
	caseAuthFailure    = "Returns error on auth failure"
	caseDuplicateName  = "Duplicate name gets counter suffix"
	caseNoDependencies = "Returns false with no dependencies"
	caseNoHostResults  = "Returns false when dep has no HostResults (unicast)"
	caseServerError    = "Returns error on server error"
	distributionUbuntu = "Ubuntu"
	errTimeout         = "timeout"
	fieldAgents        = "agents"
	fieldDistribution  = "distribution"
	fieldHostname      = "hostname"
	fieldOSInfo        = "os_info"
	fieldStdout        = "stdout"
	fieldTotal         = "total"
	hostNerd           = "nerd"
	hostStatusFailed   = "failed"
	hostStatusOK       = "ok"
	hostStatusSkipped  = "skipped"
	hostWeb01          = "web-01"
	hostWeb02          = "web-02"
	opRunCommand       = "run-command"
	stepDeploy         = "deploy"
	stepHealthCheck    = "health-check"
	stepListAgents     = "list-agents"
	taskDep            = "dep"
	taskDepA           = "dep-a"
	testToken          = "test-token"
)

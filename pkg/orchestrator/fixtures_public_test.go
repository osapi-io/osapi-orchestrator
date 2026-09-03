package orchestrator_test

// Values these tests build fixtures from. They live here rather than
// inline because each is used in four or more places, and something
// that means the same thing in every one of them should say so once.
const (
	anyAgent                = "_any"
	archX8664               = "x86_64"
	caseReturnsNonNilStep   = "Returns non-nil step"
	conditionDiskPressure   = "DiskPressure"
	conditionMemoryPressure = "MemoryPressure"
	containerC1             = "c1"
	distributionUbuntu      = "Ubuntu"
	envProd                 = "prod"
	factOSDistribution      = "os.distribution"
	fieldExitCode           = "exit_code"
	fieldResults            = "results"
	fieldStdout             = "stdout"
	fileContents            = "hello"
	hostWeb01               = "web-01"
	hostWeb02               = "web-02"
	hostWeb03               = "web-03"
	interfaceEth0           = "eth0"
	labelDatacenter         = "datacenter"
	labelEnv                = "env"
	osFamilyUbuntu          = "ubuntu"
	regionUSEast1           = "us-east-1"
	serviceNginx            = "nginx"
	stepA                   = "step-a"
	stepBadMarshal          = "bad-marshal"
	stepDeploy              = "deploy"
	stepNonexistent         = "nonexistent"
	stepRunCmd              = "run-cmd"
	taskA                   = "a"
	userAdmin               = "admin"
)

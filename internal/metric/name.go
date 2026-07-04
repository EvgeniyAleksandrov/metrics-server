package metric

type Name string

// Gauge
const (
	// MemStat
	Alloc         Name = "Alloc"
	BuckHashSys   Name = "BuckHashSys"
	Frees         Name = "Frees"
	GCCPUFraction Name = "GCCPUFraction"
	GCSys         Name = "GCSys"
	HeapAlloc     Name = "HeapAlloc"
	HeapIdle      Name = "HeapIdle"
	HeapInuse     Name = "HeapInuse"
	HeapObjects   Name = "HeapObjects"
	HeapReleased  Name = "HeapReleased"
	HeapSys       Name = "HeapSys"
	LastGC        Name = "LastGC"
	Lookups       Name = "Lookups"
	MCacheInuse   Name = "MCacheInuse"
	MCacheSys     Name = "MCacheSys"
	MSpanInuse    Name = "MSpanInuse"
	MSpanSys      Name = "MSpanSys"
	Mallocs       Name = "Mallocs"
	NextGC        Name = "NextGC"
	NumForcedGC   Name = "NumForcedGC"
	NumGC         Name = "NumGC"
	OtherSys      Name = "OtherSys"
	PauseTotalNs  Name = "PauseTotalNs"
	StackInuse    Name = "StackInuse"
	StackSys      Name = "StackSys"
	Sys           Name = "Sys"
	TotalAlloc    Name = "TotalAlloc"

	// Rangom
	RandomValue Name = "RandomValue"
)

// Counter
const PullCount Name = "PullCount"

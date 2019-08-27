package config

// Config for backup strategies.
type Config struct {
	Namespace   string
	Image       string
	Prefix      string
	CronSplit   int
	Bucket      string
	Credentials Credentials
	Resources   Resources
}

// Credentials for backup strategies.
type Credentials struct {
	ID     string
	Secret string
}

// Resources for backup tasks.
type Resources struct {
	CPU    string
	Memory string
}

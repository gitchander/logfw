package logfw

const (
	Kilobyte = 1024
	Megabyte = 1024 * Kilobyte
	Gigabyte = 1024 * Megabyte
)

type Config struct {
	FileName   string `json:"file-name,omitempty"`
	MaxSize    int64  `json:"max-size,omitempty"`
	MaxBackups int    `json:"max-backups,omitempty"`
}

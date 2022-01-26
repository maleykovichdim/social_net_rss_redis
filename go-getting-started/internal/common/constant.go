package common

import (
	_ "github.com/go-sql-driver/mysql"
	// m "go-getting-started/internal/common"
)

const (
	// UserDB     = "root3"
	// PasswordDB = "Aa123456$"

	RedisAddress  = "localhost:6379"
	RedisPassword = "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81"
)

var (
	DatabaseURL      = "root:root@/socialnet?parseTime=true"
	DatabaseURL_read = "root:root@/socialnet?parseTime=true"
	// docker version
	// databaseURL      = env("JAWSDB_URL", "root:toor@tcp(alpha:3306)/socialnet?parseTime=true")
	// databaseURL_read = env("JAWSDB_URL", "root:toor@tcp(slave:3306)/socialnet?parseTime=true")
)

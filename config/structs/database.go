package structs

import (
	"fmt"
)

// DatabaseConfiguration represents a database configuration
type DatabaseConfiguration struct {
	DatabaseName string
	Username     string
	Password     string
	Address      string `type:"optional"`
	Port         int    `type:"optional"`
}

// MariaDBConnectionString returns a connection string for MariaDB
func (c *DatabaseConfiguration) MariaDBConnectionString() (string, string) {
	return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=True",
		c.Username, c.Password, c.Address, c.Port, c.DatabaseName)
}

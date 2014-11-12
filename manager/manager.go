/*
Package manager implements the gron manager.
*/
package manager

import (
	"fmt"
)

// ReadCrontab reads the crontab
func ReadCrontab(crontabPath string) (x int, err error) {
	x, err = 0, nil

	fmt.Println("Reading crontab!")

	return x, err
}

package nullstore

import (
	"./statmsg"
	"fmt"
)

func Update(stat *statmsg.Statmsg) {
	fmt.Printf("%s %s %s %s\n", stat.Key, stat.IP, stat.Referer, stat.UA)
}

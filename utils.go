package navitas

import (
	"fmt"
	"regexp"
	"runtime"
	"time"
)

// LoadTime calculates function execution time. To use, add
// defer c.LoadTime(time.Now()) to the function body
func (n *Navitas) LoadTime(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	n.InfoLog.Println(fmt.Sprintf("Load Time: %s took %s", name, elapsed))
}

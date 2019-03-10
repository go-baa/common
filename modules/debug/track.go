package debug

import (
	"log"
	"runtime"
	"time"
)

// TrackFuncTimeUsage 记录方法调用时间
// Usage: defer TrackFuncTimeUsage(time.Now())
func TrackFuncTimeUsage(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a Function object this functions parent
	functionObject := runtime.FuncForPC(pc)

	log.Printf("%s took %s", functionObject.Name(), elapsed)
}

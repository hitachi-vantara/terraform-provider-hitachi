package common

import (
	// "fmt"
	// "strings"
	"runtime"
	"path"
)

func getSourceFileInfo(depth int) (string, string, string, int) {

	pc, fullFilePath, ln, ok := runtime.Caller(depth)

	if ok {
		cf := runtime.FuncForPC(pc).Name()
		dirPath, fileName := path.Split(fullFilePath)
		// fmt.Println(" ----------------------------- getSourceFileInfo.cf=", cf)
		// fmt.Println(" ----------------------------- dirPath=", dirPath)
		// fmt.Println(" ----------------------------- fileName=", fileName)
		var moduleName string
		if dirPath != "" {
			dirPath = dirPath[:len(dirPath)-1]
			_, moduleName = path.Split(dirPath)
		}
		fn := moduleName + "/" + fileName
		// fmt.Println(" ----------------------------- moduleName=", moduleName)
		// fmt.Println(" ----------------------------- ln=", ln)

		return dirPath, cf, fn, ln
	}

	return "", "", "", 0
}

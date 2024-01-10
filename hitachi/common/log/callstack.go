package common

import (
	"fmt"
	"runtime"
	"strings"
)

// return the file path that including the name of immediate folder and the file name with line number
func ShortFileName(fileNPathame string, fileLine int) string {
	sections := strings.Split(fileNPathame, "/")
	n := len(sections)
	if n > 0 {
		sections[n-1] = strings.Replace(sections[n-1], ".go", "", -1)
		if n >= 2 {
			sections = sections[n-2 : n]
		}
	}
	//fmt.Printf("%v\n", sections)
	shortFilePath := strings.Join(sections, "/")
	shortFilePath = fmt.Sprintf("%s:%d", shortFilePath, fileLine)
	return shortFilePath
}

func GetCallStack(index int) (runtime.Frame, error) {
	var callStack runtime.Frame
	var err error = nil
	// Caller(skip int) (pc uintptr, file string, line int, ok bool)
	pc, _, _, isOk := runtime.Caller(index)
	if !isOk {
		return callStack, fmt.Errorf(fmt.Sprintf("No call stack available for frame with index [%v]", index))
	}
	//fmt.Printf("runtime.Caller: pc=[%v], fileName=[%v], line=[%v]\n", pc, fileName, line)
	pcs := []uintptr{pc}
	frames := runtime.CallersFrames(pcs)
	var frame runtime.Frame
	for {
		fr, more := frames.Next()
		frame = fr
		//fmt.Sprintf("***** more=[%v] | function=[%s] | file=[%s] | line=[%v]\n", more, frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	frame.File = ShortFileName(frame.File, frame.Line)
	return frame, err
}

func GetCallStack2(index int) (runtime.Frame, error) {
	var callStack runtime.Frame
	var err error = nil
	callStacks := getCallStacks()
	// if len(*callStacks) >= 2 {
	// 	*callStacks = (*callStacks)[2:]
	// }
	// if len(*callStacks) > index {
	// 	callStack = (*callStacks)[index]
	// } else {
	// 	err = errors.New("No call stack for the index")
	// }

	// return callStack, err
	fmt.Print(callStacks)
	return callStack, err
}

func getCallStacks() []runtime.Frame {
	runtimeFrames := []runtime.Frame{}
	// runtimeFrames := make([]runtime.Frame, 1)
	// Ask runtime.Callers for up to 10 pcs, including runtime.Callers itself.
	// pc := make([]uintptr, 50)
	// n := runtime.Callers(0, pc)
	// fmt.Printf("n=%v\n", n)
	// fmt.Printf("runtime.Callers=%v\n", pc)
	// if n == 0 {
	// 	// No pcs available. Stop now.
	// 	// This can happen if the first argument to runtime.Callers is large.
	// 	return &runtimeFrames
	// }

	// pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	// fmt.Printf("untime.Callers=%v\n", pc)
	// frames := runtime.CallersFrames(pc)
	// for {
	// 	frame, more := frames.Next()
	// 	runtimeFrames = append(runtimeFrames, frame)
	// 	fmt.Printf("***** more:%v | function=%s | file=%s | line-%v\n", more, frame.Function, frame.File, frame.Line)
	// 	if !more {
	// 		break
	// 	}
	// }
	return runtimeFrames
}

// func GetSourceFileInfo(depth int) (string, string, string, int) {

// 	pc, fullFilePath, ln, ok := runtime.Caller(depth)

// 	if ok {
// 		cf := runtime.FuncForPC(pc).Name()
// 		dirPath, fileName := path.Split(fullFilePath)
// 		// fmtPrintln(" ----------------------------- getSourceFileInfo.cf=", cf)
// 		// fmtPrintln(" ----------------------------- dirPath=", dirPath)
// 		// fmtPrintln(" ----------------------------- fileName=", fileName)
// 		var moduleName string
// 		if dirPath != "" {
// 			dirPath = dirPath[:len(dirPath)-1]defer
// 			_, moduleName = path.Split(dirPath)
// 		}
// 		fn := moduleName + "/" + fileName
// 		//fmtPrintln(" ----------------------------- moduleName=", moduleName)

// 		return dirPath, cf, fn, ln
// 	}

// 	return "", "", "", 0
// }

package telemetry

import (
	"runtime"
	"strings"
	"sync"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
)

var (
	lastTerraformCaller string
	lastCallerLock      sync.RWMutex
)

// GetGatewayCallStackInfo returns connection type and gateway method name
func GetGatewayCallStackInfo() (string, string) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	callStackMethodNames := captureCallStack(10)
	log.WriteDebug("CALLSTACK:GetGatewayCallStackInfo: %s", strings.Join(callStackMethodNames, "\n"))

	methodStack := ""
	methodSearch := "gateway/impl."
	for _, s := range callStackMethodNames {
		if methodStack == "" && strings.Contains(s, methodSearch) {
			methodStack = s
			break
		}
	}

	return parseGatewayMethodStack(methodStack)
}

// GetTerraformCallStackInfo returns terraform resource method
// call this in gateway manager
func GetTerraformCallStackInfo() string {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	callStackMethodNames := captureCallStack(50)
	log.WriteDebug("CALLSTACK:GetResourceCallStackInfo: %s", strings.Join(callStackMethodNames, "\n"))

	resourceStack := ""
	resourceSearches := []string{
		"terraform/resource.",
		"terraform/datasource.",
		"terraform.providerConfigure",
	}

	for _, s := range callStackMethodNames {
		for _, search := range resourceSearches {
			if strings.Contains(s, search) {
				resourceStack = s
				break
			}
		}
		if resourceStack != "" {
			break
		}
	}

	caller := parseTerraformMethodStack(resourceStack)

	// ✅ Don’t cache any prefix ending with terraform.providerConfigure
	if caller != "" && !strings.HasSuffix(caller, "terraform.providerConfigure") {
		lastCallerLock.Lock()
		lastTerraformCaller = caller
		lastCallerLock.Unlock()
		log.WriteDebug("Updated lastTerraformCaller: %s", caller)
	}

	if caller != "" {
		return caller
	}

	// ✅ If not found, use last cached value
	lastCallerLock.RLock()
	cached := lastTerraformCaller
	lastCallerLock.RUnlock()

	if cached != "" {
		log.WriteDebug("Using cached lastTerraformCaller: %s", cached)
		return cached
	}

	return ""
}

// parseTerraformMethodStack
// example: resource.resourceVosbChangeUserPasswordCreate or terraform.providerConfigure
func parseTerraformMethodStack(resourceStack string) string {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("resourceStack: %s", resourceStack)

	if resourceStack == "" {
		return ""
	}

	parts := strings.Split(resourceStack, "/")
	terraformMethod := parts[len(parts)-1]
	log.WriteDebug("terraformMethod: %s", terraformMethod)
	return terraformMethod
}

// parseGatewayMethodStack
// connectionType and gateway methodName: example: "vosb" and "ChangeUserPassword"
func parseGatewayMethodStack(methodStack string) (string, string) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("methodStack: %s", methodStack)
	if methodStack == "" {
		return "", ""
	}

	parts := strings.Split(methodStack, "/")
	connectionType := parts[3]
	methodName := strings.Split(parts[len(parts)-1], ".")[2]
	return connectionType, methodName
}

// capture the call stack
func captureCallStack(depth int) []string {
	pc := make([]uintptr, depth)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	callStackMethodNames := []string{}
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		callStackMethodNames = append(callStackMethodNames, frame.Function)
		if len(callStackMethodNames) >= depth {
			break
		}
	}
	return callStackMethodNames
}

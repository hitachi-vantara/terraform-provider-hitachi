package telemetry

import (
	"runtime"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
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

	callStackMethodNames := captureCallStack(10)
	log.WriteDebug("CALLSTACK:GetResourceCallStackInfo: %s", strings.Join(callStackMethodNames, "\n"))

	resourceStack := ""
	resourceSearch1 := "terraform/resource."
	resourceSearch2 := "terraform/datasource."
	resourceSearch3 := "terraform.providerConfigure"
	for _, s := range callStackMethodNames {
		if strings.Contains(s, resourceSearch1) || strings.Contains(s, resourceSearch2) || strings.Contains(s, resourceSearch3) {
			resourceStack = s
			break
		}
	}

	return parseTerraformMethodStack(resourceStack)
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

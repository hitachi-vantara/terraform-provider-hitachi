package vssbstorage

import (
	"fmt"

	// "strconv"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
)

func GetUrl(ip string, urlPath string) string {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	url := fmt.Sprintf("https://%s/ConfigurationManager/simple/v1/%s", ip, urlPath)
	log.WriteDebug("TFDebug|url: %s", url)
	return url
}

func GetUrlWithoutVersion(ip string, urlPath string) string {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	url := fmt.Sprintf("https://%s/ConfigurationManager/simple/%s", ip, urlPath)
	log.WriteDebug("TFDebug|url: %s", url)
	return url
}

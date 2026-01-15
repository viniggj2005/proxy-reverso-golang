package shared

import "net/url"

func VerifyTlsConnection(receivedUrl string) bool {
	parsedUrl, err := url.Parse(receivedUrl)
	if err != nil {
		return false
	}
	return parsedUrl.Scheme == "https" || parsedUrl.Scheme == "wss"
}

package lib

import (
	"fmt"
	"regexp"
)

func GetServiceName (key string) (string, error) {
	pattern := `^/#/svc/(\w+)/.*$`
	is, err := regexp.MatchString(pattern, key)
	if err != nil { return "", err }
	if !is { return "", fmt.Errorf("invalid service key: %s", key)}
	rx, err := regexp.Compile(pattern)
	if err != nil { return "", err }
	groups := rx.FindStringSubmatch(key)
	serviceName := groups[1]

	return serviceName, nil
}
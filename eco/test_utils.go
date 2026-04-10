package eco

import (
	"strconv"
	"strings"
)

func getSliceKeyIndex (key string) (int, error) {
	sKey := string(key)
	iLastSlash := strings.LastIndex(sKey, "/")
	sIndex := sKey[iLastSlash+1:]
	index, err := strconv.Atoi(sIndex)
	if err != nil {
		return 0, err
	}
	return index, nil
}

package functions

import (
	"crypto/md5"
	"encoding/hex"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nxtcoder19/nthreads-backend/package/errors"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(\W|_)+`)

func CleanerNanoid(n int) (string, error) {
	id, e := nanoid.New(n)
	if e != nil {
		return "", errors.NewEf(e, "could not get nanoid()")
	}
	res := re.ReplaceAllString(id, "-")
	if strings.HasPrefix(res, "-") {
		res = "k" + res
	}
	if strings.HasSuffix(res, "-") {
		res = res + "k"
	}
	return res, nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

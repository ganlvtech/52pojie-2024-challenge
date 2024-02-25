package dynamicflag

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

func MD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func CalcFlag(uid string, secret string, t time.Time) (flag string, expiredAt time.Time) {
	flag = MD5(fmt.Sprintf("%s%s%d", uid, secret, t.Unix()/600))[:8]
	expiredAt = time.Unix(((t.Unix()/600)+1)*600, 0)
	return
}

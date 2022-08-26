package libs

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
)

/*
	MD5算法加密
*/
func GetMD5String(str string) string {
	m := md5.New()
	if _, err := io.WriteString(m, str); err != nil {
		log.Fatal(err)
	}
	res := m.Sum(nil)
	return fmt.Sprintf("%x", res)
}

package auth

import (
	"fmt"
	"time"
)

// User represents authenticated user
type User struct {
	uid       uint32
	account   string
	pwdHash   string
	loginTime time.Time
}

func (u *User) String() string {
	return fmt.Sprintf("User<%d:%s:%d>", u.uid, u.account, u.loginTime.Unix())
}

func (u *User) checkCredentials(account uint32, pwd string) error {
	// h := md5.New()
	// h.Write([]byte(pwd))
	// pwd = string(h.Sum(nil))
	if account != u.uid || pwd != u.pwdHash {
		return fmt.Errorf("credentials mismatch, want: %d:%s, got: %d:%s", u.uid, u.pwdHash, account, pwd)
	}

	return nil
}

package consts

const LoginProhibited = LoginError("login prohibited")
const LoginIncorrentAccountOrPassword = LoginError("login incorrent account or password")

type LoginError string

func (l LoginError) Error() string { return string(l) }

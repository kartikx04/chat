package auth

import (
	"github.com/gorilla/sessions"
	"github.com/kartikx04/chat/pkg"
)

var Store = sessions.NewCookieStore([]byte(pkg.LoadFile("TOKEN_SECRET")))

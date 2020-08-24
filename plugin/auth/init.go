package auth

import (
	"github.com/appootb/substratum/auth"
	"github.com/appootb/substratum/token"
)

func Init() {
	if auth.Implementor() == nil {
		auth.RegisterImplementor(auth.NewAlgorithmAuth(token.Implementor(), token.Implementor()))
	}
}

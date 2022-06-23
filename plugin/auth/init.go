package auth

import (
	"github.com/appootb/substratum/v2/auth"
	"github.com/appootb/substratum/v2/token"
)

func Init() {
	if auth.Implementor() == nil {
		auth.RegisterImplementor(auth.NewAlgorithmAuth(token.Implementor(), token.Implementor()))
	}
}

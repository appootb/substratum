package credential

import "github.com/appootb/substratum/credential"

func Init() {
	if credential.ClientImplementor() == nil {
		credential.RegisterClientImplementor(&ClientSeed{})
	}
	if credential.ServerImplementor() == nil {
		credential.RegisterServerImplementor(&ServerSeed{})
	}
}

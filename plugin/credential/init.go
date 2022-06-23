package credential

import "github.com/appootb/substratum/v2/credential"

func Init() {
	if credential.ClientImplementor() == nil {
		credential.RegisterClientImplementor(&ClientSeed{})
	}
	if credential.ServerImplementor() == nil {
		credential.RegisterServerImplementor(&ServerSeed{})
	}
}

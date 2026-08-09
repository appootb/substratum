module github.com/appootb/substratum/v2/tools

go 1.25.4

require (
	github.com/appootb/substratum/v2 v2.0.5-0.20260807102556-ab45c9beb291
	github.com/iancoleman/strcase v0.3.0
	github.com/lyft/protoc-gen-star/v2 v2.0.4
	google.golang.org/genproto/googleapis/api v0.0.0-20260807164820-c8921c73eeea
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/spf13/afero v1.3.3 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/mod v0.37.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	golang.org/x/tools v0.47.0 // indirect
)

// Always develop against the sibling framework checkout.
replace github.com/appootb/substratum/v2 => ../

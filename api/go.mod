module github.com/bitcrshr/envmgr/api

go 1.20

require (
	entgo.io/ent v0.12.3
	github.com/bitcrshr/envmgr/proto/compiled/go v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.0-rc.5
	github.com/joho/godotenv v1.5.1
	github.com/lestrrat-go/jwx/v2 v2.0.11
	github.com/lib/pq v1.10.9
	go.uber.org/zap v1.24.0
	google.golang.org/grpc v1.56.2
)

require (
	ariga.io/atlas v0.10.2-0.20230427182402-87a07dfb83bf // indirect
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/go-openapi/inflect v0.19.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/hashicorp/hcl/v2 v2.13.0 // indirect
	github.com/lestrrat-go/blackmagic v1.0.1 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.4 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/zclconf/go-cty v1.8.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace github.com/bitcrshr/envmgr/proto/compiled/go => ../proto/compiled/go

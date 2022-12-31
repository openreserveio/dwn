module github.com/openreserveio/dwn/integration-tests

go 1.19

require (
	github.com/go-resty/resty/v2 v2.7.0
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/google/uuid v1.3.0
	github.com/ipfs/go-cid v0.3.2
	github.com/ipld/go-ipld-prime v0.19.0
	github.com/multiformats/go-multibase v0.1.1
	github.com/multiformats/go-multicodec v0.7.0
	github.com/multiformats/go-multihash v0.2.1
	github.com/onsi/ginkgo/v2 v2.5.1
	github.com/onsi/gomega v1.24.1
	github.com/openreserveio/dwn/go v0.0.0
)

replace github.com/openreserveio/dwn/go v0.0.0 => ../go

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/polydawn/refmt v0.0.0-20201211092308-30ac6d18308e // indirect
	github.com/rs/zerolog v1.28.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.1.7 // indirect
)

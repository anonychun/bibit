module github.com/anonychun/bibit

go 1.27

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/aws/aws-sdk-go-v2 v1.47.0
	github.com/aws/aws-sdk-go-v2/config v1.33.5
	github.com/aws/aws-sdk-go-v2/credentials v1.20.5
	github.com/aws/aws-sdk-go-v2/service/s3 v1.113.1
	github.com/google/uuid v1.6.0
	github.com/gookit/validate/v2 v2.0.2
	github.com/jackc/pgx/v5 v5.11.0
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo/v5 v5.3.1
	github.com/oklog/ulid/v2 v2.1.2
	github.com/pressly/goose/v3 v3.28.0
	github.com/prometheus/client_golang v1.24.1
	github.com/riverqueue/river v0.47.0
	github.com/riverqueue/river/riverdriver/riverdatabasesql v0.47.0
	github.com/riverqueue/river/riverdriver/riverpgxv5 v0.47.0
	github.com/samber/do/v2 v2.1.0
	github.com/stretchr/testify v1.12.1
	github.com/uptrace/bun v1.2.18
	github.com/uptrace/bun/dialect/pgdialect v1.2.18
	github.com/uptrace/bun/extra/bundebug v1.2.18
	github.com/urfave/cli/v3 v3.12.0
	go.opentelemetry.io/otel v1.46.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.46.0
	go.opentelemetry.io/otel/exporters/prometheus v0.68.0
	go.opentelemetry.io/otel/metric v1.46.0
	go.opentelemetry.io/otel/sdk v1.46.0
	go.opentelemetry.io/otel/sdk/metric v1.46.0
	go.opentelemetry.io/otel/trace v1.46.0
	golang.org/x/crypto v0.57.0
	golang.org/x/sync v0.23.0
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.20 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.20.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.20.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.10.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.38.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.43.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.51.0 // indirect
	github.com/aws/smithy-go v1.28.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gookit/filter v1.2.3 // indirect
	github.com/gookit/goutil v0.8.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.30.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/lib/pq v1.12.3 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.3 // indirect
	github.com/prometheus/common v0.71.0 // indirect
	github.com/prometheus/otlptranslator v1.0.0 // indirect
	github.com/prometheus/procfs v0.22.0 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.5.1 // indirect
	github.com/riverqueue/river/riverdriver v0.47.0 // indirect
	github.com/riverqueue/river/rivershared v0.47.0 // indirect
	github.com/riverqueue/river/rivertype v0.47.0 // indirect
	github.com/samber/go-type-to-string v1.8.0 // indirect
	github.com/sethvargo/go-retry v0.4.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/tidwall/gjson v1.19.0 // indirect
	github.com/tidwall/match v1.2.0 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.46.0 // indirect
	go.opentelemetry.io/proto/otlp v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	golang.org/x/time v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260918162117-cecb64721679 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260918162117-cecb64721679 // indirect
)

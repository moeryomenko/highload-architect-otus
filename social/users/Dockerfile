# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.18.
ARG GO_VERSION=1.18

FROM golang:${GO_VERSION}-alpine3.15 AS base

LABEL maintainer="Maxim Eryomenko <moeryomenko@gmail.com>"

RUN apk --no-cache add git gcc musl-dev

# Precompile the entire go standard library into the first Docker cache layer: useful for other projects too!
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go install -v -installsuffix cgo -a std

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user \
    && echo '1000:x:65534:65534::/:' > /user/passwd \
    && echo '1000:x:65534:' > /user/group

FROM base as build

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /src

# Download and precompile all third party libraries, ignoring errors (some have broken tests or whatever).
COPY go.* ./

# Get and precompile third party libraries,
# See issues https://github.com/golang/go/issues/27719.
# Also reusing previous go build cache by `--mount`.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod graph | awk '$1 !~ /@/ { print $2 }' | xargs -r go get -x && \
    go list -m -f '{{ if not .Main }}{{ .Path }}/...@{{ .Version }}{{ end }}' all | tail -n +2 | \
    xargs go build -v -installsuffix cgo -i; echo done

COPY . .

# Compile! Should only compile our sources since everything else is precompiled.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir bin && \
    go build -v -installsuffix cgo -o ./bin -ldflags "-linkmode external -extldflags -static -s -w" .

FROM scratch as final

# Import the user and group files from the build stage.
COPY --from=base /user/group /user/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /src/bin /

EXPOSE 8080
USER 1000
ENTRYPOINT [ "/social" ]

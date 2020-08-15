# Go related commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test ./...
GOGET=$(GOCMD) get -u -v

# Detect the os so that we can build proper statically linked binary
OS := $(shell uname -s | awk '{print tolower($$0)}')

# Get a short hash of the git had for building images.
TAG = $$(git rev-parse --short HEAD)

# Name of actual binary to create
BINARY = app

# GOARCH tells go build which arch. to use while building a statically linked executable
GOARCH = amd64

# Setup the -ldflags option for go build here.
# While statically linking we want to inject version related information into the binary
LDFLAGS = -ldflags="$$(govvv -flags)"

.PHONY: run
run: bin #this will cause "bin" target to be build first
	./$(BINARY)-$(OS)-$(GOARCH) # Execute the binary

# bin creates a platform specific statically linked binary. Platform sepcific because if you are on
# OS-X; linux binary will not work.
.PHONY: bin
bin:
	env CGO_ENABLED=0 GOOS=$(OS) GOARCH=${GOARCH} go build -a -installsuffix cgo ${LDFLAGS} -o ${BINARY}-$(OS)-${GOARCH} . ;

# Docker build internally (within Dockerfile) triggers "make bin", which creates a "linux" binary.
.PHONY: docker
docker:
	docker build -t devstackq/$(BINARY):$(GOARCH)-$(TAG) .

# Push pushes the image to the docker repository.
.PHONY: push
push: docker
	docker push devstackq/$(BINARY):$(GOARCH)-$(TAG)

# Runs unit tests.
.PHONY: test
test:
	$(GOTEST)

# Generates a coverage report
.PHONY: cover
cover:
	${GOCMD} test -coverprofile=coverage.out ./... && ${GOCMD} tool cover -html=coverage.out

# Remove coverage report and the binary.
.SILENT: clean
.PHONY: clean
clean:
	$(GOCLEAN)
	@rm -f ${BINARY}-$(OS)-${GOARCH}
	@rm -f coverage.out

# There are much better ways to manage deps in golang, I'm going go get just for brevity
.PHONY: deps
deps:
	$(GOGET) github.com/devstackq/Forum-X/models
	$(GOGET) github.com/devstackq/Forum-X/routing


	

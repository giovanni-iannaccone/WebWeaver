GO = go
GOFLAGS = -ldflags="-s -w"

TARGET = WebWeaver

DARWIN_BIN = $(TARGET)
LINUX_BIN = $(TARGET)
WINDOWS_BIN = $(TARGET).exe

all: $(DARWIN_BIN) $(LINUX_BIN) $(WINDOWS_BIN)

$(DARWIN_BIN): ./cmd/main.go
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $@ ./cmd/main.go

$(LINUX_BIN): ./cmd/main.go
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $@ ./cmd/main.go

$(WINDOWS_BIN): ./cmd/main.go
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $@ ./cmd/main.go

clean:
	rm -f $(DARWIN_BIN) $(LINUX_BIN) $(WINDOWS_BIN)

run: $(TARGET)
	./$(TARGET)

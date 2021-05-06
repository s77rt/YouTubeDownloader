PROJECT_NAME = YouTubeDownloader

CMD_DIR = ./cmd/

BIN_DIR = ./bin/

BIN_DIR_MAC = ./bin/mac/
BIN_DIR_LINUX = ./bin/linux/
BIN_DIR_WINDOWS = ./bin/windows/

GIT_TAG = "$(shell git describe --tags)"

LD_FLAGS_MAC = "-X 'github.com/s77rt/YouTubeDownloader.Version=$(GIT_TAG)'"
LD_FLAGS_LINUX = "-X 'github.com/s77rt/YouTubeDownloader.Version=$(GIT_TAG)'"
LD_FLAGS_WINDOWS = "-X 'github.com/s77rt/YouTubeDownloader.Version=$(GIT_TAG)' -H windowsgui"

all: clean dep compile

clean:
	@echo -n "Cleaning: "
	@rm -rf $(BIN_DIR)
	@echo "[OK]"

dep:
	@echo -n "Downloading Dependencies: "
	@go get -d ./...
	@echo "[OK]"

compile:
	@echo "Compiling: "
	
	# macOS (64bit)
	@mkdir -p $(BIN_DIR_MAC)
	GOOS=darwin GOARCH=amd64 go build -ldflags $(LD_FLAGS_MAC) -o $(BIN_DIR_MAC)$(PROJECT_NAME) $(CMD_DIR)$(PROJECT_NAME)
	@zip -j $(BIN_DIR)$(PROJECT_NAME)_$(GIT_TAG)_mac.zip $(BIN_DIR_MAC)$(PROJECT_NAME)

	# Linux (64bit)
	@mkdir -p $(BIN_DIR_LINUX)
	GOOS=linux GOARCH=amd64 go build -ldflags $(LD_FLAGS_LINUX) -o $(BIN_DIR_LINUX)$(PROJECT_NAME) $(CMD_DIR)$(PROJECT_NAME)
	@zip -j $(BIN_DIR)$(PROJECT_NAME)_$(GIT_TAG)_linux.zip $(BIN_DIR_LINUX)$(PROJECT_NAME)

	# Windows (64bit)
	@mkdir -p $(BIN_DIR_WINDOWS)
	rsrc -arch amd64 -ico ico.ico -o $(CMD_DIR)$(PROJECT_NAME)/rsrc_windows_amd64.syso
	GOOS=windows GOARCH=amd64 go build -ldflags $(LD_FLAGS_WINDOWS) -o $(BIN_DIR_WINDOWS)$(PROJECT_NAME).exe $(CMD_DIR)$(PROJECT_NAME)
	@zip -j $(BIN_DIR)$(PROJECT_NAME)_$(GIT_TAG)_windows.zip $(BIN_DIR_WINDOWS)$(PROJECT_NAME).exe

	@echo "[OK]"

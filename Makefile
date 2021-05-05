PROJECT_NAME = YouTubeDownloader

CMD_DIR = ./cmd/

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
	@rm -rf $(BIN_DIR_MAC)
	@rm -rf $(BIN_DIR_LINUX)
	@rm -rf $(BIN_DIR_WINDOWS)
	@echo "[OK]"

dep:
	@echo -n "Downloading Dependencies: "
	@go get -d ./...
	@echo "[OK]"

compile:
	@echo "Compiling: "
	
	# MacOS (64bit)
	@mkdir -p $(BIN_DIR_MAC)
	GOOS=darwin GOARCH=amd64 go build -ldflags $(LD_FLAGS_MAC) -o $(BIN_DIR_MAC)$(PROJECT_NAME) $(CMD_DIR)$(PROJECT_NAME)

	# Linux (64bit)
	@mkdir -p $(BIN_DIR_LINUX)
	GOOS=linux GOARCH=amd64 go build -ldflags $(LD_FLAGS_LINUX) -o $(BIN_DIR_LINUX)$(PROJECT_NAME) $(CMD_DIR)$(PROJECT_NAME)

	# Windows (64bit)
	@mkdir -p $(BIN_DIR_WINDOWS)
	rsrc -arch amd64 -ico ico.ico -o $(CMD_DIR)$(PROJECT_NAME)/rsrc_windows_amd64.syso
	GOOS=windows GOARCH=amd64 go build -ldflags $(LD_FLAGS_WINDOWS) -o $(BIN_DIR_WINDOWS)$(PROJECT_NAME).exe $(CMD_DIR)$(PROJECT_NAME)

	@echo "[OK]"

#!/bin/bash
set -e

export DEBUG=0
export GO_BINARY_DIRECTORY=/usr/local/go/bin
export GOPROXY=https://proxy.golang.org,direct
export WORKING_DIRECTORY="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export GOROOT="/usr/local/go"
export PATH="$GO_BINARY_DIRECTORY:$PATH"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

fail() {
    echo "ERROR: $1" >&2
    exit 1
}

remove_directory() {
    local TARGET="$1"
    log "Remove directory: $TARGET"
    if [ ! -d "$TARGET" ]; then
        return 0
    fi
    rm -rf "$TARGET" 2>/dev/null || true
}

copy_file() {
    local SRC="$1"
    local DST="$2"
    if [ -f "$SRC" ]; then
        log "Copy file: $SRC to $DST"
        cp -f "$SRC" "$DST"
    fi
}

copy_icon_if_not_exists() {
    local ICON_PNG="$1"
    local ICON_PNG_IN="$2"
    if [ ! -f "$ICON_PNG" ]; then
        log "Copy icon: $ICON_PNG_IN to $ICON_PNG"
        cp -f "$ICON_PNG_IN" "$ICON_PNG"
    fi
}

create_desktop_file() {
    local APP_NAME="$1"
    local ICON_NAME="$2"
    local BIN_DIR="$3"
    local DESKTOP_FILE="$BIN_DIR/$APP_NAME.desktop"
    log "Creating desktop file: $DESKTOP_FILE"
    cat > "$DESKTOP_FILE" <<EOF
[Desktop Entry]
Type=Application
Name=$APP_NAME
Exec=$BIN_DIR/$APP_NAME
Icon=$BIN_DIR/$ICON_NAME
StartupNotify=false
Terminal=false
MimeType=
X-AppImage-Integrate=false
EOF
}

build_go_project() {
    if [ "$DEBUG" = "1" ]; then
        set -x
    fi
    local APP_NAME="$1"
    local SOURCE_DIRECTORY="$2"
    local OUTPUT_DIRECTORY="$SOURCE_DIRECTORY/bin/amd64"
    log "Starting $APP_NAME Go build"
    log "SOURCE_DIRECTORY: $SOURCE_DIRECTORY"
    remove_directory "$SOURCE_DIRECTORY/bin"
    mkdir -p "$OUTPUT_DIRECTORY"
    pushd "$SOURCE_DIRECTORY" || exit 1
    log "Running go mod tidy in $SOURCE_DIRECTORY"
    go mod tidy || (popd && exit 1)
    log "Running go build for $APP_NAME to $OUTPUT_DIRECTORY/$APP_NAME"
    go build -o "$OUTPUT_DIRECTORY/$APP_NAME" . || (popd && exit 1)
    popd
    set +x
}

log "Starting go-boost CA tools build"
log "WORKING_DIRECTORY: $WORKING_DIRECTORY"

cd "$WORKING_DIRECTORY" || fail "Failed to cd to $WORKING_DIRECTORY"

# Check and copy icon files before build
copy_icon_if_not_exists "$WORKING_DIRECTORY/server/winres/app.png" "$WORKING_DIRECTORY/server/winres/app.png.in"
copy_icon_if_not_exists "$WORKING_DIRECTORY/issueIPAddress/winres/app.png" "$WORKING_DIRECTORY/issueIPAddress/winres/app.png.in"

# Build all CA projects
build_go_project "server" "$WORKING_DIRECTORY/server"
build_go_project "issueIPAddress" "$WORKING_DIRECTORY/issueIPAddress"

# Copy all binaries to bin/amd64
BIN_DIR="$WORKING_DIRECTORY/bin/amd64"
remove_directory "$BIN_DIR"
mkdir -p "$BIN_DIR"

# Copy compiled binaries
cp -f "$WORKING_DIRECTORY/server/bin/amd64/server" "$BIN_DIR/" 2>/dev/null || true
cp -f "$WORKING_DIRECTORY/issueIPAddress/bin/amd64/issueIPAddress" "$BIN_DIR/" 2>/dev/null || true

# Copy san.json config
copy_file "$WORKING_DIRECTORY/server/san.json" "$BIN_DIR/san.json"

# Copy icon files to bin directory (using app.png name, may overwrite)
copy_file "$WORKING_DIRECTORY/server/winres/app.png" "$BIN_DIR/app.png"
copy_file "$WORKING_DIRECTORY/issueIPAddress/winres/app.png" "$BIN_DIR/app.png"

# Create desktop files
create_desktop_file "server" "app.png" "$BIN_DIR"
create_desktop_file "issueIPAddress" "app.png" "$BIN_DIR"

log "go-boost CA tools build completed successfully"
log "Output directory: $BIN_DIR"
exit 0

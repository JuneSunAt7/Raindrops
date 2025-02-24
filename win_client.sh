#!/bin/sh


FILES=("key_gens/generate_ssl.go" "cmd/client/client.go" "cmd/server.go" "cmd/data_mgmt/server_mgmt.go" "cmd/data_mgmt/client_mgmt.go" "cmd/data_mgmt/data_mgmt.go" "cmd/intelligence/run_cli_intell.go")


PLATFORMS=("windows" "linux")


for platform in "${PLATFORMS[@]}"; do
    echo "Building for $platform..."

    if [ "$platform" = "windows" ]; then
        GOOS=windows
        EXT=".exe"
    elif [ "$platform" = "linux" ]; then
        GOOS=linux
        EXT=""
    else
        echo "Unsupported platform: $platform"
        continue
    fi

    GOARCH=amd64

    for file in "${FILES[@]}"; do

        filename=$(basename "$file" .go)

        go build -o "build/$platform/${filename}${EXT}" -ldflags="-s -w" "$file"
        echo "Built: build/$platform/${filename}${EXT}"
    done
done

echo "Press enter to continue"
read name
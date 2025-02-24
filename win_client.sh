#!/bin/sh

# Определяем массив с исходными файлами
FILES=("key_gens/generate_ssl.go" "cmd/client/client.go" "cmd/server.go" "cmd/data_mgmt/server_mgmt.go" "cmd/data_mgmt/client_mgmt.go" "cmd/data_mgmt/data_mgmt.go" "cmd/intelligence/run_cli_intell.go")

# Определяем массив с целевыми операционными системами
PLATFORMS=("windows" "linux")

# Цикл по платформам
for platform in "${PLATFORMS[@]}"; do
    echo "Building for $platform..."

    # Устанавливаем GOOS в зависимости от платформы
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

    # Цикл по файлам
    for file in "${FILES[@]}"; do
        # Извлекаем имя файла без пути и расширения
        filename=$(basename "$file" .go)
        
        # Собираем приложение
        go build -o "build/$platform/${filename}${EXT}" -ldflags="-s -w" "$file"
        echo "Built: build/$platform/${filename}${EXT}"
    done
done

echo "Press enter to continue"
read name
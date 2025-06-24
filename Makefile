# Confericis Backend Makefile

# Variables
BINARY_NAME=muac-api
MAIN_FILE=./cmd/main.go
GO_CMD=go
BUILD_FLAGS=-ldflags="-s -w"
EC2_KEY=gopher.pem
EC2_USER=ec2-user
EC2_HOST=ec2-35-173-114-173.compute-1.amazonaws.com
EC2_PATH=/home/ec2-user/dev/

# Build para diferentes arquitecturas
.PHONY: build-linux build-windows build-arm64 clean deploy run test help

# Comando principal - compilar para Amazon Linux 2023 (x86_64)
build-linux:
	@echo "Compilando para Amazon Linux 2023..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO_CMD) build $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Compilación completada: $(BINARY_NAME)"

# Compilar para ARM64 (Graviton instances)
build-arm64:
	@echo "Compilando para ARM64 Linux..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO_CMD) build $(BUILD_FLAGS) -o $(BINARY_NAME)-arm64 $(MAIN_FILE)
	@echo "Compilación completada: $(BINARY_NAME)-arm64"

# Limpiar binarios
clean:
	@echo "Limpiando binarios..."
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe $(BINARY_NAME)-arm64
	@echo "Limpieza completada"

# Subir binario a EC2
upload:
	@echo "Subiendo binario a EC2..."
	scp -i "$(EC2_KEY)" $(BINARY_NAME) $(EC2_USER)@$(EC2_HOST):$(EC2_PATH)
	@echo "Binario subido exitosamente"

# Subir archivo específico - Uso: make upload-file FILE=archivo.txt
upfile:
	@if [ -z "$(FILE)" ]; then \
		echo "Error: Debes especificar el archivo con FILE=nombre_archivo"; \
		echo "Ejemplo: make upload-file FILE=ecosystem.config.cjs"; \
		exit 1; \
	fi
	@if [ ! -f "$(FILE)" ]; then \
		echo "Error: El archivo $(FILE) no existe"; \
		exit 1; \
	fi
	@echo "Subiendo archivo $(FILE) a EC2..."
	scp -i "$(EC2_KEY)" "$(FILE)" $(EC2_USER)@$(EC2_HOST):$(EC2_PATH)
	@echo "Archivo $(FILE) subido exitosamente a $(EC2_PATH)"


delete-binary:
	@echo "Eliminando binario compilado..."
	rm $(BINARY_NAME)
	@echo "Binario subido exitosamente"

# Deployment completo (compilar + subir)
deploy: build-linux upload
	@echo "Deployment completado"

# Ejecutar localmente
run:
	@echo "Ejecutando aplicación localmente..."
	$(GO_CMD) run $(MAIN_FILE)

# Verificar dependencias
deps:
	@echo "Descargando dependencias..."
	$(GO_CMD) mod download
	$(GO_CMD) mod tidy

# Desarrollo: watch y rebuild automático (requiere air)
dev:
	@echo "Iniciando modo desarrollo..."
	air

# Verificar info del proyecto
info:
	@echo "=== Información del Proyecto ==="
	@echo "Go version: $$($(GO_CMD) version)"
	@echo "GOOS: $$($(GO_CMD) env GOOS)"
	@echo "GOARCH: $$($(GO_CMD) env GOARCH)"
	@echo "Binario objetivo: $(BINARY_NAME)"
	@echo "Servidor EC2: $(EC2_HOST)"

# Conectar a EC2 (solo SSH)
ssh:
	@echo "Conectando a EC2..."
	ssh -i "$(EC2_KEY)" $(EC2_USER)@$(EC2_HOST)

# Ver logs remotos (ajustar según tu aplicación)
logs:
	@echo "Viendo logs remotos..."
	ssh -i "$(EC2_KEY)" $(EC2_USER)@$(EC2_HOST) "tail -f /var/log/confericis.log"

# Restart remoto (si usas systemd)
restart:
	@echo "Reiniciando servicio remoto..."
	ssh -i "$(EC2_KEY)" $(EC2_USER)@$(EC2_HOST) "sudo systemctl restart confericis"

# Ayuda
help:
	@echo "=== Comandos disponibles ==="
	@echo "build-linux    - Compilar para Amazon Linux 2023"
	@echo "build-windows  - Compilar para Windows"
	@echo "build-arm64    - Compilar para ARM64 Linux"
	@echo "deploy         - Compilar y subir a EC2"
	@echo "upload         - Subir binario existente a EC2"
	@echo "run-remote     - Ejecutar en EC2"
	@echo "clean          - Limpiar binarios"
	@echo "run            - Ejecutar localmente"
	@echo "deps           - Actualizar dependencias"
	@echo "ssh            - Conectar a EC2"
	@echo "info           - Ver información del proyecto"
	@echo "help           - Mostrar esta ayuda"

# Comando por defecto
.DEFAULT_GOAL := build-linux

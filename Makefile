# Confericis Backend Makefile - Configuración para Ubuntu local y EC2

# Variables
BINARY_NAME=muac-api
MAIN_FILE=./cmd/main.go
GO_CMD=go
BUILD_FLAGS=-ldflags="-s -w"

# Configuración para EC2
EC2_KEY=gopher.pem
EC2_USER=ec2-user
EC2_HOST=ec2-35-173-114-173.compute-1.amazonaws.com
EC2_PATH=/home/ec2-user/dev/

# Configuración para Ubuntu local
LOCAL_USER=nutricion
LOCAL_HOST=192.168.254.35
LOCAL_PATH=/home/nutricion/dev/

# Build para diferentes sistemas
.PHONY: build build-ubuntu build-linux build-arm64 clean deploy-local deploy-ec2 run test help

# Compilación por defecto (Ubuntu local)
build: build-ubuntu

# Compilar para Ubuntu local
build-ubuntu:
	@echo "Compilando para Ubuntu local..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO_CMD) build $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Compilación completada: $(BINARY_NAME)"

# Compilar para Amazon Linux (EC2)
build-linux:
	@echo "Compilando para Amazon Linux (EC2)..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO_CMD) build $(BUILD_FLAGS) -o $(BINARY_NAME)-amazon $(MAIN_FILE)
	@echo "Compilación completada: $(BINARY_NAME)-amazon"

# Compilar para ARM64
build-arm64:
	@echo "Compilando para ARM64..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO_CMD) build $(BUILD_FLAGS) -o $(BINARY_NAME)-arm64 $(MAIN_FILE)
	@echo "Compilación completada: $(BINARY_NAME)-arm64"

# Limpiar binarios
clean:
	@echo "Limpiando binarios..."
	rm -f $(BINARY_NAME)* $(BINARY_NAME).exe
	@echo "Limpieza completada"

# Despliegue local (Ubuntu)
deploy-local: build-ubuntu
	@echo "Moviendo binario a directorio de producción local..."
	# scp $(BINARY_NAME) $(LOCAL_USER)@$(LOCAL_HOST):$(LOCAL_PATH)
	# muestrame como quedaria ese comando aqui usando los valores no las Variables
	@echo "scp muac-api nutricion@192.168.254.35:/home/nutricion/dev/"
	@echo "------------------------------------------------------------------------"

# Despliegue en EC2
deploy-ec2: build-linux
	@echo "Subiendo binario a EC2..."
	scp -i "$(EC2_KEY)" $(BINARY_NAME)-amazon $(EC2_USER)@$(EC2_HOST):$(EC2_PATH)$(BINARY_NAME)
	@echo "Despliegue en EC2 completado"

# Subir archivo específico a servidor local
upload-local:
	@if [ -z "$(FILE)" ]; then \
		echo "Error: Especifica el archivo con FILE=nombre_archivo"; \
		exit 1; \
	fi
	@echo "Subiendo $(FILE) a servidor local..."
	scp "$(FILE)" $(LOCAL_USER)@$(LOCAL_HOST):$(LOCAL_PATH)
	@echo "Archivo $(FILE) subido a $(LOCAL_PATH)"

# Ejecutar localmente
run:
	@echo "Ejecutando aplicación localmente..."
	$(GO_CMD) run $(MAIN_FILE)

# Verificar dependencias
deps:
	@echo "Actualizando dependencias..."
	$(GO_CMD) mod tidy
	$(GO_CMD) mod download

ssh:
	@echo "Conectando a EC2..."
	ssh -i "$(EC2_KEY)" $(EC2_USER)@$(EC2_HOST)

# Ayuda actualizada
help:
	@echo "=== Comandos disponibles ==="
	@echo "build         - Compilar para Ubuntu local (por defecto)"
	@echo "build-ubuntu  - Compilar para Ubuntu local"
	@echo "build-linux   - Compilar para Amazon Linux (EC2)"
	@echo "build-arm64   - Compilar para ARM64"
	@echo "deploy-local  - Desplegar en servidor Ubuntu local"
	@echo "deploy-ec2    - Desplegar en EC2"
	@echo "upload-local  - Subir archivo al servidor local (usar FILE=archivo)"
	@echo "clean         - Limpiar binarios"
	@echo "run           - Ejecutar localmente"
	@echo "deps          - Actualizar dependencias"
	@echo "help          - Mostrar esta ayuda"

# Comando por defecto
.DEFAULT_GOAL := help
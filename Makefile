# Lista de proyectos
PROJECTS=cpu memoria io kernel
BIN_DIR=bin

# Target por defecto
all: build

# Compilar todos los binarios en la carpeta bin/
build: $(PROJECTS)

$(PROJECTS):
	@echo "Compilando $@..."
	@mkdir -p $(BIN_DIR)
	@cd $@ && go build -o ../$(BIN_DIR)/$@

# Limpiar binarios
clean:
	@echo "Limpiando binarios..."
	@rm -rf $(BIN_DIR)

.PHONY: all build clean $(PROJECTS)

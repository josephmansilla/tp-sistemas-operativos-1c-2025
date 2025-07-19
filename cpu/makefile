RUN = go run

uno:
	nano ../cpu/configs/1EG.json
	$(RUN) cpu.go 1

dos:
	nano ../cpu/configs/2EG.json
	$(RUN) cpu.go 2

tres:
	nano ../cpu/configs/3EG.json
	$(RUN) cpu.go 3

cuatro:
	nano ../cpu/configs/4EG.json
	$(RUN) cpu.go 4

# === Prueba: Corto Plazo ===
corto:
	@echo "Comienza la prueba de Corto Plazo"
	$(RUN) cpu.go 1
	$(RUN) cpu.go 2

# === Prueba: Mediano/Largo Plazo ===
lym:
	@echo "Comienza la prueba de Mediano y Largo Plazo"
	$(RUN) cpu.go 1 PLANI

# === Prueba: SWAP ===
swap:
	@echo "Comienza la prueba de SWAP"
	$(RUN) cpu.go 1 SWAP

# === Prueba: CACHE ===
cache:
	@echo "Comienza la prueba de Cache"
	$(RUN) cpu.go 1 CACHE

# === Prueba: Estabilidad General (EG) ===
eg:
	@echo "Comienza la prueba de Estabilidad General"
	$(RUN) cpu.go 1 1EG &
	$(RUN) cpu.go 2 2EG &
	$(RUN) cpu.go 3 3EG &
	$(RUN) cpu.go 4 4EG

# === Detener todos los procesos (por si quedan en background) ===
clean:
	@echo "Terminando módulos..."
	pkill -f cpu.go || true
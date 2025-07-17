RUN = go run

uno:
	nano ../cpu/configs/config.json
	$(RUN) cpu.go 1

dos:
	nano ../cpu/configs/config.json
	$(RUN) cpu.go 2

tres:
	nano ../cpu/configs/config.json
	$(RUN) cpu.go 3

cuatro:
	nano ../cpu/configs/config.json
	$(RUN) cpu.go 4

# === Prueba: Corto Plazo ===
corto:
	$(RUN) cpu.go 1 &
	$(RUN) cpu.go 2

# === Prueba: Mediano/Largo Plazo ===
lym:
	$(RUN) cpu.go 1 PLANI

# === Prueba: SWAP ===
swap:
	$(RUN) cpu.go 1 SWAP

# === Prueba: CACHE ===
cache:
	$(RUN) cpu.go 1 CACHE

# === Prueba: Estabilidad General (EG) ===
eg:
	$(RUN) cpu.go 1 1EG &
	$(RUN) cpu.go 2 2EG &
	$(RUN) cpu.go 3 3EG &
	$(RUN) cpu.go 4 4EG

# === Detener todos los procesos (por si quedan en background) ===
clean:
	@echo "Terminando módulos..."
	pkill -f cpu.go || true
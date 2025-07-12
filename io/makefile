RUN = go run

# === Prueba: Corto Plazo ===
corto:
	$(RUN) io.go DISCO

# === Prueba: Mediano/Largo Plazo ===
lym:
	$(RUN) io.go DISCO

# === Prueba: SWAP ===
swap:
	$(RUN) io.go DISCO

# === Prueba: CACHE ===
cache:
	$(RUN) io.go DISCO

# === Prueba: Estabilidad General (EG) ===
eg:
	$(RUN) io.go DISCO &
	$(RUN) io.go DISCO2 &
	$(RUN) io.go DISCO3 &
	$(RUN) io.go DISCO4

# === Detener todos los procesos (por si quedan en background) ===
clean:
	@echo "Terminando módulo..."
	pkill -f io.go || true

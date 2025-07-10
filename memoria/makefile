RUN = go run

# === Prueba: Corto Plazo ===
corto:
	$(RUN) memoria.go cortoPlazo

# === Prueba: Mediano/Largo Plazo ===
lym:
	$(RUN) memoria.go medianoLargoPlazo

# === Prueba: SWAP ===
swap:
	$(RUN) memoria.go memoriaSwap

# === Prueba: CACHE ===
cache:
	$(RUN) memoria.go memoriaCache

# === Prueba: Estabilidad General (EG) ===
eg:
	$(RUN) memoria.go estabilidadGeneral

# === Detener todos los procesos (por si quedan en background) ===
clean:
	@echo "Terminando módulo..."
	pkill -f memoria.go || true
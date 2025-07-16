RUN = go run

# === Prueba: Corto Plazo ===
corto:
	$(RUN) kernel.go PLANI_CORTO_PLAZO 0 cortoplazo

# === Prueba: Mediano/Largo Plazo ===
lym:
	$(RUN) kernel.go PLANI_LYM_IO 0

# === Prueba: SWAP ===
swap:
	$(RUN) kernel.go MEMORIA_IO 90

# === Prueba: CACHE ===
cache:
	$(RUN) kernel.go MEMORIA_BASE 256

# === Prueba: Estabilidad General (EG) ===
eg:
	$(RUN) kernel.go ESTABILIDAD_GENERAL 0

# === Detener todos los procesos (por si quedan en background) ===
clean:
	@echo "Terminando módulo..."
	pkill -f kernel.go || true

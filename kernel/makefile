RUN = go run

# === Prueba: Corto Plazo ===
corto:
	nano kernel/configs/config.json
	$(RUN) kernel.go PLANI_CORTO_PLAZO 0 cortoplazo

# === Prueba: Mediano/Largo Plazo ===
lym:
	nano ../kernel/configs/config.json
	$(RUN) kernel.go PLANI_LYM_PLAZO 0 medianoplazo

# === Prueba: SWAP ===
swap:
	nano kernel/configs/config.json
	$(RUN) kernel.go MEMORIA_IO 90 swap

# === Prueba: CACHE ===
cache:
	nano kernel/configs/config.json
	$(RUN) kernel.go MEMORIA_BASE 256 cache

# === Prueba: Estabilidad General (EG) ===
eg:
	nano kernel/configs/config.json
	$(RUN) kernel.go ESTABILIDAD_GENERAL 0

# === Detener todos los procesos (por si quedan en background) ===
clean:
	nano kernel/configs/config.json
	@echo "Terminando módulo..."
	pkill -f kernel.go || true


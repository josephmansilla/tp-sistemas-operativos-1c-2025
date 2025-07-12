RUN = go run

# === Prueba: Corto Plazo ===
corto:
	cd memoria && $(RUN) memoria.go cortoPlazo &
	sleep 4
	cd kernel && $(RUN) kernel.go PLANI_CORTO_PLAZO 0 cortoplazo &
	sleep 4
	cd cpu && $(RUN) cpu.go 1 1CP &
	cd cpu && $(RUN) cpu.go 2 2CP &
	sleep 1
	cd io && $(RUN) io.go DISCO &

# === Prueba: Mediano/Largo Plazo ===
lym:
	cd memoria && $(RUN) memoria.go medianoLargoPlazo &
	sleep 4
	cd kernel && $(RUN) kernel.go PLANI_LYM_IO 0 medianoplazo &
	sleep 4
	cd cpu && $(RUN) cpu.go 1 PLANI &
	sleep 1
	cd io && $(RUN) io.go DISCO &

# === Prueba: SWAP ===
swap:
	cd memoria && $(RUN) memoria.go memoriaSwap &
	sleep 4
	cd kernel && $(RUN) kernel.go MEMORIA_IO 90 swap &
	sleep 4
	cd cpu && $(RUN) cpu.go 1 SWAP &
	sleep 1
	cd io && $(RUN) io.go DISCO &

# === Prueba: CACHE ===
cache:
	cd memoria && $(RUN) memoria.go memoriaCache &
	sleep 4
	cd kernel && $(RUN) kernel.go MEMORIA_BASE 256 cache &
	sleep 4
	cd cpu && $(RUN) cpu.go 1 CACHE &
	sleep 1
	cd io && $(RUN) io.go DISCO &

# === Prueba: Estabilidad General (EG) ===
eg:
	cd memoria && $(RUN) memoria.go estabilidadGeneral &
	sleep 4
	cd kernel && $(RUN) kernel.go ESTABILIDAD_GENERAL 0 estabilidad &
	sleep 4
	cd cpu && $(RUN) cpu.go 1 1EG &
	cd cpu && $(RUN) cpu.go 2 2EG &
	cd cpu && $(RUN) cpu.go 3 3EG &
	cd cpu && $(RUN) cpu.go 4 4EG &
	sleep 1
	cd io && $(RUN) io.go DISCO &
	cd io && $(RUN) io.go DISCO2 &
	cd io && $(RUN) io.go DISCO3 &
	cd io && $(RUN) io.go DISCO4 &

# === Detener todos los procesos (por si quedan en background) ===
clean:
	@echo "Terminando módulos..."
	pkill -f memoria.go || true
	pkill -f kernel.go || true
	pkill -f cpu.go || true
	pkill -f io.go || true

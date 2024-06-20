# Para ejecutar en entorno local ejecutar el siguiente comando

```
go run main.go
```

### La ruta en el entorno local es realizando una petición get

```
http://localhost:8080/questions
```

### La información se parseo como se pidio en el documento de prueba


# Comando para ejecutar desde cloud run
```
gcloud run jobs deploy job-quickstart \
    --source . \
    --tasks 50 \
    --set-env-vars SLEEP_MS=10000 \
    --set-env-vars FAIL_RATE=0.1 \
    --max-retries 5 \
    --region REGION \
    --project=PROJECT_ID
```

# Ejecuta un trabajo en Cloud Run

```
gcloud run jobs execute job-quickstart --region REGION
```




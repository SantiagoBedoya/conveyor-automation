# Automatización Final

Sistema IoT de monitoreo y control para banda transportadora con sensor de gas, humedad, distancia, conteo de objetos y dashboards en tiempo real.

## Arquitectura

```
┌──────────────┐     HTTPS / JSON     ┌──────────────────┐     WebSocket     ┌────────────┐
│   ESP32      │ ──────────────────► │   Go API + WS    │ ────────────────► │ Dashboard  │
│ (firmware)   │ ◄────────────────── │  (conveyor-api)  │                   │ (HTML/JS)  │
└──────────────┘                     └────────┬─────────┘                   └────────────┘
                                              │
                                              ▼
                                       ┌──────────┐
                                       │ PostgreSQL│
                                       └──────────┘
```

## Componentes

### 1. Firmware ESP32 (`automatizacion_final.ino`)

Controla sensores y actuadores con WiFi y reporte a la API vía HTTPS.

| Componente | Pin ESP32 |
|------------|-----------|
| Gas MQ-02 | GPIO 34 |
| Humedad | GPIO 39 |
| Ultrasonido TRIG | GPIO 23 |
| Ultrasonido ECHO | GPIO 19 |
| Buzzer | GPIO 27 |
| Relay (motor banda) | GPIO 16 |
| Ventilador | GPIO 17 |
| Servo (puerta) | GPIO 18 |

**Comportamiento:**
- Reporta lecturas cada 300ms vía POST a `/api/readings`
- Gas > 500 → buzzer ON, ventilador ON
- Humedad < 3000 → buzzer ON, motor OFF, puerta OPEN
- Distancia < 10cm → incrementa contador de objetos
- En estado normal: buzzer OFF, motor ON, ventilador OFF, puerta CLOSED

### 2. Variante sin API (`no-api/`)

Firmware independiente que opera offline con la misma lógica de sensores y actuadores, sin WiFi ni reportes HTTP.

### 3. API Go (`conveyor-api/`)

Backend REST + WebSocket construido con:
- **Router:** chi
- **DB:** PostgreSQL + pgx + sqlc
- **WebSocket:** gorilla/websocket
- **Migraciones:** golang-migrate

**Endpoints:**

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/api/readings` | Crear lectura |
| `GET` | `/api/readings` | Listar (`?limit=50&offset=0`) |
| `GET` | `/api/readings/{id}` | Obtener por ID |
| `GET` | `/api/alerts` | Listar alertas |
| `GET` | `/api/alerts/active` | Alertas activas |
| `POST` | `/api/alerts/{id}/resolve` | Resolver alerta |
| `GET` | `/api/status` | Estado actual del sistema |
| `GET` | `/api/ws` | WebSocket |

### 4. Dashboard (`conveyor-api/dashboard/`)

SPA en vanilla HTML/CSS/JS + Chart.js con:
- Tarjetas de estado (banda, ventilador, buzzer, puerta)
- Lecturas en vivo (gas, humedad, distancia, objetos)
- Gráfico de líneas con toggle por sensor
- Lista de alertas activas con botón de resolver
- Tabla histórica de lecturas

## Setup rápido

```bash
# API local
cd conveyor-api
cp .env-example .env   # editar DATABASE_URL
createdb conveyor
migrate -path migrations -database "$DATABASE_URL" up
go run ./cmd/api

# Con Docker
cd conveyor-api
docker compose up --build
```

Acceder al dashboard en `http://localhost:8080`.

## Despliegue

El proyecto incluye `railway.toml` para desplegar en Railway.app. Variables de entorno requeridas: `DATABASE_URL`, `API_KEY`, `PORT`.

## Variables de entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `DATABASE_URL` | Conexión a PostgreSQL | — |
| `API_KEY` | API key para ESP32 | — |
| `PORT` | Puerto del servidor | `8080` |

## Licencia

MIT

# Stock API - Backend

## 📌 Descripción del Proyecto

Este proyecto es una **API RESTful** desarrollada en **Golang**, que permite consultar información de acciones financieras, almacenarlas en **CockroachDB**, y recomendar las mejores acciones para invertir con base en un algoritmo de análisis.

El sistema obtiene los datos desde una API externa, los almacena en la base de datos y expone endpoints para acceder a esta información. Además, proporciona una funcionalidad de recomendación basada en la evaluación de cambios en la calificación de las acciones y en la variación de sus precios objetivo.

---

## 🏗️ Arquitectura del Proyecto

El proyecto sigue una **arquitectura en capas**, dividiendo la responsabilidad en módulos separados:

1. **Base de Datos (**``**)**:

   - Conexión con CockroachDB
   - Creación de la base de datos y la tabla `stocks`
   - Inserción de datos obtenidos desde la API externa

2. **Modelos (**``**)**:

   - Define las estructuras de datos (`Stock`, `StockScore`, `APIResponse`)

3. **Servicios (**``**)**:

   - Contiene la lógica de negocio, incluyendo la puntuación y recomendación de acciones

4. **Rutas (**``**)**:

   - Define los endpoints de la API y maneja las solicitudes HTTP

5. **Main (**``**)**:

   - Punto de entrada del sistema
   - Inicializa la base de datos, servicios y API

Esta arquitectura facilita la escalabilidad y mantenimiento del código.

---

## 🔢 Lógica del Algoritmo de Recomendación

### 📊 **¿Qué evalúa el algoritmo?**

En el mercado de acciones, los inversionistas analizan varios factores antes de decidir dónde invertir. Nuestra API evalúa dos criterios principales:

1. **Cambio en la Calificación de la Acción**

   - Las acciones son calificadas por firmas de inversión con términos como `Buy`, `Neutral`, `Sell`, `Underweight`, `Overweight`.
   - Un cambio de `Sell` → `Buy` indica una mejora fuerte y recibe mayor puntuación.
   - Un cambio de `Neutral` → `Buy` es positivo, pero con menos impacto.
   - Una degradación (`Buy` → `Sell`) es penalizada con puntuaciones negativas.

2. **Cambio en el Target Price (Precio Objetivo)**

   - Indica la expectativa de crecimiento de la acción.
   - Si una acción tiene un `target_to` mayor que `target_from`, se considera positiva.
   - Se calcula el **porcentaje de crecimiento** en lugar de un valor absoluto, para evaluar proporcionalmente.

### 📈 **Fórmula de Puntuación**

Cada acción recibe un puntaje calculado con la fórmula:

\(\text{score} = (0.5 \times \text{growthScore}) + (0.3 \times \text{ratingScore}) + (0.2 \times \text{actionScore})\)

Donde:

- `growthScore` → Evaluación del crecimiento del target price basado en el porcentaje.
- `ratingScore` → Evaluación del cambio de calificación.
- `actionScore` → Evaluación de la acción realizada por el bróker (upgrade, downgrade, etc.).

Esto permite que las acciones con mejoras en su calificación y crecimiento en precio objetivo sean recomendadas primero.

---

## 🚀 Endpoints de la API

### 📌 **1. Obtener todas las acciones**

``

📌 **Descripción:** Devuelve todas las acciones almacenadas en la base de datos.

📌 **Ejemplo de Respuesta:**

```json
[
    {
        "ticker": "AMP",
        "company": "Ameriprise Financial",
        "brokerage": "Morgan Stanley",
        "action": "target raised by",
        "rating_from": "Equal Weight",
        "rating_to": "Equal Weight",
        "target_from": "$507.00",
        "target_to": "$542.00"
    }
]
```

---

### 📌 **2. Obtener las mejores recomendaciones**

``

📌 **Descripción:** Devuelve las **5 mejores acciones** para invertir según el algoritmo.

📌 **Ejemplo de Respuesta:**

```json
[
    {
        "Stock": {
            "ticker": "BSBR",
            "company": "Banco Santander (Brasil)",
            "brokerage": "The Goldman Sachs Group",
            "action": "upgraded by",
            "rating_from": "Sell",
            "rating_to": "Neutral",
            "target_from": "$4.20",
            "target_to": "$4.70"
        },
        "Score": 1.59
    }
]
```

---

## 🛠️ Instalación y Ejecución

### 📌 **1️⃣ Clonar el repositorio**

```bash
git clone <repo_url>
cd stock-api
```

### 📌 **2️⃣ Ejecutar CockroachDB con Docker**

```bash
docker-compose up -d
```

### 📌 **3️⃣ Ejecutar el servidor**

```bash
go run main.go
```

La API estará disponible en: `http://localhost:3000` 🚀

---

## 📌 Tecnologías Utilizadas

- **Golang**
- **CockroachDB**
- **Docker**
- **RESTful API**



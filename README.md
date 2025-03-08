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

Nuestra API selecciona las mejores acciones basándose en tres criterios principales:

1. **Cambio en la Calificación de la Acción**  
   - Se consideran cambios en la calificación (`Sell`, `Neutral`, `Buy`, `Overweight`, `Underweight`).
   - Se asigna una puntuación según la magnitud del cambio, por ejemplo:
     - `Sell → Buy` recibe la mayor puntuación.
     - `Buy → Sell` es penalizado con una puntuación negativa.

2. **Acción del Bróker**  
   - Se evalúan acciones como `upgraded`, `downgraded`, `target raised` o `target lowered`.
   - Acciones positivas suman puntuación, mientras que degradaciones restan.

3. **Cambio en el Precio Objetivo (Target Price)**  
   - Se calcula el porcentaje de cambio en el precio objetivo.
   - Un aumento en el `target_to` respecto a `target_from` otorga puntuación positiva.
   - Se utiliza un factor de ajuste proporcional en lugar de valores absolutos.

### 📈 **Fórmula de Puntuación**

Cada acción recibe un puntaje basado en la siguiente fórmula:

\[
\text{score} = (0.4 \times \text{ratingChangeScore}) + (0.3 \times \text{actionScore}) + (0.3 \times \text{targetPriceScore})
\]

Donde:

- `ratingChangeScore`: Evalúa el impacto del cambio en la calificación.
- `actionScore`: Evalúa la acción tomada por el bróker.
- `targetPriceScore`: Evalúa el cambio en el precio objetivo basado en porcentaje.

Las acciones con mayor puntuación se consideran mejores recomendaciones y se ordenan de mayor a menor en la respuesta de la API.

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
    }...
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
        "Score": 3.59
    }...
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



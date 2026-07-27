# TP Sistemas Operativos — The Rise of Gopher

Este proyecto es el trabajo práctico de la materia cuatrimestral **Sistemas Operativos**, centrado en la simulación de un sistema operativo distribuido, desarrollado en **Golang**.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
---

## 📄 Enunciado

Podés consultar el enunciado completo del trabajo en el siguiente enlace:  

🔗 [Consigna del TP — Google Docs](https://docs.google.com/document/d/1zoFRoBn9QAfYSr0tITsL3PD6DtPzO2sq9AtvE8NGrkc/edit?usp=sharing)

Podés consultar los tests usados para el trabajo en el siguiente enlace:

🔗 [Documentos de pruebas finales — Google Docs](https://docs.google.com/document/d/13XPliZvUBtYjaRfuVUGHWbYX8LBs8s3TDdaDa9MFr_I/edit?usp=sharing)

---

## 👨‍💻 Integrantes

| [Joseph Mansilla](https://github.com/josephmansilla) | [Ignacio Castro](https://github.com/nacho-castro) | [Santiago Torres](https://github.com/SantiagoTorres24) | [Marcelo Cabezas](https://github.com)
|:--:|:--:|:--:|:--:|
| <img src="https://avatars.githubusercontent.com/u/162230766?s=400&u=6ac208c05e9fedd414fefc12db5c38efe1c6fcd8&v=4" alt="Joseph Mansilla" width="76" height="76"> | <img src="https://avatars.githubusercontent.com/u/116680164?v=4" alt="Ignacio Castro" width="76" height="76"> | <img src="https://avatars.githubusercontent.com/u/135065796?v=4" alt="Santiago Torres" width="76" height="76"> | <img src="https://avatars.githubusercontent.com/u/143379325?v=4" alt="chelo" width="76" height="76"> |
| 🧠 Trabajó sobre el módulo **Memoria** | 🪄 Trabajó sobre el módulo **Kernel** | ⚙️ Trabajó sobre el módulo **CPU** | 🪄 Trabajó sobre el módulo **Kernel** |


![Commits en detalle](https://github.com/josephmansilla/tp-sistemas-operativos-1c-2025/graphs/contributors)
---

## Objetivos del TP

- Aplicar conceptos clave de planificación de procesos, administración de memoria y entrada/salida.
- Implementar una arquitectura distribuida con múltiples módulos comunicándose vía HTTP.
- Adquirir experiencia práctica en programación de sistemas con **Golang**.

---

## **Arquitectura del Sistema y los módulos**

El sistema está dividido en los siguientes módulos:

## - **Kernel:** planifica procesos (corto, mediano y largo plazo), administra conexiones con CPU, IO y Memoria.

![Kernel](kernel/resources/SO%202025%20KERNEL.png)

## - **CPU:** interpreta y ejecuta instrucciones, maneja TLB y caché de páginas.

<img width="972" height="594" alt="cpu" src="https://github.com/user-attachments/assets/12549952-9880-4002-b52f-1d5a6f09aae4" />

## - **Memoria + SWAP:** gestiona espacio de usuario, tablas de páginas y almacenamiento en swap.

## [VIDEO EXPLICATIVO SOBRE EL MÓDULO DE MEMORIA](https://youtu.be/twMAzy64x6Q)
<img width="1999" height="1317" alt="memoria" src="https://github.com/user-attachments/assets/510d91e1-75e2-4f57-b65b-0c271f96964d" />

[Memoria SWAP (PDF)](memoria/resources/Memoria+SWAP.pdf)
[Memoria Indexado (PDF)](memoria/resources/indexado.pdf)

## - **IO:** simula dispositivos de entrada/salida.

Todos los módulos se comunican mediante APIs HTTP, simulando un sistema operativo real distribuido.

---

## ⚙️ Tecnologías utilizadas

- 🟡 [Golang](https://go.dev/)
- 🧪 Testing con scripts y logs
- 🔌 HTTP REST APIs para la comunicación entre módulos
- 🧵 Concurrencia y sincronización

---

## 🗂 Estructura del proyecto

tp-2025/
├── cpu/
├── io/
├── kernel/
├── memoria/
├── utils/
└── scripts/ # pseudocódigos y tests




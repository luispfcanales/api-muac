# API MUAC - Sistema de Seguimiento para Niños con Desnutrición
## Descripción

API MUAC es una API REST desarrollada en Go que implementa la arquitectura hexagonal para gestionar un sistema de seguimiento de niños con desnutrición mediante la cinta MUAC (Mid-Upper Arm Circumference). Esta herramienta permite a profesionales de la salud registrar, monitorear y analizar mediciones del perímetro braquial para identificar y dar seguimiento a casos de desnutrición infantil.
## Requisitos Previos
- Go 1.16 o superior
- Git

## Configuración del Entorno de Desarrollo

### Instalación de Air (Hot Reload)

Air es una herramienta que permite la recarga en caliente de aplicaciones Go durante el desarrollo. Para instalar Air, sigue estos pasos:

1. Instala Air globalmente ejecutando:
```bash
go install github.com/cosmtrek/air@latest

## Arquitectura
La API está implementada siguiendo la arquitectura hexagonal (también conocida como puertos y adaptadores), que separa claramente:

### Dominio
- Entidades : Representan los objetos del dominio (Role, Patient, User, etc.)
- Reglas de negocio : Lógica específica del dominio
### Puertos
- Interfaces de repositorio : Definen cómo interactuar con la capa de persistencia
- Interfaces de servicio : Definen las operaciones de negocio disponibles
### Adaptadores
- Repositorios : Implementaciones concretas para acceder a la base de datos
- Manejadores HTTP : Implementaciones para exponer la API a través de HTTP
### Infraestructura
- Configuración : Gestión de variables de entorno y conexión a la base de datos
- Servidor : Configuración y gestión del servidor HTTP
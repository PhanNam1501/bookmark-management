# 🔖 Bookmark Management API 123

A robust Backend service built with **Golang** following **Clean Architecture** principles. This project provides essential services for password management, secure hashing, and system health monitoring.

[![Go Version](https://img.shields.io/badge/Go-1.25.8-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)](#)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## 🏗 Project Architecture

This project follows the **Clean Architecture** pattern to ensure high maintainability, scalability, and ease of testing by decoupling business logic from external frameworks.



### Directory Structure
* **`cmd/api`**: Application entry point. Handles initialization and startup.
* **`internal/api`**: Core server setup, routing, and Dependency Injection container.
* **`internal/handler`**: Interface layer (HTTP/REST) using the **Gin Gonic** framework.
* **`internal/service`**: Business logic layer where the primary application rules reside.
* **`docs`**: Auto-generated Swagger OpenAPI documentation.

---

## 🚀 Getting Started

### 1. Prerequisites
* **Go 1.25.8** (Installed on Windows/WSL).
* **Make** (Installed via Chocolatey or GnuWin32).
* **PowerShell** or **Git Bash**.

### 2. Installation
Clone the repository and install dependencies:
```bash
git clone [https://github.com/PhanNam1501/bookmark-management.git](https://github.com/PhanNam1501/bookmark-management.git)
cd bookmark-management
go mod tidy

# Segment3d API

Welcome to the Segment3d API project! This API provides functionality for Segment3d Mobile and Web Apps.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [Usage](#usage)
  - [Running the Server](#running-the-server)
  - [API Documentation](#api-documentation)
- [License](#license)

## Getting Started

### Prerequisites

Ensure you have the following tools installed on your machine:

- [Go 1.21.6](https://go.dev/dl/)
- [Docker](https://hub.docker.com/)

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/segment3d-app/segment3d-be.git
   cd segment3d-be
   ```

2. **Install dependencies:**

   ```bash
   go get -u ./...
   ```

3. **Install `migrate` and `sqlc`:**

   ```bash
   # Download and install the migrate CLI
   curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz | tar xvz
   sudo mv migrate.linux-amd64 /usr/local/bin/migrate

   # Download and install the sqlc CLI
   GO111MODULE=on go get github.com/kyleconroy/sqlc/cmd/sqlc@v1.10.0
   ```

   Replace `linux-amd64` with your system architecture if you are using a different operating system. Visit the [migrate GitHub releases page](https://github.com/golang-migrate/migrate/releases) and [sqlc GitHub releases page](https://github.com/kyleconroy/sqlc/releases) for the latest versions.

### Configuration

1. **Create PostgreSQL Container simply by running makefile script:**

   ```bash
    make run-container
   ```

2. **Create Database Schema:**
   ```bash
    make create-db
   ```

2. **Run Migrations:**
   ```bash
    make migrate-up
   ```

## Usage

### Running the Server

```bash
go run main.go
```

### Generate SQLC Queries

```bash
make sqlc
```

### Generate New Swagger Documentation

```bash
swag i
```

# API Documentation

To access the API documentation, visit the Swagger documentation at `/swagger` after starting the server.

## License

This project is licensed under the [MIT License](LICENSE).

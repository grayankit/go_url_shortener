# Go URL Shortener

This is a simple URL shortener written in Go.

## Prerequisites

Before running the application, you need to create a `.env` file in the root of the project.

## Configuration

The `.env` file must contain the following environment variables:

```
ADMIN_USER=your_admin_username
ADMIN_PASS=your_admin_password
```

These credentials are used to access the admin dashboard.

## Running the application

Once the `.env` file is created, you can run the application using:

```bash
go run ./cmd/server/main.go
```
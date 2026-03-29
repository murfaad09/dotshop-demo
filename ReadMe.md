# DotShop Platform

Project structure is designed to be monolethic to get start with and it will go as follows:

```
.
├── cmd
│   └── main.go       // Entry point of the application
├── internal
│   ├── config            // Configuration files
│   ├── handlers          // HTTP request handlers
│   ├── models            // Data models (using GORM)
│   ├── repositories      // Database repositories
│   ├── services          // Business logic layer
│   └── utils             // Utility functions
├── migrations            // Database migrations
├── pkg                   // Reusable packages
└── scripts               // Scripts for deployment, etc.
```

Explanation of each directory:

`cmd:` Contains the main application entry point.

`internal:` Contains the internal application code. This directory is not accessible from outside the module.

`config:` Contains configuration files for the application, such as environment variables, database configuration, etc.

`handlers:` Contains HTTP request handlers responsible for processing incoming requests and returning responses.

`models:` Contains data models representing entities in the application. We need to define our domain models using GORM here.

`repositories:` Contains database repositories responsible for interacting with the database using GORM.

`services:` Contains the business logic layer, where business logic and domain-specific operations are implemented.

`utils:` Contains utility functions used across the application.

`migrations:` Contains database migration files managed by a migration tool like goose or migrate. These files define changes to the database schema over time.

`pkg:` Contains reusable packages that can be shared across multiple applications.

`scripts:` Contains scripts for deployment, database setup, testing, etc.
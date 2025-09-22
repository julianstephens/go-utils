# go-utils

A collection of reusable Go utilities and helper functions designed to simplify common programming tasks.

## Available Packages

| Package               | Description                                                                                                                                                  |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `generic`             | Comprehensive generic utilities leveraging Go's type parameters for functional programming, slice operations, map utilities, and type-safe helpers           |
| `config`              | Reusable and idiomatic configuration management with support for environment variables, YAML/JSON files, validation, and default values                      |
| `logger`              | Unified structured logger wrapping logrus with log level control, custom formatting, and contextual logging support                                          |
| `slices`              | Generic slice utility functions for conditional selection, set operations, and element manipulation                                                          |
| `helpers`             | General utility functions including slice operations, conditional helpers, file system utilities, and struct manipulation                                    |
| `jsonutil`            | Enhanced JSON marshaling and unmarshaling with error context, formatting options, stream processing, and strict decoding support                             |
| `dbutil`              | Database utility functions and helpers for safe database interactions with connection management, query execution, transaction handling, and context support |
| `cliutil`             | Helpers and utilities for building command-line interfaces with argument parsing, interactive prompts, progress indicators, and colored output               |
| `httputil/auth`       | JWT token creation, validation, and management with role-based access control and custom claims support                                                      |
| `httputil/middleware` | Common, reusable HTTP middleware for logging, recovery, CORS, and request ID injection                                                                       |
| `httputil/request`    | HTTP request parsing utilities for JSON, form data, query parameters, and URL values                                                                         |
| `httputil/response`   | Structured HTTP response handling with extensible encoders, hooks, and status code helpers                                                                   |
| `tests`               | Shared test helpers and assertion utilities used by package tests across the repository (assertions, HTTP test helpers, and miscellaneous helpers)           |

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add support for website indexing

### Changed

- Changed protobuf schemas

## [6.2.0]

### Changed

- Enable support for latex in text generation

## [6.1.1] - 2024-02-13

### Added

- Added support for more bedrock models

## [6.1.0] - 2024-02-13

### Changed

- Changed LLM model names

### Added

- Added support for AWS Bedrock

## [6.0.1] - 2024-02-11

### Changed

- Pool embedding usage tracking to reduce database entries

## [6.0.0] - 2024-02-11

### Changed

- Changed collection proto specs

### Added

- Added GetCollection function

## [5.0.2] - 2024-02-10

### Fixed

- Fix cost calculation

## [5.0.1] - 2024-02-10

### Fixed

- Prevent cloud deployment of insecure auth service
- Fix qdrant document deletion

## [5.0.0] - 2024-02-10

### Added

- Add support for follow-up questions

## [4.3.1] - 2024-02-03

### Fixed

- Retry embedding generation on error
- Fix bug that prevented insertions of documents in the database

## [4.3.0]

### Added

- Add support for Qdrant Vector Search
- Switched to OpenAI's text-embedding-3-large embedding model

## [4.2.0] - 2023-12-23

### Added

- Implement DeleteChatMessage function

## [4.1.2] - 2023-12-19

### Fixed

- Fix page numbers bug

## [4.1.1] - 2023-12-19

### Changed

- Rename package to github.com/pzierahn/chatbot_services

### Fixed

- Fix none existing local credentials bug

## [4.1.0] - 2023-12-19

### Added

- Support for Google Gemini Pro

## [4.0.0] - 2023-12-16

### Changed

- Changed protobuf schema

## [3.0.0] - 2023-12-14

### Changed

- Use Google Cloud instead of Supabase

## [2.0.5] - 2023-11-08

### Fix

- Fix openai api issue

## [2.0.4] - 2023-11-07

### Added

- Added GTP-4 Turbo model

### Fixed

- Fix Markdown code block issue with text generation

## [2.0.3]

### Changed

- Added founding check for users

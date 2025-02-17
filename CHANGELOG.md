# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [8.0.9]

### Changed

- Update dependencies

## [8.0.8]

### Changed

- Update Anthropic Sonnet version to v2

## [8.0.7] - 2024-09-25

### Changed

- Require collection id when deleting a document from search index
- Update Google Gemini models

## [8.0.6] - 2024-08-07

### Changed

- Added openai model costs

## [8.0.5] - 2024-08-07

### Changed

- Record model id from API response

## [8.0.4] - 2024-07-19

### Added

- Added gpt-4o-mini model

## [8.0.3] - 2024-07-14

### Fixed

- Fix notion execution

## [8.0.2] - 2024-07-12

### Added

- Added model costs for openai.gpt-4o and anthropic.claude-3-5-sonnet-20240620-v1:0

### Changed

- Record model usage with model id from response

## [8.0.1] - 2024-07-12

### Changed

- Updated system prompt

## [8.0.0] - 2024-07-12

### Changed

- Replace PostgresDB with MongoDB
- Implement tool calls for better retrieval of data
- Updated multiple protobuf schemas

## [7.9.1] - 2024-06-21

### Added

- Added support for Claude 3.5 Sonnet

## [7.9.0] - 2024-06-11

### Changed

- Improved source quoting

## [7.8.2] - 2024-06-05

### Fixed

- Fix a bug that prevented the insertion of embedding vectors into the database

## [7.8.1] - 2024-06-05

### Fixed

- Fix a bug that prevented to get the current notion api key

## [7.8.0] - 2024-05-13

### Added

- Added GPT-4o

## [7.7.3] - 2024-05-01

### Added

- Trim suffix from notion filenames

### Fixed

- Fix a problem api key caching

## [7.7.2] - 2024-05-01

### Added

- Added notion api key management

## [7.7.1] - 2024-05-01

### Changed

- Update system prompt

## [7.7.0] - 2024-04-30

### Added

- Added notion service

### Removed

- Removed Anthropic package

## [7.6.4] - 2024-04-10

### Fixed

- Fixed Gemini model costs

## [7.6.3] - 2024-04-10

### Added

- Added gpt-4-turbo model

## [7.6.2] - 2024-04-09

### Fixed

- Update Gemini Pro 1.5 model pricing

## [7.6.1] - 2024-04-09

### Fixed

- Fix model selection for Gemini Models

### Updated

- Update model Gemini prices

## [7.6.0] - 2024-04-09

### Added

- Added Gemini Pro 1.5

## [7.5.0] - 2024-04-07

### Added

- Added chat batch processing

### Removed

- Removed support for Mistral AI

### Changed

- Sort models by name

## [7.4.2] - 2024-04-01

### Added

- Added support for Claude 3 Haiku

### Removed

- Removed support Pinecone

## [7.4.1] - 2024-03-09

### Fixed

- Fixed user tracking for Claude Opus

### Added

- Added pricing for Claude 3 Opus

## [7.4.0] - 2024-03-07

### Added

- Add support for Claude 3 Opus

## [7.3.1] - 2024-03-05

### Fixed

- Fix alternating roles bug for Claude 3

### Added

- Add cost support for Claude 3

## [7.3.0] - 2024-03-05

### Added

- Added support for Claude 3

## [7.2.0]

### Added

- Added Mistral AI models

## [7.1.0] - 2024-02-27

### Added

- Added crashlytics service

## [7.0.3] - 2024-02-26

### Fixed

- Fixed balance sheet

## [7.0.2] - 2024-02-23

### Changed

- Remove case sensitivity from document search

## [7.0.1] - 2024-02-23

### Fixed

- Fixed document search bug

## [7.0.0]

### Added

- Add support for website indexing
- Added support for references in chat messages

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

---
name: code-review-testing
description: Comprehensive code review, automated testing, verification, and regression prevention for web applications and backend APIs.
---

# Code Review & Testing Skill

This skill provides a systematic approach for conducting thorough code reviews, executing automated tests, validating UI/UX styling, ensuring API data integrity, and verifying zero regressions across the codebase.

## Core Review Principles

1. **Functional Correctness**:
   - Verify that business logic aligns with system models and database schemas.
   - Ensure all edge cases (null values, missing fields, token expiration, disconnected states) are properly handled with resilient fallbacks.

2. **Data Consistency & Type Safety**:
   - Validate MongoDB queries against multiple ID representations (`ObjectId`, `string`, `telegram_id`, `int`).
   - Prevent data loss during API serialization and schema transformations.
   - Maintain strict TypeScript type definitions across API clients, services, and Vue composables.

3. **UI / Design System Compliance**:
   - Verify adherence to the WEYANA Design System:
     - Primary Orange: `#F97316` / `#EA580C`
     - Light Background: `#FFFFFF`, Dark Background: `#0F0F0F` / `#181818` / `#222222`
     - Border Tokens: `#E5E7EB` (Light), `#27272A` (Dark)
     - Text Tokens: `#111111` (Primary Light), `#FFFFFF` (Primary Dark), `#6B7280` / `#A1A1AA` (Secondary)
   - Ensure no nested button or invalid HTML element hierarchies exist.

4. **Automated Verification Workflow**:
   - **Frontend**: Execute `npm run build` to ensure type checking and asset bundling compile with 0 errors.
   - **Backend**: Test FastAPI routers and dependencies for syntax, import correctness, and schema validation.
   - **Git Hygiene**: Verify clean working directory and clean commit structure before pushing to remote.

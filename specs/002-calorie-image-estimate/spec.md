# Feature Specification: Calorie Estimation from Image

**Feature Branch**: `002-calorie-image-estimate`
**Created**: 2025-12-15
**Status**: Draft
**Input**: User description: "Implement an MVP Telegram bot feature that estimates calories from a user-uploaded image via LUI only, with no storage.

User flow:
- User runs /estimate and is prompted to upload one image.
- User uploads an image.
- Bot replies with a deterministic, fixed-format result including: estimated calories and a confidence indicator.
- Provide inline buttons: Re-estimate and Cancel (minimal interaction only).

Quality + deliverables required:
- Provide a working bot name and runnable bot.
- Produce readable code with explicit quality controls (tests/lint/review gates).
- Archive all vibe coding prompts used to generate code.
- Use an LLM-based testing tool to validate the conversational flow (simulate /estimate → image → result).
- Archive LLM test prompts and output a test report with pass/fail.

Out of scope:
- Any persistence or history listing.
- Mini App and any GUI CRUD.
- Deployment automation."

## User Scenarios & Testing

### User Story 1 - Single Image Calorie Estimation (Priority: P1)

A user wants to quickly estimate the calories in their meal by taking a photo and receiving an immediate estimate through the Telegram chat interface.

**Why this priority**: This is the core value proposition of the bot - providing quick, accessible calorie estimation. Without this, the bot has no functionality.

**Independent Test**: Can be fully tested by sending the /estimate command, uploading a food image, and receiving a calorie estimate with confidence indicator. Delivers immediate value by providing actionable nutritional information.

**Acceptance Scenarios**:

1. **Given** user is in a Telegram chat with the bot, **When** user sends /estimate command, **Then** bot prompts user to upload one image
2. **Given** user has been prompted to upload an image, **When** user uploads a food image, **Then** bot analyzes the image and returns a calorie estimate with confidence indicator in a fixed format
3. **Given** user has received a calorie estimate, **When** user views the result, **Then** the message includes inline buttons for "Re-estimate" and "Cancel"

---

### User Story 2 - Re-estimation Flow (Priority: P2)

A user receives a calorie estimate but wants a fresh analysis of the same or different image without restarting the conversation flow.

**Why this priority**: Enhances user experience by allowing quick iterations without typing commands again. Secondary to the core estimation feature.

**Independent Test**: Can be tested independently by completing an initial estimation, then clicking "Re-estimate" button and uploading another image. Delivers value by streamlining repeated estimations.

**Acceptance Scenarios**:

1. **Given** user has received a calorie estimate with inline buttons, **When** user clicks "Re-estimate" button, **Then** bot prompts user to upload a new image
2. **Given** user clicked "Re-estimate" and is prompted for a new image, **When** user uploads a new food image, **Then** bot provides a fresh calorie estimate with confidence indicator

---

### User Story 3 - Cancellation Flow (Priority: P3)

A user wants to exit the estimation flow without completing it or after receiving results.

**Why this priority**: Nice-to-have for user control but not essential for MVP functionality. Users can simply stop interacting with the bot.

**Independent Test**: Can be tested by clicking "Cancel" at any point in the flow. Delivers value by providing clean exit points.

**Acceptance Scenarios**:

1. **Given** user has received a calorie estimate with inline buttons, **When** user clicks "Cancel" button, **Then** bot acknowledges cancellation and clears the current interaction
2. **Given** user is in the middle of an estimation flow, **When** user clicks "Cancel" button, **Then** bot stops the current flow and returns to idle state

---

### Edge Cases

- What happens when user uploads a non-food image (e.g., landscape, person, random object)?
- What happens when user uploads an image with no recognizable food items?
- What happens when user uploads multiple images instead of one?
- What happens when user sends text instead of an image after /estimate prompt?
- What happens when user uploads an extremely large image file?
- What happens when user uploads a corrupted or invalid image file?
- How does the bot handle network timeouts during image analysis?
- What happens if user sends /estimate command while already in an active estimation flow?

## Requirements

### Functional Requirements

- **FR-001**: System MUST provide a /estimate command that initiates the calorie estimation flow
- **FR-002**: System MUST prompt user to upload exactly one image after /estimate command is received
- **FR-003**: System MUST accept image uploads in common formats (JPEG, PNG, WebP)
- **FR-004**: System MUST analyze uploaded food images and provide a calorie estimate
- **FR-005**: System MUST include a confidence indicator with each calorie estimate (e.g., low/medium/high, or percentage)
- **FR-006**: System MUST present results in a deterministic, fixed format for consistency
- **FR-007**: System MUST provide inline keyboard buttons "Re-estimate" and "Cancel" with each result
- **FR-008**: System MUST handle "Re-estimate" button click by restarting the image upload prompt
- **FR-009**: System MUST handle "Cancel" button click by ending the current estimation flow
- **FR-010**: System MUST NOT persist any user data, images, or estimation history
- **FR-011**: System MUST NOT store uploaded images after analysis is complete
- **FR-012**: System MUST process each estimation request independently with no reference to previous requests
- **FR-013**: System MUST provide clear error messages when user uploads invalid content (non-image, corrupted file, etc.)
- **FR-014**: System MUST handle cases where no food is detected in the uploaded image
- **FR-015**: System MUST reject all images when multiple images are uploaded and prompt user to send exactly one image
- **FR-016**: System MUST analyze food images using vision AI capabilities to extract calorie information
- **FR-017**: Bot MUST have a unique, descriptive name that reflects its calorie estimation purpose

### Key Entities

This feature operates without persistence, so there are no stored entities. All operations are stateless and ephemeral:

- **Estimation Session**: Temporary in-memory state representing a single /estimate flow (not persisted)
  - User ID (from Telegram)
  - Current flow state (awaiting image, processing, showing results)
  - Uploaded image (held in memory during analysis only)
  - Analysis result (calorie estimate + confidence)

## Success Criteria

### Measurable Outcomes

- **SC-001**: Users can complete a full estimation flow (command → upload → result) in under 10 seconds under normal conditions
- **SC-002**: System provides calorie estimates with confidence indicators for at least 90% of food images submitted
- **SC-003**: Bot responds to /estimate command within 1 second
- **SC-004**: Image analysis and response delivery completes within 8 seconds of image upload
- **SC-005**: Automated conversational flow tests pass with 100% success rate for defined user scenarios
- **SC-006**: Code passes all linting and quality gates before deployment
- **SC-007**: All vibe coding prompts and LLM test prompts are archived in the repository
- **SC-008**: Test report shows clear pass/fail results for each tested scenario

## Assumptions

- Users will primarily upload images containing recognizable food items
- Calorie estimates do not need to be nutritionist-level accurate; reasonable approximations are acceptable for MVP
- Bot will operate in English language only for MVP
- Users understand that estimates are approximate and not medical/nutritional advice
- Single-user testing flow is sufficient for MVP; no multi-user concurrency stress testing required
- Google Gemini Vision API (via ai.dev/Google AI Studio) has been selected as the vision AI provider for image analysis
- Google Gemini Vision API will be accessible with reasonable rate limits for MVP testing
- Telegram Bot API v3 (github.com/tucnak/telebot) supports inline keyboard buttons and image handling
- Fixed-format response means consistent structure (e.g., "Estimated Calories: X kcal | Confidence: Y%")

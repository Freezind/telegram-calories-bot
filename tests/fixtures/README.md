# Test Fixtures

This directory contains test assets used by the automated testing suite.

## Required Files

### food.jpg

**Purpose:** Test image for calorie estimation scenarios (S2, S3, S4)

**Requirements:**
- File format: JPEG, PNG, or WebP
- Content: A clear photo of food items
- Recommended: Simple meal with 2-3 recognizable food items (e.g., grilled chicken with vegetables)
- Size: < 5MB (Telegram file size limit for photos)

**How to add:**
1. Find or take a photo of a simple meal
2. Save it as `food.jpg` in this directory
3. The test runner will use this image for estimate scenarios

**Example foods that work well:**
- Grilled chicken with rice and broccoli
- Sandwich with chips
- Pasta with tomato sauce
- Salad with protein

**Note:** The LLM judge validates response structure only (foods list, calories number, confidence level), not calorie accuracy. Any food image will work for testing purposes.

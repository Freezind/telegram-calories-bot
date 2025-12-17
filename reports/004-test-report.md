# Test Report: LLM-Based Bot Testing

**Generated:** 2025-12-17T09:40:45Z
**Test Duration:** 26s
**Total Scenarios:** 5
**Passed:** 1
**Failed:** 4

---

## Summary

| Scenario | Verdict | Rationale |
|----------|---------|-----------||
| /start Welcome Message | ❌ FAIL | Test execution error: LLM judge evaluation failed: gemini API call failed: Error 503, Message: The model is overloaded. Please try again later., Status: UNAVAILABLE, Details: [] |
| /estimate + Image Upload | ✅ PASS | The evidence clearly shows a prompt was sent with the text 'Please send a food image', which directly matches the expected behavior of the bot prompting for an image after the /estimate command. |
| Re-estimate Button Message Preservation | ❌ FAIL | Test execution error: LLM judge evaluation failed: gemini API call failed: Error 503, Message: The model is overloaded. Please try again later., Status: UNAVAILABLE, Details: [] |
| Cancel Button Message Preservation | ❌ FAIL | Test execution error: LLM judge evaluation failed: gemini API call failed: Error 503, Message: The model is overloaded. Please try again later., Status: UNAVAILABLE, Details: [] |
| Mini App Page Load | ❌ FAIL | Test execution error: failed to GET Mini App URL: Get "https://lesson-teachers-any-processed.trycloudflare.com": dial tcp: lookup lesson-teachers-any-processed.trycloudflare.com: no such host |

---

## Scenario Details

### S1: /start Welcome Message

**Duration:** 18.285s

**Captured Evidence:**

Evidence 1 (bot_message):
```json
{
  "type": "bot_message",
  "timestamp": "2025-12-17T18:40:19.383562+09:00",
  "data": {
    "message_id": 1,
    "message_sent": true,
    "text": "👋 Welcome to Calorie Estimation Bot!\n\nI help you estimate the calories in your food by analyzing images.\n\n**How to use:**\n1. Send /estimate command\n2. Upload a photo of your food\n3. Receive calorie estimate with confidence indicator\n\n**Features:**\n• 🍽️ Instant calorie estimation\n• 📊 Confidence indicators (Low/Medium/High)\n• 🔄 Re-estimate with different images\n• ❌ Cancel anytime\n\nReady to start? Send /estimate to begin!"
  }
}
```

**Error:**
```
LLM judge evaluation failed: gemini API call failed: Error 503, Message: The model is overloaded. Please try again later., Status: UNAVAILABLE, Details: []
```

**LLM Judge Verdict:** ❌ FAIL
**Rationale:** Test execution error: LLM judge evaluation failed: gemini API call failed: Error 503, Message: The model is overloaded. Please try again later., Status: UNAVAILABLE, Details: []

---

### S2: /estimate + Image Upload

**Duration:** 2.689s

**Captured Evidence:**

Evidence 1 (estimate_prompt):
```json
{
  "type": "estimate_prompt",
  "timestamp": "2025-12-17T18:40:37.669446+09:00",
  "data": {
    "prompt_sent": true,
    "prompt_text": "📸 Please send a food image for calorie estimation"
  }
}
```

**LLM Judge Verdict:** ✅ PASS
**Rationale:** The evidence clearly shows a prompt was sent with the text 'Please send a food image', which directly matches the expected behavior of the bot prompting for an image after the /estimate command.

---

### S3: Re-estimate Button Message Preservation

**Duration:** 2.666s

**Captured Evidence:**

Evidence 1 (message_preservation):
```json
{
  "type": "message_preservation",
  "timestamp": "2025-12-17T18:40:40.358465+09:00",
  "data": {
    "new_prompt_sent": true,
    "new_prompt_text": "📸 Please send another food image",
    "previous_message_deleted": false,
    "previous_message_id": 1
  }
}
```

**Error:**
```
LLM judge evaluation failed: gemini API call failed: Error 503, Message: The model is overloaded. Please try again later., Status: UNAVAILABLE, Details: []
```

**LLM Judge Verdict:** ❌ FAIL
**Rationale:** Test execution error: LLM judge evaluation failed: gemini API call failed: Error 503, Message: The model is overloaded. Please try again later., Status: UNAVAILABLE, Details: []

---

### S4: Cancel Button Message Preservation

**Duration:** 2.734s

**Captured Evidence:**

Evidence 1 (message_preservation):
```json
{
  "type": "message_preservation",
  "timestamp": "2025-12-17T18:40:43.024987+09:00",
  "data": {
    "cancellation_sent": true,
    "cancellation_text": "Estimation canceled. Use /estimate to start again.",
    "previous_message_deleted": false,
    "previous_message_id": 2
  }
}
```

**Error:**
```
LLM judge evaluation failed: gemini API call failed: Error 503, Message: The model is overloaded. Please try again later., Status: UNAVAILABLE, Details: []
```

**LLM Judge Verdict:** ❌ FAIL
**Rationale:** Test execution error: LLM judge evaluation failed: gemini API call failed: Error 503, Message: The model is overloaded. Please try again later., Status: UNAVAILABLE, Details: []

---

### S5: Mini App Page Load

**Duration:** 6ms

**Error:**
```
failed to GET Mini App URL: Get "https://lesson-teachers-any-processed.trycloudflare.com": dial tcp: lookup lesson-teachers-any-processed.trycloudflare.com: no such host
```

**LLM Judge Verdict:** ❌ FAIL
**Rationale:** Test execution error: failed to GET Mini App URL: Get "https://lesson-teachers-any-processed.trycloudflare.com": dial tcp: lookup lesson-teachers-any-processed.trycloudflare.com: no such host

---

## Test Completion

**Test run completed at:** 2025-12-17T09:40:45Z
**Overall result:** ❌ FAIL (4 scenario(s) failed)

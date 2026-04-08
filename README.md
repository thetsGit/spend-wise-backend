# SpendWise Backend

Email expense analyzer and SaaS subscription discovery service. Accepts email data via JSON upload, analyzes the content using AI (OpenAI), extracts categorized spending and SaaS signals, then serves the results through a REST API.

## Context & Approach

### Starting Point

My background is frontend heavy (React, Vue, React Native, TypeScript) with limited backend exposure. This project required technologies I had zero or minimal experience with:

| Area | Prior Experience |
|------|-----------------|
| Go | None, first time writing Go |
| AI / prompt engineering | Limited, have been using AI as a chat user but never integrated APIs |
| PostgreSQL | Some experience with relational databases |
| Docker / Docker Compose | Some experience |
| HTTP servers | Some experience in JS/PHP |

### How I Approached It

Given the knowledge gaps, I leaned heavily into problem-based learning. Instead of trying to "learn Go" as a whole, I broke the assignment down into small problems (connect to Postgres, build an HTTP handler, call an external API) and figured out each one by asking the right questions, drawing analogies from what I already know in JS/TS/PHP, and working through it step by step.

For example, understanding Go structs was easier once I mapped them to TypeScript interfaces. Understanding `error` returns clicked once I compared it to try/catch. Connection pooling in pgx made sense once I thought about it like a database client in Node. This kind of analogical thinking helped me move faster than going through tutorials.

I also made good use of AI tools and documentation along the way. They were helpful for closing syntax gaps, sanity-checking architecture decisions, and speeding up the parts where I already understood the "what" and "why" but needed help with the "how" in a new language. The design decisions, tradeoffs, and architecture choices are my own.

Rather than jumping into code, I invested time upfront in research and architecture planning:

1. **Schema design first.** Modeled the relationship between raw emails and AI-extracted spending/SaaS data before writing any code
2. **Preset-driven validation.** Defined all constrained field values (categories, signal types, billing cycles, confidence levels) as a single source of truth, used across prompts, validation, and normalization
3. **Prompt engineering.** Structured prompts following [OpenAI's prompt engineering best practices](https://platform.openai.com/docs/guides/prompt-engineering) with clear role, explicit constraints, exact output format, and edge case handling
4. **Go learning.** Mapped Go concepts to familiar JS/TS patterns (structs = interfaces, packages = modules, `error` return = try/catch) to move quickly without getting stuck on syntax

## Architecture

```
[JSON Upload] → [Validate] → [PostgreSQL]
                                  ↓
                             [Build Prompt]
                                  ↓
                            [OpenAI API]
                                  ↓
                          [Parse AI Response]
                                  ↓
                     [Normalize + Score + Insert]
                                  ↓
                        [spending + saas_discovery]
                                  ↓
                          [REST API ← Frontend]
```

### Data Model

```
email (raw input, source of truth)
  ├── spending (AI-extracted transactions)
  └── saas_discovery (AI-detected SaaS signals)
```

One email can produce a spending record, a SaaS detection, both, or neither. They are sibling tables linked through `email_id`, not directly related to each other.

### Processing Pipeline (Upload Endpoint)

```
1. Parse JSON body, validate emails, reject invalid ones
2. Insert valid emails to DB, skip duplicates (ON CONFLICT)
3. Build AI prompt with all emails (batch, single API call)
4. Call OpenAI, parse JSON response
5. For each result:
   - Normalize fields against presets (e.g AI says "Food Delivery", we map to "food_delivery")
   - Calculate confidence score (rule-based, not trusting AI self-assessment blindly)
   - Insert spending/saas records
6. Return upload summary
```

## Key Design Decisions

### 1. Preset-Driven Validation

All constrained fields are defined as Go maps, acting as the single source of truth:

```go
var SpendingCategories = map[string]bool{
    "food_delivery": true, "travel": true, "software": true, ...
}
```

These presets are dynamically injected into AI prompts so the LLM knows the valid values. They are also used for runtime normalization where invalid AI output falls back gracefully. The AI is constrained but the system does not break if it returns unexpected values.

### 2. Dual Confidence Scoring

Each spending and SaaS record stores two confidence values:
- `ai_confidence`: what the LLM self-reported
- `confidence`: our own rule-based score computed from field completeness

This way the system's confidence is decoupled from the AI provider. If the AI says "high" but critical fields are missing, our score will reflect the actual reality.

### 3. Batch AI Processing

All emails are sent to OpenAI in a single prompt, extracting both spending and SaaS data at once. This means one API call instead of N×2 calls, which reduces latency, cost, and API rate limit risk.

### 4. Duplicate Detection

Composite unique constraint on `(sender, recipient, subject, date)`. Re-uploading the same file is safe since duplicates are silently skipped via `ON CONFLICT DO NOTHING`.

### 5. AI Prompt Design

Structured prompting approach following [OpenAI's best practices](https://platform.openai.com/docs/guides/prompt-engineering):

- **Role**: "You are an email analyzer that performs two tasks"
- **Constraints**: preset values listed explicitly per field
- **Edge cases**: use `null` for unknown fields
- **Important rules**: "An email can produce BOTH a spending record AND a SaaS signal" (discovered during testing that AI was skipping spending records for SaaS invoices)
- **Output format**: exact JSON schema provided at the end of the prompt
- **Guard**: "Respond with ONLY a valid JSON array. No markdown, no explanation."

## Project Structure

```
├── cmd/server/main.go           # Entry point, bootstrap, router, server
├── internal/
│   ├── ai/
│   │   ├── ai.go                # OpenAI HTTP client
│   │   └── models.go            # Request/response types (unexported)
│   ├── config/config.go         # Environment config with defaults
│   ├── database/
│   │   ├── connection.go        # pgx pool connection with ping verification
│   │   └── queries.go           # All DB operations (insert, select, aggregate)
│   ├── handlers/
│   │   ├── handlers.go          # Handler struct and constructor
│   │   ├── helpers.go           # JSON response helpers
│   │   ├── upload.go            # POST /api/emails/upload (the main pipeline)
│   │   ├── spending.go          # GET /api/spending + summary
│   │   └── saas.go              # GET /api/saas + summary
│   ├── models/
│   │   ├── models.go            # All structs (Email, Spending, SaaS, API types)
│   │   └── methods.go           # Validation, scoring, AI response parsing
│   ├── presets/
│   │   ├── presets.go           # Constrained field values (single source of truth)
│   │   └── normalizers.go       # Field normalization functions
│   ├── prompts/prompts.go       # AI prompt builder (injects presets dynamically)
│   └── utils/utils.go           # Generic normalize and map keys helpers
├── migrations/v001_init.sql     # Database schema (auto-runs on first boot)
├── sample_emails.json           # Test dataset from assignment
├── sample_emails_2.json         # Additional test dataset
├── docker-compose.yaml          # Postgres + API services
├── Dockerfile                   # Multi-stage build (builder then alpine)
└── .env.example                 # Environment template
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/emails/upload` | Upload email JSON, triggers AI analysis pipeline |
| `GET` | `/api/spending` | List transactions. Filters: `category`, `start_date`, `end_date` |
| `GET` | `/api/spending/summary` | Aggregated spend by category with totals |
| `GET` | `/api/saas` | List detected SaaS tools. Filters: `product_name`, `signal_type` |
| `GET` | `/api/saas/summary` | Total estimated monthly SaaS spend and tool count |

All responses follow a consistent envelope:

```json
{
  "status": "success",
  "status_code": 200,
  "message": "Success",
  "data": { ... }
}
```

## Setup

### Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose
- An OpenAI API key ([platform.openai.com](https://platform.openai.com/))

### Quick Start

```bash
# 1. Clone
git clone https://github.com/thetsGit/spend-wise-be.git
cd spend-wise-be

# 2. Configure environment
cp .env.example .env
# Fill in your values (see Environment Variables below)

# 3. Start everything
docker compose up --build

# 4. Test the pipeline
curl -X POST http://localhost:8000/api/emails/upload \
  -H "Content-Type: application/json" \
  -d @sample_emails.json
```

### Environment Variables

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=spend_wise
DB_HOST=db
DB_PORT=5432
HTTP_PORT=8000
ALLOWED_ORIGINS=http://localhost:5173

OPEN_AI_URL=https://api.openai.com/v1/chat/completions
OPEN_AI_MODEL=gpt-4o-mini
OPEN_AI_API_KEY=your_key_here

MAX_UPLOAD_SIZE_KB=5120
```

### Test Endpoints

```bash
# Upload emails
curl -X POST http://localhost:8000/api/emails/upload \
  -H "Content-Type: application/json" \
  -d @sample_emails.json

# Query spending
curl http://localhost:8000/api/spending
curl http://localhost:8000/api/spending/summary
curl "http://localhost:8000/api/spending?category=software"

# Query SaaS
curl http://localhost:8000/api/saas
curl http://localhost:8000/api/saas/summary

# Check DB directly
docker compose exec db psql -U postgres -d spend_wise \
  -c "SELECT merchant, amount, category, confidence FROM spending;"
docker compose exec db psql -U postgres -d spend_wise \
  -c "SELECT product_name, signal_type, billing_cycle, estimated_cost, confidence FROM saas_discovery;"
```

## What I Would Improve Given More Time

- **Transaction rollback.** Wrap email inserts in a DB transaction so partial failures do not leave orphan records. Right now if the 5th email insert fails, the first 4 are already saved with no way to undo
- **Reprocessing endpoint.** Something like `POST /api/emails/reprocess` to re-run AI on already stored emails when prompts get improved. Since raw emails are stored as source of truth, this would be straightforward
- **Presets API.** A `GET /api/presets` endpoint so frontend can dynamically build filter dropdowns from the backend source of truth instead of hardcoding the same values
- **Modular prompt templates.** Right now spending and SaaS extraction share one combined prompt. Ideally each analyzer (spending, SaaS, and potentially new ones like job alerts or subscription renewals) would have its own template with a shared base, so extending does not require touching existing prompts
- **Multi-currency support.** An `amount_usd` column computed at insert time using exchange rates, so aggregation across different currencies actually makes sense
- **Monitor fallback frequency.** Track how often normalization falls back to `other` or `unknown`. If a certain value keeps appearing, that is a signal to add it as a new preset instead of lumping it into the catch-all
- **Rate limiting** on both the API endpoints and outgoing AI provider calls to prevent abuse and budget overruns
- **LLM provider adapter.** Abstract the AI call behind an interface so swapping between OpenAI, Anthropic, Gemini, or any other provider is just a config change. Right now the OpenAI integration is directly coupled
- **Wider input support.** Currently only JSON upload is supported. Would like to add CSV parsing, and eventually direct email integration via IMAP or Gmail API so users do not need to export emails manually
- **Store raw AI response.** Save the full AI output per email for debugging and prompt iteration. When the prompt changes, having the old responses helps compare results

## Time Breakdown

| Phase | Time (approx) |
|-------|---------------|
| Research and architecture (schema design, Go basics, prompt strategy) | ~3 hrs |
| Backend: DB setup, migrations, Docker Compose | ~1 hr |
| Backend: Models, presets, normalizers, validators | ~1.5 hrs |
| Backend: AI integration (provider setup, prompt building, response parsing) | ~2 hrs |
| Backend: Upload pipeline (main handler wiring everything together) | ~1.5 hrs |
| Backend: GET endpoints, filters, summaries | ~1 hr |
| Backend: Docker, CORS, config, error handling polish | ~1 hr |
| Frontend: API layer, views, components (separate repo) | ~3 hrs |
| Testing and prompt refinement | ~1 hr |
| README and submission | ~1 hr |
| **Total** | **~16 hrs** |

The extended time mostly reflects learning Go from scratch and navigating AI provider billing issues (started with Anthropic, then tried OpenAI, then Gemini, and finally back to OpenAI with purchased credits). An equivalent project in TypeScript/Node.js would have probably taken around 4 to 5 hours.

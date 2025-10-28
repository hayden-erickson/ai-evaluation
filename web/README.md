# Habit Streaks UI

A modern, minimal React (Vite + TypeScript) frontend for the Habits API.

This UI lets users:

- Login/Register via phone number + password.
- Create, edit, delete habits.
- View current streaks with a one-day skip allowed.
- Create, edit, delete logs for each habit (optionally with durations).

## Tech

- React 18 + TypeScript
- Vite dev server with proxy to API at `/api`
- No extra UI libraries; custom minimal CSS

## Prerequisites

- Node.js 18+
- Go 1.21+

## Local Development

1) Start the Go API (port 8080 by default):

- Terminal A (repo root: `new-ui`):
  - `go run main.go`

2) Start the web app:

- Terminal B (web folder):
  - `npm install`
  - `npm run dev`

The app will open at http://localhost:5173 and proxy API calls from `/api/*` to `http://localhost:8080/*`.

You can override the API base by creating a `.env` file in `web/` with:

```
VITE_API_BASE=/api
```

## Production Build

- From `web/`:
  - `npm run build`
  - Output in `web/dist/`

### Deploying the Web UI to AWS (S3 + CloudFront)

1) Build the UI: `npm run build`
2) Create an S3 bucket for static website hosting (private, Block Public Access ON).
3) Upload `dist/` contents to S3.
4) Create a CloudFront distribution with S3 as origin.
5) Configure default root object to `index.html`.
6) Set appropriate cache policies (e.g., cache static assets aggressively, `index.html` with low TTL).

### Hosting the API

- Run the Go API behind an Application Load Balancer or on ECS/EC2 with HTTPS via ACM.
- Recommended: expose the API under the same domain using a reverse proxy path of `/api` to avoid CORS in production.
  - Example: CloudFront behavior `/api/*` -> API origin; default behavior -> S3.
  - Ensure the API honors `X-Forwarded-*` headers if needed.

### CORS and CSP Notes

- In development, the Vite proxy avoids CORS issues.
- If you deploy the UI and API on different domains in production, enable CORS on the API accordingly.
- The API currently sets a basic Content-Security-Policy. When serving the frontend from S3/CloudFront and the API on a different origin, adjust CSP and CORS to allow requests to the API origin.

## Environment Variables (web)

- `VITE_API_BASE` (default: `/api`) — Change if you are not using a reverse proxy and need to point directly to an API origin like `https://api.example.com`.

## Troubleshooting

- TypeScript errors like "Cannot find module 'react'" typically mean `npm install` hasn't been run.
- If API calls fail with 401, ensure you have registered and logged in. The token is stored in `localStorage`.
- If your habit requires duration (created with duration_seconds), logs must include `duration_seconds` per backend validation.

# AI Booking Travel - Car Rental Brokerage Aggregation

This is a complete B2B and B2C Car Rental Brokerage application. It aggregates car rentals from multiple suppliers (like Hertz and Flex), enabling searches, bookings, and reservations.

- **Agents (B2B)**: Operate under offices and organizations. They enjoy lower markups, and billing is deferred (managed through monthly invoicing to the office/organization).
- **Customers (B2C)**: Pay upfront via credit card before a booking can be placed. User creation and authentication happen via Phone Number + OTP at the time of checkout.

## Architecture

- **Backend**: Built in Go using the **[Encore](https://encore.dev/)** framework. It boasts an autonomous infrastructure and automatic distributed tracing. Uses Postgres (`sqlc`) for DB interaction.
- **Frontend**: Next.js App Router (React), utilizing Next-Intl for localization (Hebrew and English), React Query with generated Encore clients, Shadcn UI & Tailwind CSS. **Payload CMS** is also integrated natively inside the frontend to handle site content via the database.

---

## Prerequisites

Before running the application locally, ensure you have the following tools installed:

1. **[Docker](https://www.docker.com/products/docker-desktop/)**: Must be running in the background. Encore uses Docker to automatically provision local PostgreSQL databases and pub/sub.
2. **[Encore CLI](https://encore.dev/docs/install)**: The backend development tool.
   - macOS/Linux: `curl -L https://encore.dev/install.sh | bash`
   - or via Homebrew: `brew install encoredev/tap/encore`
3. **[Go](https://go.dev/doc/install)** (1.22+ required).
4. **[Node.js](https://nodejs.org/) & [pnpm](https://pnpm.io/installation)**: Core frontend tools. Install `pnpm` natively or via npm: `npm i -g pnpm`.

---

## Local Development

### 1. Backend

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. **Setup Secrets**: Look into the app code, numerous packages expect secrets. You can set them for your local environment using the `encore secret set` command:
   ```bash
   # Core & Notifications
   encore secret set --type local SecretKey
   encore secret set --type local emailPassword
   encore secret set --type local smsToken
   encore secret set --type local FirstAdminEmail
   encore secret set --type local FirstAdminPassword

   # Integrations
   encore secret set --type local IcountPassword
   encore secret set --type local translationToken

   # Broker Suppliers
   encore secret set --type local flexAgentCode
   encore secret set --type local flexPassword
   encore secret set --type local hertzAgentDutyCode
   encore secret set --type local hertzVendorNumber
   encore secret set --type local hertzCodeContext
   ```
3. **Run the Backend**:
   Make sure Docker is running, then run the app:
   ```bash
   encore run
   ```
   *Note: If you change database schemas, generate new `sqlc` models using `make gen` from the `/backend` directory.*

4. **(Optional) Gotenberg PDF Generator**:
   If you are testing PDF generation, you will need a running Gotenberg container. Open a new terminal and run:
   ```bash
   docker run --rm -p 8080:3000 gotenberg/gotenberg:8-chromium
   ```

### 2. Frontend & Payload CMS

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. **Environment Variables**:
   Create a `.env` file in the root of the frontend folder and supply the following variables needed by Payload CMS (and potentially others):
   ```env
   # Example values
   PAYLOAD_SECRET=your_super_secret_payload_string
   DATABASE_URL=postgres://user:password@localhost:5432/your_local_db
   ```
   *(Note: Encore spins up dynamic ports for DBs. You may need to obtain the local connection string of Payload's designated database by running `encore db conn-uri <dbname>` inside `/backend` if it uses an Encore DB, or provide an external Neon DB string).*
3. **Install Dependencies**:
   ```bash
   pnpm install
   ```
4. **Run the Development Server**:
   ```bash
   pnpm dev
   ```

*(Payload CMS will automatically start up as part of the Next.js process).*

---

## AI Agents & Development Tools

This project is configured to work with AI coding assistants and relies on the Model Context Protocol (MCP) to supply backend intelligence to your assistant.

- **Project Conventions**: `CLAUDE.md` is the entry point for AI assistants. Detailed codebase rules, Go/Encore practices and frontend routing conventions live in `context/project-conventions.md`.
- **Encore MCP Server**: Configured in `.vscode/mcp.json`, the project uses the `encore-mcp` server (invoked via `encore mcp run --app=ai-booking-travel-bo22`). This enables your IDE's AI assistant to automatically understand the Encore application architecture, call local endpoints, query databases, and read local API specs.

---

## Environments & CI/CD

Currently, the application relies on the **Local** environment. As the project evolves, the deployment strategy will follow this structure:

### Environments

- **Stage**: This environment updates automatically on every merge to the `dev` branch.
- **Production**: This environment updates automatically on every merge to the `main` branch.

### Deployment Architecture

- **Backend**: Hosted on **Encore Cloud**, which automatically analyzes the codebase and provisions the necessary cloud infrastructure (databases, pub/sub, cron jobs, serverless functions, etc.).
- **Frontend & CMS**: Deployed via **Vercel** for edge execution and scalable serverless rendering of the Next.js App Router.
- **Frontend Database**: **Neon DB** (Serverless Postgres) is used by Payload CMS in deployed environments.
- **CI/CD Pipeline**: Handled via **GitHub Actions** (to orchestrate the respective deployment flows to Encore and Vercel simultaneously).
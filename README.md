
# Transaction Service

This is a RESTful web service for managing transactions. It supports storing transactions, retrieving them by ID, finding all transactions of a specific type, and calculating sums for linked transactions.

## Prerequisites

- Docker and Docker Compose installed on your system.

## Running the Project

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd transaction-service
   ```

2. Start the application:
   ```bash
   docker-compose up --build
   ```

3. The service will be available at `http://localhost:8080`.

## API Endpoints with `curl` Examples

1. **Create a Transaction**
   ```bash
   curl -X PUT http://localhost:8080/transactionservice/transaction/10 \
   -H "Content-Type: application/json" \
   -d '{"amount": 5000, "type": "cars"}'
   ```

2. **Get a Transaction**
   ```bash
   curl -X GET http://localhost:8080/transactionservice/transaction/10
   ```

3. **Get Transactions by Type**
   ```bash
   curl -X GET http://localhost:8080/transactionservice/types/cars
   ```

4. **Get Transaction Sum**
   ```bash
   curl -X GET http://localhost:8080/transactionservice/sum/10
   ```

## Database Schema

The service uses the following PostgreSQL table for storing transactions:

```sql
CREATE TABLE public.transactions (
	id int8 NOT NULL,
	amount numeric(20, 2) NOT NULL,
	"type" varchar(255) NOT NULL,
	parent_id int8 NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	CONSTRAINT transactions_pkey PRIMARY KEY (id),
	CONSTRAINT fk_transactions_children FOREIGN KEY (parent_id) REFERENCES public.transactions(id)
);

CREATE INDEX idx_transactions_parent_id ON public.transactions USING btree (parent_id);
CREATE INDEX idx_transactions_type ON public.transactions USING btree (type);
```

### Schema Details
- **`id`**: Unique identifier for each transaction.
- **`amount`**: The transaction amount (up to 20 digits with 2 decimal places).
- **`type`**: Type/category of the transaction (e.g., "cars", "shopping").
- **`parent_id`**: Optional parent transaction ID, allowing hierarchical relationships.
- **`created_at`** & **`updated_at`**: Timestamps for tracking record creation and updates.


## Environment Variables

These environment variables are pre-configured in `docker-compose.yml`:
- `DB_HOST`: Database host (default: `db`)
- `DB_USER`: Database username (default: `postgres`)
- `DB_PASSWORD`: Database password (default: `postgres`)
- `DB_NAME`: Database name (default: `transaction_db`)
- `DB_PORT`: Database port (default: `5432`)

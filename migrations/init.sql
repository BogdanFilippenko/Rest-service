CREATE TABLE IF NOT EXISTS subscription (
id SERIAL PRIMARY KEY,
service_name VARCHAR(300) NOT NULL,
price INTEGER NOT NULL,
user_id UUID NOT NULL,
start_data TEXT NOT NULL,
end_data TEXT

);
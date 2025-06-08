-- Drop existing tables if they exist
DROP TABLE IF EXISTS tweets;

-- Create tweets table
CREATE TABLE tweets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id for faster lookups
CREATE INDEX idx_tweets_user_id ON tweets(user_id);

-- Create index on created_at for chronological queries
CREATE INDEX idx_tweets_created_at ON tweets(created_at); 
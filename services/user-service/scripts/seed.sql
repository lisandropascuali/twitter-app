-- Insert mock users
INSERT INTO users (id, username) VALUES
    ('123e4567-e89b-12d3-a456-426614174000', 'john_doe'),
    ('223e4567-e89b-12d3-a456-426614174000', 'jane_smith'),
    ('323e4567-e89b-12d3-a456-426614174000', 'bob_wilson'),
    ('423e4567-e89b-12d3-a456-426614174000', 'alice_johnson'),
    ('523e4567-e89b-12d3-a456-426614174000', 'charlie_brown')
ON CONFLICT (id) DO NOTHING;

-- Insert mock follow relationships
INSERT INTO user_follows (follower_id, followed_id) VALUES
    -- John follows Jane and Bob
    ('123e4567-e89b-12d3-a456-426614174000', '223e4567-e89b-12d3-a456-426614174000'),
    ('123e4567-e89b-12d3-a456-426614174000', '323e4567-e89b-12d3-a456-426614174000'),
    
    -- Jane follows Alice and Charlie
    ('223e4567-e89b-12d3-a456-426614174000', '423e4567-e89b-12d3-a456-426614174000'),
    ('223e4567-e89b-12d3-a456-426614174000', '523e4567-e89b-12d3-a456-426614174000'),
    
    -- Bob follows John and Alice
    ('323e4567-e89b-12d3-a456-426614174000', '123e4567-e89b-12d3-a456-426614174000'),
    ('323e4567-e89b-12d3-a456-426614174000', '423e4567-e89b-12d3-a456-426614174000'),
    
    -- Alice follows Jane and Charlie
    ('423e4567-e89b-12d3-a456-426614174000', '223e4567-e89b-12d3-a456-426614174000'),
    ('423e4567-e89b-12d3-a456-426614174000', '523e4567-e89b-12d3-a456-426614174000'),
    
    -- Charlie follows John and Bob
    ('523e4567-e89b-12d3-a456-426614174000', '123e4567-e89b-12d3-a456-426614174000'),
    ('523e4567-e89b-12d3-a456-426614174000', '323e4567-e89b-12d3-a456-426614174000')
ON CONFLICT (follower_id, followed_id) DO NOTHING; 
-- Clear existing data
TRUNCATE TABLE tweets;

-- Insert sample tweets
INSERT INTO tweets (user_id, content) VALUES
    (1, 'Just launched my new project! #coding #golang'),
    (1, 'Learning about microservices architecture today. Very interesting stuff!'),
    (2, 'Beautiful day for coding! ☀️'),
    (2, 'Just finished implementing a new feature. Time for a coffee break! ☕'),
    (3, 'Working on some exciting new features for our application!'),
    (3, 'The power of Go and PostgreSQL together is amazing!'),
    (4, 'Debugging is like being the detective in a crime movie where you are also the murderer.'),
    (4, 'Code review time! Always learning something new from the team.'),
    (5, 'Just deployed our latest update. Everything running smoothly!'),
    (5, 'Remember: The best code is the code you don''t have to write.'); 
ALTER TABLE posts
ADD COLUMN user_id BIGINT;

-- Filling the existing rows with user_id
UPDATE posts SET user_id = 1;

-- Making the column to be NOT NULL
ALTER TABLE posts
ALTER COLUMN user_id SET NOT NULL;

ALTER TABLE posts
ADD CONSTRAINT posts_user_id_fk
FOREIGN KEY (user_id) REFERENCES users (id)
ON DELETE CASCADE;

CREATE INDEX idx_posts_user_id ON posts (user_id)
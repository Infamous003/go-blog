UPDATE posts SET subtitle = '' WHERE subtitle IS NULL;

ALTER TABLE posts 
ALTER COLUMN subtitle SET NOT NULL,
ALTER COLUMN subtitle SET DEFAULT '';

-- updating the already existning posts slugs with a unique value
-- concatenating 'post-' with their ids, which will be unique for all slugs
UPDATE posts
SET slug = 'post-' || id
WHERE slug IS NULL OR slug = '';

UPDATE posts SET slug = '' WHERE slug IS NULL;

ALTER TABLE posts
ALTER COLUMN slug SET NOT NULL;
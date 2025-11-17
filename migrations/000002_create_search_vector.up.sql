-- Adding the search vector column
ALTER TABLE posts ADD COLUMN search_vector tsvector;

-- Creating a trigger function that updates the search_vector to have vector on title, subtitle and content
-- title gets the highest priority with 'A', and content gets lowest with 'C'
CREATE FUNCTION posts_search_vector_update() RETURNS trigger AS $$
BEGIN
  NEW.search_vector :=
    setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(NEW.subtitle, '')), 'B') ||
    setweight(to_tsvector('english', coalesce(NEW.content, '')), 'C');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- coalesce, makes it so that instead of NULL, an empty string is passed


-- Creating the trigger(when the trigger function should work)
CREATE TRIGGER posts_search_vector_update
BEFORE INSERT OR UPDATE ON posts
FOR EACH ROW EXECUTE FUNCTION posts_search_vector_update();

-- Backfilling the old rows
UPDATE posts SET search_vector = 
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(subtitle, '')), 'B') ||
    setweight(to_tsvector('english', coalesce(content, '')), 'C');

-- Creating a GIN index
CREATE INDEX posts_search_vector_idx
ON posts USING GIN (search_vector);
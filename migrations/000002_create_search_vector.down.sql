DROP INDEX IF EXISTS posts_search_vector_idx;

DROP TRIGGER IF EXISTS posts_search_vector_update ON posts;

DROP FUNCTION IF EXISTS posts_search_vector_update();

ALTER TABLE posts DROP COLUMN IF EXISTS search_vector;

-- we first drop the index, then the trigger, followed by the trigger function and then the column
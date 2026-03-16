-- +goose Up
ALTER TABLE courses 
ALTER COLUMN teacher_id DROP NOT NULL;

ALTER TABLE courses 
DROP CONSTRAINT courses_teacher_id_fkey;

ALTER TABLE courses 
ADD CONSTRAINT courses_teacher_id_fkey 
FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE courses 
DROP CONSTRAINT courses_teacher_id_fkey;

ALTER TABLE courses 
ADD CONSTRAINT courses_teacher_id_fkey 
FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE;

ALTER TABLE courses 
ALTER COLUMN teacher_id SET NOT NULL;

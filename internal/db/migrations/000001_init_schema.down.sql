-- 1. Удаление таблиц
DROP TABLE users CASCADE;
DROP TABLE families CASCADE;
DROP TABLE family_members CASCADE;
DROP TABLE shopping_lists CASCADE;
DROP TABLE shopping_items CASCADE;
DROP TABLE tasks CASCADE;
DROP TABLE family_events CASCADE;
DROP TABLE expenses CASCADE;
DROP TABLE expense_splits CASCADE;

-- 2. Удаление типов
DROP TYPE role_members;
DROP TYPE priority_task;
DROP TYPE status_task;
DROP TYPE split_type_expenses;
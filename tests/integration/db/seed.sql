INSERT INTO monsters (
    id,name, monster_def_id, experience
) VALUES (1,'Schiggo', 1, 0);
INSERT INTO monsters (
    id,name, monster_def_id, experience
) VALUES (2,'Bisa', 1, 0);
INSERT INTO monsters (
    id,name, monster_def_id, experience
) VALUES (3,'Glumander', 1, 0);
INSERT INTO monsters (
    id,name, monster_def_id, experience
) VALUES (4,'Schiggo', 2, 0);


INSERT INTO users (id,username, external_id ) VALUES (1,'felix', '12345678');
INSERT INTO inventory_items (user_id, item_def_id, quantity) VALUES (1,'stone',100);

-- Add an old job that needs to be canceled
-- INSERT INTO jobs (user_id, job_def_id, started_at, updated_at, job_type) VALUES (1, 'stoneBar', '2022-01-01 00:00:00', '2022-01-01 00:00:00', 'processing');
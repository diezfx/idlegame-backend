BEGIN;
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    external_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE jobs(
        id SERIAL PRIMARY KEY,
        job_def_id TEXT NOT NULL,
        user_id integer NOT NULL,
        started_at TIMESTAMPTZ NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL,
        job_type TEXT NOT NULL,

        CONSTRAINT fk_user
          FOREIGN KEY(user_id) 
          REFERENCES users(id)
);

CREATE TABLE job_monsters (
    monster_id integer PRIMARY KEY,
    job_id integer NOT NULL,
    CONSTRAINT fk_job
      FOREIGN KEY(job_id) 
	  REFERENCES jobs(id)
);

CREATE TABLE gathering_jobs (
    job_id SERIAL PRIMARY KEY,
    gathering_type TEXT NOT NULL,
    CONSTRAINT fk_job
      FOREIGN KEY(job_id) 
	    REFERENCES jobs(id)
);



CREATE TABLE monsters(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    monster_def_id integer NOT NULL,
    experience integer NOT NULL
);


CREATE TABLE inventory_items(
    user_id integer NOT NULL,
    item_def_id TEXT NOT NULL,
    quantity integer NOT NULL,

    PRIMARY KEY(user_id, item_def_id),
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
      REFERENCES users(id)
);



END;

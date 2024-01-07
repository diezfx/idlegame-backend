BEGIN;
CREATE TABLE jobs(
        id SERIAL PRIMARY KEY,
        started_at TIMESTAMP NOT NULL,
        job_type TEXT
);

CREATE TABLE job_monsters (
    monster_id integer PRIMARY KEY,
    job_id integer NOT NULL,
    CONSTRAINT fk_job
      FOREIGN KEY(job_id) 
	  REFERENCES jobs(id)
);

CREATE TABLE woodcutting_jobs (
    job_id SERIAL PRIMARY KEY,
    tree_type TEXT NOT NULL,
    CONSTRAINT fk_job
      FOREIGN KEY(job_id) 
	  REFERENCES jobs(id)
);

END;

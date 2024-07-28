-- +goose Up

CREATE TABLE "events" (
                          "id" serial PRIMARY KEY,
                          "title" varchar NOT NULL,
                          "time" time NOT NULL,
                          "created_at" timestamptz DEFAULT (now())
);

INSERT INTO "events" ("title", "time")
VALUES
    ('Morning Yoga', '06:00:00'),
    ('Team Standup Meeting', '09:00:00'),
    ('Brunch with Clients', '11:00:00'),
    ('Product Launch Review', '14:00:00'),
    ('Workshop: Innovation Strategies', '16:00:00'),
    ('Evening Run', '18:00:00'),
    ('Dinner with Partners', '19:30:00'),
    ('Networking Event', '20:30:00'),
    ('Late Night Coding', '22:00:00'),
    ('Midnight Meditation', '23:59:00');

-- +goose Down

DROP TABLE "events";

-- +goose Up

CREATE TABLE "events" (
                          "id" serial PRIMARY KEY,
                          "title" varchar NOT NULL,
                          "time" timestamptz NOT NULL,
                          "created_at" timestamptz DEFAULT (now())
);

INSERT INTO "events" ("title", "time")
VALUES
    ('Morning Yoga', '2024-08-06 06:00:00+00'),
    ('Team Standup Meeting', '2024-08-06 09:00:00+00'),
    ('Brunch with Clients', '2024-08-06 11:00:00+00'),
    ('Product Launch Review', '2024-08-06 14:00:00+00'),
    ('Workshop: Innovation Strategies', '2024-08-06 16:00:00+00'),
    ('Evening Run', '2024-08-06 18:00:00+00'),
    ('Dinner with Partners', '2024-08-06 19:30:00+00'),
    ('Networking Event', '2024-08-06 20:30:00+00'),
    ('Late Night Coding', '2024-08-06 22:00:00+00'),
    ('Midnight Meditation', '2024-08-06 23:59:00+00');

-- +goose Down

DROP TABLE "events";

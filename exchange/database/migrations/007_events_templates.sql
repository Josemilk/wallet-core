CREATE TABLE IF NOT EXISTS event_templates (
 id UUID PRIMARY KEY,
 slug TEXT NOT NULL UNIQUE,
 title TEXT NOT NULL,
 event_type TEXT NOT NULL CHECK (event_type IN ('price_prediction','sports_result','sports_score','seasonal_challenge','skill_quiz')),
 season_key TEXT NOT NULL,
 icon TEXT NOT NULL,
 description TEXT NOT NULL,
 active BOOLEAN NOT NULL DEFAULT TRUE,
 sort_order INTEGER NOT NULL DEFAULT 0,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS events (
 id UUID PRIMARY KEY,
 template_id UUID NOT NULL REFERENCES event_templates(id),
 title TEXT NOT NULL,
 asset TEXT,
 sport TEXT,
 league TEXT,
 event_data JSONB NOT NULL DEFAULT '{}'::jsonb,
 entry_price NUMERIC(78,0) NOT NULL DEFAULT 0 CHECK (entry_price >= 0),
 prize_amount NUMERIC(78,0) NOT NULL CHECK (prize_amount >= 0),
 prize_asset TEXT NOT NULL,
 starts_at TIMESTAMPTZ NOT NULL,
 ends_at TIMESTAMPTZ NOT NULL,
 status TEXT NOT NULL CHECK (status IN ('draft','scheduled','open','locked','settling','settled','cancelled')),
 created_by UUID NOT NULL REFERENCES users(id),
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 CHECK (ends_at > starts_at)
);
CREATE INDEX IF NOT EXISTS events_carousel_idx ON events(status, starts_at, ends_at);
CREATE INDEX IF NOT EXISTS events_template_idx ON events(template_id, season_key);

CREATE TABLE IF NOT EXISTS event_entries (
 id UUID PRIMARY KEY,
 event_id UUID NOT NULL REFERENCES events(id),
 user_id UUID NOT NULL REFERENCES users(id),
 selection JSONB NOT NULL,
 entry_amount NUMERIC(78,0) NOT NULL CHECK (entry_amount >= 0),
 prize_amount NUMERIC(78,0) NOT NULL DEFAULT 0 CHECK (prize_amount >= 0),
 status TEXT NOT NULL CHECK (status IN ('entered','won','lost','refunded','cancelled')),
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 UNIQUE(event_id,user_id)
);

INSERT INTO event_templates(id,slug,title,event_type,season_key,icon,description,sort_order)
VALUES
('00000000-0000-0000-0000-000000000101','crypto-price-prediction','Crypto Price Prediction','price_prediction','core','chart','Predict whether a selected asset reaches the published target condition before the event closes.',1),
('00000000-0000-0000-0000-000000000102','football-result','Football Result','sports_result','football','football','Predict the result of a scheduled football match using the published fixture.',2),
('00000000-0000-0000-0000-000000000103','basketball-result','Basketball Result','sports_result','basketball','basketball','Predict the result of a scheduled basketball game.',3),
('00000000-0000-0000-0000-000000000104','tennis-match','Tennis Match','sports_result','tennis','tennis','Predict the winner of a scheduled tennis match.',4),
('00000000-0000-0000-0000-000000000105','seasonal-quiz','Seasonal Skill Quiz','skill_quiz','seasonal','quiz','Answer a fixed set of prepared questions before the event closes.',5)
ON CONFLICT(slug) DO NOTHING;

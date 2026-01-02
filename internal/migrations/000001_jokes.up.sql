CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS jokes (
    id bigserial primary key,
    external_id varchar(50) unique not null,
    joke_url varchar unique not null,
    content text not null,
    created_at timestamp not null default current_timestamp
);
CREATE INDEX IF NOT EXISTS idx_jokes_content ON jokes USING gin(to_tsvector('simple', content));

CREATE TABLE IF NOT EXISTS users (
    id bigserial primary key,
    hashed_pw bytea not null,
    created_at timestamp not null default current_timestamp,
    email citext unique not null
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea primary key,
    user_id bigint not null references users on delete cascade,
    expires_at timestamp not null
);

-- seed database
INSERT INTO jokes (external_id, joke_url, content)
VALUES
('c-3yrrglr0ouxifeo2rzsw', 'https://api.chucknorris.io/jokes/c-3yrrglr0ouxifeo2rzsw', 'In the movie The Matrix, Chuck Norris is the Matrix. If you pay close attention in the green falling code scenes, you can make out the faint texture of his beard.'),
('7w5sd_KSTgGwsVCQc3U8KA', 'https://api.chucknorris.io/jokes/7w5sd_KSTgGwsVCQc3U8KA', 'Chuck Norris was origionally to be cast as Neo in The Matrix, but when they told him about the powers he would recieve, he laughed and said im Chuck Norris! He then proceeded to roundhouse kick the Warner Brother representatives and flying off to drink some beer and kick some random person''s ass'),
('KzX7anzHQMKJuViPUmBMVg', 'https://api.chucknorris.io/jokes/KzX7anzHQMKJuViPUmBMVg', 'In the Matrix, how can you tell it from the real world? Well, in the real world, all leading roles played by different actors who always get the girl are actually Chuck Norris'),
('-jMF0mF7Rcu6tUiIZ49pkw', 'https://api.chucknorris.io/jokes/-jMF0mF7Rcu6tUiIZ49pkw', 'Chuck Norris took the red pill, the blue pill, 4 horse tranquilizers and a handful of rat poison, washed it down with a bottle of bourbon, and roundhouse kicked Morpheus out of the Matrix. Norris then unplugged himself out ofthe Matrix by simply flexing, broke into Zion and roundhouse kicked Morpheus again, just because he''s Chuck Norris.');
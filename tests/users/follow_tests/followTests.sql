INSERT INTO users (
    id,
    email,
    country_code,
    verified
) VALUES (
    1,
    'erdajtsopjani.tech@gmail.com',
    'RKS',
    TRUE
);


INSERT INTO users (
    id,
    email,
    country_code,
    verified
) VALUES (
    2,
    'testexample@gmail.com',
    'RKS',
    TRUE
);

INSERT INTO users (
    id,
    email,
    country_code,
    verified
) VALUES (
    3,
    'testmailmail@mail.com',
    'RKS',
    TRUE
);

INSERT INTO user_tokens(
    user_id,
    token
) VALUES (
    1,
    'token1'
);

INSERT INTO user_tokens(
    user_id,
    token
) VALUES (
    2,
    'token2'
);

INSERT INTO user_tokens(
    user_id,
    token
) VALUES (
    3,
    'token3'
);

INSERT INTO follows(
    follower_id,
    following_id
) VALUES (
    2,
    3
);

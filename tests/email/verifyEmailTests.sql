INSERT INTO users (
    id,
    email,
    country_code,
    verified
) VALUES (
    2,
    'erdajtsopjani.tech@gmail.com',
    'RKS',
    TRUE
);

INSERT INTO verification_Tokens (
    token,
    user_id
) VALUES (
    '123456',
    2
);

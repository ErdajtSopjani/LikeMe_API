INSERT INTO users (
    id,
    email,
    country_code,
    verified,
    created_at
) VALUES (
    1,
    'erdajtsopjani.tech@gmail.com',
    'RKS',
    TRUE,
    '2024-09-17 00:13:33'
);

INSERT INTO user_tokens (
    user_id,
    token,
    created_at
) VALUES (
    1,
    'jUy2Iti6p3GqQxp0TjwrGA==',
    '2024-09-17 00:13:33'
);

INSERT INTO user (
    id,
    email,
    country_code,
    verified,
) VALUES (
    2,
    'test@mail.com',
    'RKS',
    TRUE
);

INSERT INTO user_tokens (
    user_id,
    token,
    created_at
) VALUES (
    2,
    'gsa2I2kja3GqQxp0TKhj1A==',
    '2024-09-17 00:13:33'
);

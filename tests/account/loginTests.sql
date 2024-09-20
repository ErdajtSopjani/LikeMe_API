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

INSERT INTO "two_factors" (
    "code",
    "user_id",
    "created_at",
    "expires_at"
) VALUES (
    5692124, -- valid code: 5692124 invalid code: 162508
    '2',
    '2024-09-17 00:13:33',
    '2028-09-17 00:43:33'
) RETURNING "id";

INSERT INTO "two_factors" (
    "code",
    "user_id",
    "created_at",
    "expires_at"
) VALUES (
    162508, -- valid code: 152508 invalid code: 162508
    '2',
    '2020-09-17 00:13:33',
    '2022-09-17 00:43:33'
) RETURNING "id";

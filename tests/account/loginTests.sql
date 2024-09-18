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
    CURRENT_TIMESTAMP
);

INSERT INTO "two_factors" (
    "code",
    "user_id",
    "created_at",
    "expires_at"
) VALUES (
    '5692124', -- valid code: 5692124 invalid code: 162508
    '1',
    '2024-09-17 00:13:33',
    '2028-09-17 00:43:33'
) RETURNING "id";

INSERT INTO "two_factors" (
    "code",
    "user_id",
    "created_at",
    "expires_at"
) VALUES (
    '162508', -- valid code: 152508 invalid code: 162508
    '1',
    '2020-09-17 00:13:33',
    '2022-09-17 00:43:33'
) RETURNING "id";

INSERT INTO "two_factors" (
    "code",
    "user_id",
    "created_at",
    "expires_at"
) VALUES (
    '112308', -- valid code: 152508 invalid code: 162508
    '4',
    '2024-09-17 00:13:33',
    '2028-09-17 00:43:33'
) RETURNING "id";

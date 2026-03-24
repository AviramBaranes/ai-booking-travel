create unlogged table available_plans_snapshots (
    id bigserial primary key,
    created_at timestamptz not null default now (),
    driver_age text not null,
    pickup_date text not null,
    pickup_time text not null,
    return_date text not null,
    return_time text not null,
    country_code text not null,
    plans json not null
);

create index available_plans_snapshots_created_at_idx on available_plans_snapshots (created_at);
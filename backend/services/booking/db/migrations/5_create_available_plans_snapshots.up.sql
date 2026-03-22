create table
    available_plans_snapshots (
        id bigserial primary key,
        created_at timestamptz not null default now (),
        plans json not null
    );

create index available_plans_snapshots_created_at_idx on available_plans_snapshots (created_at);
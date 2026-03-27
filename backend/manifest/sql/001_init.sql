create table if not exists users (
    id bigserial primary key,
    username varchar(64) not null unique,
    password varchar(255) not null,
    role varchar(32) not null default 'admin',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists reports (
    id bigserial primary key,
    trade_date varchar(16) not null,
    stock_code varchar(16) not null,
    stock_name varchar(64) not null,
    summary text not null,
    risk_level varchar(32) not null,
    recommendation varchar(32) not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists model_configs (
    id bigserial primary key,
    base_url varchar(255) not null,
    api_key varchar(512) not null,
    model varchar(128) not null,
    checked_at timestamptz null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_reports_trade_date on reports(trade_date);
create index if not exists idx_reports_stock_code on reports(stock_code);

insert into users (username, password, role)
values ('admin', 'admin123456', 'admin')
on conflict (username) do nothing;

insert into reports (trade_date, stock_code, stock_name, summary, risk_level, recommendation)
values
('2026-03-27', '600519', '贵州茅台', '高位震荡，短期关注量价配合与机构资金延续性。', 'medium', 'watch'),
('2026-03-27', '300750', '宁德时代', '趋势仍强，但需关注板块轮动与估值消化节奏。', 'medium', 'hold')
on conflict do nothing;

create table sent_alerts if not exists {
    id varchar(255) primary key,
    club_id varchar(20) not null,
    club_name varchar(255) not null,
    court_id varchar(20) not null,
    court_name varchar(255) not null,
    slot_date date not null,
    start_time date not null,
    end_time date not null,
    created_at timestamp not null default NOW(),
    updated_at timestamp not null default NOW()
}
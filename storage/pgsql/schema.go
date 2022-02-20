package pgsql

const schema = `
create table if not exists users (
	id bigserial not null,
	login varchar (255) unique not null,
	password varchar (255) not null,
	balance decimal (16,2) not null default 0,
	outcome decimal (16,2)  not null default 0
);

create table if not exists orders (
    id varchar (255) not null primary key,
    user_id integer not null,
    status varchar (50) not null default 'NEW',
    accrual decimal (16,2) not null default 0,
    uploaded_at timestamp(0) with time zone not null default current_timestamp
);

create table if not exists transactions (
    id bigserial not null,
    user_id integer not null,
    amount decimal (16, 2),
    created_at timestamp(0) with time zone not null default current_timestamp
);
`

// Migrate exec sql code
func (s Storage) Migrate() error {
	_, err := s.db.Exec(schema)

	return err
}

func (s Storage) Drop() {
	s.db.Exec("truncate table transactions;")
	s.db.Exec("truncate table orders;")
	s.db.Exec("truncate table users;")
}

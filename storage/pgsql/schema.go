package pgsql

const schema = `
create table if not exists users (
id bigserial not null,
login varchar(255) unique not null,
password varchar(255) not null
);
`

// Migrate exec sql code
func (s Storage) Migrate() error {
	_, err := s.db.Exec(schema)

	return err
}

// Clear truncates all tables for testing.
func (s Storage) Clear() {
	s.db.Exec("truncate table users;")
}

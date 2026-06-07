CREATE TABLE users (
    id varchar(36) PRIMARY KEY,
    display_name varchar(255) NOT NULL,
    username varchar(255) UNIQUE NOT NULL,
    password varchar(2048) NOT NULL,
    email varchar(2048) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

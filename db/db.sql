CREATE TABLE "User" (
    user_id bigserial primary key,
    user_name varchar(40) NOT NULL,
    user_pw varchar(50) NOT NULL,
    date_joined timestamp default NULL
);

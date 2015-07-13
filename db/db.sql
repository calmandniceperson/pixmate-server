/* CREATE SCHEMA imgturtle FOR USER/ROLE imgturtle */
CREATE SCHEMA IF NOT EXISTS imgturtle AUTHORIZATION imgturtle;

CREATE TABLE IF NOT EXISTS imgturtle.User (
    user_id bigserial not null,
    user_name varchar(40) not null,
    user_pw varchar(50) not null,
    date_joined timestamp default null,
    PRIMARY KEY(user_id, user_name)
);

CREATE TABLE IF NOT EXISTS imgturtle.Img (
  image_id varchar(100) primary key,
  image_title varchar(100) not null,
  image_path varchar(100) not null,
  image_desc varchar(400),
  date_uploaded timestamp not null,
  uploader_id bigserial not null,
  uploader_name varchar(40) not null,
  FOREIGN KEY (uploader_id, uploader_name) REFERENCES imgturtle.User(user_id, user_name)
);

CREATE TABLE IF NOT EXISTS imgturtle.Comment (
  comment_id bigserial primary key,
  image_id varchar(100) REFERENCES imgturtle.Img(image_id),
  comment_text varchar(100) not null,
  creator_id bigserial not null,
  creator_name varchar(40) not null,
  FOREIGN KEY (creator_id, creator_name) REFERENCES imgturtle.User(user_id, user_name)
);

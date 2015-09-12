/* CREATE SCHEMA imgturtle FOR USER/ROLE imgturtle */
CREATE SCHEMA imgturtle AUTHORIZATION imgturtle;

/*DROP TABLE IF EXISTS imgturtle.User;*/
CREATE TABLE IF NOT EXISTS imgturtle.User (
    user_id text default md5(random()::text),
    user_name text not null,
    user_pw bytea not null,
    user_email text not null,
    user_hash text not null,
    user_p_pic_path text,
    date_joined timestamp default now(),
    PRIMARY KEY(user_id, user_name)
);

CREATE TABLE IF NOT EXISTS imgturtle.Following (
  id serial PRIMARY KEY,
  user1_id text not null,
  user1_name text not null,
  user2_id text not null,
  user2_name text not null,
  status int not null default 1,
  FOREIGN KEY (user1_id, user1_name) REFERENCES imgturtle.User(user_id, user_name),
  FOREIGN KEY (user2_id, user2_name) REFERENCES imgturtle.User(user_id, user_name)
);

/*DROP TABLE IF EXISTS imgturtle.Img;*/
CREATE TABLE IF NOT EXISTS imgturtle.Img (
  image_id text primary key /*default md5(random()::text)*/,
  image_title text not null,
  image_path text not null,
  image_f_ext text not null, /* file extension */
  image_desc text,
  date_uploaded timestamp not null default now(),
  uploader_id text,
  uploader_name text,
  FOREIGN KEY (uploader_id, uploader_name) REFERENCES imgturtle.User(user_id, user_name)
);

/*DROP TABLE IF EXISTS imgturtle.Comment;*/
CREATE TABLE IF NOT EXISTS imgturtle.Comment (
  comment_id bigserial primary key,
  image_id text REFERENCES imgturtle.Img(image_id),
  comment_text text not null,
  creator_id text not null,
  creator_name text not null,
  FOREIGN KEY (creator_id, creator_name) REFERENCES imgturtle.User(user_id, user_name)
);

/* CREATE SCHEMA pixmate FOR USER/ROLE pixmate */
CREATE SCHEMA pixmate AUTHORIZATION pixmate;

/*DROP TABLE IF EXISTS pixmate.User;*/
CREATE TABLE IF NOT EXISTS pixmate.User (
    user_name text not null,
    user_pw bytea not null,
    user_email text not null,
    user_hash text not null,
    user_p_pic_path text,
    date_joined timestamp default now(),
    private int default 1, /* 0=public, 1=friends*/
    PRIMARY KEY(user_name)
);

CREATE TABLE IF NOT EXISTS pixmate.Following (
  user1_name text not null,
  user2_name text not null,
  status int not null default 1, /* 1=not yet accepted, 0=accepted*/
  PRIMARY KEY(user1_name, user2_name),
  FOREIGN KEY (user1_name) REFERENCES pixmate.User(user_name),
  FOREIGN KEY (user2_name) REFERENCES pixmate.User(user_name)
);

/*DROP TABLE IF EXISTS pixmate.Img;*/
CREATE TABLE IF NOT EXISTS pixmate.Img (
  image_id text primary key /*default md5(random()::text)*/,
  image_title text not null,
  image_path text not null,
  image_f_ext text not null, /* file extension */
  image_desc text,
  date_uploaded timestamp not null default now(),
  uploader_name text,
  FOREIGN KEY (uploader_name) REFERENCES pixmate.User(user_name)
);

/*DROP TABLE IF EXISTS pixmate.Comment;*/
CREATE TABLE IF NOT EXISTS pixmate.Comment (
  comment_id bigserial primary key,
  image_id text REFERENCES pixmate.Img(image_id),
  comment_text text not null,
  creator_name text not null,
  FOREIGN KEY (creator_name) REFERENCES pixmate.User(user_name)
);

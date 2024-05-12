create table novel_user (
	id INT auto_increment,
	email varchar(64) not null comment '이메일',
	credential varchar(128) comment '비밀번호',
	name varchar(16) not null comment '이름',
	created_at timestamp not null default NOW() comment '생성날짜',
	udpated_at timestamp not null default NOW() comment '최근접속 날짜',
	provider varchar(16) not null comment '로그인 플랫폼',
	primary key (id),
	unique key (email, provider)
);

create table tag (
	id INT auto_increment,
	name varchar(255) not null,
	primary key (id)
);

create table novel (
	id INT auto_increment,
	novel_user_id INT not null comment '유저 아이디',
	tag_id INT not null comment '태그 아이디',
	title varchar(255) not null comment '제목',
	description varchar(1024) not null comment '설명',
	total_views INT not null default 0 comment '총조회수',
	total_recommandation INT not null default 0 comment '추천수',
	create_at timestamp not null default NOW() comment '생성날짜',
	last_updated_at timestamp not null default NOW() comment '최종업데이트',
	last_updated_episode INT not null default 0 comment '최종화수',
	primary key (id),
	foreign key(novel_user_id) references novel_user(id),
	foreign key(tag_id) references tag(id)
);

create table novel_content (
	id int auto_increment,
	novel_id int not null comment '소설 아이디',
	content text not null comment '내용',
	views int not null default 0 comment '조회수',
	created_at timestamp not null default NOW() comment '업로드날짜',
	content_length int not null default 0 comment '화',
	recommandation int not null default 0 comment '추천수',
	primary key (id),
	foreign key(novel_id) references novel(id)
);

create table novel_comment (
	id int auto_increment,
	novel_content_id int not null comment '소설내용 아이디',
	novel_user_id int not null comment '유저 아이디',
	comment_id int not null comment '상위 댓글 아이디',
	content varchar(1024) not null,
	primary key (id),
	foreign key(novel_content_id) references novel_content(id),
	foreign key(novel_user_id) references novel_user(id),
	foreign key(comment_id) references novel_comment(id)
);
create table buy_wagers
(
	id int auto_increment,
	wager_id int null,
	buying_price float not null,
	bought_at datetime not null,
	constraint buy_wager_id_uindex
		unique (id),
	constraint buy_wager_wagers_id_fk
		foreign key (wager_id) references wagers (id)
			on delete cascade
);

alter table buy_wagers
	add primary key (id);


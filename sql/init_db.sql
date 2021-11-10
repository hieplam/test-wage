create table wagers
(
	id int auto_increment,
	total_wager_value int not null,
	odds int not null,
	selling_percentage int not null,
	selling_price float not null,
	current_selling_price float null,
	percentage_sold int null,
	amount_sold float null,
	placed_at datetime not null,
	constraint wagers_id_uindex
		unique (id)
);

alter table wagers
	add primary key (id);

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


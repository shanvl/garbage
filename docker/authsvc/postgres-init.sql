-- create role enum type
do
$$
    begin
        if not exists(select 1 from pg_type where typname = 'role') then
            create type role as enum ('admin', 'member');
        end if;
    end
$$;

-- create users table
create table if not exists users
(
    id             varchar(50) primary key,
    active         bool        not null,
    activate_token text        not null,
    email          varchar(50) not null,
    first_name     varchar(25) not null,
    last_name      varchar(25) not null,
    password_hash  text        not null,
    role           role        not null,
    text_search    tsvector generated always as (to_tsvector('simple', first_name || ' ' || last_name || ' ' ||
                                                                       email)) stored
);

create index if not exists users_text_search_idx on users using gin (text_search);
create index if not exists users_active_idx on users (active);
create index if not exists users_email_idx on users (email);
create index if not exists users_activate_token_idx on users (activate_token);

-- create clients table. Clients in a sense of browsers, apps etc
create table if not exists clients
(
    id            varchar(50) primary key,
    refresh_token text not null
);

-- create user_client table
create table if not exists user_client
(
    user_id   varchar(50) not null,
    client_id varchar(50) not null,
    primary key (user_id, client_id),
    foreign key (client_id) references clients (id)
        on delete cascade
        on update cascade,
    foreign key (user_id) references users (id)
        on delete cascade
        on update cascade
);

-- populate users
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('fb670cf2-9461-4568-990c-e9cf322f182f', 'Giralda', 'Brognot', 'gbrognot0@printfriendly.com', 'admin', true, '', 'uIO8dT5D');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('9c0d8fd1-3b79-41fd-a1dd-2f8c40e8655c', 'Hodge', 'Thiolier', 'hthiolier1@people.com.cn', 'member', true, '', 'mWr7Hn');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('fbfbedfb-41ce-45dc-b8d5-f60bbb3bb9ec', 'Leif', 'Petris', 'lpetris2@moonfruit.com', 'admin', true, '', 'BhtRpAM1lL0u');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('861a7287-a8b0-4bb1-94a7-5c37c4597945', 'Rolando', 'Brame', 'rbrame3@istockphoto.com', 'member', true, '', 'QZQEH7J');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('7b23adc3-aedd-457a-b448-da1d7141e05e', 'Gertrudis', 'McKelvie', 'gmckelvie4@netscape.com', 'member', true, '', 'gsTd0XNioqW');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('af0248ca-7b00-4d90-9cb8-344ec2e1db86', 'Zitella', 'Barthot', 'zbarthot5@simplemachines.org', 'admin', true, '', '2sPXb9fhoY');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('9a353d70-fd4f-4053-a29d-012927a83ab1', 'Eduard', 'Strand', 'estrand6@ameblo.jp', 'member', true, '', 'Zwlh4nULNbg');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('be89197a-c953-46cb-9a94-9465fb8f5558', 'Kip', 'Sawell', 'ksawell7@hc360.com', 'member', false, 'd9a202ea-6d57-4449-a7a3-e2afbee510e9', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c45de003-4372-4361-a34f-e2ffbd82b1ce', 'Davidde', 'Gissop', 'dgissop8@yellowbook.com', 'admin', true, '', '1RW8pF');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('39a69617-2c76-4939-9cc4-ccaa0d7d1c6c', 'Maisie', 'Gowans', 'mgowans9@bloglovin.com', 'admin', true, '', 'e8DLNGwUB');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('3585b2aa-2772-49f2-b001-29afaac5267c', 'Jeffrey', 'Cansdill', 'jcansdilla@google.co.uk', 'member', true, '', 'x3R6nt');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('13adb135-bf5d-41b0-917d-bf967ab6cdac', 'Patrizio', 'Kingswell', 'pkingswellb@bravesites.com', 'member', false, '3303bbd2-5615-4559-9abb-b64c02e60834', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('8bf59554-5540-450f-91b4-dc116ea25a88', 'Halley', 'Feavearyear', 'hfeavearyearc@51.la', 'admin', true, '', 'ST1vvlNfn6');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('a9058af0-20dc-4ec3-b697-f6c8ccce068d', 'Lela', 'Kindred', 'lkindredd@chron.com', 'admin', true, '', '55QMijhC');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('94cb6da0-24a2-4adc-9b79-796dc7f86671', 'Kristoforo', 'Ianson', 'kiansone@pen.io', 'admin', true, '', 'wUxiIn5IanFE');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('2d04b20c-2eef-4da0-b820-47394d64f936', 'Raddy', 'Applegarth', 'rapplegarthf@360.cn', 'member', false, '2564dc50-2e83-4020-a08a-0dbe46973a4f', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('49cd401c-6e74-42ae-9d73-499f9c585121', 'Lazar', 'Zuenelli', 'lzuenellig@ebay.com', 'member', true, '', '9EgAh8');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('82a27152-0384-49cb-956a-7c68790c2907', 'Janeva', 'Peckham', 'jpeckhamh@de.vu', 'admin', true, '', 'Rl2zzLiZ');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('89123e6a-6fc3-431d-a686-27cb7e0fbbf3', 'Rey', 'Biggans', 'rbiggansi@purevolume.com', 'admin', true, '', 'iY5LzTCM');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('016c486c-e66c-44e1-8aca-25426e367a4b', 'Rey', 'Chisnall', 'rchisnallj@hugedomains.com', 'member', false, '3ab6b0d5-12da-4c5d-89e5-3978896174c0', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('6e93d7d4-0a1c-43b9-95a7-39356e8479a0', 'Viviene', 'McNalley', 'vmcnalleyk@techcrunch.com', 'member', true, '', '7Pw7sIfWiYN');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('ffaa3762-77f3-4079-ba9c-06de10981d68', 'Rudd', 'Schulkins', 'rschulkinsl@skype.com', 'member', false, '2c832426-a168-4dd9-9242-186f85f50511', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('8a4e1756-873b-4ec2-9033-219c33733000', 'Rubie', 'Danis', 'rdanism@tumblr.com', 'admin', true, '', 'K0BkhMF');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('6ce5b5b0-7363-4511-aa9c-023968d49759', 'Jerald', 'Corrao', 'jcorraon@photobucket.com', 'member', true, '', 'ex21zLeXx');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('315c993a-40db-470b-9ba2-b147c0ab5e71', 'Elihu', 'Pauly', 'epaulyo@cocolog-nifty.com', 'member', true, '', 'OzPIa06FofN');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('33d828eb-a2f6-4093-a82b-70aa738c23f3', 'Merry', 'Zwicker', 'mzwickerp@squidoo.com', 'admin', true, '', 'bfkfQFd8Xt');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('7ac0ebea-bb0e-4106-844a-cf598e2b1807', 'Yance', 'Fulop', 'yfulopq@yelp.com', 'admin', true, '', 'KGzX8IpVvznT');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('3d5ebe9e-7b76-4ade-87d2-3c889c30f3d1', 'Elfreda', 'Crosser', 'ecrosserr@state.tx.us', 'admin', true, '', '5TN1Ox1j');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('6ae3456a-110e-499f-8890-b361042e1788', 'Matthias', 'Castana', 'mcastanas@linkedin.com', 'member', false, '0651ffbb-17b2-41ff-8f0d-5175e0287565', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('1923e608-880e-4676-a800-841c1f65e8e9', 'Gregg', 'Twentyman', 'gtwentymant@mlb.com', 'admin', true, '', '77AuN2p');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('ce05731b-2579-4a9b-a019-843f26081085', 'Jecho', 'Askie', 'jaskieu@g.co', 'member', false, '1353645d-6f11-4873-9965-748f7d9bfa8d', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('37c99d56-e09d-4986-9e16-fb077adc2b1e', 'Gunther', 'Dami', 'gdamiv@spiegel.de', 'admin', true, '', '3YVjdkLKP');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c5d87819-9592-4ccc-be14-09cb6de85d38', 'Leah', 'Georgeau', 'lgeorgeauw@blogs.com', 'admin', true, '', 'IkDsNTb');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('eae51c3b-90af-483a-a69c-d0acbeff0b4b', 'Guss', 'Bassingden', 'gbassingdenx@acquirethisname.com', 'member', false, '1c2d005e-82d1-426c-9bad-a840f055a423', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('5bab59f9-8442-40de-9462-fbc5a4863016', 'Vere', 'Hastin', 'vhastiny@sohu.com', 'admin', true, '', 'ROCWPNsMpu');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('33824cb8-8e00-4c2e-8daf-88a49ddf129d', 'Freddie', 'Kimble', 'fkimblez@fotki.com', 'member', false, 'ca2dfd78-194f-4bc4-9ba9-ce7076515248', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('cd552fd6-8dc4-4ee6-a72c-c646392047a1', 'Milicent', 'Furness', 'mfurness10@ftc.gov', 'member', false, '382306ba-1c47-4e59-865f-876b39e699ef', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('40678e6b-03ac-41de-94ee-c0a898d01356', 'Hadlee', 'McCahey', 'hmccahey11@arstechnica.com', 'member', false, '7b937d83-0f92-45dd-ba9f-ca16de16f0c9', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c55df9cd-dcfc-49ac-afa2-e30e26c5abdd', 'Jo-ann', 'Edgeley', 'jedgeley12@unesco.org', 'admin', true, '', 'p0cO7A4dv3M');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('d44e0ff4-7f94-4c9f-8733-a516465f1965', 'Nickie', 'Dives', 'ndives13@ihg.com', 'member', true, '', '4C1kCHdyb');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('08a6e711-ab56-4108-88fc-6cbb7b614f85', 'Thane', 'Gyse', 'tgyse14@ucoz.ru', 'member', false, 'e54268d1-9850-431f-89a7-e8bb164c400b', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('1345febb-32c5-4081-b292-49e10627f396', 'Kareem', 'Newvell', 'knewvell15@wikispaces.com', 'member', true, '', 'dSPIXu2yz');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('8ea02dd0-efcd-4546-ba22-26b46a5d1382', 'Natividad', 'Mc Menamin', 'nmcmenamin16@boston.com', 'member', true, '', 'yaTG9MlVyyt');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('7fd3971b-2744-428f-a129-0139ff315402', 'Winne', 'Awton', 'wawton17@upenn.edu', 'admin', true, '', '4KtjCwb3');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c2622cbf-04ce-4d3e-8f28-581823feef3c', 'Brianna', 'Imlach', 'bimlach18@multiply.com', 'member', false, '54a288b8-d37f-4708-a7ef-f14acf918505', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('9d2da61b-07c9-4e1b-ad03-68a021beb89b', 'Harald', 'Blaksley', 'hblaksley19@columbia.edu', 'member', true, '', 'n2k8QEeQs6');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('4b917589-b5ec-4298-8742-e0940325651d', 'Daren', 'McCloud', 'dmccloud1a@skype.com', 'member', false, '7cc881b6-c642-4dd6-a3e2-f4b6976ab4b8', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('68d4532b-ee51-4b43-9c2e-c971027887c3', 'Graehme', 'Buglass', 'gbuglass1b@chronoengine.com', 'member', true, '', 'MERweR78');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('eb07cd1c-479a-44ed-897a-c8c3421bdd05', 'Tammy', 'Pitt', 'tpitt1c@technorati.com', 'admin', true, '', 'Ufc3emTtw');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('48d5c10d-91ee-4632-b47e-ed2092748eab', 'Brocky', 'Arons', 'barons1d@addtoany.com', 'admin', true, '', 'CeTVQV47xD');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('21d65da8-d0aa-40c4-9dc2-fe3bc7b9eaf2', 'Kippy', 'Woolhouse', 'kwoolhouse1e@networkadvertising.org', 'member', true, '', 'GjOf1Sh');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('4a815610-bea1-4074-a869-459de2811c36', 'Michele', 'Chamberlen', 'mchamberlen1f@ovh.net', 'admin', true, '', 'oBuokt');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('99e89af2-f5b7-47f7-9c80-fd475d6c82f5', 'Hobart', 'Gagg', 'hgagg1g@stanford.edu', 'admin', true, '', 'Si1SkCLzBY8w');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('691839eb-4acb-4668-8345-8127f877d8bb', 'Meggy', 'D''Ambrogi', 'mdambrogi1h@ameblo.jp', 'admin', true, '', 'aswlfMJ6P');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('b39d6049-377b-4a62-adcd-73009074bd80', 'Taffy', 'Ryal', 'tryal1i@topsy.com', 'member', false, 'c94172ed-2e9d-468f-b2b4-a1fc1cb5fcbf', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('0112b573-e59a-46b6-8745-e77a75be1ac6', 'Anastasie', 'Salters', 'asalters1j@newsvine.com', 'admin', true, '', 'b1KDqG');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('b475df13-f92f-42b6-9bfd-644dce9e69e9', 'Haroun', 'Kosiada', 'hkosiada1l@dell.com', 'member', true, '', '0RgLQ9');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('f8eab169-6c17-4ee9-a325-ff422e8b45c6', 'Rhodia', 'Meindl', 'rmeindl1m@indiatimes.com', 'admin', true, '', 'l6Siw8G');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('b51b2711-6b98-4bac-91bc-4b23fd88eb5d', 'Shina', 'Harfoot', 'sharfoot1n@craigslist.org', 'member', false, '5b173178-757b-4273-be22-f6a3a99ba222', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('cd4cbbad-3208-4d1f-bc16-0ea4214195a7', 'Gardiner', 'Mattacks', 'gmattacks1o@nature.com', 'member', true, '', 'AEIQtQc');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c07d053c-1e8b-4e23-8f28-fecc1987ae73', 'Zachery', 'Foxall', 'zfoxall1p@discuz.net', 'admin', true, '', 'ZYyYOS');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('0fb8a95b-8788-42a5-8d38-a752548ffb2e', 'Ruth', 'Marxsen', 'rmarxsen1q@yahoo.co.jp', 'member', true, '', 'ehiybUn97Wx');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('962b7d41-d05d-4179-9157-1928986fc982', 'Yale', 'Orhrt', 'yorhrt1r@w3.org', 'member', false, 'edabd0e3-9b15-44c3-837f-8e502fd4c0b9', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('b0ff4c3f-0e8c-4ab7-8246-006fa9eaa806', 'Malinda', 'Scritch', 'mscritch1s@merriam-webster.com', 'admin', true, '', '5qjPufvdF');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('71db2582-111c-4e00-97f9-62e6a3a189f8', 'Rog', 'Schrir', 'rschrir1t@yelp.com', 'admin', true, '', 'dKXuve');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('f32de303-080f-486e-93b0-f49942d9e10f', 'Hartley', 'Tourle', 'htourle1u@simplemachines.org', 'admin', true, '', 'vJ4FAGxS');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('eee09fff-4a27-48e9-afb9-3265018be5a9', 'Gearard', 'Tadgell', 'gtadgell1v@a8.net', 'admin', true, '', 'bz0P1Gfb');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('229e40fd-89fe-4ee7-b25d-a79d167398f9', 'Tiertza', 'Bowstead', 'tbowstead1w@google.co.jp', 'member', false, '611e29f9-37c2-4b91-a920-07625f01e2ae', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('07f56109-315b-42c8-ae20-9122c27854f2', 'Charlot', 'Akester', 'cakester1x@examiner.com', 'admin', true, '', 'XqfGkC');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('1f916341-7551-4a94-9ca9-69006db4dd70', 'Denise', 'Highman', 'dhighman1y@xrea.com', 'member', true, '', 'aIb1vLEv');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('2edda33c-c225-4511-bd2f-2ced358f8782', 'Brendis', 'Najafian', 'bnajafian1z@smh.com.au', 'member', true, '', 'Zb9SlKBTmAhg');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('e71961a3-7019-4321-8e87-80ca241c48b2', 'Eloisa', 'Pray', 'epray20@independent.co.uk', 'admin', true, '', 'DrWLHeS');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('cbba2119-7450-426f-92b8-c235e7b6329f', 'Scotti', 'Merrywether', 'smerrywether21@last.fm', 'admin', true, '', 'cJ7b9Ba');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('728f1a83-30b4-44b6-affb-63615d7387f2', 'Ilyse', 'Joannet', 'ijoannet22@nature.com', 'member', false, 'f4341b2a-45b4-42e4-abc4-b1c4989186fe', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('7c6bf5b5-ffd4-4415-baa7-cf6ab3bfd2d2', 'Tonia', 'Bertolin', 'tbertolin23@soundcloud.com', 'admin', true, '', 'me81DEP3JO');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('b99a5b60-2818-4c1d-a076-10675050961d', 'Friedrick', 'MacWilliam', 'fmacwilliam24@hexun.com', 'admin', true, '', '1hBrBV49m');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('ab1184bf-f8aa-48e0-a072-09ba2a7cc505', 'Vale', 'Riediger', 'vriediger25@bizjournals.com', 'admin', true, '', 'JPCgCpcg42kI');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('05600353-e6bd-4015-86f9-e93df06bf5f7', 'Juliann', 'Silversmidt', 'jsilversmidt26@earthlink.net', 'admin', true, '', 'P3yxfiTbR');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('b1c1424b-1cd8-482e-96d6-d6205a319f14', 'Melisent', 'Illem', 'millem27@nih.gov', 'admin', true, '', 'CB4TrKjxW4C');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c59ae5dc-2290-4ecb-8581-147b100fb6bd', 'Yvonne', 'Chatfield', 'ychatfield28@google.de', 'admin', true, '', 'LdmW1w6bZvu');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('a58f0ce1-0076-4fcf-98b2-968cbe59c6e7', 'Beilul', 'Addyman', 'baddyman29@wikia.com', 'member', false, '4998050b-0779-4902-ae4c-a0bb3eb0dabb', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('dbd9344e-6a1a-456a-928b-40cbae14393e', 'Vivian', 'Camacke', 'vcamacke2a@berkeley.edu', 'member', true, '', 'kPD4Ms');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('413b5941-bd97-4b47-8aa7-d2a38d78e78c', 'Rudie', 'Willers', 'rwillers2b@privacy.gov.au', 'admin', true, '', 'tL8dtlJpp');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('e55f8d74-0387-41f2-8883-d65bcf59a7a3', 'Crystal', 'Rudloff', 'crudloff2c@amazon.co.uk', 'member', false, '0db8cb39-502a-49ec-a7b0-b3eaabeba412', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('8def1081-4be1-4ba7-9f19-570c79dea275', 'Sven', 'MacWhan', 'smacwhan2d@pbs.org', 'admin', true, '', 'VhNqIGb8');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('444f8d9e-742e-4a58-9b17-9e7c7e278b9f', 'Delphinia', 'Bramhill', 'dbramhill2e@unc.edu', 'admin', true, '', 'o7LOZnMR');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('05661f5e-fc4d-479f-bb93-ac943bcd7c4f', 'Flossy', 'Gamage', 'fgamage2f@is.gd', 'member', true, '', '69Z3oQ6nP');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('8505f5d8-6a3f-4ade-808b-f5763262dc9f', 'Matty', 'Dibden', 'mdibden2g@redcross.org', 'member', false, '1ef8ff4f-d3c8-42ee-921a-a0e7cc86b799', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('9eb24ecf-a3b5-4c13-a550-7cd583c87af8', 'Cesar', 'Millgate', 'cmillgate2h@ustream.tv', 'admin', true, '', 'U3Q4GlJb');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c7dba66b-e6d8-4bbd-81c2-f4dfe518eec7', 'Moshe', 'Deetlof', 'mdeetlof2i@usda.gov', 'member', false, '3d6d2e67-ea29-4585-8c5d-1a52e436a1b7', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('868d6cdd-10ed-4e6b-887a-6e0e4579fb63', 'Arman', 'Udall', 'audall2j@liveinternet.ru', 'member', true, '', 'HaWaBr12El');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('7677f1f8-62c4-4ce4-a749-acb85d88a744', 'Orelia', 'Sharkey', 'osharkey2k@fastcompany.com', 'admin', true, '', 'rkb0L7zIY2i');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('4f164770-3afd-415e-96b9-ecc4dc62fd5c', 'Ferdinand', 'Brouwer', 'fbrouwer2l@ucla.edu', 'member', false, '4dcc5a4e-e888-434a-a676-711a2f4ccd3e', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('ba5ec7da-65f8-43e7-9c37-5daa495335ef', 'Matthus', 'Scuse', 'mscuse2m@whitehouse.gov', 'admin', true, '', 'byuByW4kZm');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('69d49a2f-e4a8-426c-af55-ad52e4d5b98c', 'Arlyne', 'Wannan', 'awannan2n@wisc.edu', 'admin', true, '', 'QKURz2');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('1b3060c6-bcb2-411b-a14c-a205188744fe', 'Tobit', 'Sebborn', 'tsebborn2o@webnode.com', 'member', true, '', '9U1RfD8fjP');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('ceba8349-74ab-4654-8235-6587375b4460', 'Carmencita', 'Costell', 'ccostell2p@cam.ac.uk', 'admin', true, '', '3toMzy');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('1094984d-0fd0-4527-8696-b77ca47c2443', 'Onfre', 'Binfield', 'obinfield2q@indiegogo.com', 'admin', true, '', 'fhmngcWGBh');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('39998450-b891-43a3-8aa1-89ec0a143171', 'Walden', 'Fetteplace', 'wfetteplace2r@deviantart.com', 'admin', true, '', 'O29T77BgtJ');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('cbdded07-dbc8-43b1-a79b-4033db679913', 'Edwin', 'Headings', 'eheadings2s@e-recht24.de', 'member', true, '', 'XCgGPH60');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('2c877d97-202a-465d-9fda-6ed90291582e', 'Ulrike', 'Brittan', 'ubrittan2t@exblog.jp', 'admin', true, '', 'cQbCd5QqwJx');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('a126e0ca-acaa-44ac-b557-d933c6896ff5', 'Josie', 'Crumby', 'jcrumby2u@ning.com', 'admin', true, '', 'JgkMcbWwe');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('94791c00-17ae-4362-b8a9-3a873dff7eee', 'Gordon', 'Attarge', 'gattarge2v@sogou.com', 'admin', true, '', 'Bu9PTD2rBgx');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('05af3564-6d55-4e91-816d-3868bef33567', 'Tamarra', 'Kitchaside', 'tkitchaside2w@google.ru', 'member', false, '1630e9d9-45d9-42be-89d8-dbec1ae176c6', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('4d0da606-98d2-4aaa-9570-70c368a7602f', 'Haily', 'Dalgarnocht', 'hdalgarnocht2x@indiatimes.com', 'admin', true, '', 'orbZuOaL');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('e16d6f0e-8f24-4cd3-a3db-d2324b8fff90', 'Wyndham', 'Pea', 'wpea2y@rambler.ru', 'member', false, '8b00813e-2bda-4183-b11c-d0001087a775', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('9b5b3394-7ce4-4bce-b955-0ef6e66a74ea', 'Sergio', 'Leist', 'sleist2z@ibm.com', 'member', false, '16c08c12-4771-47f8-b55e-8cbf19475c71', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('832857c0-5449-4a5b-a93a-af4e91955112', 'Laurence', 'Maestrini', 'lmaestrini30@furl.net', 'member', false, 'e6e0f8af-3e60-49ac-bdf4-19dbfb352239', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('358c86d6-e3f0-411a-ba25-5d3e35eb04a4', 'Della', 'Spurge', 'dspurge31@sciencedirect.com', 'admin', true, '', '6W3Ivs');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('a48c2751-d541-4b72-8415-bed5912d0770', 'Valeria', 'Wombwell', 'vwombwell32@europa.eu', 'admin', true, '', 'KFujxXKm');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('6657b3b8-8454-487e-864b-a9e12346fe3a', 'Fields', 'Godbehere', 'fgodbehere33@cloudflare.com', 'member', false, 'c4f8f1b6-a24a-4a93-ab20-5236c24fa8b2', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('7ac08123-c79a-4800-b790-fc102e91b070', 'Uri', 'Broadbere', 'ubroadbere34@samsung.com', 'admin', true, '', 'iAoYR4');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('d15e3b87-bb9d-42bf-a3da-492dd9025d90', 'Morganne', 'Stentiford', 'mstentiford35@miitbeian.gov.cn', 'admin', true, '', 'CFNRQJk');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('663b54f5-6332-428e-af47-509455528cb5', 'Gus', 'Dewfall', 'gdewfall36@amazonaws.com', 'admin', true, '', 'k2l0Kqa8VAtO');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('6187cf96-3b7f-4c72-be61-c6d372c8efea', 'Davin', 'Jans', 'djans37@earthlink.net', 'member', false, '21bac83a-bc85-4cb6-91b1-b9cf67cfd7a2', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c4ef8f6c-75cf-4336-aa1b-77c9a91383db', 'Almeta', 'Sperring', 'asperring38@fema.gov', 'admin', true, '', 'RBg87jeZYQI');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('ee1e8552-1a2a-458a-810e-e4c2858aabed', 'Maynord', 'Kroger', 'mkroger39@statcounter.com', 'admin', true, '', 'AMrWVxHZQk0U');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('ced30db8-f25e-4a1b-9b3e-876ae06fb336', 'Emmye', 'Stickford', 'estickford3a@un.org', 'admin', true, '', 'QXPShvJ4AkSx');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('a5027977-7440-488f-921f-d58c85233968', 'Aileen', 'Meekings', 'ameekings3b@statcounter.com', 'admin', true, '', '17uecCznlBp');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('d864a8c0-97e7-4d10-8cb4-bdcb6af16116', 'Isiahi', 'Wraighte', 'iwraighte3c@admin.ch', 'member', false, '292389c5-183f-428e-b6bc-284d212cb1c9', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('31d51cfa-c7d7-4f32-9b99-a664c456a3af', 'Pavel', 'Rouzet', 'prouzet3d@ifeng.com', 'member', false, '3012baff-90b2-4915-bb84-23b153d4f01d', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c401a5c4-dd51-4b98-b778-1e48f5c9524b', 'Teirtza', 'McGonnell', 'tmcgonnell3e@usgs.gov', 'member', true, '', 'FvtXMdI');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('ee5ccd2d-c0e6-4f72-9686-6a308f9198bc', 'Lulita', 'Greenstead', 'lgreenstead3f@zimbio.com', 'member', false, 'fda2874d-3db0-4e2e-b7f0-6fbbff0cc2f1', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('0432e763-41d5-41de-945d-eec69004e78e', 'Fransisco', 'Radenhurst', 'fradenhurst3g@xinhuanet.com', 'member', true, '', 'DOlW4S7gWdvE');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('5f8a2758-264d-4db2-bd29-3ae650a575aa', 'Margarette', 'Ubach', 'mubach3h@ow.ly', 'member', true, '', 'uRp8KoKhPrFS');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('37dd4fd3-8851-4ca2-96db-c7a46bbca996', 'Cordy', 'Keaton', 'ckeaton3i@ow.ly', 'admin', true, '', '6SEu5NiQ');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('fc0610ba-ba34-41b4-9678-a04f55e287d6', 'Amabel', 'Louden', 'alouden3j@gnu.org', 'admin', true, '', 'EBfEzAzxmN');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('fb9d95d1-f2ac-4957-990d-d466206f42b8', 'Ketty', 'Perrygo', 'kperrygo3k@reuters.com', 'admin', true, '', 'QH84GSd');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('69930a4d-f096-4bcd-a960-63af69aaca18', 'Antonietta', 'Bowld', 'abowld3l@theguardian.com', 'admin', true, '', 'CXRDOz1Mp6e');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('0736d580-9186-4d25-8560-86146056980c', 'Shaylah', 'Du Hamel', 'sduhamel3m@google.com.br', 'admin', true, '', 'IA9NzI8d3qv');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('be9061b4-4c68-4d6c-9c63-46dc5a30e357', 'Joachim', 'Rosenthal', 'jrosenthal3n@blogtalkradio.com', 'member', true, '', 'V0HwReXbcd15');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('d3a77ce1-e15d-4ee9-96c4-9c21a98f5854', 'Felipe', 'O''Heagertie', 'foheagertie3o@parallels.com', 'admin', true, '', 'I0ymazpPwZtt');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('15d2932d-c2d6-4999-8041-b775b8dcb84a', 'Alfonse', 'Christofe', 'achristofe3p@prnewswire.com', 'admin', true, '', 'pTBXSOmGIb03');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('7b4fd856-3774-427c-9a71-de79dfe66425', 'Dani', 'Raff', 'draff3q@aboutads.info', 'admin', true, '', 'FeaxWE8TUuVn');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('3f5268e5-190d-4592-817a-6bb8831d05b2', 'Miquela', 'Kitchiner', 'mkitchiner3r@google.com.br', 'admin', true, '', 'fOEEAx');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('12ce4427-65b2-4d68-b9ea-102794439f88', 'Tuesday', 'Mumbeson', 'tmumbeson3s@domainmarket.com', 'member', true, '', 'frzV91r');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('4df1e307-0898-43e9-8438-659bec7a344a', 'Murdoch', 'Matchitt', 'mmatchitt3t@yahoo.com', 'admin', true, '', 'JPKNBvU');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('d7915b66-5b2c-4125-984c-ff921b820f9a', 'Pren', 'Hedgeley', 'phedgeley3u@howstuffworks.com', 'member', true, '', 'Cb7ubrWBGS');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('3683d3a1-907e-4126-aef8-9794db7ac168', 'Rod', 'Sabine', 'rsabine3v@wordpress.org', 'member', false, '4f55d589-f08c-48d5-b32d-b24b3eaa3613', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c1109c07-83a6-43e1-86e5-93486ed003a2', 'Lincoln', 'Avard', 'lavard3w@tinyurl.com', 'member', true, '', 'Oj1ogqKRS0O');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('240788bf-c274-4b9f-b6a0-f6c5466916d3', 'Dominik', 'Bigglestone', 'dbigglestone3x@hubpages.com', 'admin', true, '', 'LVILFRj1nV');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('ebf733d4-683c-41e8-9300-e831e441b277', 'Clemmie', 'Maffei', 'cmaffei3y@prlog.org', 'member', true, '', 'hYjxrfFCL');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c9a67c8e-3adc-4a7c-8203-27cf825c7d0b', 'Katharine', 'Kettoe', 'kkettoe3z@theatlantic.com', 'member', true, '', 'hs1nc5dWuL9Q');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('17f842cf-9ef2-4d78-8143-d08cbcf5416d', 'Leodora', 'Calven', 'lcalven40@flickr.com', 'member', false, '840b99bf-10aa-4a03-a7be-6ceae1e20067', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('35382e30-1059-4af7-a313-d04efae4b9b2', 'Gerti', 'Haggath', 'ghaggath41@quantcast.com', 'member', true, '', 'gvNnSlK');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('b4a60352-5cf1-4c4a-bc9f-88d1015034b2', 'Merl', 'Bosse', 'mbosse42@blog.com', 'admin', true, '', 'ynkkJB1J');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('8f41e8d0-88cc-4c45-81f4-f278ed0e97b6', 'Tammara', 'Losano', 'tlosano43@51.la', 'admin', true, '', 'a4sR9LU');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('0fe021e0-f857-40c7-b883-4cc3f92d9e2e', 'Amy', 'Lillecrop', 'alillecrop44@storify.com', 'admin', true, '', 'wzwu9e9gkxAw');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('5ca81022-3522-4b7a-b18d-c02f0c979d32', 'Devina', 'Hiom', 'dhiom45@canalblog.com', 'member', true, '', '2pojdBtbF');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('57704abd-d2fb-402d-8e26-f084d621ba34', 'Eula', 'Statton', 'estatton46@marriott.com', 'member', true, '', 'McogYrrI');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('f828dff6-44b4-46b1-8e30-64b866d713b5', 'Chelsie', 'Robion', 'crobion47@printfriendly.com', 'member', true, '', 'mc2o1Q9');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('a302db77-0761-4408-9075-666f612c3b37', 'Merv', 'Sitlinton', 'msitlinton48@howstuffworks.com', 'admin', true, '', 'HkZDP1OnQn');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('cce2e970-db73-4cc5-92a5-6cfc2fb3d771', 'Allianora', 'Aust', 'aaust49@live.com', 'member', false, '718852d9-eb4b-475f-b8b9-4b6b3ec131b1', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('947f7803-e8f9-469a-acdd-6754b25dba93', 'Feliza', 'Champain', 'fchampain4a@list-manage.com', 'member', false, '823328c4-2fbc-4b99-b38e-5a3126834cca', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('bf74d7df-b808-4e07-9108-f09b64d9eb10', 'Sarine', 'Rowlands', 'srowlands4b@timesonline.co.uk', 'admin', true, '', 'BoI8boK6Lo4');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('ad406f9b-6a76-4dd2-bbc9-3c07979165a6', 'Kathe', 'Tudbald', 'ktudbald4c@ebay.co.uk', 'admin', true, '', 'h5mAQS22Hu2');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('41cea5f3-2639-4a8a-8a00-0cd4fb1efa8f', 'Alfy', 'Kopke', 'akopke4d@diigo.com', 'member', false, '6d5649d7-7216-4974-820a-73837235224f', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('336143b2-43f3-468a-bfc0-db670fa1af49', 'Hinze', 'Perulli', 'hperulli4e@google.com.hk', 'member', false, '4419a66b-777c-42c7-be9e-74717af15b25', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('747b22d5-29f8-4e6d-b245-ae56d22906fa', 'Bette', 'Gayne', 'bgayne4f@webeden.co.uk', 'admin', true, '', 'zSbtqzD');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('69182df5-bf42-43d1-8ace-be567a50d9fa', 'Babara', 'Dowzell', 'bdowzell4g@usatoday.com', 'admin', true, '', 'KJL60qU');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('522f1508-f8d8-478a-bfdb-a14cf0bf5400', 'Lars', 'Gherardi', 'lgherardi4h@statcounter.com', 'member', true, '', 'lfGj1hX');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('804a79aa-0683-46cd-903e-3f6cd6a4b6de', 'Seth', 'Hembling', 'shembling4i@hexun.com', 'admin', true, '', 'BmHt4p6QVJgi');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('026408ee-6a90-405b-9a94-7685125ee413', 'Stephanus', 'Tight', 'stight4j@bbb.org', 'member', true, '', 'bdsSln');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('e147e0b4-ec01-4906-86c1-539f525f6f52', 'Glynda', 'Hardey', 'ghardey4k@vk.com', 'admin', true, '', 'jJcYrLdKy');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('5f25821a-0fe7-4926-b7fe-c8b2cd2c0235', 'Charlena', 'Cowl', 'ccowl4l@tmall.com', 'member', false, '1a051fc4-f483-4762-83ef-f3b1a6d67c77', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('377406bd-0479-4acb-b92c-dc03b26010d4', 'Thatch', 'Eskriet', 'teskriet4m@nhs.uk', 'member', false, 'f4c2ce6f-0988-487a-ac58-befea6869413', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('39de33e5-0488-456f-b813-29f6370d9d88', 'Bert', 'Moynham', 'bmoynham4n@arizona.edu', 'member', false, 'fdbb3e1f-0cd3-440e-9eb9-8e2f42f01460', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('7d31f1ee-5572-496c-8ae9-eb82bb803ea8', 'Claus', 'Ewles', 'cewles4o@is.gd', 'member', true, '', 'if7lLDIwk');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c1ffee5e-ce3b-453a-945b-9e4883a854e8', 'Jeremie', 'Prangnell', 'jprangnell4p@eepurl.com', 'member', false, '5334226a-0c6e-48cf-8108-58a4cdd0e445', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('416a41e6-e185-4e9d-ace7-914350678df0', 'Esther', 'Kearton', 'ekearton4q@forbes.com', 'member', false, 'b1bbdf23-a41a-49d8-9020-d89232875728', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('6bae441a-1d05-4829-ba57-02f59eaed160', 'Northrop', 'Dobbins', 'ndobbins4r@technorati.com', 'member', true, '', '7YGayBAbZ3Km');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('721ba45a-c1fb-4a39-9125-5de8c09c2a4a', 'Kipper', 'Dupey', 'kdupey4s@skype.com', 'admin', true, '', 'itoftbQ');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c13d7073-2913-4fe1-8e43-6832d0c7a7a5', 'Aimil', 'Hartless', 'ahartless4t@is.gd', 'member', true, '', 'QWCRsd');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('de688ade-57f9-41c3-bfe1-b0932b8893e7', 'Kelci', 'Piechnik', 'kpiechnik4u@techcrunch.com', 'member', false, 'd03df90b-dc76-44d0-808d-aa6f8058cd88', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('d0603efa-74f1-4d8c-8f87-51322348db95', 'Marve', 'Crookshanks', 'mcrookshanks4v@ihg.com', 'admin', true, '', 'dbMeNP');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('52e2d47d-2cbf-4889-88a4-9cca4a435bff', 'Ignacius', 'Absalom', 'iabsalom4w@exblog.jp', 'member', true, '', 'g5JAAZO');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('fa572a5b-836c-4c9e-a800-e4f127c52912', 'Chucho', 'Brauner', 'cbrauner4x@opera.com', 'admin', true, '', 'tK6b9VJ3ik');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('d381d5a9-7e3b-4f37-8f7e-dad2471d256f', 'Wilbur', 'Dumini', 'wdumini4y@dell.com', 'member', false, 'bc99c3b1-17b4-4e7e-92b4-a7327afd6cff', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('36ee2254-84db-4bcc-9eed-6727546fb3f0', 'Cody', 'Astlatt', 'castlatt4z@nhs.uk', 'admin', true, '', 'ghT6oy');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('7f35e553-31a1-4979-b355-b456bbff3463', 'Harmonie', 'Clibbery', 'hclibbery50@gizmodo.com', 'admin', true, '', '2yZzI3');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('2a7e8720-f606-4b97-91b7-84ccaa31f8d2', 'Ruthy', 'Fawdry', 'rfawdry51@hugedomains.com', 'member', false, '67cf87a2-fbd4-4f96-8cf0-f6879f175f03', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('415f9c3a-9fac-4794-8524-dfdd40e9ed06', 'Dev', 'Augie', 'daugie52@bloomberg.com', 'admin', true, '', 'o1b0Ia');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('94367041-2430-435f-962c-5bd040d93cc7', 'Colby', 'Estabrook', 'cestabrook53@studiopress.com', 'member', true, '', '7KdZLubWxe9');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('4718a8d4-8e7b-40a8-9b28-5a0f67fa1fef', 'Almeria', 'Joannidi', 'ajoannidi54@chronoengine.com', 'admin', true, '', '7FCI0OG5');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('70963e6a-5275-4be7-89c7-1aebde5f4be2', 'Jeanelle', 'Mallabund', 'jmallabund55@drupal.org', 'admin', true, '', 'cJQE24xx');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('108fee32-0b6f-4776-83f3-6cf6ad58a0b2', 'Shem', 'Hughs', 'shughs56@theguardian.com', 'admin', true, '', '3qwY3Tsv');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('aaaa974a-0f74-4850-b2b7-a88c0a07ffcb', 'Abbi', 'McMichael', 'amcmichael57@biglobe.ne.jp', 'member', true, '', 'pz7mltT5d');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('6493955b-8bdf-46d4-b3ab-c1626ebd6989', 'Rebecka', 'Ivakhnov', 'rivakhnov58@printfriendly.com', 'member', true, '', 'hiZWKHa');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('13990d56-dd6d-4702-b2d0-905d7bf9cb36', 'Conrade', 'Matthiae', 'cmatthiae59@wired.com', 'admin', true, '', 'ppSPgAR');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('ee77815a-cc8d-41fc-a285-0f138d59a4a4', 'Filberte', 'Stieger', 'fstieger5a@bing.com', 'admin', true, '', 'gBiUQXG');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('74510ed5-fc23-4d04-af73-34928c719fed', 'Alison', 'Langmuir', 'alangmuir5b@g.co', 'member', false, '33b4538c-ec31-419e-b754-cf19ec5439f4', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('a9969ed8-a411-430b-8664-d1b5c08416d3', 'Zechariah', 'Copper', 'zcopper5c@latimes.com', 'member', false, 'cca08f7d-a26c-4b94-bd76-dec9f1f6085a', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('1d63d5f9-97aa-4d33-a9b1-55cc6ad8017b', 'Gian', 'Imeson', 'gimeson5d@ucoz.com', 'member', false, '8b0a5b45-c009-4b92-b215-7356270f277b', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('8d99925a-ad55-4a62-91f6-bcf069084a50', 'Noah', 'Barkhouse', 'nbarkhouse5e@lulu.com', 'admin', true, '', 'WrTACy7');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('68bf6c3f-e863-42fc-bd8f-21bef6a18564', 'Winni', 'Poskitt', 'wposkitt5f@sun.com', 'member', false, '4e359c35-9929-485e-ad07-2b8f6d3aec88', '');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('18f17bc3-3d1d-4314-87c4-35a39850b121', 'Lynda', 'Egerton', 'legerton5g@cnet.com', 'admin', true, '', 'sW9vGBxvNUi');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('760a62de-d9ab-4763-a6d1-375be4babd56', 'Hedda', 'Kingaby', 'hkingaby5h@unc.edu', 'admin', true, '', 'NstmsWD');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('c2864860-90d5-44f2-b027-878afaa80f89', 'Ysabel', 'Wolfendell', 'ywolfendell5i@themeforest.net', 'admin', true, '', '4p1x2PjGa');
insert into users (id, first_name, last_name, email, role, active, activate_token, password_hash) values ('986d7dbc-212b-4836-8c31-4759d32e7e46', 'Rosemonde', 'Tatton', 'rtatton5j@technorati.com', 'admin', true, '', 'BSEqnZF7');

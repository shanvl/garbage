/*
 A script to be executed at the start of the docker postgres image, which populates a db with sample data.
 It creates 900 pupils and 30 events. Each pupil has 80% chance to attend an event and 50% chance to bring
 some recyclables to that event
 */

-- create resource enum type
do
$$
    begin
        if not exists(select 1 from pg_type where typname = 'resource') then
            create type resource as enum ('paper', 'plastic', 'gadgets');
        end if;
    end
$$;

-- create event table
create table if not exists event
(
    id                varchar(25) primary key,
    date              date        not null,
    name              varchar(25) not null,
    resources_allowed resource[]  not null check (cardinality(resources_allowed) >= 1)
);

create index if not exists event_resources_allowed on event using gin (resources_allowed);
create index if not exists event_id_date_name_idx on event (id, date, name);
create index if not exists event_id_name_date_idx on event (id, name, date);

-- create pupil table
create table if not exists pupil
(
    id                varchar(25) not null primary key,
    class_letter      char        not null,
    class_date_formed date        not null,
    first_name        varchar(25) not null,
    last_name         varchar(25) not null,
    text_search       tsvector generated always as (to_tsvector('simple', first_name || ' ' || last_name || ' ' ||
                                                                          extract(year from class_date_formed)::text ||
                                                                          class_letter ||
                                                                          ' ' ||
                                                                          class_letter || ' ' ||
                                                                          extract(year from class_date_formed)::text))
                          stored
);

create index if not exists pupil_text_search_idx on pupil using gin (text_search);
create index if not exists pupil_class_name_idx on pupil (class_date_formed, class_letter, last_name, first_name);

-- create resources table
create table if not exists resources
(
    pupil_id varchar(25) not null,
    event_id varchar(25) not null,
    paper    float4      not null default 0,
    plastic  float4      not null default 0,
    gadgets  float4      not null default 0,
    primary key (pupil_id, event_id),
    foreign key (event_id) references event (id)
        on delete cascade
        on update cascade,
    foreign key (pupil_id) references pupil (id)
        on delete cascade
        on update cascade
);

create index if not exists resources_event_id_gadgets_idx on resources (event_id, gadgets desc nulls last);
create index if not exists resources_event_id_paper_idx on resources (event_id, paper desc nulls last);
create index if not exists resources_event_id_plastic_idx on resources (event_id, plastic desc nulls last);


-- insert 900 pupils to the pupil table
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yrmxdjuotf', 'Vaughan', 'Enrdigo', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gwfothlucj', 'Fan', 'O''Tuohy', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('otxmuzyesl', 'Dale', 'Corey', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bsfihpudro', 'Shelden', 'Lawling', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('aecwnilqxr', 'Shanan', 'Wildsmith', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pldyceawgb', 'Pamela', 'Domeney', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fbrvwichkn', 'Germaine', 'Newsome', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xkezwovtdc', 'Northrop', 'Unger', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('chagrdkxsl', 'Gonzales', 'Sonner', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eyouafgtnm', 'Deane', 'Gilroy', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gctwzjidaq', 'Abby', 'Hearons', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ajpcbivefh', 'Katinka', 'Lightbowne', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kvyarjwhxq', 'Judas', 'Konert', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mwsrvnckup', 'Lucian', 'Sheers', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('flkaznpgxd', 'Roby', 'Telling', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ozgmikdsyc', 'Saxe', 'Humphery', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('iemtlqyubg', 'Rik', 'Parysiak', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ytdhapxwos', 'Evangelia', 'Oertzen', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jqinoyhekg', 'Davina', 'Rowthorn', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tijwglmqex', 'Luz', 'Stolworthy', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gyukimepso', 'Duffie', 'Lisciardelli', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nivlqdmzgj', 'Hillary', 'Adamovicz', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jpcrlogqvt', 'Morissa', 'Leon', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('chqpukltaw', 'Shanta', 'Huskinson', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('znritbmahc', 'Mattie', 'Tschirasche', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mdawrkzqny', 'Angus', 'Arundale', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zsrhckofml', 'Nolana', 'Sidgwick', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pbxknitucs', 'Kurt', 'Nealand', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tswjgaelqi', 'Lexy', 'Jaine', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qulhtgbyjz', 'Torin', 'Amiable', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gnutyeorza', 'Allie', 'Adamsky', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jrbaizxdwo', 'Abbie', 'Trevain', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vjzsautbhn', 'Jackelyn', 'Helliar', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fdtayheuvp', 'Gordon', 'Placidi', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mcriwtlzyf', 'Crosby', 'Beckenham', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qbngxtmoes', 'Dean', 'O''Sheerin', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uhyvowdrfe', 'Casey', 'Bonsale', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('prfsjzavmi', 'Cirstoforo', 'Elsy', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qgmkpshctv', 'Zared', 'Champniss', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fbynxjaztl', 'Baudoin', 'Cannavan', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ehnpbtvqkz', 'Ermanno', 'Ducham', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('oqtgnezpfh', 'Corbet', 'Telezhkin', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bnxqrwgzkd', 'Jerrilee', 'Kauschke', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bqctndeauf', 'Whit', 'Bleier', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gpyskbcahl', 'Dexter', 'Collumbell', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vlorqywshj', 'Jeniffer', 'Ind', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yfbrqgcnje', 'Solly', 'Rijkeseis', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('arwbdjgcyp', 'Rhonda', 'Shingles', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vfgomtkbhi', 'Jess', 'Labbet', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fryjeuvhqn', 'Grannie', 'Nowlan', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('selxkzwamo', 'Powell', 'Henken', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ojyadrsekw', 'Wolfy', 'Rossant', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ifxscmpyvn', 'Noni', 'Moncreiffe', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cyrbqzvpaj', 'Fritz', 'Karlmann', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yrganofqcj', 'Gerrie', 'Pennington', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jkrmistclp', 'Sheeree', 'Mingo', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hvqbsozutx', 'Wells', 'Christauffour', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lecswzqvkp', 'Andrej', 'Pendlenton', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zyxugrfpmo', 'Lorine', 'Oaks', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uvkntzqsjh', 'Ernie', 'Ferran', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hyltfizkrg', 'Saraann', 'Belchem', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('evwhtdkxnl', 'Jason', 'Begwell', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cijauhxmfb', 'Channa', 'Gallie', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qsyhpbejut', 'Raina', 'Ledeker', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('byfxudjlnv', 'Demetris', 'Fentem', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dkqoihrvse', 'Herrick', 'Dreini', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jvgwhfptir', 'Graehme', 'Shewan', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('owzkyxbinr', 'Kai', 'Halfacree', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jusfhcmtnq', 'Giulia', 'Woltman', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jlemthrbfw', 'Corene', 'Yitzowitz', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ymrzpectvu', 'Jaynell', 'Trobridge', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fivegruqst', 'Ardelis', 'Stilly', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uvlsrtkjzb', 'Iris', 'Conrad', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zcqpnfaiyt', 'Rikki', 'Minter', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xyohludpvc', 'Oona', 'Thrush', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('thmiuraedj', 'Sande', 'Dreelan', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('iokjvqpsyg', 'Petronille', 'Arons', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sahenfovyc', 'Nerte', 'Jellings', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wqmrkhfalc', 'Gayelord', 'Phelipeaux', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bltkozupyv', 'Larissa', 'Trafford', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('phqlczjbst', 'Kris', 'Lippiett', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mgxcazhdwi', 'Giraldo', 'Windrass', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jnvbwkzslq', 'Joela', 'Markus', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qcgodtxewh', 'Shandy', 'Kynaston', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('oqadbvhrup', 'Rickert', 'Drinkeld', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zgjmfctdbq', 'Edi', 'Donett', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('txoszykrle', 'Garrott', 'Arnull', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('disvochtax', 'Rheta', 'Hayzer', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bfljayopzh', 'Ermin', 'Pegram', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xlawyzitmd', 'Ash', 'MacGorley', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ygzsoxamdu', 'Chrissie', 'Edmenson', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bogusdqamf', 'Fleming', 'Thompkins', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dsxnthcybz', 'Reube', 'Wickes', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kpydqbezcn', 'Sylvester', 'Mumberson', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mnycbezgqw', 'Grover', 'Beadman', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('igzkuwjbms', 'Shirleen', 'Delamaine', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('istwmarvzn', 'Zorana', 'Garrity', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('slmizapnvq', 'Nils', 'Durbyn', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kshnpawulb', 'Pall', 'Brown', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wgftdisvhl', 'Lisette', 'Capper', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lvqtgkcxij', 'Joya', 'MacCheyne', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('deunzygbai', 'Hardy', 'Coatts', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cpayilkben', 'Andrea', 'Glasscock', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('naxpqlgkwh', 'Berky', 'Berston', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eybtmquszo', 'Johan', 'Jorge', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ozxsuhnjcb', 'Ariela', 'Ranaghan', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jioybwnpzx', 'Felicdad', 'Bisset', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('buwxgkfzdr', 'Parsifal', 'Gibbings', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mczvwkybxd', 'Kyrstin', 'Handke', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fhsioebdzc', 'Kath', 'Gonnely', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('crwglxidfy', 'Tiena', 'Bilsland', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dalbsgmovn', 'Carey', 'Swanston', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lfueaciokw', 'Major', 'Losemann', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rdklcwjfng', 'Winifield', 'Bradnick', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mpwbyzxrta', 'Marybelle', 'Kersey', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cgqopwreis', 'Caprice', 'Finby', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uitnfbmgxv', 'Bonnee', 'Rayson', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wetpudqibj', 'Stanislaw', 'Carlos', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cjuzqlrdoa', 'Vachel', 'Codner', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dhqwvzoixf', 'Ellerey', 'Pickerin', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ohmctpdvea', 'Darin', 'Farnorth', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('faqzetwdkx', 'Rance', 'Mannooch', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ieufqdzybj', 'Kimble', 'Hincks', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wmdngpyxez', 'Brana', 'Carvilla', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kvogjurbhx', 'Ailsun', 'Sabater', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jciybxqzlw', 'Agnola', 'Joreau', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jxlceaowky', 'Zena', 'Lethbrig', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qmhzjbfkce', 'Otis', 'Thoumasson', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rsondhuwke', 'Jackqueline', 'Eade', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mtrweshzbx', 'Dorthea', 'Cullinan', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vwxysfezlt', 'Maxwell', 'Sapseed', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wuztfxgsvc', 'Moe', 'Hembling', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qxhiluajcg', 'Alvina', 'Shovell', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eqnkiblzph', 'Ignaz', 'Beaudry', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('oxgqhnltbz', 'Glenda', 'Chadwyck', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('neslywdzva', 'Wenona', 'Chitson', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('brzsmxjqhf', 'Rodie', 'Twiddell', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tuwoscfhgd', 'Travus', 'Coatts', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qmlrfsbjgy', 'Munroe', 'Rotherforth', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bmvyceloru', 'Ryon', 'Crangle', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nwcrtmeupb', 'Ailey', 'Kaasman', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('haivmkpneo', 'Jamill', 'Casperri', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gsytoqnbir', 'Bertrand', 'Gleeton', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zqrolvntje', 'Rabi', 'Gladbach', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('flhaevkonc', 'Ring', 'Batchan', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mkdtrgoehl', 'Yankee', 'Ancliff', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gwzxtqyeij', 'Robin', 'Laidlow', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nxawikjmfr', 'Luciana', 'De Benedetti', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wigphoqxzm', 'Blinni', 'O''Dowling', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xmtzesjyaq', 'Jourdan', 'O'' Mulderrig', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uncaxgovjq', 'Jesse', 'Clarson', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jkihzqyucw', 'Renee', 'Kinchlea', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uyxabdcwjf', 'Cy', 'Metzing', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rcowqtenhi', 'Blaine', 'Broadis', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nmvyksogdh', 'Adrea', 'Tabbitt', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zmebqnjlou', 'Maryann', 'Gonoude', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fvdsmhoyxq', 'Starlin', 'Hinrichs', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cunshrqyxv', 'Anneliese', 'Robshaw', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zxigklretp', 'Tomasina', 'Robilliard', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('prtjcdakxo', 'Nita', 'Berney', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qniuoavwjf', 'Jozef', 'Elgie', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gofkcpmewb', 'Zenia', 'Cracknall', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('risnufodwg', 'Phil', 'Aylwin', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gcewjixpkb', 'Whitby', 'Dunphy', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fegkcpnzvw', 'Dur', 'Hardwin', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tsxqzgkvyw', 'Chloette', 'Kears', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cmkwfydaqr', 'Maribelle', 'Hargreves', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pnizmxwrou', 'Shayna', 'McKerlie', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xeoisfrtmv', 'Emmalee', 'Craker', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hynfxwdaui', 'Consalve', 'McCaughey', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rsayxotgve', 'Barbe', 'Elcoat', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mrwdjqzxou', 'Hamil', 'Prangley', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yulczibpat', 'Haydon', 'Digweed', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fqmzgcboue', 'Pavlov', 'Fydoe', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('isqdxbmhaj', 'Norman', 'Mair', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gtkbiuceza', 'Gunar', 'Riste', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pmxziydcjv', 'Karlee', 'Faustian', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vrxzosutcn', 'Katerina', 'Stansbury', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('znyrbimujs', 'Anabel', 'Iacovaccio', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('craxhpnzke', 'Alexia', 'Ellcock', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xpztrylsvg', 'Rickie', 'Trevarthen', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('egqctnpviz', 'Pooh', 'Cane', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qlboxiwuhr', 'Kristine', 'Spittles', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ftjbhmrcek', 'Elli', 'Domenici', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zhybpjxgds', 'Vanessa', 'Witts', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fevqworsxa', 'Oliviero', 'Petru', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('girowvujhd', 'Chadwick', 'Leek', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bljtwynzqc', 'Washington', 'Gronou', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kfjbrqzhxn', 'Philomena', 'Scruby', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gniymvurdt', 'Abbott', 'Leachman', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tqzypuckjh', 'Mal', 'Manton', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tvoxsqwndj', 'Hayyim', 'Dowyer', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('szgdriubny', 'Allie', 'Gawkroge', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dgkflvhrji', 'Konstantin', 'Lemm', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eboqjantdk', 'Malanie', 'Peregrine', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fzmakvotlw', 'Nev', 'Bougen', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qtnfcxbhoy', 'Lanni', 'Gellett', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kdmoliyefn', 'Constantia', 'Mounter', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ycsavgdwlq', 'Carena', 'Exroll', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pyltdmbfrz', 'Joleen', 'Headrick', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qrphuwozyv', 'Milty', 'Goldson', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lpbznjrhfq', 'Denyse', 'Geleman', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fmanliwvqe', 'Delmore', 'McQuillen', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('funghtqbro', 'Dori', 'Hynes', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mvyhnwrdpz', 'Rudolf', 'Fredy', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pcqfiwvjhg', 'Christiano', 'Euels', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ydhfmgcowi', 'Launce', 'Hembry', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fvqrbinejt', 'Kerrill', 'Pollastro', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rklcmaszhj', 'Aaren', 'Doog', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ljbvgauxsm', 'Erwin', 'Parsley', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mdnwkyfqxz', 'Shelley', 'Andrelli', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('peudkqsial', 'Christiane', 'Gunner', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fhekoxazjq', 'Coletta', 'Mccaull', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yfzsrqligh', 'Lillian', 'Teresi', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mplywturqs', 'Marlo', 'Astlatt', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gslnjztvhe', 'Audry', 'Dann', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tqiubspref', 'Ellsworth', 'Arrault', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vdmyanzopi', 'Jade', 'Lethby', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gelzwhbscu', 'Janice', 'Dealy', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('asweyikmdg', 'Merilyn', 'Suthworth', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bnmkutyxjh', 'Lita', 'Lisle', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eyhcqmzfpu', 'Tonia', 'Creamen', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sctozlbjfe', 'Giovanna', 'Isworth', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('flvzunihcj', 'Tobe', 'Shilstone', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dwyzlaqkec', 'Arline', 'Bentame', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('inotvgwdkq', 'Sharia', 'Scalera', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rzgwujsvkb', 'Marie', 'Grieswood', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rabhekwijc', 'Molly', 'Whitmarsh', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cfvljywrzh', 'Bryana', 'Worwood', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jrufbwpdmt', 'Jarad', 'Whitehair', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ebtfqckyxp', 'Noll', 'Swapp', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pjdueahmcr', 'Dulcy', 'Retallack', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lyxqpmtndr', 'Federica', 'Duckworth', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dsblrugtvi', 'Estelle', 'Phipps', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qkfisuoxny', 'Kristo', 'Haysman', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('foyjxtkzcd', 'Mora', 'Leisk', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nfskgrimeo', 'Marja', 'Aloway', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jlvbufqcwa', 'Howard', 'Rogers', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('plygsuediz', 'Ealasaid', 'O''Sesnane', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yvukgxcbqm', 'Nicki', 'Rimell', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('puozajlwbx', 'Shirleen', 'Elwin', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('drieljqncw', 'Thor', 'Dumbreck', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xcydjglmhv', 'Shae', 'Readwing', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('alkpotnisg', 'Kettie', 'Hazleton', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('weoydjfqng', 'Clare', 'Willder', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kyinjlgfsa', 'Inger', 'Edie', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('geazsklvxq', 'Elmira', 'Koenraad', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sglvurwhtn', 'Alyce', 'Brussell', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('coabjeikhv', 'Lisbeth', 'Divina', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('umdlbchonp', 'Desi', 'Geldeford', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pcoaigefvj', 'Faye', 'Himsworth', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rhlkcutoiy', 'Emeline', 'Guidera', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ujfeaopimg', 'Marcia', 'Regardsoe', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jmgrbontpx', 'Karena', 'Grigorian', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bnrwugjmzc', 'Blithe', 'Cordeix', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('egywarhpfd', 'Gun', 'Francke', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rqjonilxub', 'Benito', 'Potte', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dewqtpbvkf', 'Bartholomeus', 'Dishmon', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pqumzrnwef', 'Anastassia', 'Brimming', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pqflgxsnzy', 'Johna', 'Gennings', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ytmualqkxh', 'Benedikt', 'Smillie', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('thnifgczyj', 'Carmen', 'Crebo', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('amyeqvzbhk', 'Myrna', 'Bezants', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jvraxzsncd', 'Mahmud', 'Trimme', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xkgjafvswc', 'Dacia', 'Whodcoat', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rtikjpgonl', 'Ethan', 'O''Carmody', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('usiwgfrphq', 'Bryn', 'Keppie', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lekgphbwsq', 'Hartwell', 'Brownlea', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('egrbxqwoun', 'Evelyn', 'Hazeley', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gafqrolxjk', 'Bart', 'Nast', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bjnzthxuqo', 'Penny', 'Playle', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zdfpywrgol', 'Luis', 'MacParlan', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('etaqsvixun', 'Meris', 'Garron', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hoivzamrdp', 'Salmon', 'Farmar', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('isyhorzfux', 'Erny', 'Ambler', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xpryjomeab', 'Pansy', 'Kincla', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('btmyzjcfdv', 'Tallia', 'Pollins', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('grmxsadhol', 'Camala', 'Berni', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('iyodwubahv', 'Charlie', 'Westfalen', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qzoagwrjks', 'Issi', 'Oager', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('oasgcdnbrl', 'Kelila', 'Edwicke', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dpsmyxfbjz', 'Christine', 'Rotlauf', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tmwlbajynz', 'Corabella', 'Parkin', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eufksbngma', 'Bryan', 'Lefley', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nqchetkyoi', 'Tremayne', 'Oram', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zpfbnmcljq', 'Nickolai', 'Bonhan', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zqrvfepauj', 'Matelda', 'Denerley', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gscptheoyz', 'Jillene', 'Brothers', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tvdepbjoaz', 'Jacky', 'Minker', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('duopvtkhei', 'Berny', 'Seakes', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xzmjfkvido', 'Merna', 'Sparshutt', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fixuncgdkq', 'Margret', 'Aves', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vgfjdultix', 'Noel', 'Spelsbury', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bknoewhfji', 'Anette', 'Sadlier', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('iyfnweqoht', 'Franky', 'Brattan', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('whdspctxag', 'Verile', 'Cuer', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xpdeyumrqo', 'Rollo', 'Arundel', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mdjuvswtiy', 'Mozelle', 'Fussell', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uxleidmyon', 'Ekaterina', 'Swanbourne', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('asytbcqkiv', 'Pennie', 'McMurdo', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('autbcwhqgo', 'Appolonia', 'Mollitt', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gdnacpkfws', 'Aline', 'Reedie', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lkaxtpgyvi', 'Gay', 'Bullas', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kezxrhlwib', 'Brockie', 'Turland', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kafzlsiwvh', 'Annie', 'Allanby', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ictghxoraq', 'Deny', 'Thrasher', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('shnxdguoba', 'Loraine', 'Cackett', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dhaqcibmxg', 'Kendal', 'Garnsworthy', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('stjpuzlxfn', 'Hunfredo', 'Lawton', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('raoxblzypj', 'Mirilla', 'Ventura', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('otgferalyv', 'Thomasin', 'Frame', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('osycehinvu', 'Victoria', 'Kettles', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xvmbsrpuzd', 'Alina', 'Gratrex', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gzmouhqcwy', 'Jeanna', 'D''Hooghe', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pachsivjko', 'Brana', 'Hellmore', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gfvlrujhcy', 'Elfreda', 'Beuscher', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hbjlkicdys', 'Issi', 'Challener', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rvutelyzbm', 'Pippa', 'Tapenden', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mbroekidtv', 'Julee', 'McFayden', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fwkcdzopln', 'Nelia', 'Goodbur', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('selbnvcmkj', 'Rabbi', 'Abdy', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rmzfyhawxb', 'Dominick', 'Lampart', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tplikrsnxj', 'Meredith', 'Perotti', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hifnocmdaz', 'Siward', 'Lutty', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vrqjaxslut', 'Harrison', 'Rowling', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wgunotlkhe', 'Randi', 'Ivanchenkov', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rboyndvptf', 'Lilian', 'Murfin', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qujbigprvd', 'Brnaba', 'Coole', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bynihqpsaj', 'Sib', 'Stickland', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fitclhvbpz', 'Rosabel', 'Dalla', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hybwzietrk', 'Georas', 'Wild', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('icuesvodyr', 'Mercedes', 'Cavanagh', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('clhdseujor', 'Paddie', 'Offield', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gtozelcura', 'Moselle', 'Machon', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ibzpjnorgv', 'Deanna', 'Burgisi', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('efltjzvbmx', 'Fionna', 'Lynas', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gdutxmhzrp', 'Carey', 'Muckersie', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('swcjkyoebg', 'Micheal', 'Freeberne', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pufhyeozxm', 'Eunice', 'Warnock', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('juofsadzrb', 'Leonidas', 'Plante', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lkenbuhfro', 'Zonda', 'Pedrielli', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wfhiqdeybo', 'Kizzee', 'Tetford', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kvxswbgrci', 'Aida', 'Quantick', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rlqkoiumah', 'Elijah', 'Snowsill', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xvqaeoisbg', 'Zelma', 'Dearnley', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zvwrqdetbs', 'Rebe', 'Barnwell', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cehybzjtfa', 'Chris', 'Shelp', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sdqvtceywm', 'Clyde', 'Bleas', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mlrtazwdbc', 'Ewell', 'Bletsoe', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('oaxnieruby', 'Renate', 'Chaston', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rohjqaiczf', 'Penelope', 'Bole', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ryipbowdau', 'Evie', 'Gocke', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('entkomalvs', 'Yorgos', 'Coggeshall', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vgzjkfpteq', 'Joellen', 'Delacoste', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tdzpwymrlk', 'Karel', 'Rainville', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wskzgcrhqp', 'Kristien', 'Popland', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mewgrphqsy', 'Deedee', 'Bagster', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qsphewbfrx', 'Tamra', 'McNirlan', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bxnhrgtsjz', 'Symon', 'Stoffels', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cewyzrijgt', 'Wanda', 'Birdsey', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qbiwfm', 'Grannie', 'Fynan', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gcevyhrdxo', 'Arluene', 'Tull', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hwoixkrebz', 'Silvan', 'Housecroft', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('aeszkripco', 'Lissi', 'Briamo', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hatqozlecy', 'Malchy', 'Ramalho', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('luibxwetyr', 'Babb', 'McGrail', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jkvrthzefw', 'Wrennie', 'Robotham', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ylwfzrpcqn', 'Gwynne', 'Credland', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('brvmcfyzqx', 'Caron', 'Corteis', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pxjlzwokby', 'Luz', 'Culkin', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('upbwdnlagx', 'Hymie', 'Pooley', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('segmktulwi', 'Mia', 'Toun', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('itycxpajvn', 'Inger', 'Chaize', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('samqnuydir', 'Kristin', 'Caukill', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('byazcvunhk', 'Mitchell', 'Pancoust', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gyhxscjmfv', 'Tandy', 'Seres', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tvdgfzlqwo', 'Mose', 'Nimmo', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wofpsbgunl', 'Wolfy', 'Iwanczyk', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('abgsmyjwuf', 'Aidan', 'Balham', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('whtzdvflgj', 'Marley', 'Fockes', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gyczkpnhrw', 'Giustino', 'Delahunty', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ijvtpwchgd', 'Melody', 'Raftery', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xcinaymtdg', 'Leonardo', 'Jinkinson', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kxfqjupgav', 'Zackariah', 'Borman', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vzcemtqowy', 'Merola', 'Coulsen', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ryjivdqnkh', 'Mitzi', 'Lalevee', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pdtriyojkn', 'Vinni', 'Brokenshaw', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ndpgwzetja', 'Michele', 'Rioch', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xctsivjraq', 'Anderson', 'Pilling', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dfktyrzmwh', 'Ted', 'Housen', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('smkfiazuwr', 'Broderic', 'Wrixon', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jatozskivw', 'Mina', 'Cogin', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('enlbcsgyrq', 'Horatio', 'Lanahan', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('elgavkpzjx', 'Zuzana', 'Warwicker', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bsrwnqfmvi', 'Baxie', 'Borrie', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('darexzisyh', 'Ilsa', 'Stoke', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cwjqnepuhy', 'Mirabel', 'Krishtopaittis', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ycisruabgp', 'Phillis', 'McKinney', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gpfnzvmsqb', 'Izzy', 'Addams', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ukcwpeynvh', 'Parry', 'Vosse', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('luntbkdavr', 'Tyrone', 'Rowe', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lraovxujie', 'Darla', 'Cuerdall', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kuenhoxtry', 'Veronique', 'Janks', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('obzlywxacp', 'Alyse', 'Tallon', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qmpgjiayfo', 'Freddy', 'Pearn', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cuojdlevyt', 'Hillary', 'Olley', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('abvngzqfxc', 'Arleta', 'Fathers', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nybdphzgik', 'Parrnell', 'Asaaf', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('latvfhjbid', 'Rhiamon', 'Ginley', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rflehatmxj', 'Lotty', 'Bennet', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('revhkxuytq', 'Hanan', 'Dymoke', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('asifuywvzg', 'Norry', 'Box', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('havixwzvzg', 'Sidonnie', 'Rupel', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tyksrojwfu', 'Sheila', 'Rapinett', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kgfamtocqh', 'Lenette', 'Rens', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nltexvjois', 'Tom', 'Connock', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gdtzleiyfm', 'Wes', 'Skatcher', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yjvopgxecu', 'Stewart', 'Overill', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mbcxqfzywi', 'Powell', 'Shord', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('joretmvsxc', 'Hebert', 'Farrall', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hjuvenbqxi', 'Onida', 'Pile', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('muohyqabgf', 'Karil', 'Folland', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ynczdhwvfs', 'Allister', 'MacDonell', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fpychuvrbg', 'Windy', 'Issitt', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('acltobiwpe', 'Arluene', 'Samuel', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('oaceptgquy', 'Gilly', 'Shreenan', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sgmtixqulk', 'Godart', 'Delahunty', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vmljwrazpx', 'Herrick', 'Spriggs', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jebzgwtapf', 'Jennica', 'Sawrey', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qubzkwdjls', 'Brina', 'Callery', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nlwevfbsmo', 'Barbaraanne', 'Clemitt', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kotaluegsc', 'Shani', 'Doles', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zwrkmqxitc', 'Alasdair', 'Mazzey', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('adhqmouvbe', 'Chadwick', 'Mandifield', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gnqhsypzxw', 'Pamella', 'Vila', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('swxhqrlgdn', 'Peyton', 'Gothup', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hmtuwqdfjc', 'Lizbeth', 'Logan', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jcenxviasr', 'Miller', 'Thirkettle', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jgnehsbavf', 'Brad', 'Gueny', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bmrwgyakxd', 'Elvyn', 'Sacco', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pjgtceqdao', 'Reed', 'Cubitt', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jqhxvanrkl', 'Gail', 'Keppel', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yogepqdhsm', 'Hesther', 'Woolston', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vcrkontbae', 'Pebrook', 'Grisard', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('emptdqhkur', 'Doro', 'Dulake', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hlwkntqyxd', 'Hershel', 'Dowley', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gzowrpnvsf', 'Lisle', 'Ivakhnov', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bkzefdxcyv', 'Anstice', 'Rowly', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hvlznacumt', 'Maximilien', 'Harrigan', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nfovdhqxwt', 'Johanna', 'Ortas', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sirdumgefp', 'Pamella', 'Stitle', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kmtbruhqyp', 'Caddric', 'Tuckwell', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('daexszyvmn', 'Nicholas', 'Sentinella', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gzifptdeyo', 'Nero', 'Shillitto', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sjivhckrme', 'Shelly', 'Naisbitt', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('orghkiqwfl', 'Lionel', 'Ivison', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sqbhjzdlkn', 'Pascale', 'Muddiman', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jipwlkdagf', 'Eliza', 'Fortman', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wblhpkozry', 'Dayle', 'Bleesing', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jwpxfsdmzl', 'Jarid', 'Frostdyke', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sqcypdmjou', 'Mannie', 'Mace', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lcaqekpnjr', 'Shellie', 'Dakers', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dkjzgesitp', 'Merwin', 'Kightly', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rmzosjxugw', 'Suzann', 'Farrimond', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('avdbpycstx', 'Bonnee', 'McPolin', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mqpbukrywh', 'Adlai', 'Barnson', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uavjqfezmi', 'James', 'Rashleigh', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('irvadjfemy', 'Saidee', 'Syson', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xtzweorkfs', 'Erasmus', 'Roebuck', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('euzylixpnf', 'Lorry', 'Derwin', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fykjacrsub', 'Matthew', 'Balcombe', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bakzuwpisj', 'Goddart', 'Platts', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gawslyvorq', 'Iorgo', 'Deex', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sudrtecyxw', 'Julia', 'Rustadge', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xhqfpznier', 'Mord', 'Dewing', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fqyxtsiplj', 'Kary', 'Scamaden', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mrwcsuklgp', 'Isabel', 'Jaffa', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wtumlrzboc', 'Leandra', 'Corston', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('khbeizanyd', 'Rosalia', 'Millan', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jhqdzigofu', 'Manfred', 'Freshwater', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xloyanfzbh', 'Patin', 'Steen', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ltbdgepsun', 'Sherlocke', 'Djokic', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yxtzndwbfp', 'Carey', 'Cashman', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uzvgayfljm', 'Stanislaus', 'Strathman', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ecypqshizk', 'Waverley', 'Chainey', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sbemlzfrvh', 'Obed', 'Feare', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vtbnsfapzy', 'Tripp', 'Solway', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sbtcfvxaqn', 'Lonni', 'Dennistoun', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gfnrjuqawo', 'Lissy', 'Barford', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('raumdzvchl', 'Benjy', 'Oganian', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lamjpunswf', 'Janessa', 'Goshawk', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ofkywmghqu', 'Linette', 'Ibell', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ynxdouregz', 'Wald', 'Woodfin', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bvofrjzylw', 'Richie', 'Gelder', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dmecuibgtp', 'Jared', 'Marquand', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gxhoqtwmzu', 'Othilia', 'Maffione', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wcavkfuseh', 'Rosalinda', 'Gallardo', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('trldpfaoyv', 'Randell', 'Ralfe', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('amrvgfbxtl', 'Ettie', 'Phillipson', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fxicdtjoql', 'Patric', 'Spilling', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('devtkpumrs', 'Corry', 'Lympany', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qzgwrbvxad', 'Klaus', 'Tidball', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('niqtrsaxbm', 'Allard', 'Orhtmann', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dkolhawrcb', 'Loren', 'Gonzalvo', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ohasurpmde', 'Justus', 'Pourvoieur', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('iwcxpoqsyf', 'Raine', 'Bulbeck', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tfqolbumay', 'Der', 'Wykes', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yopaqfenlw', 'Terri', 'Milam', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('edaszgfubi', 'Earle', 'Lenz', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pcmuxizaql', 'Eleen', 'Goble', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yifsqulxgn', 'Chickie', 'Krollman', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wnvyqrjxbe', 'Abraham', 'Portt', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qhtodgemvk', 'Emmit', 'Parriss', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kpvjldqbfw', 'Ramsey', 'Booker', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wyfteiblvs', 'Talia', 'Ticehurst', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qbofteduhi', 'Gates', 'Breche', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jfexwchgbp', 'Lacey', 'Mazey', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('apcfqivnzo', 'Tania', 'Martindale', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vbcmdtywgj', 'Ferrel', 'Storm', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tuheryjasx', 'Zabrina', 'Abrahamian', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hxiawztcsj', 'Poppy', 'Dyet', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fvcelmgajx', 'Robinia', 'Mosley', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pynrejutgm', 'Hieronymus', 'Speeks', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uzoyrntgqv', 'Celle', 'Menezes', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qvduaypriz', 'Taddeusz', 'Farncombe', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hbmkjqsnoa', 'Dorie', 'Le Conte', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fbjlucvhqm', 'Clevey', 'Corradeschi', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hvgfapojcq', 'Kincaid', 'Brunker', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fcvtzpwksx', 'Nanette', 'Lembrick', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nbratpixwy', 'Alexandro', 'Baptiste', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hiyxqatlcr', 'Kaylyn', 'Blasio', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lkwcvturbx', 'Packston', 'Bette', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wxvzmjpfql', 'Roth', 'McDyer', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('blzycveinx', 'Aeriela', 'Dominiak', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uwrczmjqad', 'Martino', 'Emberton', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lciomkfueg', 'Tarrah', 'Tinson', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zqejnupxao', 'Rodd', 'Radoux', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zprxktmvey', 'Annissa', 'Stannis', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xsbjwkhyou', 'Gill', 'Tarply', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cgryfuetdb', 'Merv', 'Hourigan', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wkpsetvcdr', 'Vincenz', 'Hurran', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('iyhkrtgpev', 'Beitris', 'Mardall', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gvwekyldtx', 'Emery', 'Broek', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('egnkjlyvpa', 'Sibylle', 'Welbrock', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bsotnzjwfa', 'Augy', 'Rennison', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eyhwfqinlb', 'Austen', 'Sommers', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gzxuvtkidn', 'Happy', 'Cofax', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eocpsrwzmq', 'Madeline', 'Dinneges', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hmrytulozb', 'Myrvyn', 'Bodicam', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wciuyajftq', 'Frieda', 'Siss', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('omqpbnsrcj', 'Mallory', 'Trussman', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xmcyjbautk', 'Ware', 'De la Yglesia', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('szcpigbleo', 'Denny', 'Lillie', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jqyvefcxld', 'Matty', 'Ickovicz', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pxakfgdije', 'Stephi', 'Eschalotte', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eirulvhcox', 'Suzann', 'Tweddle', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fxulwkpcni', 'Rosemonde', 'Baselli', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ingkyrpqsw', 'Ardith', 'Garlant', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('knjdyzeufr', 'Darbee', 'Phaup', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wrmuihvdse', 'Rab', 'Hards', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qxjumrzspl', 'Sebastiano', 'Ciementini', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('urbifkspvn', 'Linus', 'Torres', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('aqmodkptwf', 'Chiquita', 'Daly', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fwjrvedbyz', 'Christa', 'Mavin', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pgjaxtekzn', 'Roarke', 'Strelitzer', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('twygxbckjv', 'Noellyn', 'Tremble', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('aihjoyfgzt', 'Aviva', 'Northcliffe', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mkuedpwjhz', 'Yuri', 'Gulc', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tbfikycwul', 'Riobard', 'Drillingcourt', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xnsubfiowe', 'Dew', 'Barnaby', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uekbwzhiyn', 'Berry', 'Lyngsted', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('txzaphbysk', 'Nicolas', 'Klimkovich', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ymugxidnzo', 'Kellie', 'Grewer', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tknwiedprl', 'Arthur', 'Van Rembrandt', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ltbuzgdihk', 'Cecil', 'Baynon', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nlutgzmkoq', 'Katusha', 'Jardine', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('biyjoxadcs', 'Alastair', 'Podd', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yjveiparql', 'Nelson', 'Phinnis', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ytxbdjgpui', 'Janie', 'Greening', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('auzwktflsx', 'Bertha', 'Esslement', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mlxgwiqpky', 'Philipa', 'Dubble', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pbhvjocsxu', 'Dani', 'Foster-Smith', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ecswqvkugn', 'Georgetta', 'Wyse', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zdxwfgorqp', 'Tatum', 'Buchett', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qixmvegctu', 'Ailene', 'Bentz', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('inexmgrvcl', 'Darrell', 'Grayne', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bgctmvdjlr', 'Verla', 'Eul', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('petawqvnxf', 'Amii', 'Fairhead', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yfcwetarkx', 'Jaquith', 'Copin', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kzagwqurps', 'Carr', 'Jessop', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zohlduwyai', 'Ealasaid', 'Raxworthy', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vxtolpzand', 'Agnesse', 'Killeen', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('onaifbhxsw', 'Alidia', 'Fullylove', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ioabdnqryc', 'Yorker', 'Vlahos', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xrynhfktzb', 'Harris', 'Bradneck', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vuheoacnbw', 'Bianka', 'Kaliszewski', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ulypwmtkbi', 'Flinn', 'Abele', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zxuvsrnglt', 'Elysha', 'Rabbatts', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pdrbnwaxsq', 'Morna', 'Leyzell', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pekismqzal', 'Andrea', 'Edmondson', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xdulzmybjt', 'Basilio', 'Bullan', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nezvqojusl', 'Katine', 'McKleod', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('oldmgfhxzr', 'Leah', 'Buie', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dqfurkoeyx', 'Tildie', 'Bursnall', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('facmxevutd', 'Dora', 'Wiggam', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('snvjkbxled', 'Yoshi', 'Goodboddy', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eqwoabxlgs', 'Lalo', 'Gioan', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ihrclkzqeo', 'Kahaleel', 'Bruck', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gmwblnkdjo', 'Nissie', 'Robjents', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vugdtbwjke', 'Dannie', 'Ellar', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vhkybomdwf', 'Matthus', 'Tregido', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qypjcghvzt', 'Miquela', 'Anthill', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ljfznrivgy', 'Happy', 'Tather', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('amojqiutsb', 'Kania', 'Juniper', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vtdxuzfbro', 'Job', 'Pattinson', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gnyvpekcou', 'Ruperta', 'Bleasdille', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wxkuhvcnbq', 'Tiertza', 'Sweetlove', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('exzurbjldw', 'Atlanta', 'Balazs', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jvefkporys', 'Garth', 'Imos', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nefyatzpqd', 'Free', 'Mattiessen', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dgtqloxwjk', 'Niki', 'Robic', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fomgzrlues', 'Cosimo', 'Medlicott', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tdnlzsbhmw', 'Gunther', 'Forker', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eqiljzpduh', 'Elmira', 'Talloe', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qepdtfzbjy', 'Jazmin', 'Rootham', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wmxlaofijq', 'Dana', 'Caro', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kijprgewud', 'Duky', 'Tripe', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qusgyfnlxb', 'Cordell', 'Mapes', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gxmkqucvpr', 'Rriocard', 'Diano', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rxhsmdwjqv', 'Clarita', 'Darrell', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vamednlukq', 'Wyatt', 'Leakner', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rvbplunwhe', 'Lonnard', 'Ligerton', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hcbotzmdsi', 'Ambrosi', 'Viner', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('czavigkebq', 'Elinore', 'Loving', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hracgpebxk', 'Ibbie', 'Bilbee', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mawqvposcr', 'Ronnica', 'Cran', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gnvmctkhuq', 'Agata', 'Duckett', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cngbdamytq', 'Verna', 'Zanni', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wmhoxitqgr', 'Clyve', 'Stallebrass', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tnqhmflwzy', 'Prue', 'Knatt', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lgtdaumncq', 'Giuseppe', 'Filson', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('myfvtpjnqx', 'Torry', 'Bickerton', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zdvnlsyukf', 'Merrile', 'Madle', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('aohkzbymew', 'Brigitta', 'Ensley', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nwpstxbgku', 'Flint', 'Egalton', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cdjqowbmpf', 'Harv', 'Bennetto', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('imxejtonah', 'Hilliary', 'Devoy', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rbzudoqcni', 'Spenser', 'Ludy', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wsnxeikrpt', 'Adorne', 'Turl', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fturkxgsqi', 'Ody', 'Alvaro', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jotieycavd', 'Dion', 'Rembaud', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pcsbzavitn', 'Erinna', 'Borham', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ouvklyrajd', 'Dasha', 'Hayen', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bsmohtvaqk', 'Duke', 'Santora', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tivfcnazrb', 'Jere', 'Garrick', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vobnqmtrzx', 'Ilsa', 'Thewles', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('auxjnvmpkb', 'Jackie', 'Darnbrook', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zmwdsygcau', 'Hilarius', 'Skeat', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nfxsakgtzy', 'Gloria', 'Morrill', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tywgbxiuoz', 'Theo', 'Peasee', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dtikozcpnq', 'Danya', 'Parfrey', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xouzmaqswg', 'Massimo', 'Dallimare', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rmkcbwuqoi', 'Abramo', 'Kleinsinger', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lwjhbmroie', 'Annetta', 'Lardez', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xfntpobevi', 'Georas', 'Twidale', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ohjvufedmt', 'Wally', 'Trusty', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yepjzbhacw', 'Baryram', 'Mirams', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eiwxafjhkq', 'Analiese', 'Brandts', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sayofvmcbx', 'Linoel', 'Calltone', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tvzqickwed', 'Bogey', 'Bottom', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ruznykjmlw', 'Hyacinth', 'Finnimore', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tbwnadkirc', 'Gerrie', 'Karlolczak', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('koyswgctrf', 'Gabriel', 'Zeal', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xpkndvlazi', 'Elizabeth', 'Rolston', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vizmwelunt', 'Vittoria', 'Cluett', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ayiqsvwfct', 'Bayard', 'Gasparro', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ocenupjkwt', 'Doralyn', 'Onraet', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rqmwoifzal', 'Codie', 'Leidecker', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fubgmnvecp', 'Marice', 'Keyser', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fmtzuahnbk', 'Adan', 'Matchett', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nhzlyxkpaf', 'Pollyanna', 'McInility', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('izobrqjxmt', 'Xerxes', 'Barefoot', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('oityrvhmkl', 'Linnet', 'Dunkinson', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('iodhxpewbu', 'Mikel', 'Chetwin', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('onvziptseh', 'Pearce', 'Hunsworth', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zjkrevnxua', 'Britt', 'Mohan', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pkxjbgtqrm', 'Fonzie', 'Gurrado', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('locyqahvwr', 'Earle', 'Minelli', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jmspxirckq', 'Philis', 'Munkton', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yfwbrchsqj', 'Amil', 'Hearthfield', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qlwyvhirex', 'Tammy', 'Ramstead', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vburtwxnhg', 'Belva', 'Mion', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fvqcbjdrxn', 'Gustavo', 'Stainton', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tpweonkdhv', 'Tani', 'Moulding', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jofcliwpaq', 'Arda', 'Summerhayes', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('eslzackhyj', 'Gayle', 'Rummings', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ozsrndqubx', 'Penny', 'Boskell', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('axzjnrsfmk', 'Nananne', 'Golagley', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zqmsryebnt', 'Johanna', 'Cantu', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('binrtxfzwm', 'Gal', 'Billison', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xsofbgvcti', 'Mick', 'Slainey', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cloishxdyn', 'Pembroke', 'Midgley', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jhyxgfiskl', 'Nadia', 'Newlove', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ovxdpakfrs', 'Decca', 'Iacovuzzi', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vmgrdqsice', 'Cathryn', 'Chilles', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fkvodhsqyu', 'Bar', 'Correa', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pvogwfibry', 'Alec', 'Akaster', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hofupbtcmk', 'Gabi', 'Hirth', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wuopsyfqde', 'Tedd', 'Elan', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fdamiovsqz', 'Axe', 'Lesor', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pobfslrncd', 'Carley', 'Garroway', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bxuegflynw', 'Anallise', 'Das', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cykajgqdpf', 'Neddy', 'Libby', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fdvwchyatp', 'Stanislaw', 'Burcombe', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gwheyokftx', 'Archer', 'Sooper', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lfztgayuxj', 'Etan', 'Duffitt', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nayhkudjvz', 'Douglass', 'Gillison', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uiocsbqexh', 'Bud', 'Layland', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uoikxpgtes', 'Angil', 'Eddy', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rcdtqbyxvz', 'Roanne', 'Scholling', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gwibjphyxl', 'Alexine', 'Leppingwell', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('orkltncuqi', 'Norma', 'Yearby', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vyizfqhlax', 'Quint', 'Penvarne', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jampodcwln', 'Nero', 'Dineen', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vaobtwfzxl', 'Nana', 'Reckus', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('adgyperkji', 'Loise', 'Brusin', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bydxufhrmn', 'Dulcinea', 'Saunter', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qjbmodhufz', 'Jere', 'Boatwright', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('seaykmlzfc', 'Cinnamon', 'Good', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kfmjcurbzx', 'Zea', 'Soppeth', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lpaqriudkz', 'Lauritz', 'Midgely', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('clijstbyek', 'Scarlett', 'Pizzie', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qbagcsekyw', 'Samuel', 'Petegrew', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hrqxvzituf', 'Ileane', 'Goodhew', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hgrcxlemdi', 'Mia', 'Hardy-Piggin', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('oltraynvxq', 'Conni', 'Fennelow', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xricwvdsgl', 'Baryram', 'Gusney', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vtofzqknhm', 'Chelsae', 'Bess', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cslepraukg', 'Julina', 'Rown', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hypcsvaolx', 'Sharyl', 'Santore', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('afiyvmzgwb', 'Kial', 'Rapaport', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tmdgkcensa', 'Peggi', 'Duddy', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zjgdvcqumk', 'Jacqui', 'Filippo', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('csetqymxnb', 'Carolina', 'Gathercole', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kumabtojrx', 'Saree', 'Scotney', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gkdheviaoy', 'Moses', 'Farreil', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zivsdgcynu', 'Eleanore', 'Bugge', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zyasmhjqeg', 'Lynnet', 'Harteley', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zwkhtnlrdp', 'Amaleta', 'McLemon', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('munsigovpb', 'Brier', 'Gillard', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kqrmnjiozw', 'Pedro', 'Hairs', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kuqsgwimbp', 'Trixie', 'Shambroke', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rzjnyvcloh', 'Rodrick', 'Jemmett', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uzgphkbfem', 'Kara-lynn', 'Rumford', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jwgurbfztm', 'Zsa zsa', 'Guilloneau', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lwsjrautbd', 'Charley', 'Hambright', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pyujwthfbl', 'Nikos', 'Shufflebotham', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cqbrydeuzk', 'Stefano', 'Djurdjevic', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fybizprndk', 'Kat', 'Jefferson', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xkbwmgfsnz', 'Ivett', 'Webburn', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('stcyxbfkdp', 'Luciano', 'Blaskett', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qjzbpvlmkx', 'Mile', 'Cochrane', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qliopmvygz', 'Granthem', 'Yeld', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bokztpdjsh', 'Krista', 'Wythill', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wxbhjtcksl', 'Rosie', 'Shevlane', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zposjdvamf', 'Corenda', 'Pendell', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('iskvhulxzc', 'Clevey', 'Wofenden', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wgaihzkpue', 'Lindi', 'Carayol', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tqmogfvlkz', 'Iseabal', 'Puttick', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('smkbpajxgo', 'Klara', 'Learmonth', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('npjwgysxhu', 'Nathalia', 'Greir', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xvolbmfzwk', 'Glynda', 'McCue', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kgapjrismc', 'Darsey', 'Heinsen', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sjnazdhwub', 'Daphene', 'Bravery', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pwghdjvuaf', 'Evered', 'Henrys', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zhimgcksrw', 'Templeton', 'Tothacot', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bmsicvpakn', 'Mil', 'Wilmington', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gkdhzxsqyc', 'Mabel', 'Ripon', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tznmegqipw', 'Ignaz', 'Tremmil', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fuhmsoaxtw', 'Hattie', 'Scutchin', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gxmlfjhbni', 'Kaja', 'Blindmann', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lqygmkosfd', 'Colly', 'Riggott', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qfisgrnuow', 'Shell', 'O''Ruane', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lkfyimcjha', 'August', 'Mahedy', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pnbhmzlqyo', 'Haslett', 'Bickerdike', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dfukirvejg', 'Onfroi', 'Presman', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('evnqhpwjti', 'Gonzalo', 'Maypowder', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xbhtpnaslq', 'Laurena', 'Sahlstrom', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vlrtxkoqac', 'Aldus', 'Denzey', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wuqjivrgdz', 'Isaiah', 'Mingardi', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wenmlsozhv', 'Robyn', 'Popescu', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vrlazidumt', 'Ingra', 'Aylin', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zivxeduoft', 'Killian', 'Dorr', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fyjlsgioqh', 'Bradan', 'Doubleday', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fznbygkqlx', 'Corey', 'Corkill', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gtrvdncwxi', 'Adara', 'Mustard', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nbogwzvhec', 'Jarid', 'Harlock', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ryohtdvpgx', 'Merrili', 'Fawdery', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ijploenxym', 'Felisha', 'Prior', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('riyqjuobpt', 'Lorine', 'Sarsons', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ikjhvngrbz', 'Udall', 'Sandcraft', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fxulopdsyz', 'Tracee', 'Chavey', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hleqyajgou', 'Clarabelle', 'Palfreeman', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('oviejugzfx', 'Erin', 'Winborn', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hcjbdfqokv', 'Timi', 'Shann', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kumsnfjixp', 'Rusty', 'Driffill', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hkwajodglu', 'Dolores', 'Rook', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fsjdlgunph', 'Stavro', 'Joderli', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ybupzcorwj', 'Diena', 'Goalley', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gmazjhryfo', 'Annetta', 'Penberthy', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('avjntmkqrx', 'Maighdiln', 'Barkaway', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pedtvynclo', 'Dulcinea', 'Brafferton', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xncdpkqjau', 'Veradis', 'Hatherley', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('iagwyurpnv', 'Rebeca', 'Houlahan', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wlfmcynvsx', 'Kylen', 'Howgego', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('sqzrmpjuvf', 'Judi', 'Freegard', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ipnwkvlbfh', 'Perice', 'Spykings', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('gdjbinrmyp', 'Jermain', 'Hickinbottom', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('izowpvtmhe', 'Eula', 'Jaskiewicz', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xatluiocey', 'Paolo', 'Brien', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xnkfpodrch', 'Ivan', 'Gorman', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('inokqatsum', 'Cissy', 'Adds', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jsclqedmyn', 'Dulcea', 'Ozelton', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ikvqstfhgx', 'Michelle', 'Lakeman', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ihkxgopuct', 'Courtnay', 'Bridgland', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bkrhczmvfw', 'Merline', 'Frizell', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rdpycgnuvi', 'Caressa', 'Carlens', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vfablkoyex', 'Kienan', 'Cleary', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qrgudizxks', 'Myles', 'Yanne', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('jdipshyvrl', 'Maryanna', 'Leggin', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hbrwvtumsz', 'Cheslie', 'Garrod', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('inmusbkgra', 'Nanon', 'Slyde', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mpgnvkqtrc', 'Darby', 'Radclyffe', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('cvxeuapksd', 'Ring', 'Bleesing', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('piuqfolcgn', 'Darsie', 'Pilcher', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tagivsuyjp', 'Lizzy', 'Crinidge', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uzqgaelbjt', 'Baron', 'Jennaroy', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kpnuxyzaol', 'Brynn', 'Giacobini', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('miswfkdpbq', 'Jodi', 'Dane', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dfkvoeuaqm', 'Ingelbert', 'Emblem', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('awoguyezbr', 'Zulema', 'Tattershaw', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lizbqstyme', 'Nellie', 'Newnham', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zoflhaqbtw', 'Noami', 'Shadfourth', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('esbzxqgntm', 'Nathanil', 'Thorpe', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dfpjxgqcub', 'Claudio', 'Reeder', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mueryndvtg', 'Natalie', 'Fairnie', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bnhcrdsvgw', 'Germaine', 'Hawkswood', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('voyawciqeb', 'Bartlett', 'Hunn', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yezbnuitrw', 'Martynne', 'Smithson', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ntqoupkiwg', 'Arleen', 'Woolaghan', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tezvkhxacp', 'Kelwin', 'Baff', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('hsbxgdflmr', 'Dalia', 'Bartaloni', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pmnrotyuqa', 'Mathilde', 'Enevoldsen', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('nashblvkyp', 'Boycey', 'Sutcliff', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('deikwxoczb', 'Simeon', 'Sjollema', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('yuvozgklrn', 'Nevsa', 'McPheat', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bigyrlkmfd', 'Edi', 'Pull', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('dfwbspnljq', 'Marcelia', 'Erwin', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('fpcrghuwez', 'Otho', 'Reddy', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('giwkhdbaox', 'Danice', 'McGirr', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qyhkuerxdt', 'Alonzo', 'Southern', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wmsxltzvfh', 'Onofredo', 'Godmer', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('qndxvpwmtg', 'Doy', 'Yule', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zupyfntigs', 'Pasquale', 'Joanaud', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xbaqnzyimv', 'Devina', 'Simmill', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('rfvskheupd', 'Niki', 'Waszczyk', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ognhtxvqmj', 'Emlynne', 'Krolik', '2018-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('crmtgiasno', 'Ivy', 'Illyes', '2019-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ervobtycxn', 'Carey', 'Huggon', '2009-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bsymwqhjot', 'Rici', 'Spofford', '2010-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('usmikpceah', 'Dynah', 'Wolvey', '2011-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('brwkimetox', 'Emelda', 'Vye', '2012-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('htscbjpadw', 'Danella', 'Loades', '2013-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xwhcyikamq', 'Darci', 'Deniset', '2014-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xcyjdeiwgq', 'Obediah', 'Eade', '2015-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tprsvmbdxw', 'Abbey', 'Wrey', '2016-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('lqktwexcmo', 'Amabel', 'Bownas', '2017-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('swqxkvufme', 'Cammy', 'Turbefield', '2018-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('loepwzfdcn', 'Gilligan', 'Crosfeld', '2019-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ojfhxqkzdt', 'Madeleine', 'Warn', '2009-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vpzuymbrnf', 'Milt', 'Melsom', '2010-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('kozbtxsadp', 'Mark', 'Nason', '2011-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('tiomdxbgsf', 'Tawsha', 'Greenrodd', '2012-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('uvazhjdmlq', 'Sancho', 'Stockbridge', '2013-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('vmkcgwtura', 'Cyrus', 'Seak', '2014-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xtcprfuyva', 'Carree', 'Syms', '2015-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('wscpohzynv', 'Reece', 'Biernacki', '2016-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('mobyqurhjw', 'Travers', 'Heiden', '2017-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('pozdkrlusw', 'Virgie', 'Sinclaire', '2018-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ketoqszfxu', 'Xylia', 'Dayley', '2019-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('zibelvqhoj', 'Netta', 'Colaton', '2009-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('txovmwgyhj', 'Adler', 'Pulhoster', '2010-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('xjcatwzyhr', 'Myrilla', 'Jeroch', '2011-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('emptdlhvbx', 'Shina', 'Spelman', '2012-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('iwzaycqflo', 'Burl', 'Addionizio', '2013-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ljagvtzymc', 'Sasha', 'Jellybrand', '2014-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('edmjblpcwo', 'Odilia', 'Carass', '2015-09-01', 'a');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('bnvplfigyh', 'Marna', 'Birrell', '2016-09-01', 'b');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('ynatikdxwl', 'Sax', 'Gerin', '2017-09-01', 'c');
insert into pupil (id, first_name, last_name, class_date_formed, class_letter)
values ('svhlfpzuad', 'Bourke', 'Simmill', '2018-09-01', 'a');

-- insert 30 events to the event table
insert into event (id, date, name, resources_allowed)
values ('enbucokdvp', '2016-01-20', 'Public-key secondary', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('jexudlkomc', '2016-02-15', 'Triple-buffered', '{gadgets, paper}');
insert into event (id, date, name, resources_allowed)
values ('dqfvbtcwmo', '2016-03-12', 'Networked intranet', '{paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('chwojmnvad', '2016-04-05', 'Self-enabling portal', '{plastic}');
insert into event (id, date, name, resources_allowed)
values ('eivxtdykpb', '2016-05-17', 'Diverse background', '{gadgets, plastic}');
insert into event (id, date, name, resources_allowed)
values ('skyjlgriwz', '2017-09-15', 'Stand-alone', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('cqsbntoihw', '2017-10-11', 'Advanced user-facing', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('tbyqixcupe', '2017-11-11', 'Cross-group executive', '{gadgets, plastic}');
insert into event (id, date, name, resources_allowed)
values ('fsbgjwcqpr', '2017-12-20', 'Open-architected', '{gadgets, paper}');
insert into event (id, date, name, resources_allowed)
values ('lcqhmxbzka', '2018-01-10', 'Balanced algorithm', '{gadgets}');
insert into event (id, date, name, resources_allowed)
values ('nxmqsiluhz', '2018-02-19', 'Fully-configurable', '{paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('hcrtnasmeg', '2018-03-25', 'Monitored uniform', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('tupmhyqkwe', '2018-04-08', 'Reverse-engineered', '{gadgets, plastic}');
insert into event (id, date, name, resources_allowed)
values ('awsujqglvz', '2018-05-23', 'Phased encoding', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('nbghxfwjzc', '2018-10-26', 'Inverse structure', '{paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('jpsqlhrovz', '2018-12-09', 'Expanded toolset', '{paper}');
insert into event (id, date, name, resources_allowed)
values ('dpskbjhxet', '2019-01-12', 'Digitized extranet', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('gxlwktivau', '2019-02-03', 'Fully-configurable hub', '{gadgets, plastic}');
insert into event (id, date, name, resources_allowed)
values ('jfqsekarbd', '2019-03-25', 'Optional flexibility', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('bklzqpoihu', '2019-04-11', 'De-engineered hub', '{gadgets, paper}');
insert into event (id, date, name, resources_allowed)
values ('sfcyxzvtoa', '2019-10-22', 'Configurable hardware', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('avmceihjpu', '2019-12-21', 'Total circuit', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('zbvksgadcm', '2020-02-25', 'Self-enabling coronavirus', '{gadgets, plastic}');
insert into event (id, date, name, resources_allowed)
values ('rqblmjutph', '2020-03-15', 'Coronable panic', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('svhbzdogqu', '2020-04-24', 'Virus corona', '{paper}');
insert into event (id, date, name, resources_allowed)
values ('sdovdvyuiw', '2020-05-09', 'Configurable structure', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('tdgljkiozk', '2020-09-15', 'Inverse algorithm', '{paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('pdhanqomso', '2020-11-11', 'Cross-group intranet', '{gadgets}');
insert into event (id, date, name, resources_allowed)
values ('cgxdgdfqgs', '2021-01-28', 'Balanced uniform', '{gadgets, paper, plastic}');
insert into event (id, date, name, resources_allowed)
values ('mufirelawx', '2021-02-10', 'Advanced portal', '{gadgets, paper, plastic}');

-- Insert pupils and events interactions into the resource table. Each pupil has 80% chance to attend an event and 50%
-- chance to bring some recyclables to that event
insert into resources(pupil_id, event_id, paper, plastic, gadgets)
select pupil_id,
       event_id,
       case
           when 'paper' = any (resources_allowed) and random() < 0.5 then (random() * 100)::numeric(9, 3)
           else 0
           end,
       case
           when 'plastic' = any (resources_allowed) and random() < 0.5 then (random() * 100)::numeric(9, 3)
           else 0
           end,
       case
           when 'gadgets' = any (resources_allowed) and random() < 0.5 then (random() * 100)::numeric(9, 3)
           else 0
           end
from (
         select e.id as event_id, p.id as pupil_id, e.resources_allowed as resources_allowed
         from event as e
                  inner join pupil as p on p.class_date_formed <= e.date
             and p.class_date_formed >= e.date - interval '11 years'
             and e.date <= now()
             and random() < 0.8
     ) as t
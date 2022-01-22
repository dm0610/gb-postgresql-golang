Практические задания

1.Развернуть сервер PostgreSQL в Docker.
<p>Ответ: 
<p>

```
~/Documents/GeekBrains/GB-Postgres/gb-postgresql-golang/HW01$ ./01.sh 
261d38607dd77c42a4ffd42283f0f6273b452baeb30b37f62b221b65f98857fa

~/Documents/GeekBrains/GB-Postgres/gb-postgresql-golang/HW01$ docker ps
CONTAINER ID   IMAGE           COMMAND                  CREATED         STATUS         PORTS                                       NAMES
261d38607dd7   postgres:13.1   "docker-entrypoint.s…"   7 seconds ago   Up 6 seconds   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp   postgres

~/Documents/GeekBrains/GB-Postgres/gb-postgresql-golang/HW01$ ls -alt
total 20
-rw-rw-r--  1 dmvstrelnikov    dmvstrelnikov 2784 янв 16 18:17 readme.md
drwx------ 19 systemd-coredump root          4096 янв 16 18:14 mntdata
-rwxrwxr-x  1 dmvstrelnikov    dmvstrelnikov  228 янв 16 18:13 01.sh
drwxrwxr-x  3 dmvstrelnikov    dmvstrelnikov 4096 янв 16 18:13 .
drwxrwxr-x  4 dmvstrelnikov    dmvstrelnikov 4096 янв 16 18:10 ..
dmvstrelnikov@dmvstrelnikov-VirtualBox:~/Documents/GeekBrains/GB-Postgres/gb-postgresql-golang/HW01$ docker exec -it postgres bash
root@261d38607dd7:/# echo $PGDATA
/var/lib/postgresql/data
root@261d38607dd7:/# psql -d postgres -U postgres -h localhost -p 5432
psql (13.1 (Debian 13.1-1.pgdg100+1))
Type "help" for help.
postgres=# 
```
<p>2.Создать пользователя и базу данных.
<p>Ответ:
<p>

```
postgres=# create database users;
CREATE DATABASE
postgres=# create user techuser  password 'tEchpassW!!';
CREATE ROLE
postgres=# grant all privileges on database users to techuser;
GRANT
postgres=# \q
root@261d38607dd7:/# psql -U techuser -d users
psql (13.1 (Debian 13.1-1.pgdg100+1))
Type "help" for help.
users=> \conninfo
You are connected to database "users" as user "techuser" via socket in "/var/run/postgresql" at port "5432".
users=> 
```
<p>3.В базе из пункта 2 создать таблицу: не менее трёх столбцов различных типов. SQL-запрос на создание таблицы добавить в текстовый файл class1_hometask.txt.
<p>Ответ:
<p>

```
users=> CREATE TABLE personaldata
users-> ( user_id        serial PRIMARY KEY,
users(>   user_firstname char(50) NOT NULL,
users(>   user_lastname  char(50) NOT NULL,
users(>   user_email     char(50)
users(> );
CREATE TABLE
```
<p>4.В таблицу из пункта 3 вставить не менее трёх строк. SQL-запрос на вставку добавить в текстовый файл class1_hometask.txt.
<p>Ответ:
<p>

```
users=> INSERT INTO personaldata (user_firstname, user_lastname, user_email) VALUES 
('Artem', 'Ivanov', 'a.ivanov1@mail.ru' );
INSERT 0 1
users=> select * from personaldata;
 user_id |                   user_firstname                   |                   user_lastname                    |                     user_email                     
---------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------
       1 | Artem                                              | Ivanov                                             | a.ivanov1@mail.ru                                 
(1 row)

users=> INSERT INTO personaldata (user_firstname, user_lastname, user_email) VALUES 
users-> ('Dmitry', 'Semenov', 'd.semenov3@mail.ru' ),
users-> ('Ivan', 'Demchenko', 'i.demchenko1@mail.ru' );
INSERT 0 2
users=> select * from personaldata;
 user_id |                   user_firstname                   |                   user_lastname                    |                     user_email                     
---------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------
       1 | Artem                                              | Ivanov                                             | a.ivanov1@mail.ru                                 
       2 | Dmitry                                             | Semenov                                            | d.semenov3@mail.ru                                
       3 | Ivan                                               | Demchenko                                          | i.demchenko1@mail.ru                              
(3 rows)

users=> 
```

<p>5.Используя мета-команды psql, вывести список всех сущностей в базе данных из пункта 2. Полученный список сущностей добавить в текстовый файл class1_hometask.txt.
<p>Ответ
<p>

```
\! clear
users-> \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 users     | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =Tc/postgres         +
           |          |          |            |            | postgres=CTc/postgres+
           |          |          |            |            | techuser=CTc/postgres
(4 rows)

users-> \d
                    List of relations
 Schema |           Name           |   Type   |  Owner   
--------+--------------------------+----------+----------
 public | personaldata             | table    | techuser
 public | personaldata_user_id_seq | sequence | techuser
(2 rows)

users-> \dt
            List of relations
 Schema |     Name     | Type  |  Owner   
--------+--------------+-------+----------
 public | personaldata | table | techuser
(1 row)

users-> \dv
Did not find any relations.
users-> \dm
Did not find any relations.
users-> \di
                      List of relations
 Schema |       Name        | Type  |  Owner   |    Table     
--------+-------------------+-------+----------+--------------
 public | personaldata_pkey | index | techuser | personaldata
(1 row)

users-> \dn
  List of schemas
  Name  |  Owner   
--------+----------
 public | postgres
(1 row)

users-> \dT
     List of data types
 Schema | Name | Description 
--------+------+-------------
(0 rows)

users-> \x
Expanded display is on.
users-> \dT
(0 rows)

users-> \set
AUTOCOMMIT = 'on'
COMP_KEYWORD_CASE = 'preserve-upper'
DBNAME = 'users'
ECHO = 'none'
ECHO_HIDDEN = 'off'
ENCODING = 'UTF8'
ERROR = 'false'
FETCH_COUNT = '0'
HIDE_TABLEAM = 'off'
HISTCONTROL = 'none'
HISTSIZE = '500'
HOST = '/var/run/postgresql'
IGNOREEOF = '0'
LASTOID = '0'
LAST_ERROR_MESSAGE = 'syntax error at or near "rable"'
LAST_ERROR_SQLSTATE = '42601'
ON_ERROR_ROLLBACK = 'off'
ON_ERROR_STOP = 'off'
PORT = '5432'
PROMPT1 = '%/%R%x%# '
PROMPT2 = '%/%R%x%# '
PROMPT3 = '>> '
QUIET = 'off'
ROW_COUNT = '3'
SERVER_VERSION_NAME = '13.1 (Debian 13.1-1.pgdg100+1)'
SERVER_VERSION_NUM = '130001'
SHOW_CONTEXT = 'errors'
SINGLELINE = 'off'
SINGLESTEP = 'off'
SQLSTATE = '00000'
USER = 'techuser'
VERBOSITY = 'default'
VERSION = 'PostgreSQL 13.1 (Debian 13.1-1.pgdg100+1) on x86_64-pc-linux-gnu, compiled by gcc (Debian 8.3.0-6) 8.3.0, 64-bit'
VERSION_NAME = '13.1 (Debian 13.1-1.pgdg100+1)'
VERSION_NUM = '130001'
users-> 

```

<p>6.Придумать проект, над которым вы будете работать в последующих уроках. Примером может быть любой сайт/приложение, которым вы пользуетесь (YouTube, Medium, AirBnB, Booking, Twitter, Instagram, Facebook, LinkedIn). Это может быть соц. сеть, блог, книга рецептов, база данных авиаперелетов, мессенджер, система бронирования бань и саун и т.п.
<p>Ответ:
<p>

```
Локальный небольшой аналог heroku.
```

<p>7.Кратко (не более 10 предложений) описать суть проекта и основной use-case в файле schema.sql (описывать как sql комментарий в начале файла).
<p>Ответ:
<p>

```
1. Будет несколько сервисов. Гитхаб, dockerhub, jenkins, k8s(yandex, например)
2. Я создаю в базе в таблице projects новую строку(id,prName) - новый проект. 
3. В других таблицах, отвечающих за сервисы создаю строки с id,name,expirationDate и т.п.
4. Сферические микросервисы в вакууме мониторят базу и после появления нового "инстанса" создают его реальную копию на заданных ресурсах.
5. Пока не представляю как это реализовать, но всю жизнь мечтал создать своё облако :) 
```

<p>8.Разработать структуру базы данных, которая будет фундаментом для выбранного проекта (не менее трёх таблиц, не более 10 таблиц). В структуре базы данных желательно иметь логические связи между сущностями (не менее одной связи). Команды на создание таблиц описать в файле schema.sql.
<p>Ответ:
<p>

```
create database projects;

CREATE TABLE projects (
 id SERIAL PRIMARY KEY,
 name TEXT NOT NULL,
 owner_email  TEXT NOT NULL UNIQUE
);

CREATE TABLE services (
 id SERIAL PRIMARY KEY,
 service_name TEXT NOT NULL UNIQUE,
 service_address TEXT NOT NULL UNIQUE,
 service_techuser TEXT NOT NULL UNIQUE,
 service_passwd TEXT NOT NULL UNIQUE
);

CREATE TABLE instances (
 id SERIAL PRIMARY KEY,
 service_id INTEGER NOT NULL UNIQUE,
 instance_name TEXT NOT NULL UNIQUE,
 instance_expire_time TEXT NOT NULL UNIQUE,
 project_id INTEGER NOT NULL UNIQUE,
 FOREIGN KEY (service_id) REFERENCES services(id),
 FOREIGN KEY (project_id) REFERENCES projects(id)
);

projects=# \dt+
                               List of relations
 Schema |   Name    | Type  |  Owner   | Persistence |    Size    | Description 
--------+-----------+-------+----------+-------------+------------+-------------
 public | instances | table | postgres | permanent   | 8192 bytes | 
 public | projects  | table | postgres | permanent   | 8192 bytes | 
 public | services  | table | postgres | permanent   | 8192 bytes | 
(3 rows)

projects=# 
```
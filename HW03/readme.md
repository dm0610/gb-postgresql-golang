<p>Практические задания

<p>1. Составить 3–5 типовых запросов к данным в созданном проекте БД. Описать их в файле queries.sql.
<p>Ответ:

```
SELECT instance_name, instance_description FROM instances WHERE instance_name LIKE '%jenkins' ORDER BY instance_name;


SELECT 
 (SELECT title FROM projects WHERE projects.id = instances.project_id)
   AS project_title,
 instance_name,
 (SELECT owner_email FROM projects WHERE projects.id = instances.project_id)
   AS owner_email,
 (SELECT service_title FROM services WHERE services.id = instances.service_id)
   AS service_title,   
 (SELECT service_address FROM services WHERE services.id = instances.service_id)
   AS service_address   
FROM instances;

SELECT projects.id AS "Project Id", 
       projects.title AS "Project Title", 
       instances.instance_name AS "Instance Name",
       services.service_title AS "Service Title",
       services.service_address AS "Service Address"
    FROM projects 
        JOIN instances
            ON projects.id = instances.project_id
        JOIN services
            ON services.id = instances.service_id;

```

<p>2. Предложить, на каких полях можно создать индексы для ускорения запросов из п. 1. Создать требуемые индексы (не более трёх). Команды на создание индексов описать в файле schema.sql (редактировать файл из прошлого урока).
<p>Ответ:

```
projects=# create index services_index on services (service_title, service_address, service_techuser, service_passwd);
CREATE INDEX
projects=# create index instances_index on instances (instance_name);
CREATE INDEX
projects=# create index projects_index on projects (title);
CREATE INDEX
projects=# 

```

<p>3. Для каждого индекса из п. 2 показать анализ запроса до/после добавления индекса, оценить занимаемый индексом объём диска. Отчёт представить в файле analysis.txt.
<p>Ответ:

```
Размеры индексов:
---
projects=#  select pg_relation_size('services_index');
 pg_relation_size 
------------------
            16384
(1 row)

projects=#  select pg_relation_size('instances_index');
 pg_relation_size 
------------------
            16384
(1 row)

projects=#  select pg_relation_size('projects_index');
 pg_relation_size 
------------------
            16384
(1 row)

projects=# 

Информация об индексе:
projects=# \di+ services_index
                                     List of relations
 Schema |      Name      | Type  |  Owner   |  Table   | Persistence | Size  | Description 
--------+----------------+-------+----------+----------+-------------+-------+-------------
 public | services_index | index | postgres | services | permanent   | 16 kB | 
(1 row)

projects=# \di+ projects_index
                                     List of relations
 Schema |      Name      | Type  |  Owner   |  Table   | Persistence | Size  | Description 
--------+----------------+-------+----------+----------+-------------+-------+-------------
 public | projects_index | index | postgres | projects | permanent   | 16 kB | 
(1 row)

projects=# \di+ instances_index
                                      List of relations
 Schema |      Name       | Type  |  Owner   |   Table   | Persistence | Size  | Description 
--------+-----------------+-------+----------+-----------+-------------+-------+-------------
 public | instances_index | index | postgres | instances | permanent   | 16 kB | 
(1 row)

---
Делаю explain analize запросов.
projects=# explain analyze SELECT instance_name, instance_description FROM instances WHERE instance_name LIKE '%jenkins' ORDER BY instance_name;
                                                QUERY PLAN                                                
----------------------------------------------------------------------------------------------------------
 Sort  (cost=1.24..1.24 rows=1 width=24) (actual time=0.063..0.065 rows=3 loops=1)
   Sort Key: instance_name
   Sort Method: quicksort  Memory: 25kB
   ->  Seq Scan on instances  (cost=0.00..1.23 rows=1 width=24) (actual time=0.024..0.045 rows=3 loops=1)
         Filter: (instance_name ~~ '%jenkins'::text)
         Rows Removed by Filter: 15
 Planning Time: 0.357 ms
 Execution Time: 0.098 ms
(8 rows)

projects=# 
projects=# 
projects=# explain analyze SELECT 
projects-#  (SELECT title FROM projects WHERE projects.id = instances.project_id)
projects-#    AS project_title,
projects-#  instance_name,
projects-#  (SELECT owner_email FROM projects WHERE projects.id = instances.project_id)
projects-#    AS owner_email,
projects-#  (SELECT service_title FROM services WHERE services.id = instances.service_id)
projects-#    AS service_title,   
projects-#  (SELECT service_address FROM services WHERE services.id = instances.service_id)
projects-#    AS service_address   
projects-# FROM instances;
                                                      QUERY PLAN                                                       
-----------------------------------------------------------------------------------------------------------------------
 Seq Scan on instances  (cost=0.00..77.23 rows=18 width=138) (actual time=0.046..0.559 rows=18 loops=1)
   SubPlan 1
     ->  Seq Scan on projects  (cost=0.00..1.04 rows=1 width=13) (actual time=0.003..0.003 rows=1 loops=18)
           Filter: (id = instances.project_id)
           Rows Removed by Filter: 2
   SubPlan 2
     ->  Seq Scan on projects projects_1  (cost=0.00..1.04 rows=1 width=19) (actual time=0.005..0.005 rows=1 loops=18)
           Filter: (id = instances.project_id)
           Rows Removed by Filter: 2
   SubPlan 3
     ->  Seq Scan on services  (cost=0.00..1.07 rows=1 width=8) (actual time=0.004..0.005 rows=1 loops=18)
           Filter: (id = instances.service_id)
           Rows Removed by Filter: 5
   SubPlan 4
     ->  Seq Scan on services services_1  (cost=0.00..1.07 rows=1 width=20) (actual time=0.003..0.004 rows=1 loops=18)
           Filter: (id = instances.service_id)
           Rows Removed by Filter: 5
 Planning Time: 0.821 ms
 Execution Time: 0.598 ms
(19 rows)

projects=# 
projects=# explain analyze SELECT projects.id AS "Project Id", 
projects-#        projects.title AS "Project Title", 
projects-#        instances.instance_name AS "Instance Name",
projects-#        services.service_title AS "Service Title",
projects-#        services.service_address AS "Service Address"
projects-#     FROM projects 
projects-#         JOIN instances
projects-#             ON projects.id = instances.project_id
projects-#         JOIN services
projects-#             ON services.id = instances.service_id;
                                                     QUERY PLAN                                                      
---------------------------------------------------------------------------------------------------------------------
 Hash Join  (cost=2.20..3.58 rows=18 width=59) (actual time=0.046..0.080 rows=18 loops=1)
   Hash Cond: (instances.service_id = services.id)
   ->  Hash Join  (cost=1.07..2.36 rows=18 width=39) (actual time=0.022..0.043 rows=18 loops=1)
         Hash Cond: (instances.project_id = projects.id)
         ->  Seq Scan on instances  (cost=0.00..1.18 rows=18 width=26) (actual time=0.006..0.010 rows=18 loops=1)
         ->  Hash  (cost=1.03..1.03 rows=3 width=21) (actual time=0.008..0.009 rows=3 loops=1)
               Buckets: 1024  Batches: 1  Memory Usage: 9kB
               ->  Seq Scan on projects  (cost=0.00..1.03 rows=3 width=21) (actual time=0.004..0.006 rows=3 loops=1)
   ->  Hash  (cost=1.06..1.06 rows=6 width=36) (actual time=0.018..0.019 rows=6 loops=1)
         Buckets: 1024  Batches: 1  Memory Usage: 9kB
         ->  Seq Scan on services  (cost=0.00..1.06 rows=6 width=36) (actual time=0.010..0.012 rows=6 loops=1)
 Planning Time: 0.493 ms
 Execution Time: 0.119 ms
(13 rows)

projects=# 


```
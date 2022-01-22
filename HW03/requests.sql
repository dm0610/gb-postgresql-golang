explain analyze SELECT instance_name, instance_description FROM instances WHERE instance_name LIKE '%jenkins' ORDER BY instance_name;


explain analyze SELECT 
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

explain analyze SELECT projects.id AS "Project Id", 
       projects.title AS "Project Title", 
       instances.instance_name AS "Instance Name",
       services.service_title AS "Service Title",
       services.service_address AS "Service Address"
    FROM projects 
        JOIN instances
            ON projects.id = instances.project_id
        JOIN services
            ON services.id = instances.service_id;

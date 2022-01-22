INSERT INTO projects (
    title, 
    owner_email) VALUES 
('Customers Feedback', 'a.ivanov1@mail.ru'),
('New Cloud', 's.atremonov2@mail.ru'),
('Shared Data', 'n.semenov3@mail.ru');

INSERT INTO services (
    service_title, 
    service_address, 
    service_techuser, 
    service_passwd) VALUES 
('jenkins', 'jenkins.mycompany.ru', 'jenkins-tech-user', 'q23Fsdf79wsew3g'),
('harbor', 'harbor.mycompany.ru', 'harbor-tech-user', 'q23Fsd239wsew3g'),
('nexus', 'nexus.mycompany.ru', 'nexus-tech-user', '2q3Fsdf79wsew3g'),
('gitlab', 'gitlab.mycompany.ru', 'gitlab-tech-user', 'z3Gs3df79wsew3g'),
('sonarqube', 'sonarqube.mycompany.ru', 'sonarqube-tech-user', 'z42Asdf79wsew3g'),
('kubernetes', 'k8s.mycompany.ru', 'k8s-tech-user', 'Aqwe3Fsdf79wsew3g');


INSERT INTO instances (
    project_id,
    service_id, 
    instance_name,
    instance_description,
    instance_expire_time) VALUES 
(1, 1, 'cf_jenkins', 'Description 1', '-1'),
(1, 2, 'cf_harbor', 'Description 2', '-1'),
(1, 3, 'cf_nexus', 'Description 3', '-1'),
(1, 4, 'cf_gitlab', 'Description 4', '-1'),
(1, 5, 'cf_sonarqube', 'Description 5', '-1'),
(1, 6, 'cf_k8s', 'Description 6', '-1');

INSERT INTO instances (
    project_id,
    service_id, 
    instance_name, 
    instance_description,    
    instance_expire_time) VALUES 
(2, 1, 'nc_jenkins', 'Description 7', '-1'),
(2, 2, 'nc_harbor', 'Description 8', '-1'),
(2, 3, 'nc_nexus', 'Description 9', '-1'),
(2, 4, 'nc_gitlab', 'Description 10', '-1'),
(2, 5, 'nc_sonarqube', 'Description 11', '-1'),
(2, 6, 'nc_k8s', 'Description 12', '-1');

INSERT INTO instances (
    project_id,
    service_id, 
    instance_name, 
    instance_description,    
    instance_expire_time) VALUES 
(3, 1, 'sd_jenkins', 'Description 13', '-1'),
(3, 2, 'sd_harbor', 'Description 14', '-1'),
(3, 3, 'sd_nexus', 'Description 15', '-1'),
(3, 4, 'sd_gitlab', 'Description 16', '-1'),
(3, 5, 'sd_sonarqube', 'Description 17', '-1'),
(3, 6, 'sd_k8s', 'Description 18', '-1');


INSERT INTO customers(name, phone)
VALUES ('Salom1', '+992000000001'),
       ('Salom2', '+992000000002'),
       ('Salom3', '+992000000003'),
       ('Salom4', '+992000000004'),
       ('Salom5', '+992000000005'),
       ('Salom6', '+992000000006'),
       ('Salom7', '+992000000007'),
       ('Salom8', '+992000000008'),
       ('Salom9', '+992000000009'),
       ('Salom10', '+992000000010')
ON CONFLICT DO NOTHING
RETURNING id;
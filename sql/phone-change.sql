UPDATE customers
set phone='+992000000011'
where id = 11 RETURNING  id, name, active;
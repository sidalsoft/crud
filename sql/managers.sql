SELECT id,
       name,
       salary*1000 as salary,
       plan*1000 as plan,
       COALESCE((SELECT sum(price * qty)
                 FROM sale_positions,
                      sales
                 WHERE sale_id = sales.id
                   and sales.manager_id = m.id), 0) as total
from managers as m
order by total desc;
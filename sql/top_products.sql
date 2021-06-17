SELECT id,
       name,
       (SELECT sum(price * qty)
        FROM sale_positions
        WHERE product_id = p.id) as total
FROM products as p
where name != 'Tea'
ORDER BY total desc
limit 3;
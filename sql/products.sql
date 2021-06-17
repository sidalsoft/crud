SELECT id, name, qty
FROM products
WHERE active = true
ORDER BY qty desc
LIMIT 10;
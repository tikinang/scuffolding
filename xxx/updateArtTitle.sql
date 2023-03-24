UPDATE art AS a JOIN customer AS c ON a.customer_id = c.id SET a.title = CONCAT(c.name, ': ', a.id);

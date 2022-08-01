SELECT *
FROM order_assets oa 
INNER JOIN matched_orders mo 
ON oa.order_id = mo.order_id 
WHERE id = ?
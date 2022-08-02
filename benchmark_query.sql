SELECT *
FROM order_assets oa 
INNER JOIN matched_orders mo 
ON oa.order_id = mo.order_id 
WHERE address = '0x8068a2c7735060589ab03685e220b322b5ec9a71'
and id = 6919056368701220851

id = 1443635317331776148
address = 0x70bd60f625f6dd082ae1f59b80dc78cfa8b47f18
INSERT INTO permissions (id, code)
VALUES
  (1, 'listing:create'),
  (2, 'listing:delete'),
  (3, 'listing:edit'),
  (4, 'listing:moderate');

SELECT setval(pg_get_serial_sequence('permissions', 'id'), (SELECT MAX(id) FROM permissions));

INSERT INTO roles (id, name)
VALUES
  (1, 'Administrator'),
  (2, 'User');

SELECT setval(pg_get_serial_sequence('roles', 'id'), (SELECT MAX(id) FROM roles));

INSERT INTO roles_permissions (role_id, permission_id)
SELECT roles.id, permissions.id
FROM roles
CROSS JOIN permissions
WHERE name = 'Administrator'
UNION
SELECT roles.id, permissions.id
FROM roles
CROSS JOIN (
  SELECT id
  FROM permissions
  WHERE code IN ('listing:create', 'listing:delete', 'listing:edit')
) AS permissions
WHERE name = 'User';

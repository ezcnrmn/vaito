INSERT INTO categories (name)
VALUES
  ('Services'),
  ('Electronics'),
  ('Personal items'),
  ('Transport'),
  ('Pets'),
  ('Home & Garden');

INSERT INTO listing_statuses (name)
VALUES
  ('Draft'),
  ('Moderation'),
  ('Active'),
  ('Inactive');

INSERT INTO permissions (code)
VALUES
  ('listing:create'),
  ('listing:delete'),
  ('listing:edit'),
  ('listing:moderate');

INSERT INTO roles (name)
VALUES
  ('Administrator'),
  ('User');

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

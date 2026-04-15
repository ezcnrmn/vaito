INSERT INTO categories (id, name)
VALUES
  (1, 'Services'),
  (2, 'Electronics'),
  (3, 'Personal items'),
  (4, 'Transport'),
  (5, 'Pets'),
  (6, 'Home & Garden');

SELECT setval(pg_get_serial_sequence('categories', 'id'), (SELECT MAX(id) FROM categories));

INSERT INTO listing_statuses (id, name)
VALUES
  (1, 'Draft'),
  (2, 'Moderation'),
  (3, 'Active'),
  (4, 'Inactive');

SELECT setval(pg_get_serial_sequence('listing_statuses', 'id'), (SELECT MAX(id) FROM listing_statuses));

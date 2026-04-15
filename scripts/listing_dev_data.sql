INSERT INTO listings (title, description, category_id, user_id, created_at, published_at, version, status_id, price) 
VALUES 
(
  'Продам iPhone 3G',
  'Полный комплект: Зарядное устройство. Коробка. Наушники. Чехол в подарок.',
  (SELECT id FROM categories WHERE name='Electronics'),
  2,
  NOW(),
  NULL,
  1,
  (SELECT id FROM listing_statuses WHERE name='Draft'),
  1500
),
(
  'Клетка для попугая',
  'Ширина клетки 30см, высота 40см, глубина 20см',
  (SELECT id FROM categories WHERE name='Pets'),
  2,
  NOW() - INTERVAL '2 days',
  NOW(),
  3,
  (SELECT id FROM listing_statuses WHERE name='Active'),
  1000
),
(
  'Услуги курьера на авто',
  'Доставлю ваши посылки по городу в течение часа. Бережное обращение.',
  (SELECT id FROM categories WHERE name='Services'),
  2,
  NOW() - INTERVAL '1 day',
  NOW(),
  3,
  (SELECT id FROM listing_statuses WHERE name='Active'),
  500
),
(
  'Игровая приставка Sony PlayStation 5',
  'Почти новая, два джойстика в комплекте. Играли всего пару раз.',
  (SELECT id FROM categories WHERE name='Electronics'),
  1,
  NOW() - INTERVAL '5 hours',
  NOW(),
  1,
  (SELECT id FROM listing_statuses WHERE name='Draft'),
  45000
),
(
  'Кожаная куртка (размер L)',
  'Натуральная кожа, состояние отличное. Носил один сезон.',
  (SELECT id FROM categories WHERE name='Personal items'),
  2,
  NOW() - INTERVAL '3 days',
  NOW(),
  3,
  (SELECT id FROM listing_statuses WHERE name='Active'),
  7000
),
(
  'Велосипед горный',
  '21 скорость, дисковые тормоза. Требуется небольшая настройка переключателей.',
  (SELECT id FROM categories WHERE name='Transport'),
  2,
  NOW() - INTERVAL '10 days',
  NULL,
  4,
  (SELECT id FROM listing_statuses WHERE name='Inactive'),
  12000
),
(
  'Котенок мейн-кун',
  'Мальчик, 3 месяца. Привит, приучен к лотку. Очень ласковый.',
  (SELECT id FROM categories WHERE name='Pets'),
  1,
  NOW() - INTERVAL '12 hours',
  NOW(),
  3,
  (SELECT id FROM listing_statuses WHERE name='Active'),
  15000
),
(
  'Набор садовых инструментов',
  'Лопата, грабли, секатор. Все из нержавеющей стали.',
  (SELECT id FROM categories WHERE name='Home & Garden'),
  2,
  NOW() - INTERVAL '4 days',
  NOW(),
  3,
  (SELECT id FROM listing_statuses WHERE name='Active'),
  3500
),
(
  'Репетитор по математике',
  'Подготовка к экзаменам, разбор сложных тем. Опыт работы 5 лет.',
  (SELECT id FROM categories WHERE name='Services'),
  2,
  NOW() - INTERVAL '1 hour',
  NOW(),
  3,
  (SELECT id FROM listing_statuses WHERE name='Active'),
  1000
),
(
  'Ноутбук для учебы',
  'Компактный, легкий, заряд держит около 5 часов. Состояние на 4.',
  (SELECT id FROM categories WHERE name='Electronics'),
  2,
  NOW() - INTERVAL '6 days',
  NULL,
  1,
  (SELECT id FROM listing_statuses WHERE name='Draft'),
  20000
),
(
  'Зимние сапоги 38 размер',
  'Теплые, на натуральном меху. Новые, не подошел размер.',
  (SELECT id FROM categories WHERE name='Personal items'),
  2,
  NOW() - INTERVAL '2 days',
  NOW(),
  5,
  (SELECT id FROM listing_statuses WHERE name='Active'),
  4000
),
(
  'Диван угловой',
  'Большой, раскладной. Самовывоз. Есть небольшие потертости на подлокотнике.',
  (SELECT id FROM categories WHERE name='Home & Garden'),
  1,
  NOW() - INTERVAL '8 days',
  NOW(),
  3,
  (SELECT id FROM listing_statuses WHERE name='Active'),
  8000
);

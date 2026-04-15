DELETE FROM users WHERE email IN ('admin@test.com', 'user@test.com');

-- Пароль для всех тестовых пользователей: 12345678
INSERT INTO users (id, name, email, password_hash, role_id, created_at, version) 
VALUES 
(1, 'Admin', 'admin@test.com', '\x24326124313224744f486d706559746a726a774c456c6468357231644f626c62692f68416d6664526d544e386c6e386536783158384f334935423843', (SELECT id FROM roles WHERE name='Administrator'), NOW(), 1),
(2, 'User', 'user@test.com', '\x24326124313224744f486d706559746a726a774c456c6468357231644f626c62692f68416d6664526d544e386c6e386536783158384f334935423843', (SELECT id FROM roles WHERE name='User'), NOW(), 1);


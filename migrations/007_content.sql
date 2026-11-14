ALTER TABLE services ADD COLUMN description TEXT NOT NULL DEFAULT '';
ALTER TABLE barbers ADD COLUMN title TEXT NOT NULL DEFAULT '';
ALTER TABLE barbers ADD COLUMN photo TEXT NOT NULL DEFAULT '';

UPDATE services SET name = 'Мужская стрижка' WHERE name = 'Стрижка';
UPDATE services SET name = 'Моделирование бороды' WHERE name = 'Борода';

INSERT INTO services (name, price, duration_min, description) VALUES
    ('Мужская стрижка', 1800, 60, 'Консультация, мытьё головы, стрижка и укладка'),
    ('Стрижка машинкой', 1200, 30, 'Одна или несколько насадок, окантовка'),
    ('Моделирование бороды', 1000, 30, 'Форма, контур, уход маслом'),
    ('Стрижка и борода', 2500, 90, 'Полный образ за один визит'),
    ('Королевское бритьё', 1500, 45, 'Опасная бритва, горячие полотенца'),
    ('Детская стрижка', 1200, 45, 'Для детей до 12 лет')
ON CONFLICT (name) DO UPDATE SET
    price = EXCLUDED.price,
    duration_min = EXCLUDED.duration_min,
    description = EXCLUDED.description;

INSERT INTO barbers (name, title) VALUES
    ('Иван', 'Старший барбер, 8 лет опыта'),
    ('Пётр', 'Барбер, специалист по бороде'),
    ('Артём', 'Барбер-стилист')
ON CONFLICT (name) DO UPDATE SET title = EXCLUDED.title;

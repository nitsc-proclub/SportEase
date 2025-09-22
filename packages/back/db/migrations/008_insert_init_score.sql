INSERT INTO team_points (class_id, init_score) VALUES
(11, 20),
(12, 95),
(13, 20),
(21, 90),
(22, 160),
(23, 100),
(31, 205),
(32, 165),
(33, 35),
(41, 240),
(42, 70),
(43, 140),
(51, 15),
(52, 85),
(53, 10),
(6, 200)
ON DUPLICATE KEY UPDATE
init_score = VALUES(init_score);

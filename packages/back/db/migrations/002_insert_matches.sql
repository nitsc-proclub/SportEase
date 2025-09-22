-- データベースの初期化（必要に応じてコメントアウトを解除）

-- ----------------------------------------------------------------
-- 1. トーナメント情報の登録
-- ----------------------------------------------------------------
INSERT INTO tournaments (name, sport, weather_condition) VALUES
('バレーボール', 'volleyball', 'any'),
('卓球（晴天時）', 'table_tennis', 'sunny'),
('卓球（雨天時）', 'table_tennis', 'rainy'),
('卓球（雨天時・敗者戦左側）', 'table_tennis', 'rainy'),
('卓球（雨天時・敗者戦右側）', 'table_tennis', 'rainy'),
('8人制サッカー', 'soccer', 'any'),
('学年対抗リレーAブロック', 'relay', 'any'),
('学年対抗リレーBブロック', 'relay', 'any');

-- ----------------------------------------------------------------
-- 2. バレーボール: チームと試合の登録 (P.8参照)
-- ----------------------------------------------------------------
SET @volleyball_tour_id = (SELECT id FROM tournaments WHERE name = 'バレーボール');

-- チーム登録
INSERT INTO teams (class_id, name, tournament_id) VALUES
(6, '専・教', @volleyball_tour_id), (43, 'IE4', @volleyball_tour_id), (51, 'IS5', @volleyball_tour_id), (42, 'IT4', @volleyball_tour_id),
(32, 'IT3', @volleyball_tour_id), (22, 'IT2', @volleyball_tour_id), (11, '1-1', @volleyball_tour_id), (23, 'IE2', @volleyball_tour_id),
(31, 'IS3', @volleyball_tour_id), (21, 'IS2', @volleyball_tour_id), (41, 'IS4', @volleyball_tour_id), (53, 'IE5', @volleyball_tour_id),
(12, '1-2', @volleyball_tour_id), (13, '1-3', @volleyball_tour_id), (33, 'IE3', @volleyball_tour_id), (52, 'IT5', @volleyball_tour_id);

-- 試合枠の作成（決勝から順に作成）
INSERT INTO matches (tournament_id, round, match_number_in_round) VALUES
(@volleyball_tour_id, 4, 16), -- 決勝
(@volleyball_tour_id, 4, 15), -- 3位決定戦
(@volleyball_tour_id, 3, 14), -- 準決勝
(@volleyball_tour_id, 3, 13), -- 準決勝
(@volleyball_tour_id, 2, 12), -- 準々決勝
(@volleyball_tour_id, 2, 11), -- 準々決勝
(@volleyball_tour_id, 2, 10), -- 準々決勝
(@volleyball_tour_id, 2, 9);  -- 準々決勝

-- 各試合枠のIDを変数に格納
SET @v_match_16_id = (SELECT id FROM matches WHERE tournament_id = @volleyball_tour_id AND match_number_in_round = 16);
SET @v_match_15_id = (SELECT id FROM matches WHERE tournament_id = @volleyball_tour_id AND match_number_in_round = 15);
SET @v_match_14_id = (SELECT id FROM matches WHERE tournament_id = @volleyball_tour_id AND match_number_in_round = 14);
SET @v_match_13_id = (SELECT id FROM matches WHERE tournament_id = @volleyball_tour_id AND match_number_in_round = 13);
SET @v_match_12_id = (SELECT id FROM matches WHERE tournament_id = @volleyball_tour_id AND match_number_in_round = 12);
SET @v_match_11_id = (SELECT id FROM matches WHERE tournament_id = @volleyball_tour_id AND match_number_in_round = 11);
SET @v_match_10_id = (SELECT id FROM matches WHERE tournament_id = @volleyball_tour_id AND match_number_in_round = 10);
SET @v_match_9_id = (SELECT id FROM matches WHERE tournament_id = @volleyball_tour_id AND match_number_in_round = 9);

-- 準決勝の勝者は決勝へ、敗者は3位決定戦へ進む (ここでは勝者の進路のみ設定)
UPDATE matches SET next_match_id = @v_match_16_id WHERE id IN (@v_match_14_id, @v_match_13_id);

-- 準々決勝の勝者は準決勝へ
UPDATE matches SET next_match_id = @v_match_14_id WHERE id IN (@v_match_12_id, @v_match_11_id);
UPDATE matches SET next_match_id = @v_match_13_id WHERE id IN (@v_match_10_id, @v_match_9_id);

-- 1回戦の試合を登録（チームと次の試合IDを設定）
INSERT INTO matches (tournament_id, round, match_number_in_round, team1_id, team2_id, next_match_id) VALUES
(@volleyball_tour_id, 1, 1, (SELECT id FROM teams WHERE name = '専・教' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IE4' AND tournament_id = @volleyball_tour_id), @v_match_9_id),
(@volleyball_tour_id, 1, 2, (SELECT id FROM teams WHERE name = 'IS5' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IT4' AND tournament_id = @volleyball_tour_id), @v_match_9_id),
(@volleyball_tour_id, 1, 3, (SELECT id FROM teams WHERE name = 'IT3' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IT2' AND tournament_id = @volleyball_tour_id), @v_match_10_id),
(@volleyball_tour_id, 1, 4, (SELECT id FROM teams WHERE name = '1-1' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IE2' AND tournament_id = @volleyball_tour_id), @v_match_10_id),
(@volleyball_tour_id, 1, 5, (SELECT id FROM teams WHERE name = 'IS3' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IS2' AND tournament_id = @volleyball_tour_id), @v_match_11_id),
(@volleyball_tour_id, 1, 6, (SELECT id FROM teams WHERE name = 'IS4' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IE5' AND tournament_id = @volleyball_tour_id), @v_match_11_id),
(@volleyball_tour_id, 1, 7, (SELECT id FROM teams WHERE name = '1-2' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = '1-3' AND tournament_id = @volleyball_tour_id), @v_match_12_id),
(@volleyball_tour_id, 1, 8, (SELECT id FROM teams WHERE name = 'IE3' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IT5' AND tournament_id = @volleyball_tour_id), @v_match_12_id);


-- ----------------------------------------------------------------
-- 3. 卓球（晴天時）: チームと試合の登録 (P.13参照)
-- ----------------------------------------------------------------
SET @tt_sunny_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（晴天時）');

INSERT INTO teams (class_id, name, tournament_id) VALUES
(12, '1-2', @tt_sunny_tour_id), (53, 'IE5', @tt_sunny_tour_id), (13, '1-3', @tt_sunny_tour_id), (43, 'IE4', @tt_sunny_tour_id),
(51, 'IS5', @tt_sunny_tour_id), (23, 'IE2', @tt_sunny_tour_id), (21, 'IS2', @tt_sunny_tour_id), (41, 'IS4', @tt_sunny_tour_id),
(33, 'IE3', @tt_sunny_tour_id), (22, 'IT2', @tt_sunny_tour_id), (32, 'IT3', @tt_sunny_tour_id), (42, 'IT4', @tt_sunny_tour_id),
(11, '1-1', @tt_sunny_tour_id), (31, 'IS3', @tt_sunny_tour_id), (52, 'IT5', @tt_sunny_tour_id), (6, '専・教', @tt_sunny_tour_id);

-- 試合枠の作成（決勝から順に作成）
INSERT INTO matches (tournament_id, round, match_number_in_round) VALUES
(@tt_sunny_tour_id, 4, 16), -- 決勝
(@tt_sunny_tour_id, 4, 15), -- 3位決定戦
(@tt_sunny_tour_id, 3, 14), -- 準決勝
(@tt_sunny_tour_id, 3, 13), -- 準決勝
(@tt_sunny_tour_id, 2, 12), -- 準々決勝
(@tt_sunny_tour_id, 2, 11), -- 準々決勝
(@tt_sunny_tour_id, 2, 10), -- 準々決勝
(@tt_sunny_tour_id, 2, 9);  -- 準々決勝

-- 各試合枠のIDを変数に格納
SET @tts_match_16_id = (SELECT id FROM matches WHERE tournament_id = @tt_sunny_tour_id AND match_number_in_round = 16);
SET @tts_match_15_id = (SELECT id FROM matches WHERE tournament_id = @tt_sunny_tour_id AND match_number_in_round = 15);
SET @tts_match_14_id = (SELECT id FROM matches WHERE tournament_id = @tt_sunny_tour_id AND match_number_in_round = 14);
SET @tts_match_13_id = (SELECT id FROM matches WHERE tournament_id = @tt_sunny_tour_id AND match_number_in_round = 13);
SET @tts_match_12_id = (SELECT id FROM matches WHERE tournament_id = @tt_sunny_tour_id AND match_number_in_round = 12);
SET @tts_match_11_id = (SELECT id FROM matches WHERE tournament_id = @tt_sunny_tour_id AND match_number_in_round = 11);
SET @tts_match_10_id = (SELECT id FROM matches WHERE tournament_id = @tt_sunny_tour_id AND match_number_in_round = 10);
SET @tts_match_9_id = (SELECT id FROM matches WHERE tournament_id = @tt_sunny_tour_id AND match_number_in_round = 9);

-- 準決勝の勝者は決勝へ
UPDATE matches SET next_match_id = @tts_match_16_id WHERE id IN (@tts_match_14_id, @tts_match_13_id);
-- 準々決勝の勝者は準決勝へ
UPDATE matches SET next_match_id = @tts_match_14_id WHERE id IN (@tts_match_12_id, @tts_match_11_id);
UPDATE matches SET next_match_id = @tts_match_13_id WHERE id IN (@tts_match_10_id, @tts_match_9_id);


INSERT INTO matches (tournament_id, round, match_number_in_round, team1_id, team2_id, next_match_id) VALUES
(@tt_sunny_tour_id, 1, 1, (SELECT id FROM teams WHERE name = '1-2' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IE5' AND tournament_id = @tt_sunny_tour_id), @tts_match_9_id),
(@tt_sunny_tour_id, 1, 2, (SELECT id FROM teams WHERE name = '1-3' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IE4' AND tournament_id = @tt_sunny_tour_id), @tts_match_9_id),
(@tt_sunny_tour_id, 1, 3, (SELECT id FROM teams WHERE name = 'IS5' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IE2' AND tournament_id = @tt_sunny_tour_id), @tts_match_10_id),
(@tt_sunny_tour_id, 1, 4, (SELECT id FROM teams WHERE name = 'IS2' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IS4' AND tournament_id = @tt_sunny_tour_id), @tts_match_10_id),
(@tt_sunny_tour_id, 1, 5, (SELECT id FROM teams WHERE name = 'IE3' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IT2' AND tournament_id = @tt_sunny_tour_id), @tts_match_11_id),
(@tt_sunny_tour_id, 1, 6, (SELECT id FROM teams WHERE name = 'IT3' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IT4' AND tournament_id = @tt_sunny_tour_id), @tts_match_11_id),
(@tt_sunny_tour_id, 1, 7, (SELECT id FROM teams WHERE name = '1-1' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IS3' AND tournament_id = @tt_sunny_tour_id), @tts_match_12_id),
(@tt_sunny_tour_id, 1, 8, (SELECT id FROM teams WHERE name = 'IT5' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = '専・教' AND tournament_id = @tt_sunny_tour_id), @tts_match_12_id);


-- ----------------------------------------------------------------
-- 4. 卓球（雨天時）: チームと試合の登録 (P.15参照)
-- ----------------------------------------------------------------
SET @tt_rainy_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（雨天時）');

INSERT INTO teams (class_id, name, tournament_id) VALUES
(12, '1-2', @tt_rainy_tour_id), (53, 'IE5', @tt_rainy_tour_id), (13, '1-3', @tt_rainy_tour_id), (43, 'IE4', @tt_rainy_tour_id),
(51, 'IS5', @tt_rainy_tour_id), (23, 'IE2', @tt_rainy_tour_id), (21, 'IS2', @tt_rainy_tour_id), (41, 'IS4', @tt_rainy_tour_id),
(33, 'IE3', @tt_rainy_tour_id), (22, 'IT2', @tt_rainy_tour_id), (32, 'IT3', @tt_rainy_tour_id), (42, 'IT4', @tt_rainy_tour_id),
(11, '1-1', @tt_rainy_tour_id), (31, 'IS3', @tt_rainy_tour_id), (52, 'IT5', @tt_rainy_tour_id), (6, '専・教', @tt_rainy_tour_id);

-- 試合枠の作成（決勝から順に作成）
INSERT INTO matches (tournament_id, round, match_number_in_round) VALUES
(@tt_rainy_tour_id, 4, 20), -- 決勝
(@tt_rainy_tour_id, 4, 19), -- 3位決定戦
(@tt_rainy_tour_id, 3, 18), -- 準決勝
(@tt_rainy_tour_id, 3, 17), -- 準決勝
(@tt_rainy_tour_id, 2, 12), -- 準々決勝
(@tt_rainy_tour_id, 2, 11), -- 準々決勝
(@tt_rainy_tour_id, 2, 10), -- 準々決勝
(@tt_rainy_tour_id, 2, 9);  -- 準々決勝

-- 各試合枠のIDを変数に格納
SET @ttr_match_20_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 20);
SET @ttr_match_19_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 19);
SET @ttr_match_18_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 18);
SET @ttr_match_17_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 17);
SET @ttr_match_12_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 12);
SET @ttr_match_11_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 11);
SET @ttr_match_10_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 10);
SET @ttr_match_9_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 9);

-- 準決勝の勝者は決勝へ
UPDATE matches SET next_match_id = @ttr_match_20_id WHERE id IN (@ttr_match_18_id, @ttr_match_17_id);
-- 準々決勝の勝者は準決勝へ
UPDATE matches SET next_match_id = @ttr_match_18_id WHERE id IN (@ttr_match_12_id, @ttr_match_11_id);
UPDATE matches SET next_match_id = @ttr_match_17_id WHERE id IN (@ttr_match_10_id, @ttr_match_9_id);

-- 1回戦の試合を登録（要項のトーナメント表に合わせて組み合わせを修正）
INSERT INTO matches (tournament_id, round, match_number_in_round, team1_id, team2_id, next_match_id) VALUES
(@tt_rainy_tour_id, 1, 1, (SELECT id FROM teams WHERE name = '1-2' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IE5' AND tournament_id = @tt_rainy_tour_id), @ttr_match_9_id),
(@tt_rainy_tour_id, 1, 2, (SELECT id FROM teams WHERE name = '1-3' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IE4' AND tournament_id = @tt_rainy_tour_id), @ttr_match_9_id),
(@tt_rainy_tour_id, 1, 3, (SELECT id FROM teams WHERE name = 'IS5' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IE2' AND tournament_id = @tt_rainy_tour_id), @ttr_match_10_id),
(@tt_rainy_tour_id, 1, 4, (SELECT id FROM teams WHERE name = 'IS2' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IS4' AND tournament_id = @tt_rainy_tour_id), @ttr_match_10_id),
(@tt_rainy_tour_id, 1, 5, (SELECT id FROM teams WHERE name = 'IE3' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IT2' AND tournament_id = @tt_rainy_tour_id), @ttr_match_11_id),
(@tt_rainy_tour_id, 1, 6, (SELECT id FROM teams WHERE name = 'IT3' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IT4' AND tournament_id = @tt_rainy_tour_id), @ttr_match_11_id),
(@tt_rainy_tour_id, 1, 7, (SELECT id FROM teams WHERE name = '1-1' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IS3' AND tournament_id = @tt_rainy_tour_id), @ttr_match_12_id),
(@tt_rainy_tour_id, 1, 8, (SELECT id FROM teams WHERE name = 'IT5' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = '専・教' AND tournament_id = @tt_rainy_tour_id), @ttr_match_12_id);

-- 敗者戦1 (左側ブロック)
SET @tt_rainy_loser1_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（雨天時・敗者戦左側）');
INSERT INTO matches (tournament_id, round, match_number_in_round) VALUES
(@tt_rainy_loser1_tour_id, 2, 15), -- 敗者戦1 決勝
(@tt_rainy_loser1_tour_id, 1, 13), -- 敗者戦1 1回戦(第1試合)
(@tt_rainy_loser1_tour_id, 1, 14); -- 敗者戦1 1回戦(第2試合)

SET @ttrl1_match_15_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_loser1_tour_id AND match_number_in_round = 15);
SET @ttrl1_match_13_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_loser1_tour_id AND match_number_in_round = 13);
SET @ttrl1_match_14_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_loser1_tour_id AND match_number_in_round = 14);
UPDATE matches SET next_match_id = @ttrl1_match_15_id WHERE id IN (@ttrl1_match_13_id, @ttrl1_match_14_id);

-- 敗者戦2 (右側ブロック)
SET @tt_rainy_loser2_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（雨天時・敗者戦右側）');
INSERT INTO matches (tournament_id, round, match_number_in_round) VALUES
(@tt_rainy_loser2_tour_id, 2, 16), -- 敗者戦2 決勝
(@tt_rainy_loser2_tour_id, 1, 13), -- 敗者戦2 1回戦(第1試合)
(@tt_rainy_loser2_tour_id, 1, 14); -- 敗者戦2 1回戦(第2試合)

SET @ttrl2_match_16_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_loser2_tour_id AND match_number_in_round = 16);
SET @ttrl2_match_13_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_loser2_tour_id AND match_number_in_round = 13);
SET @ttrl2_match_14_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_loser2_tour_id AND match_number_in_round = 14);
UPDATE matches SET next_match_id = @ttrl2_match_16_id WHERE id IN (@ttrl2_match_13_id, @ttrl2_match_14_id);

-- ----------------------------------------------------------------
-- 5. 8人制サッカー: チームと試合の登録 (P.21参照)
-- ----------------------------------------------------------------
SET @soccer_tour_id = (SELECT id FROM tournaments WHERE name = '8人制サッカー');

INSERT INTO teams (class_id, name, tournament_id) VALUES
(31, 'IS3', @soccer_tour_id), (23, 'IE2', @soccer_tour_id), (11, '1-1', @soccer_tour_id), (21, 'IS2', @soccer_tour_id),
(41, 'IS4', @soccer_tour_id), (52, 'IT5', @soccer_tour_id), (51, 'IS5', @soccer_tour_id), (6, '専・教', @soccer_tour_id),
(12, '1-2', @soccer_tour_id), (13, '1-3', @soccer_tour_id), (33, 'IE3', @soccer_tour_id), (42, 'IT4', @soccer_tour_id),
(32, 'IT3', @soccer_tour_id), (43, 'IE4', @soccer_tour_id), (22, 'IT2', @soccer_tour_id), (53, 'IE5', @soccer_tour_id);

INSERT INTO matches (tournament_id, round, match_number_in_round) VALUES
(@soccer_tour_id, 4, 16), -- 決勝
(@soccer_tour_id, 4, 15), -- 3位決定戦
(@soccer_tour_id, 3, 14), -- 準決勝
(@soccer_tour_id, 3, 13), -- 準決勝
(@soccer_tour_id, 2, 12), -- 準々決勝
(@soccer_tour_id, 2, 11), -- 準々決勝
(@soccer_tour_id, 2, 10), -- 準々決勝
(@soccer_tour_id, 2, 9);  -- 準々決勝

SET @s_match_16_id = (SELECT id FROM matches WHERE tournament_id = @soccer_tour_id AND match_number_in_round = 16);
SET @s_match_15_id = (SELECT id FROM matches WHERE tournament_id = @soccer_tour_id AND match_number_in_round = 15);
SET @s_match_14_id = (SELECT id FROM matches WHERE tournament_id = @soccer_tour_id AND match_number_in_round = 14);
SET @s_match_13_id = (SELECT id FROM matches WHERE tournament_id = @soccer_tour_id AND match_number_in_round = 13);
SET @s_match_12_id = (SELECT id FROM matches WHERE tournament_id = @soccer_tour_id AND match_number_in_round = 12);
SET @s_match_11_id = (SELECT id FROM matches WHERE tournament_id = @soccer_tour_id AND match_number_in_round = 11);
SET @s_match_10_id = (SELECT id FROM matches WHERE tournament_id = @soccer_tour_id AND match_number_in_round = 10);
SET @s_match_9_id = (SELECT id FROM matches WHERE tournament_id = @soccer_tour_id AND match_number_in_round = 9);

UPDATE matches SET next_match_id = @s_match_16_id WHERE id IN (@s_match_14_id, @s_match_13_id);
UPDATE matches SET next_match_id = @s_match_14_id WHERE id IN (@s_match_12_id, @s_match_11_id);
UPDATE matches SET next_match_id = @s_match_13_id WHERE id IN (@s_match_10_id, @s_match_9_id);

INSERT INTO matches (tournament_id, round, match_number_in_round, team1_id, team2_id, next_match_id) VALUES
(@soccer_tour_id, 1, 1, (SELECT id FROM teams WHERE name = 'IS3' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IE2' AND tournament_id = @soccer_tour_id), @s_match_9_id),
(@soccer_tour_id, 1, 2, (SELECT id FROM teams WHERE name = '1-1' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IS2' AND tournament_id = @soccer_tour_id), @s_match_9_id),
(@soccer_tour_id, 1, 3, (SELECT id FROM teams WHERE name = 'IS4' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IT5' AND tournament_id = @soccer_tour_id), @s_match_10_id),
(@soccer_tour_id, 1, 4, (SELECT id FROM teams WHERE name = 'IS5' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = '専・教' AND tournament_id = @soccer_tour_id), @s_match_10_id),
(@soccer_tour_id, 1, 5, (SELECT id FROM teams WHERE name = '1-2' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = '1-3' AND tournament_id = @soccer_tour_id), @s_match_11_id),
(@soccer_tour_id, 1, 6, (SELECT id FROM teams WHERE name = 'IE3' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IT4' AND tournament_id = @soccer_tour_id), @s_match_11_id),
(@soccer_tour_id, 1, 7, (SELECT id FROM teams WHERE name = 'IT3' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IE4' AND tournament_id = @soccer_tour_id), @s_match_12_id),
(@soccer_tour_id, 1, 8, (SELECT id FROM teams WHERE name = 'IT2' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IE5' AND tournament_id = @soccer_tour_id), @s_match_12_id);
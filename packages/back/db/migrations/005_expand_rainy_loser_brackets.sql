-- 卓球（雨天時） 敗者復活戦の試合枠を各ブロック2試合（準決勝×2）+ 決勝に整備
-- 既存構成:
--   左ブロック: round2=15(決勝), round1=13(1試合のみ)
--   右ブロック: round2=16(決勝), round1=14(1試合のみ)
-- 本マイグレーションで追加:
--   左ブロック: round1=14 を追加し、13/14 の勝者を 15 へリンク
--   右ブロック: round1=13 を追加し、13/14 の勝者を 16 へリンク

SET @loser_left_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（雨天時・敗者戦左側）');
SET @loser_right_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（雨天時・敗者戦右側）');

-- 左ブロック: R1 No.14 を追加（存在しない場合のみ）
-- 存在しない場合のみ追加
INSERT INTO matches (tournament_id, round, match_number_in_round)
SELECT @loser_left_tour_id, 1, 14
WHERE NOT EXISTS (
  SELECT 1 FROM matches WHERE tournament_id = @loser_left_tour_id AND round = 1 AND match_number_in_round = 14
);
SET @left_r1_14_id = (SELECT id FROM matches WHERE tournament_id = @loser_left_tour_id AND round = 1 AND match_number_in_round = 14);

-- 左ブロック決勝(15) を取得
SET @left_final_id = (SELECT id FROM matches WHERE tournament_id = @loser_left_tour_id AND round = 2 AND match_number_in_round = 15);

-- 左ブロックの既存R1 No.13 を取得
SET @left_r1_13_id = (SELECT id FROM matches WHERE tournament_id = @loser_left_tour_id AND round = 1 AND match_number_in_round = 13);

-- 左ブロック: 13/14 の勝者が 15 へ進むように設定
UPDATE matches SET next_match_id = @left_final_id WHERE id IN (@left_r1_13_id, @left_r1_14_id);

-- 右ブロック: R1 No.13 を追加（存在しない場合のみ）
-- 存在しない場合のみ追加
INSERT INTO matches (tournament_id, round, match_number_in_round)
SELECT @loser_right_tour_id, 1, 13
WHERE NOT EXISTS (
  SELECT 1 FROM matches WHERE tournament_id = @loser_right_tour_id AND round = 1 AND match_number_in_round = 13
);
SET @right_r1_13_id = (SELECT id FROM matches WHERE tournament_id = @loser_right_tour_id AND round = 1 AND match_number_in_round = 13);

-- 右ブロック決勝(16) を取得
SET @right_final_id = (SELECT id FROM matches WHERE tournament_id = @loser_right_tour_id AND round = 2 AND match_number_in_round = 16);

-- 右ブロックの既存R1 No.14 を取得
SET @right_r1_14_id = (SELECT id FROM matches WHERE tournament_id = @loser_right_tour_id AND round = 1 AND match_number_in_round = 14);

-- 右ブロック: 13/14 の勝者が 16 へ進むように設定
UPDATE matches SET next_match_id = @right_final_id WHERE id IN (@right_r1_13_id, @right_r1_14_id);

-- 1回戦の敗者 → 敗者復活戦 へのリンクも4試合×2ブロックに分割
SET @rainy_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（雨天時）');

-- 雨天本戦 1〜4番の敗者 → 左ブロック R1:13,14
SET @r_m1 = (SELECT id FROM matches WHERE tournament_id = @rainy_tour_id AND round = 1 AND match_number_in_round = 1);
SET @r_m2 = (SELECT id FROM matches WHERE tournament_id = @rainy_tour_id AND round = 1 AND match_number_in_round = 2);
SET @r_m3 = (SELECT id FROM matches WHERE tournament_id = @rainy_tour_id AND round = 1 AND match_number_in_round = 3);
SET @r_m4 = (SELECT id FROM matches WHERE tournament_id = @rainy_tour_id AND round = 1 AND match_number_in_round = 4);

UPDATE matches SET loser_next_match_id = @left_r1_13_id WHERE id IN (@r_m1, @r_m2);
UPDATE matches SET loser_next_match_id = @left_r1_14_id WHERE id IN (@r_m3, @r_m4);

-- 雨天本戦 5〜8番の敗者 → 右ブロック R1:13,14
SET @r_m5 = (SELECT id FROM matches WHERE tournament_id = @rainy_tour_id AND round = 1 AND match_number_in_round = 5);
SET @r_m6 = (SELECT id FROM matches WHERE tournament_id = @rainy_tour_id AND round = 1 AND match_number_in_round = 6);
SET @r_m7 = (SELECT id FROM matches WHERE tournament_id = @rainy_tour_id AND round = 1 AND match_number_in_round = 7);
SET @r_m8 = (SELECT id FROM matches WHERE tournament_id = @rainy_tour_id AND round = 1 AND match_number_in_round = 8);

UPDATE matches SET loser_next_match_id = @right_r1_13_id WHERE id IN (@r_m5, @r_m6);
UPDATE matches SET loser_next_match_id = @right_r1_14_id WHERE id IN (@r_m7, @r_m8);



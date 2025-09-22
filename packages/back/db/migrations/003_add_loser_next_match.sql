-- matchesテーブルに敗者進路カラム追加
ALTER TABLE matches ADD COLUMN loser_next_match_id INT NULL AFTER next_match_id;

-- バレーボール: 準決勝の敗者は3位決定戦へ進む
SET @volleyball_tour_id = (SELECT id FROM tournaments WHERE name = 'バレーボール');
SET @v_match_15_id = (SELECT id FROM matches WHERE tournament_id = @volleyball_tour_id AND match_number_in_round = 15);
SET @v_match_14_id = (SELECT id FROM matches WHERE tournament_id = @volleyball_tour_id AND match_number_in_round = 14);
SET @v_match_13_id = (SELECT id FROM matches WHERE tournament_id = @volleyball_tour_id AND match_number_in_round = 13);
UPDATE matches SET loser_next_match_id = @v_match_15_id WHERE id IN (@v_match_14_id, @v_match_13_id);

-- 卓球（晴天時）: 準決勝の敗者は3位決定戦へ進む
SET @tt_sunny_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（晴天時）');
SET @tts_match_15_id = (SELECT id FROM matches WHERE tournament_id = @tt_sunny_tour_id AND match_number_in_round = 15);
SET @tts_match_14_id = (SELECT id FROM matches WHERE tournament_id = @tt_sunny_tour_id AND match_number_in_round = 14);
SET @tts_match_13_id = (SELECT id FROM matches WHERE tournament_id = @tt_sunny_tour_id AND match_number_in_round = 13);
UPDATE matches SET loser_next_match_id = @tts_match_15_id WHERE id IN (@tts_match_14_id, @tts_match_13_id);

-- 卓球（雨天時）: 準決勝の敗者は3位決定戦へ進む
SET @tt_rainy_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（雨天時）');
SET @ttr_match_19_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 19);
SET @ttr_match_18_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 18);
SET @ttr_match_17_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 17);
UPDATE matches SET loser_next_match_id = @ttr_match_19_id WHERE id IN (@ttr_match_18_id, @ttr_match_17_id);

-- サッカー: 準決勝の敗者は3位決定戦へ進む
SET @soccer_tour_id = (SELECT id FROM tournaments WHERE name = '8人制サッカー');
SET @s_match_15_id = (SELECT id FROM matches WHERE tournament_id = @soccer_tour_id AND match_number_in_round = 15);
SET @s_match_14_id = (SELECT id FROM matches WHERE tournament_id = @soccer_tour_id AND match_number_in_round = 14);
SET @s_match_13_id = (SELECT id FROM matches WHERE tournament_id = @soccer_tour_id AND match_number_in_round = 13);
UPDATE matches SET loser_next_match_id = @s_match_15_id WHERE id IN (@s_match_14_id, @s_match_13_id);
